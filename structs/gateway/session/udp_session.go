package session

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway/payload"
)

var ntpEpoch = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

type UdpEventFunc func(UdpSession, payload.VoicePacket) error
type UdpDiscoveryEventFunc func(UdpSession, payload.DiscoveryPacket) error

type udpSession struct {
	Session
	mu     *sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	udpData     gateway.UdpData
	encryption  gateway.TransportEncryptionMode
	discovered  bool
	secretKey   [32]byte
	sentPackets int
	sentBytes   int

	eventHandler *udpEventHandler

	closeGroup     structs.SyncGroup
	connectReady   chan struct{}
	discoveryReady chan struct{}
	speakingReady  chan struct{}
}

type UdpSession interface {
	Write(data []byte, binary bool)
	Connect(gateway string, udp bool) error
	Exit(graceful bool) error
	Error(err error)
	Discover() error
	KeepAlive()
	ResetSentData()
	IsDiscovered() bool
	GetUdpData() *gateway.UdpData
	SetUdpData(data gateway.UdpData)
	SetSecretKey(key [32]byte)
	GetSecretKey() [32]byte
	SetEncryption(encryption gateway.TransportEncryptionMode)
	GetEncryption() gateway.TransportEncryptionMode
	GetDiscoveryReady() <-chan struct{}
	GetSpeakingReady() <-chan struct{}
	GetConnectReady() <-chan struct{}
	CloseConnectReady()
	CloseDiscoveryReady()
	CloseSpeakingReady()
}

func NewUdpSession() UdpSession {
	u := &udpSession{
		mu:             &sync.Mutex{},
		Session:        NewSession(),
		eventHandler:   NewEventHandler[udpEventHandler](),
		closeGroup:     *structs.NewSyncGroup(),
		connectReady:   make(chan struct{}),
		discoveryReady: make(chan struct{}),
		speakingReady:  make(chan struct{}),
	}
	u.ctx, u.cancel = context.WithCancel(context.Background())

	u.closeGroup.AddChannel("connectReady")
	u.closeGroup.AddChannel("discoveryReady")
	u.closeGroup.AddChannel("speakingReady")

	u.SetListenFunc(u.validateEvent)
	u.SetHandleFunc(u.handleEvent)
	u.SetPayloadDecoders(&payload.DiscoveryPacket{}, &payload.VoicePacket{})
	u.SetEventDecoders(&payload.DiscoveryPacket{}, &payload.VoicePacket{}, &payload.SenderReportPacket{})
	u.SetErrorHandlers(map[error]func(){
		net.ErrClosed: func() {
			u.Exit(false)
		},
		io.EOF: func() {
			u.Exit(false)
		},
		io.ErrUnexpectedEOF: func() {
			u.Exit(false)
		},
	})
	u.SetValidCloseErrors(io.EOF, io.ErrUnexpectedEOF, net.ErrClosed)

	return u
}

func (u *udpSession) Write(data []byte, binary bool) {
	u.Session.Write(data, binary)

	u.mu.Lock()
	u.sentPackets++
	u.sentBytes += len(data)
	u.mu.Unlock()
}

func (u *udpSession) Exit(graceful bool) error {
	defer u.cancel()
	u.CloseConnectReady()
	u.CloseDiscoveryReady()
	u.CloseSpeakingReady()

	return u.Session.Exit(graceful)
}

func (u *udpSession) Discover() error {
	return u.eventHandler.HandleEvent(u, &payload.DiscoveryPacket{PacketType: 0})
}

