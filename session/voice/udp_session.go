package voice

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

var ntpEpoch = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

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
	PacketSent()
	BytesSent(bytes int)
	ResetSentData()
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
	connectionReady chan struct{}
	discoveryReady  chan struct{}
	speakingReady   chan struct{}
	readChan        chan []byte
	writeChan       chan []byte
	errorChan       chan error

	// context for managing lifecycle
	ctx    context.Context
	cancel context.CancelFunc

	// for managing the UDP packets sent
	sentPackets int
	sentBytes   int
}

var _ UdpSession = (*udpSession)(nil)

func NewUdpSession() UdpSession {
	var u udpSession
	u.Mu = &sync.Mutex{}
	u.EventHandler = NewUdpEventHandler()
	u.connectionReady = make(chan struct{})
	u.discoveryReady = make(chan struct{})
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

	<-u.discoveryReady

	// start keeping the conn alive once we discover
	go u.keepAlive()
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

	close(u.readChan)
	close(u.writeChan)
	close(u.errorChan)

	*u = *NewUdpSession().GetSession()
	return nil
}

func (u *udpSession) Write(data []byte) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	if len(u.writeChan) < cap(u.writeChan) {
		u.writeChan <- data
	} else {
		u.errorChan <- fmt.Errorf("failed to write data to udp write channel")
	}
}

func (u *udpSession) keepAlive() {
	// send a keep alive packet every 5 seconds
	keepAlive := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-u.ctx.Done():
			keepAlive.Stop()
			u.ResetSentData()
			return
		case <-keepAlive.C:
			if u.Conn != nil {
				// send keep alive packet
				senderReport := voice.SenderReportPacket{}
				senderReport.Version = 2
				senderReport.Padding = false
				senderReport.ReceptionReportCount = 0
				senderReport.PacketType = 200
				senderReport.Length = 0
				senderReport.SSRC = uint32(u.GetConnData().SSRC)
				senderReport.NTPTimestamp = constructNTPTimestamp(time.Now().UTC())
				senderReport.RTPTimestamp = constructRTPTimestamp(senderReport.NTPTimestamp)
				senderReport.SenderPacketCount = uint32(u.sentPackets)
				senderReport.SenderOctetCount = uint32(u.sentBytes)

				packet, err := senderReport.MarshalBinary()
				if err != nil {
					u.errorChan <- fmt.Errorf("error marshalling keep alive packet: %v", err)
					keepAlive.Stop()
					break
				} else {
					_, err := u.Conn.Write(packet)
					if err != nil {
						u.errorChan <- fmt.Errorf("error sending keep alive packet: %v", err)
						keepAlive.Stop()
						break
					}
					fmt.Printf("packets sent: %d, bytes sent: %d\n", u.sentPackets, u.sentBytes)
				}
			}
		}
	}
}

func (u *udpSession) discover() error {
	return u.GetEventHandler().HandleSendDiscoveryEvent(u, voice.DiscoveryPacket{})
}

// listen will listen for incoming packets and handle them accordingly
func (u *udpSession) listen() {
	for {
		select {
		case <-u.ctx.Done():
			return
		case msg := <-u.readChan:
			// Try to decode as DiscoveryPacket
			// only need to do this once, can ignore duplicated discovery packets
			if u.discoveryReady != nil {
				discoveryPacket, err := decodeDiscoveryPacket(msg)
				if discoveryPacket != nil && err == nil {
					if err := u.EventHandler.HandleReceiveDiscoveryEvent(u, *discoveryPacket); err != nil {
						u.errorChan <- fmt.Errorf("error handling received discovery packet: %v", err)
					}
					continue
				}
			}

			// Try to decode as VoicePacket
			voicePacket, err := decodeBinaryPacket(msg)
			if voicePacket != nil && err == nil {
				// currently no implementation for receiving voice packets
				// could potentially build something to record audio, but thats kind of spyware-y
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
			if err != nil && !errors.Is(err, net.ErrClosed) {
				fmt.Println(err)
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

func (u *udpSession) PacketSent() {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.sentPackets++
}

func (u *udpSession) BytesSent(bytes int) {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.sentBytes += bytes
}

func (u *udpSession) ResetSentData() {
	u.Mu.Lock()
	defer u.Mu.Unlock()
	u.sentPackets = 0
	u.sentBytes = 0
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

func constructNTPTimestamp(t time.Time) uint64 {
	// Calculate the number of seconds since the NTP epoch
	seconds := uint32(t.Sub(ntpEpoch).Seconds())

	// Calculate the fractional part
	fraction := uint32((t.Sub(ntpEpoch).Nanoseconds() % 1e9) * (1 << 32) / 1e9)

	// Combine the integer and fractional parts into a 64-bit unsigned fixed-point number
	ntpTimestamp := uint64(seconds)<<32 | uint64(fraction)

	return ntpTimestamp
}

func constructRTPTimestamp(ntpTimestamp uint64) uint32 {
	// Extract the integer part of the NTP timestamp (first 32 bits)
	ntpSeconds := uint32(ntpTimestamp >> 32)

	// Extract the fractional part of the NTP timestamp (last 32 bits)
	ntpFraction := uint32(ntpTimestamp & 0xFFFFFFFF)

	// Convert the NTP timestamp to RTP timestamp using the clock rate
	rtpTimestamp := ntpSeconds*uint32(48000) + uint32((uint64(ntpFraction)*uint64(48000))>>32)

	return rtpTimestamp
}
