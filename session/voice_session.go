package session

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	sendevents "github.com/Carmen-Shannon/simple-discord/gateway/send_events"
	voice_gateway "github.com/Carmen-Shannon/simple-discord/gateway/voice"
	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/voice"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var ErrConnClosed = errors.New("use of closed network connection")

type VoiceSession struct {
	// Mutex for thread safety
	Mu *sync.Mutex

	// setting up a one-time closure of the connectReady channel
	once sync.Once

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
	BotData *structs.BotData

	// UDP Connection details
	UdpConn *voice.UdpData

	// checking if the voice session is ready to connect
	isConnectReady bool

	// channels
	connectReady   chan struct{}
	heartbeatReady chan struct{}
	stopHeartbeat  chan struct{}
	resumeReady    chan struct{}
	readChan       chan []byte
	writeChan      chan []byte
	errorChan      chan error
}

func NewVoiceSession() *VoiceSession {
	var vs VoiceSession
	vs.Mu = &sync.Mutex{}
	vs.once = sync.Once{}
	vs.SetEventHandler(NewVoiceEventHandler())
	vs.SetConnected(false)
	vs.SetConnectReady(false)
	vs.heartbeatReady = make(chan struct{})
	vs.stopHeartbeat = make(chan struct{})
	vs.connectReady = make(chan struct{})
	vs.readChan = make(chan []byte)
	vs.writeChan = make(chan []byte, 4096)
	vs.errorChan = make(chan error)

	return &vs
}

func (s *Session) NewVoiceSession() {
	vs := NewVoiceSession()
	vs.SetBotData(s.GetBotData())
	s.SetVoiceSession(vs)
}

func (v *VoiceSession) Exit() error {
	if v.stopHeartbeat != nil {
		close(v.stopHeartbeat)
		v.stopHeartbeat = nil
	}

	if v.GetConnected() {
		if err := v.GetConn().Close(websocket.StatusNormalClosure, "disconnect"); err != nil {
			return fmt.Errorf("error closing voice websocket: %v", err)
		}
	}

	*v = *NewVoiceSession()
	return nil
}

func (v *VoiceSession) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if v.GetConnected() {
		return nil
	}
	if v.GetResumeURL() == nil {
		return errors.New("voice session requires a resume URL to connect")
	}

	conn, _, err := websocket.Dial(ctx, *v.GetResumeURL()+"?v=8", nil)
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

// when a session is disconnected and can be resumed, use this
func (v *VoiceSession) ResumeSession() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if v.GetGuildID() == nil || v.GetSessionID() == nil || v.GetToken() == nil || v.GetSequence() == nil {
		return errors.New("voice session requires a guild ID, session ID, and token to resume")
	}

	v.heartbeatReady = make(chan struct{})
	v.stopHeartbeat = make(chan struct{})
	v.connectReady = make(chan struct{})
	v.readChan = make(chan []byte)
	v.writeChan = make(chan []byte, 4096)
	v.errorChan = make(chan error)

	v.SetConnected(false)
	v.SetConnectReady(false)

	conn, _, err := websocket.Dial(ctx, *v.GetResumeURL()+"?v=8", nil)
	if err != nil {
		return err
	}
	v.SetConn(conn)

	go v.listen()
	go v.handleRead()
	go v.handleWrite()
	go v.handleError()

	// send the resume payload
	var resumePayload voice.VoicePayload
	resumePayload.OpCode = voice.Resume
	resumePayload.Data = sendevents.VoiceResumeEvent{
		ServerID:  *v.GetGuildID(),
		SessionID: *v.GetSessionID(),
		Token:     *v.GetToken(),
		SeqAck:    *v.GetSequence(),
	}

	if err := v.EventHandler.HandleEvent(v, resumePayload); err != nil {
		return err
	}

	<-v.resumeReady
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer close(v.readChan)

	var buffer bytes.Buffer
	v.Conn.SetReadLimit(-1)

	for {
		_, bytes, err := v.Conn.Read(ctx)
		if err != nil {
			if err == io.EOF {
				break
			} else if errors.As(err, &websocket.CloseError{Code: websocket.StatusNormalClosure}) {
				v.Exit()
				return
			}
			v.errorChan <- fmt.Errorf("error reading from voice websocket: %v", err)
			break
		}

		buffer.Write(bytes)

		for {
			decoder := json.NewDecoder(&buffer)
			var msg json.RawMessage
			startOffset := buffer.Len()

			if err := decoder.Decode(&msg); err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					// incomplete message
					if startOffset <= buffer.Len() {
						buffer.Truncate(startOffset)
					}
					break
				}

				// If JSON decoding fails, attempt to decode as binary
				payload, err := v.decodeBinaryMessage(buffer.Bytes())
				if payload == nil && err == nil {
					// incomplete message
					if startOffset <= buffer.Len() {
						buffer.Truncate(startOffset)
					}
					break
				} else if err != nil {
					v.errorChan <- fmt.Errorf("error decoding binary message: %v", err)
					buffer.Reset()
					break
				}

				payloadBytes, err := payload.MarshalBinary()
				if err != nil {
					v.errorChan <- fmt.Errorf("error marshalling binary message: %v", err)
					buffer.Reset()
					break
				}

				v.readChan <- payloadBytes
				buffer.Reset()
				continue
			}

			v.readChan <- msg

			var remainingData []byte
			if int(decoder.InputOffset()) > len(buffer.Bytes()) {
				buffer.Reset()
			} else {
				remainingData = buffer.Bytes()[decoder.InputOffset():]
			}

			if len(remainingData) > 0 {
				buffer.Reset()
				buffer.Write(remainingData)
			} else {
				buffer.Reset()
			}
		}
	}
}

// has a retry mechanism with a delay of 2 seconds
// after 3 retries, give up and go home
func (v *VoiceSession) handleWrite() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer close(v.writeChan)

	retryCount := 0
	maxRetries := 3
	retryDelay := time.Second * 2

	for data := range v.writeChan {
		for {
			var msg json.RawMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				v.errorChan <- fmt.Errorf("error unmarshalling data: %v", err)
				break
			}

			if err := wsjson.Write(ctx, v.Conn, msg); err != nil {
				if strings.Contains(err.Error(), ErrConnClosed.Error()) {
					fmt.Println("voice conn closed, returning from handleWrite")
					return
				}
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
		log.Printf("voice error: %v\n", err)
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

func (v *VoiceSession) ClearChannelID() {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.ChannelID = nil
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

func (v *VoiceSession) SetBotData(botData *structs.BotData) {
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

func (v *VoiceSession) GetBotData() *structs.BotData {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.BotData
}

func (v *VoiceSession) GetUdpConn() *voice.UdpData {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.UdpConn
}

func (v *VoiceSession) IsConnectReady() bool {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	return v.isConnectReady
}

func (v *VoiceSession) SetConnectReady(ready bool) {
	v.Mu.Lock()
	defer v.Mu.Unlock()
	v.isConnectReady = ready

	if v.isConnectReady && !v.Connected {
		v.once.Do(func() {
			close(v.connectReady)
		})
	}
}
