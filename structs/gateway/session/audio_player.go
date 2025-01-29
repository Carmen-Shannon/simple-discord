package session

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway/payload"
	"github.com/Carmen-Shannon/simple-discord/util/crypto"
)

type audioPlayer struct {
	mu     *sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	session UdpSession

	speakingFunc       func(bool) error
	selectProtocolFunc func() error

	connected bool
	playing   bool

	audioResource AudioResource
}

type AudioPlayer interface {
	Play(path string) error
	Connect() error
	Exit()
	IsConnected() bool
	IsPlaying() bool
	GetSession() UdpSession
}

func NewAudioPlayer(speakingFunc func(bool) error, selectProtocolFunc func() error) AudioPlayer {
	a := &audioPlayer{
		mu:                 &sync.Mutex{},
		session:            NewUdpSession(),
		audioResource:      NewAudioResource(),
		speakingFunc:       speakingFunc,
		selectProtocolFunc: selectProtocolFunc,
	}
	a.ctx, a.cancel = context.WithCancel(context.Background())
	return a
}

func (a *audioPlayer) Play(path string) error {
	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return err
		}
	} else if a.IsPlaying() {
		return errors.New("audio is already playing")
	} else if a.IsConnected() {
		a.mu.Lock()
		a.audioResource.Exit()
		a.audioResource = NewAudioResource()
		a.mu.Unlock()
	}

	go func() {
		if err := a.audioResource.RegisterFile(path); err != nil {
			a.session.Error(err)
			return
		}
	}()

	go a.playAudio()
	return nil
}

func (a *audioPlayer) Connect() error {
	gateway := fmt.Sprintf("%s:%d", a.session.GetUdpData().Address, a.session.GetUdpData().Port)
	if err := a.session.Connect(gateway, true); err != nil {
		return err
	}

	if err := a.session.Discover(); err != nil {
		return err
	}

	select {
	case <-a.ctx.Done():
		return nil
	case <-a.session.GetDiscoveryReady():
	}

	if err := a.selectProtocolFunc(); err != nil {
		return err
	}

	select {
	case <-a.ctx.Done():
		return nil
	case <-a.session.GetSpeakingReady():
	}

	go a.session.KeepAlive()
	a.mu.Lock()
	a.connected = true
	a.mu.Unlock()
	return nil
}

func (a *audioPlayer) Exit() {
	defer a.cancel()
	defer a.audioResource.Exit()

	if a.IsConnected() {
		a.session.Exit(true)
		a.mu.Lock()
		a.connected = false
		a.mu.Unlock()
	}
}

func (a *audioPlayer) IsConnected() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.connected
}

func (a *audioPlayer) IsPlaying() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.playing
}

func (a *audioPlayer) GetSession() UdpSession {
	return a.session
}

func (a *audioPlayer) playAudio() {
	if !a.IsPlaying() {
		if err := a.speakingFunc(true); err != nil {
			a.session.Error(err)
			return
		}

		time.Sleep(250 * time.Millisecond)
		a.mu.Lock()
		a.playing = true
		a.mu.Unlock()
	}

	sendChan := make(chan []byte)

	frameSize := 960
	frameTime := 20 * time.Millisecond

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer close(sendChan)
		if err := a.prepAudio(sendChan, frameSize); err != nil {
			a.session.Error(err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.sendAudio(frameTime, sendChan); err != nil {
			a.session.Error(err)
			return
		}
	}()

	wg.Wait()
	a.mu.Lock()
	a.playing = false
	a.mu.Unlock()

	if a.IsConnected() {
		if err := a.speakingFunc(false); err != nil {
			a.session.Error(err)
			return
		}
	}
}

func (a *audioPlayer) prepAudio(sendChan chan []byte, frameSize int) error {
	var seq uint16
	var timestamp uint32
	var nonce uint32
	timestamp = uint32((time.Now().Unix() / 4) - 1)

	encryption := a.session.GetEncryption()
	secretKey := a.session.GetSecretKey()
	ssrc := a.session.GetUdpData().SSRC

	header := payload.NewRtpHeader(seq, timestamp, uint32(ssrc))
	stream := a.audioResource.GetStream()

	for {
		select {
		case <-a.ctx.Done():
			return nil
		case encoded, ok := <-stream:
			if !ok {
				for i := 0; i < 5; i++ {
					select {
					case <-a.ctx.Done():
						return nil
					default:
						header.Seq = seq
						header.Timestamp = timestamp

						packet, err := a.encrypt(payload.SilenceFrame, *header, encryption, nonce, secretKey)
						if err != nil {
							return err
						}

						select {
						case sendChan <- packet.Bytes():
						case <-a.ctx.Done():
							return nil
						}

						seq++
						nonce++
						timestamp += uint32(frameSize)
					}
				}
				return nil
			} else if !a.connected {
				return nil
			}

			header.Seq = seq
			header.Timestamp = timestamp
			packet, err := a.encrypt(encoded, *header, encryption, nonce, secretKey)
			if err != nil {
				return err
			}

			select {
			case sendChan <- packet.Bytes():
			case <-a.ctx.Done():
				return nil
			}

			seq++
			nonce++
			timestamp += uint32(frameSize)
		}
	}
}

func (a *audioPlayer) sendAudio(frameTime time.Duration, receiveChan chan []byte) error {
	ticker := time.NewTicker(frameTime)
	defer ticker.Stop()
	defer a.session.ResetSentData()
	for {
		select {
		case <-a.ctx.Done():
			return nil
		case msg, ok := <-receiveChan:
			if !ok {
				return nil
			} else if !a.connected {
				return nil
			}

			a.session.Write(msg, false)
			<-ticker.C
		}
	}
}

func (a *audioPlayer) encrypt(packet []byte, rtpHeader payload.RTPHeader, encryptionMode gateway.TransportEncryptionMode, nonce uint32, secretKey [32]byte) (bytes.Buffer, error) {
	var encryptedAudio bytes.Buffer
	headerBytes, err := rtpHeader.MarshalBinary()
	if err != nil {
		return encryptedAudio, err
	}

	switch encryptionMode {
	case gateway.AEAD_AES256_GCM:
		encryptedFrame, err := crypto.EncryptAESGCM(packet, secretKey[:], headerBytes, nonce)
		if err != nil {
			return encryptedAudio, err
		}
		encryptedAudio.Write(encryptedFrame)
	case gateway.AEAD_XCHACHA20_POLY1305:
		chachaNonce := make([]byte, 24)
		copy(chachaNonce[:12], headerBytes)
		encryptedFrame, err := crypto.EncryptXChaCha20Poly1305(packet, secretKey[:], chachaNonce)
		if err != nil {
			return encryptedAudio, err
		}
		encryptedAudio.Write(encryptedFrame)
	default:
		return encryptedAudio, errors.New("unsupported encryption mode")
	}

	return encryptedAudio, nil
}
