package voice

import (
	"bytes"
	"errors"
	"sync"
	"time"

	"github.com/Carmen-Shannon/simple-discord/structs/voice"
	"github.com/Carmen-Shannon/simple-discord/util/crypto"
	"github.com/Carmen-Shannon/simple-discord/util/ffmpeg"
)

type AudioPlayer interface {
	IsPlaying() bool
	IsConnected() bool
	GetUdpSession() UdpSession
	SetUdpSession(session UdpSession)
	GetVoiceSession() VoiceSession
	SetVoiceSession(v VoiceSession)
	Connect() error
	Play(filepath string) error
}

type audioPlayer struct {
	mu           sync.Mutex
	playing      bool
	connected    bool
	session      UdpSession
	voiceSession VoiceSession
}

var _ AudioPlayer = (*audioPlayer)(nil)

func NewAudioPlayer() AudioPlayer {
	return &audioPlayer{
		mu:           sync.Mutex{},
		playing:      false,
		connected:    false,
		session:      nil,
		voiceSession: nil,
	}
}

// Connect requires a VoiceSession to be passed in order to connect to the UDP server, since forming a UDP connection is reliant on the voice gateway
// we need to send voice gateway ops
func (a *audioPlayer) Connect() error {
	if a.session == nil || a.voiceSession == nil {
		return errors.New("sessions not initialized")
	}

	if err := a.session.Connect(); err != nil {
		return err
	}

	// send select protocol payload
	var selectProtocolPayload voice.VoicePayload
	selectProtocolPayload.OpCode = voice.SelectProtocol
	if err := a.voiceSession.GetEventHandler().HandleEvent(a.voiceSession, selectProtocolPayload); err != nil {
		return err
	}

	<-a.session.GetSession().speakingReady

	a.mu.Lock()
	a.connected = true
	a.mu.Unlock()
	return nil
}

// using this method to test
func (a *audioPlayer) Play(path string) error {
	if a.GetUdpSession() == nil || a.GetVoiceSession() == nil {
		return errors.New("sessions not initialized")
	}

	if !a.IsConnected() {
		if err := a.Connect(); err != nil {
			return err
		}
	}

	if a.playing {
		return errors.New("already playing")
	}

	audio := NewAudio()
	go func() {
		if err := audio.RegisterFile(path); err != nil {
			return
		}
	}()

	// Stream the encoded audio
	go a.play(audio)

	return nil
}

func (a *audioPlayer) play(audio Audio) {
	if !a.IsPlaying() {
		if err := a.voiceSession.Speaking(); err != nil {
			a.session.GetSession().errorChan <- err
			return
		}

		// TODO: fix this garbage with something better
		time.Sleep(100 * time.Millisecond) // should be enough time for the speaking event to process before modifying the playing state

		a.mu.Lock()
		a.playing = true
		a.mu.Unlock()
	}

	packetChannel := make(chan []byte)
	done := make(chan struct{})
	sampleSize := 960 // 48khz sample rate (48000 samples/second) * 0.02 seconds (20ms frame time) = 960 samples per frame
	frameDuration := 20 * time.Millisecond

	go func() {
		err := a.prepAudio(audio, packetChannel, sampleSize)
		if err != nil {
			a.session.GetSession().errorChan <- err
			close(packetChannel)
			return
		}
		close(packetChannel)
	}()

	go func() {
		err := a.sendAudio(frameDuration, packetChannel)
		if err != nil {
			a.session.GetSession().errorChan <- err
			close(done)
			return
		}
		close(done)
	}()

	<-done

	if err := a.voiceSession.Speaking(); err != nil {
		a.session.GetSession().errorChan <- err
		return
	}

	time.Sleep(100 * time.Millisecond) // should be enough time for the speaking event to process before modifying the playing state

	a.mu.Lock()
	a.playing = false
	a.mu.Unlock()
}

func (a *audioPlayer) sendAudio(frameTime time.Duration, packetChan chan []byte) error {
	ticker := time.NewTicker(frameTime)
	defer ticker.Stop()
	for {
		select {
		case <-a.GetUdpSession().GetSession().ctx.Done():
			return nil
		case msg, ok := <-packetChan:
			if !ok {
				return nil
			} else if !a.IsConnected() {
				return nil
			}

			a.GetUdpSession().Write(msg)

			a.GetUdpSession().PacketSent()
			a.GetUdpSession().BytesSent(len(msg)) // packet size

			<-ticker.C
		}
	}
}

