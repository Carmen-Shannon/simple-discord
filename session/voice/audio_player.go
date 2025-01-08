package voice

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	a.connected = true
	return nil
}

// using this method to test
func (a *audioPlayer) Play(path string) error {
	if a.session == nil || a.voiceSession == nil {
		return errors.New("sessions not initialized")
	}

	if !a.connected {
		if err := a.Connect(); err != nil {
			return err
		}
	}

	if a.playing {
		return errors.New("already playing")
	}

	audio := NewAudio()
	if err := audio.RegisterFile(path); err != nil {
		return err
	}

	// Debugging function to save audio data to a file
	// rootDir, err := os.Getwd()
	// if err != nil {
	// 	return fmt.Errorf("failed to get current directory: %w", err)
	// }
	// savePath := filepath.Join(rootDir, "local", "test-encoded.opus")
	// if err := saveAudioData(audio.GetData(), savePath); err != nil {
	// 	return fmt.Errorf("failed to save audio data: %w", err)
	// }

	// Stream the encoded audio
	go a.play(audio)

	return nil
}

func saveAudioData(data []byte, filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (a *audioPlayer) play(audio Audio) {
	if !a.IsPlaying() {
		if err := a.voiceSession.Speaking(); err != nil {
			a.session.GetSession().errorChan <- err
			return
		}

		time.Sleep(100 * time.Millisecond) // should be enough time for the speaking event to process before modifying the playing state

		a.mu.Lock()
		a.playing = true
		a.mu.Unlock()
	}

	packetChannel := make(chan []byte, 100)
	done := make(chan struct{})

	// Get the metadata of the audio file
	metadata := audio.GetMetadata()
	// set up audio stream parameters
	frameDuration := 20 * time.Millisecond
	frameSize := 320 // 20ms frame duration with a 128kbps bitrate means 320 bytes per frame

	fmt.Println("audio data length:", metadata.DurationMs/1000)
	fmt.Println("total bytes:", len(audio.GetData()))
	fmt.Println("frame size:", frameSize)
	fmt.Println("encryption mode:", string(a.GetUdpSession().GetConnData().Mode))
	fmt.Println("address:", a.GetUdpSession().GetConnData().Address)
	fmt.Println("port:", a.GetUdpSession().GetConnData().Port)

	go func() {
		// incrementals
		var seq uint16
		var timestamp uint32
		var nonceCounter uint32

		encoded := audio.GetData()
		encryptionMode := a.GetUdpSession().GetConnData().Mode
		secretKey := a.GetUdpSession().GetSession().SecretKey
		ssrc := a.GetUdpSession().GetConnData().SSRC

		for len(encoded) > 0 {
			// Set up the RTP header
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

			rtpHeaderBytes, err := rtpHeader.MarshalBinary()
			if err != nil {
				a.session.GetSession().errorChan <- err
				close(done)
				return
			}

			frame := encoded[:min(len(encoded), frameSize)]
			encryptedAudio, err := a.encryptAudio(frame, rtpHeader, encryptionMode, nonceCounter, secretKey)
			if err != nil {
				a.session.GetSession().errorChan <- err
				close(done)
				return
			}

			payload := bytes.Buffer{}
			if _, err := payload.Write(rtpHeaderBytes); err != nil {
				a.session.GetSession().errorChan <- err
				close(done)
				return
			}
			if _, err := payload.Write(encryptedAudio.Bytes()); err != nil {
				a.session.GetSession().errorChan <- err
				close(done)
				return
			}

			// Log the packet details for debugging
			fmt.Printf("Sending packet: Seq=%d, Timestamp=%d, Nonce=%d, PayloadSize=%d\n", seq, timestamp, nonceCounter, len(payload.Bytes()))

			// Send the packet to the packetChannel
			packetChannel <- payload.Bytes()

			// Move to the next frame
			encoded = encoded[min(len(encoded), frameSize):]

			// Increment the nonce counter, sequence, and timestamp
			nonceCounter++
			seq++
			timestamp += uint32(frameSize)
		}

		// Send five frames of silence (0xF8, 0xFF, 0xFE)
		silenceFrame := []byte{0xF8, 0xFF, 0xFE}
		for i := 0; i < 5; i++ {
			packetChannel <- silenceFrame
		}

		close(packetChannel)
	}()

	var decryptedData bytes.Buffer

	go func() {
		for {
			select {
			case packet, ok := <-packetChannel:
				if !ok {
					close(done)
					break
				}

				// Decrypt the packet
				decryptedAudio, err := a.decryptPayload(packet, a.GetUdpSession().GetConnData().Mode)
				if err != nil {
					a.session.GetSession().errorChan <- err
					close(done)
					return
				}

				// Write the decrypted audio data to the buffer
				decryptedData.Write(decryptedAudio)

				// Send the packet over the network
				a.session.Write(packet)
				// update the udp session stats
				a.GetUdpSession().PacketSent()
				a.GetUdpSession().BytesSent(len(packet) - 12) // packet size - the 12 byte RTP header

				// Wait for the frame duration
				time.Sleep(frameDuration)
			case <-done:
				return
			}
		}
	}()

	<-done

	if err := a.voiceSession.Speaking(); err != nil {
		a.session.GetSession().errorChan <- err
		return
	}

	time.Sleep(100 * time.Millisecond) // should be enough time for the speaking event to process before modifying the playing state

	// Save the decrypted audio data to a file
	if err := saveDecryptedAudio(decryptedData.Bytes(), "local/decrypted_audio.opus"); err != nil {
		fmt.Println("Failed to save decrypted audio:", err)
	}

	a.mu.Lock()
	a.playing = false
	a.mu.Unlock()
}

func saveDecryptedAudio(data []byte, filename string) error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Construct the absolute path to the local directory
	absPath := filepath.Join(cwd, "local", filename)

	// Create the local directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// Write the file to the local directory
	return os.WriteFile(absPath, data, 0644)
}

