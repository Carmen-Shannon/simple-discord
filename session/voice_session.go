package session

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	voice_gateway "github.com/Carmen-Shannon/simple-discord/gateway/voice"
	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/voice"
	"golang.org/x/net/websocket"
)

type VoiceSession struct {
	// Mutex for thread safety
	Mu sync.Mutex

	// Websocket connection
	Conn *websocket.Conn

	// Heartbeat response time
	HeartbeatACK *int

	// Latest sequence number
	Sequence *int

	// Custom event handler for voice gateway
	EventHandler *VoiceEventHandler

	// Session ID of the voice session
	SessionID *string

	// Token for the voice session
	Token *string

	// Guild ID (server the session is for)
	GuildID *structs.Snowflake

	// Channel ID in the guild
	ChannelID *structs.Snowflake

	// Gateway URL to resume/connect to
	ResumeURL *string

	// Is the voice session connected
	Connected bool

	// Bot details
	BotData *structs.Bot

	// UDP Connection details
	UdpConn *voice.UdpData

	// channels
	connectReady   chan struct{}
	heartbeatReady chan struct{}
	stopHeartbeat  chan struct{}
	readChan       chan []byte
	writeChan      chan []byte
	errorChan      chan error
}

func (s *Session) NewVoiceSession() {
	var vs VoiceSession
	vs.SetEventHandler(NewVoiceEventHandler())
	vs.SetConnected(false)
	vs.connectReady = make(chan struct{})
	vs.heartbeatReady = make(chan struct{})
	vs.stopHeartbeat = make(chan struct{})
	vs.readChan = make(chan []byte)
	vs.writeChan = make(chan []byte, 4096)
	vs.errorChan = make(chan error)

	s.SetVoiceSession(&vs)
}

func (v *VoiceSession) Connect() error {
	if v.GetConnected() {
		return nil
	}
	if v.GetResumeURL() == nil {
		return errors.New("voice session requires a resume URL to connect")
	}

	conn, err := websocket.Dial(*v.GetResumeURL()+"?v=8", "", "http://localhost")
	if err != nil {
		return err
	}
	v.SetConn(conn)

	go v.listen()
	go v.handleRead()
	go v.handleWrite()
	go v.handleError()

	// send identify payload
	var identifyPayload voice.VoicePayload
	identifyPayload.OpCode = voice.Identify

	if err := v.EventHandler.HandleEvent(v, identifyPayload); err != nil {
		return err
	}

	<-v.heartbeatReady
	v.SetConnected(true)
	return nil
}

// writes messages as raw bytes to the writeChan
func (v *VoiceSession) Write(data []byte) {
	if len(v.writeChan) < cap(v.writeChan) {
		v.writeChan <- data
	} else {
		v.errorChan <- fmt.Errorf("failed to write data to write channel")
	}
}

// listens for new messages sent to the readChan and parses them before submitting them to the EventHandler
func (v *VoiceSession) listen() {
	for msg := range v.readChan {
		var payload voice.VoicePayload
		// first try the message as a json payload, then try as a binary payload
		if err := json.Unmarshal(msg, &payload); err != nil {
			var binaryPayload voice.BinaryVoicePayload
			if err := binaryPayload.UnmarshalBinary(msg); err != nil {
				v.errorChan <- fmt.Errorf("error parsing message: %v", err)
				continue
			}

			// TODO: implement binary receive event function
			binaryPayload.Data, err = voice_gateway.NewBinaryReceiveEvent(binaryPayload)
			if err != nil {
				v.errorChan <- fmt.Errorf("error parsing event: %v", err)
				continue
			}

			if err := v.EventHandler.HandleBinaryEvent(v, binaryPayload); err != nil {
				v.errorChan <- fmt.Errorf("error handling binary event: %v", err)
				continue
			}
		}

		var err error
		payload.Data, err = voice_gateway.NewReceiveEvent(payload)
		if err != nil {
			v.errorChan <- fmt.Errorf("error parsing event: %v\n%s", err, payload.ToString())
			continue
		}

		if err := v.EventHandler.HandleEvent(v, payload); err != nil {
			v.errorChan <- fmt.Errorf("error handling event: %v\n%s", err, payload.ToString())
			continue
		}
	}
}