func (a *audioPlayer) prepAudio(audio Audio, sendChan chan []byte, frameSize int) error {
	// done channel
	done := make(chan struct{})
	var once sync.Once
	// incrementals
	var seq uint16
	var timestamp uint32
	var nonceCounter uint32
	timestamp = uint32((time.Now().Unix() / 4) - 1)

	// encoded := audio.GetData()
	encryptionMode := a.GetUdpSession().GetConnData().Mode
	secretKey := a.GetUdpSession().GetSession().SecretKey
	ssrc := a.GetUdpSession().GetConnData().SSRC

	// TODO: add a NewRTPHeader method to voice.RTPHeader that can take a seq, timestamp, and ssrc
	rtpHeader := voice.RTPHeader{
		Version:     2,
		Padding:     false,
		Extension:   false,
		CSRCCount:   0,
		Marker:      false,
		PayloadType: 120,
		Seq:         seq,
		Timestamp:   timestamp,
		SSRC:        uint32(ssrc),
		CSRC:        []uint32{},
	}

	// Silence frame for the end-of-stream silence
	silenceFrame := []byte{0xF8, 0xFF, 0xFE}

	for {
		select {
		case <-a.GetUdpSession().GetSession().ctx.Done():
			once.Do(func() { close(done) })
			return nil
		case encoded, ok := <-audio.GetStream():
			if !ok {
				once.Do(func() { close(done) })
				break
			}
			rtpHeader.Seq = seq
			rtpHeader.Timestamp = timestamp
			packet, err := a.encryptAudio(encoded, rtpHeader, encryptionMode, nonceCounter, secretKey)
			if err != nil {
				a.session.GetSession().errorChan <- err
				return nil
			}

			sendChan <- packet.Bytes()

			nonceCounter++
			seq++
			timestamp += uint32(frameSize)
		case <-done:
		silenceLoop:
			for i := 0; i < 5; i++ {
				select {
				case <-a.GetUdpSession().GetSession().ctx.Done():
					break silenceLoop
				default:
					rtpHeader.Seq = seq
					rtpHeader.Timestamp = timestamp

					packet, err := a.encryptAudio(silenceFrame, rtpHeader, encryptionMode, nonceCounter, secretKey)
					if err != nil {
						a.session.GetSession().errorChan <- err
						return err
					}

					sendChan <- packet.Bytes()

					// Increment the nonce counter, sequence, and timestamp
					nonceCounter++
					seq++
					timestamp += uint32(frameSize)
				}
			}
			a.GetUdpSession().ResetSentData()
			return nil
		}
	}
}

func (a *audioPlayer) encryptAudio(unencryptedAudio []byte, header voice.RTPHeader, encryptionMode voice.TransportEncryptionMode, nonceCounter uint32, secretKey [32]byte) (bytes.Buffer, error) {
	var encryptedAudio bytes.Buffer
	headerBytes, err := header.MarshalBinary()
	if err != nil {
		return encryptedAudio, err
	}

	switch encryptionMode {
	// TODO: add support for chacha
	case voice.AEAD_XCHACHA20_POLY1305:
		chachaNonce := make([]byte, 24)
		copy(chachaNonce[:12], headerBytes)
		encryptedFrame, err := crypto.EncryptXChaCha20Poly1305(unencryptedAudio, secretKey[:], chachaNonce)
		if err != nil {
			return encryptedAudio, err
		}
		encryptedAudio.Write(encryptedFrame)
	case voice.AEAD_AES256_GCM:
		encryptedFrame, err := crypto.EncryptAESGCM(unencryptedAudio, secretKey[:], headerBytes, nonceCounter)
		if err != nil {
			return encryptedAudio, err
		}
		encryptedAudio.Write(encryptedFrame)
	default:
		// Handle other encryption modes if necessary
		return encryptedAudio, errors.New("unsupported encryption mode, not sure how you got here")
	}

	return encryptedAudio, nil
}

func (a *audioPlayer) IsPlaying() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.playing
}

func (a *audioPlayer) IsConnected() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.connected
}

func (a *audioPlayer) GetUdpSession() UdpSession {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.session
}

func (a *audioPlayer) SetUdpSession(session UdpSession) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.session = session
}

func (a *audioPlayer) GetVoiceSession() VoiceSession {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.voiceSession
}

func (a *audioPlayer) SetVoiceSession(v VoiceSession) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.voiceSession = v
}

// Audio is meant to store the metadata and PCM audio
type Audio interface {
	RegisterFile(path string) error
	RegisterStream()
	GetMetadata() voice.AudioMetadata
	GetStream() chan []byte
}

type audio struct {
	mu       sync.Mutex
	metadata voice.AudioMetadata
	stream   chan []byte
}

var _ Audio = (*audio)(nil)

// NewAudio will initialize a new Audio instance, setting up the stream for audio and the mutex for thread safety
func NewAudio() Audio {
	return &audio{
		mu:     sync.Mutex{},
		stream: make(chan []byte),
	}
}

// RegisterFile is a temporary function while I work on a more permanent solution, for now this will take a file path to any audio file (should in theory handle mp4 m4a etc) and process the audio into PCM.
// The PCM data will be encoded per-frame to Opus frames and sent to the stream channel.
func (a *audio) RegisterFile(path string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	metadata, err := ffmpeg.ConvertFileToOpus(path, true, a.stream)
	if err != nil {
		return errors.New("failed to convert file to Opus: " + err.Error())
	}

	a.metadata = *metadata
	close(a.stream)
	return nil
}

func (a *audio) RegisterStream() {
	a.mu.Lock()
	defer a.mu.Unlock()

	// // Process the stream and convert to PCM
	// audioReader, metadata, err := ffmpeg.StreamBytesToPCM(a.stream, true)
	// if err != nil {
	// 	fmt.Println("failed to process stream:", err)
	// 	return
	// }

	// a.metadata = *metadata

	// // Process the PCM stream and convert to Opus
	// audioReader, err = ffmpeg.StreamPCMToOpus(a.stream)
	// if err != nil {
	// 	fmt.Println("failed to process stream:", err)
	// 	return
	// }

	// a.data = audioReader
}

func (a *audio) GetMetadata() voice.AudioMetadata {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.metadata
}

func (a *audio) GetStream() chan []byte {
	return a.stream
}
