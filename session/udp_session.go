package session

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Carmen-Shannon/simple-discord/structs/voice"
)

type UdpSession interface {
	SetConnData(udpConn *voice.UdpData)
	GetConnData() *voice.UdpData
	IsConnected() bool
	Connect() error
	Exit() error
	GetSession() *udpSession
	SetSecretKey(key [32]byte)
	SetEncryptionMode(mode voice.TransportEncryptionMode)
	Write(data []byte)
	GetEventHandler() *UdpEventHandler
	SetEventHandler(handler *UdpEventHandler)
}

type udpSession struct {
	// mutex for thread safety
	Mu *sync.Mutex
	// udp connection details
	ConnData *voice.UdpData
	// udp connection
	Conn *net.UDPConn
	// secret key for encoding audio frames
	SecretKey [32]byte
	// Encryption mode
	EncryptionMode voice.TransportEncryptionMode
	// Event handler
	EventHandler *UdpEventHandler

	// channels
	discovered    chan struct{}
	speakingReady chan struct{}
	readChan      chan []byte
	writeChan     chan []byte
	errorChan     chan error

	// context for managing lifecycle
	ctx    context.Context
	cancel context.CancelFunc
}

var _ UdpSession = (*udpSession)(nil)

func NewUdpSession() UdpSession {
	var u udpSession
	u.Mu = &sync.Mutex{}
	u.EventHandler = NewUdpEventHandler()
	u.discovered = make(chan struct{})
	u.speakingReady = make(chan struct{})
	u.readChan = make(chan []byte)
	u.writeChan = make(chan []byte, 2048)
	u.errorChan = make(chan error)
	u.ctx, u.cancel = context.WithCancel(context.Background())

	return &u
}

func (u *udpSession) Connect() error {
	if u.Conn != nil {
		if err := u.Exit(); err != nil {
			return err
		}
	} else if u.GetConnData() == nil {
		return errors.New("please set up UDP connection data before connecting")
	}

	connData := u.GetConnData()
	addr := fmt.Sprintf("%s:%d", connData.Address, connData.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("error resolving UDP address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("error dialing UDP connection: %v", err)
	}
	u.Conn = conn

	go u.listen()
	go u.handleWrite()
	go u.handleRead()
	go u.handleError()

	if err := u.discover(); err != nil {
		return fmt.Errorf("error discovering external IP: %v", err)
	}

	return nil
}

func (u *udpSession) Exit() error {
	if u.Conn != nil {
		if err := u.Conn.Close(); err != nil {
			if !errors.Is(err, net.ErrClosed) {
				return fmt.Errorf("error closing voice UDP: %v", err)
			}
		}
	}

	u.cancel()
	time.Sleep(1 * time.Second) //arbitrary sleep to allow for cleanup
	close(u.readChan)
	close(u.writeChan)
	close(u.errorChan)

	u.Conn = nil
	return nil
}

func (u *udpSession) Write(data []byte) {
	if len(u.writeChan) < cap(u.writeChan) {
		u.writeChan <- data
	} else {
		u.errorChan <- fmt.Errorf("failed to write data to udp write channel")
	}
}

func (u *udpSession) discover() error {
	return u.GetEventHandler().HandleSendDiscoveryEvent(u, voice.DiscoveryPacket{})
}

func (u *udpSession) listen() {
	for {
		select {
		case <-u.ctx.Done():
			return
		case msg := <-u.readChan:
			// Try to decode as DiscoveryPacket
			discoveryPacket, err := decodeDiscoveryPacket(msg)
			if discoveryPacket != nil && err == nil {
				if err := u.EventHandler.HandleReceiveDiscoveryEvent(u, *discoveryPacket); err != nil {
					u.errorChan <- fmt.Errorf("error handling received discovery packet: %v", err)
				}
				continue
			}

			// Try to decode as VoicePacket
			voicePacket, err := decodeBinaryPacket(msg)
			if voicePacket != nil && err == nil {
				fmt.Println("HANDLING A VOICE PACKET")
				fmt.Println(voicePacket.ToString())
				continue
			}

			// If neither, send error
			if err != nil {
				u.errorChan <- fmt.Errorf("error decoding packet: %v", err)
			}
		}
	}
}

func (u *udpSession) handleWrite() {
	for {
		select {
		case <-u.ctx.Done():
			return
		case data := <-u.writeChan:
			if _, err := u.Conn.Write(data); err != nil {
				u.errorChan <- fmt.Errorf("error writing to UDP connection: %v", err)
			}
		}
	}
}

func (u *udpSession) handleRead() {
	data := make([]byte, 2048) // Allocate a single byte slice of 2048 bytes

	for {
		select {
		case <-u.ctx.Done():
			return
		default:
			n, err := u.Conn.Read(data)
			if err != nil {
				u.errorChan <- fmt.Errorf("error reading from UDP connection: %v", err)
				return
			}

			// Try to decode as DiscoveryPacket
			discoveryPacket, err := decodeDiscoveryPacket(data[:n])
			if discoveryPacket != nil && err == nil {
				u.readChan <- data[:n]
				continue
			}

			// Try to decode as VoicePacket
			voicePacket, err := decodeBinaryPacket(data[:n])
			if voicePacket != nil && err == nil {
				u.readChan <- data[:n]
				continue
			}

			// If neither, send error
			if err != nil {
				u.errorChan <- fmt.Errorf("error decoding packet: %v", err)
			}
		}
	}
}

func (u *udpSession) handleError() {
	for {
		select {
		case <-u.ctx.Done():
			return
		case err := <-u.errorChan:
			log.Printf("udp session error: %v\n", err)
		}
	}
}

func (u *udpSession) GetSession() *udpSession {
	return u
}

func (u *udpSession) SetConnData(udpConn *voice.UdpData) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.ConnData = udpConn
}

func (u *udpSession) GetConnData() *voice.UdpData {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	return u.ConnData
}

func (u *udpSession) IsConnected() bool {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	return u.Conn != nil
}

func (u *udpSession) SetSecretKey(key [32]byte) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.SecretKey = key
}

func (u *udpSession) SetEncryptionMode(mode voice.TransportEncryptionMode) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.EncryptionMode = mode
}

func (u *udpSession) GetEventHandler() *UdpEventHandler {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	return u.EventHandler
}

func (u *udpSession) SetEventHandler(handler *UdpEventHandler) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.EventHandler = handler
}

func decodeBinaryPacket(data []byte) (*voice.VoicePacket, error) {
	var packet voice.VoicePacket
	err := packet.UnmarshalBinary(data)
	if err != nil && errors.Is(err, io.EOF) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error decoding voice packet: %v", err)
	}
	return &packet, nil
}

func decodeDiscoveryPacket(data []byte) (*voice.DiscoveryPacket, error) {
	var packet voice.DiscoveryPacket
	err := packet.UnmarshalBinary(data)
	if err != nil && errors.Is(err, io.EOF) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error decoding discovery packet: %v", err)
	}
	return &packet, nil
}