func (a *audioPlayer) encryptAudio(unencryptedAudio []byte, header voice.RTPHeader, encryptionMode voice.TransportEncryptionMode, nonceCounter uint32, secretKey [32]byte) (bytes.Buffer, error) {
	var encryptedAudio bytes.Buffer
	headerBytes, err := header.MarshalBinary()
	if err != nil {
		return encryptedAudio, err
	}

	switch encryptionMode {
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

func (a *audioPlayer) decryptPayload(encryptedPacket []byte, encryptionMode voice.TransportEncryptionMode) ([]byte, error) {
	// extract the RTP header from the first 12 bytes of the packet
	rtpHeader := voice.RTPHeader{}
	if err := rtpHeader.UnmarshalBinary(encryptedPacket[:12]); err != nil {
		fmt.Println("failed to unmarshal RTP header:", err)
		return nil, err
	}
	// extract the nonce from the last 4 bytes of the packet
	nonce := binary.BigEndian.Uint32(encryptedPacket[len(encryptedPacket)-4:])
	noncePadding := make([]byte, 12)
	binary.BigEndian.PutUint32(noncePadding[:4], nonce)
	var decryptedAudio bytes.Buffer
	switch encryptionMode {
	case voice.AEAD_XCHACHA20_POLY1305:
		// Decrypt the packet using XChaCha20-Poly1305
		// This is not implemented in this snippet
	case voice.AEAD_AES256_GCM:
		decryptedFrame, err := crypto.DecryptAESGCM(encryptedPacket[12:len(encryptedPacket)-4], a.GetUdpSession().GetSession().SecretKey[:], noncePadding, encryptedPacket[:12])
		if err != nil {
			fmt.Println("failed to decrypt frame:", err)
			return nil, err
		}
		decryptedAudio.Write(decryptedFrame)
	default:
	}

	return decryptedAudio.Bytes(), nil
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

func (a *audioPlayer) SetVoiceSession(v VoiceSession) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.voiceSession = v
}

// Audio is meant to store the metadata and PCM audio
type Audio interface {
	RegisterFile(path string) error
	RegisterStream()
	GetData() []byte
	GetMetadata() voice.AudioMetadata
}

type audio struct {
	mu       sync.Mutex
	metadata voice.AudioMetadata
	data     io.Reader
	buffer   []byte
	stream   chan []byte
}

var _ Audio = (*audio)(nil)

func NewAudio() Audio {
	return &audio{
		mu:     sync.Mutex{},
		stream: make(chan []byte, 100),
	}
}

func (a *audio) RegisterFile(path string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	audio, metadata, err := ffmpeg.ConvertFileToOpus(path, true)
	if err != nil {
		return errors.New("failed to convert file to Opus: " + err.Error())
	}

	a.data = bufio.NewReader(bytes.NewReader(audio))
	a.metadata = *metadata
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

func (a *audio) GetData() []byte {
	a.mu.Lock()
	defer a.mu.Unlock()

	// If the buffer is already populated, return it
	if a.buffer != nil {
		return a.buffer
	}

	// Read from the data reader and store in the buffer
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(a.data)
	if err != nil {
		fmt.Println("failed to read data:", err)
		return nil
	}

	a.buffer = buf.Bytes()
	return a.buffer
}

func (a *audio) GetMetadata() voice.AudioMetadata {
	return a.metadata
}