func (u *udpSession) KeepAlive() {
	// send a keep alive packet every 5 seconds
	keepAlive := time.NewTicker(5 * time.Second)
	defer keepAlive.Stop()
	for {
		select {
		case <-u.ctx.Done():
			u.ResetSentData()
			return
		case <-keepAlive.C:
			ntpTimestamp := u.constructNTPTimestamp(time.Now().UTC())
			senderReport := payload.NewSenderReportPacket(
				ntpTimestamp,
				u.constructRTPTimestamp(ntpTimestamp),
				uint32(u.udpData.SSRC),
				uint32(u.sentPackets),
				uint32(u.sentBytes),
			)

			packet, err := senderReport.Marshal()
			if err != nil {
				u.Error(err)
				return
			}

			u.Write(packet, false)
			fmt.Printf("packets sent: %d, bytes sent: %d\n", u.sentPackets, u.sentBytes)
		}
	}
}

func (u *udpSession) IsDiscovered() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.discovered
}

func (u *udpSession) ResetSentData() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.sentPackets = 0
	u.sentBytes = 0
}

func (u *udpSession) GetUdpData() *gateway.UdpData {
	u.mu.Lock()
	defer u.mu.Unlock()
	return &u.udpData
}

func (u *udpSession) SetUdpData(data gateway.UdpData) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.udpData = data
}

func (u *udpSession) SetSecretKey(key [32]byte) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.secretKey = key
}

func (u *udpSession) GetSecretKey() [32]byte {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.secretKey
}

func (u *udpSession) SetEncryption(encryption gateway.TransportEncryptionMode) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.encryption = encryption
}

func (u *udpSession) GetEncryption() gateway.TransportEncryptionMode {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.encryption
}

func (u *udpSession) GetDiscoveryReady() <-chan struct{} {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.discoveryReady
}

func (u *udpSession) GetSpeakingReady() <-chan struct{} {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.speakingReady
}

func (u *udpSession) GetConnectReady() <-chan struct{} {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.connectReady
}

func (u *udpSession) CloseConnectReady() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.closeGroup.CloseChannels["connectReady"].Do(func() {
		close(u.connectReady)
	})
}

func (u *udpSession) CloseDiscoveryReady() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.closeGroup.CloseChannels["discoveryReady"].Do(func() {
		close(u.discoveryReady)
	})
	u.discovered = true
}

func (u *udpSession) CloseSpeakingReady() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.closeGroup.CloseChannels["speakingReady"].Do(func() {
		close(u.speakingReady)
	})
}

func (u *udpSession) handleEvent(p payload.Payload) error {
	return u.eventHandler.HandleEvent(u, p)
}

func (u *udpSession) validateEvent(p payload.Payload) (any, error) {
	dp, ok := p.(*payload.DiscoveryPacket)
	if !ok {
		vp, ok := p.(*payload.VoicePacket)
		if !ok {
			sr, ok := p.(*payload.SenderReportPacket)
			if !ok {
				return nil, errors.New("invalid udp payload type - validate error: " + p.ToString())
			}
			return sr, nil
		}
		return vp, nil
	}
	return dp, nil
}

func (u *udpSession) constructNTPTimestamp(t time.Time) uint64 {
	// Calculate the number of seconds since the NTP epoch
	seconds := uint32(t.Sub(ntpEpoch).Seconds())

	// Calculate the fractional part
	fraction := uint32((t.Sub(ntpEpoch).Nanoseconds() % 1e9) * (1 << 32) / 1e9)

	// Combine the integer and fractional parts into a 64-bit unsigned fixed-point number
	ntpTimestamp := uint64(seconds)<<32 | uint64(fraction)

	return ntpTimestamp
}

func (u *udpSession) constructRTPTimestamp(ntpTimestamp uint64) uint32 {
	// Extract the integer part of the NTP timestamp (first 32 bits)
	ntpSeconds := uint32(ntpTimestamp >> 32)

	// Extract the fractional part of the NTP timestamp (last 32 bits)
	ntpFraction := uint32(ntpTimestamp & 0xFFFFFFFF)

	// Convert the NTP timestamp to RTP timestamp using the clock rate
	rtpTimestamp := ntpSeconds*uint32(48000) + uint32((uint64(ntpFraction)*uint64(48000))>>32)

	return rtpTimestamp
}