// reads frames from the gateway in increments of 1024 bytes
// dynamically resizes the buffer array to fit the full message and writes the message to the readChan
func (v *VoiceSession) handleRead() {
	defer close(v.readChan)

	var buffers [][]byte

	for {
		tempBuffer := make([]byte, 1024)
		n, err := v.Conn.Read(tempBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			v.errorChan <- err
			break
		}

		buffers = append(buffers, tempBuffer[:n])

		for {
			combinedBuffer := bytes.Join(buffers, nil)
			decoder := json.NewDecoder(bytes.NewReader(combinedBuffer))
			var msg json.RawMessage
			startOffset := len(combinedBuffer)

			if err := decoder.Decode(&msg); err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					// incomplete message
					break
				}

				// If JSON decoding fails, attempt to decode as binary
				payload, err := v.decodeBinaryMessage(combinedBuffer)
				if payload == nil && err == nil {
					// incomplete message
					break
				} else if err != nil {
					v.errorChan <- fmt.Errorf("error decoding message: %v", err)
					buffers = nil
					break
				}

				payloadBytes, err := payload.MarshalBinary()
				if err != nil {
					v.errorChan <- fmt.Errorf("error marshalling binary message: %v", err)
					buffers = nil
					break
				}

				v.readChan <- payloadBytes
				continue
			}

			v.readChan <- msg

			offset := startOffset - int(decoder.InputOffset())
			if offset > 0 && offset <= len(combinedBuffer) {
				remainingData := combinedBuffer[decoder.InputOffset():]
				buffers = [][]byte{remainingData}
			} else {
				buffers = nil
			}
		}
	}
}

// reads from the writeChan and writes the message to thSe gateway
// has a retry mechanism with a delay of 2 seconds
// after 3 retries, give up and go home
func (v *VoiceSession) handleWrite() {
	defer close(v.writeChan)

	retryCount := 0
	maxRetries := 3
	retryDelay := time.Second * 2

	for data := range v.writeChan {
		for {
			if _, err := v.Conn.Write(data); err != nil {
				if retryCount < maxRetries {
					retryCount++
					log.Printf("write error: %v, retrying %d/%d", err, retryCount, maxRetries)
					time.Sleep(retryDelay)
					continue
				} else {
					v.errorChan <- fmt.Errorf("write error after %d retries: %v", maxRetries, err)
					// TODO: implement reconnect logic
					// if err := s.ReconnectSession(); err != nil {
					// 	s.errorChan <- fmt.Errorf("error resuming session: %v", err)
					// 	s.Exit()
					// 	break
					// }
					return
				}
			}
			retryCount = 0 // Reset retry count on successful write
			break
		}
	}
}

// reads from the errorChan and logs the error
func (v *VoiceSession) handleError() {
	defer close(v.errorChan)

	for err := range v.errorChan {
		log.Printf("error: %v\n", err)
	}
}

// decodeBinaryMessage attempts to decode the message using binary big-endian encoding
func (v *VoiceSession) decodeBinaryMessage(data []byte) (*voice.BinaryVoicePayload, error) {
	buf := bytes.NewReader(data)
	var payload voice.BinaryVoicePayload

	// Read OpCode
	if err := binary.Read(buf, binary.BigEndian, &payload.OpCode); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read opcode: %w", err)
	}

	// Read remaining bytes into Payload
	payloadData := make([]byte, buf.Len())
	if _, err := buf.Read(payloadData); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, nil
		}
		return &payload, fmt.Errorf("failed to read payload: %w", err)
	}

	payload.Data = payloadData
	return &payload, nil
}

func (v *VoiceSession) SetConn(conn *websocket.Conn) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.Conn = conn
}

func (v *VoiceSession) SetHeartbeatACK(heartbeatACK int) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.HeartbeatACK = &heartbeatACK
}

func (v *VoiceSession) SetSequence(sequence int) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.Sequence = &sequence
}

func (v *VoiceSession) SetEventHandler(eventHandler *VoiceEventHandler) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.EventHandler = eventHandler
}

func (v *VoiceSession) SetSessionID(id string) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.SessionID = &id
}

func (v *VoiceSession) SetToken(token string) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.Token = &token
}

func (v *VoiceSession) SetGuildID(guildID structs.Snowflake) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.GuildID = &guildID
}

func (v *VoiceSession) SetChannelID(channelID structs.Snowflake) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.ChannelID = &channelID
}

func (v *VoiceSession) SetResumeURL(resumeURL string) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.ResumeURL = &resumeURL
}

func (v *VoiceSession) SetConnected(connected bool) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.Connected = connected
}

func (v *VoiceSession) SetBotData(botData *structs.Bot) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.BotData = botData
}

func (v *VoiceSession) SetUdpConn(udpConn *voice.UdpData) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.UdpConn = udpConn
}

func (v *VoiceSession) GetConn() *websocket.Conn {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.Conn
}

func (v *VoiceSession) GetHeartbeatACK() *int {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.HeartbeatACK
}

func (v *VoiceSession) GetSequence() *int {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.Sequence
}

func (v *VoiceSession) GetEventHandler() *VoiceEventHandler {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.EventHandler
}

func (v *VoiceSession) GetSessionID() *string {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.SessionID
}

func (v *VoiceSession) GetToken() *string {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.Token
}

func (v *VoiceSession) GetGuildID() *structs.Snowflake {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.GuildID
}

func (v *VoiceSession) GetChannelID() *structs.Snowflake {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.ChannelID
}

func (v *VoiceSession) GetResumeURL() *string {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.ResumeURL
}

func (v *VoiceSession) GetConnected() bool {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.Connected
}

func (v *VoiceSession) GetBotData() *structs.Bot {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.BotData
}

func (v *VoiceSession) GetUdpConn() *voice.UdpData {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.UdpConn
}
