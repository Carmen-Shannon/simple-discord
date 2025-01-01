package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	gateway "github.com/Carmen-Shannon/simple-discord/gateway"
	requestutil "github.com/Carmen-Shannon/simple-discord/gateway/request_util"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	gateway_structs "github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"github.com/Carmen-Shannon/simple-discord/util"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"golang.org/x/net/websocket"
)

// handles establishing a new session with the discord gateway
func NewSession(token string, intents []structs.Intent) (*Session, error) {
	var sess Session
	sess.SetToken(&token)
	sess.SetEventHandler(NewEventHandler())
	sess.SetIntents(intents)
	sess.SetServers(make(map[string]structs.Server))
	sess.helloReceived = make(chan struct{})
	sess.stopHeartbeat = make(chan struct{})
	sess.replyAck = make(chan struct{})
	sess.readChan = make(chan []byte)
	sess.writeChan = make(chan []byte, 4096)
	sess.errorChan = make(chan error)

	// set up a dummy voice session
	sess.NewVoiceSession()

	ws, err := sess.dialer(nil, "/?v=10&encoding=json")
	if err != nil {
		return nil, err
	}
	sess.SetConn(ws)

	// set up the goroutines to listen, read, write, and handle errors
	go sess.listen()
	go sess.handleRead()
	go sess.handleWrite()
	go sess.handleError()

	// stop here until the HELLO event is receieved
	<-sess.helloReceived

	var identifyData gateway_structs.Payload
	identifyData.OpCode = gateway_structs.Identify

	// let her rip tater chip
	if err := sess.EventHandler.HandleEvent(&sess, identifyData); err != nil {
		return nil, err
	}

	return &sess, nil
}

type Session struct {
	// Mutex for thread safety
	Mu sync.Mutex

	// Websocket connection
	Conn *websocket.Conn

	// Heartbeat response time
	HeartbeatACK *int

	// Latest sequence number
	Sequence *int

	// Custom event handler
	EventHandler *EventHandler

	// Session ID
	ID *string

	// Gateway URL to resume with
	ResumeURL *string

	// Bot token
	Token *string

	// Intents to subscribe to
	Intents []structs.Intent

	// Servers the bot is in
	Servers map[string]structs.Server

	// Bot details
	BotData *structs.Bot

	// Voice session
	Voice *VoiceSession

	// various channels, self explanatory what each one does
	helloReceived chan struct{}
	stopHeartbeat chan struct{}
	replyAck      chan struct{}
	readChan      chan []byte
	writeChan     chan []byte
	errorChan     chan error
}

func (s *Session) JoinVoice(guildID, channelID structs.Snowflake) error {
	if s.Voice == nil {
		return fmt.Errorf("voice session not initialized")
	}
	s.Voice.SetGuildID(guildID)
	s.Voice.SetChannelID(channelID)

	// send the voice state update payload
	var voiceStateUpdate gateway_structs.Payload
	voiceStateUpdate.OpCode = gateway_structs.VoiceStateUpdate

	if err := s.EventHandler.HandleEvent(s, voiceStateUpdate); err != nil {
		return err
	}

	<- s.Voice.connectReady

	if err := s.Voice.Connect(); err != nil {
		return err
	}
	return nil
}

// closes the hearbeat and websocket connection
func (s *Session) Exit() error {
	if s.stopHeartbeat != nil {
		close(s.stopHeartbeat)
	}
	// Close the connection
	if err := s.Conn.Close(); err != nil {
		return fmt.Errorf("error closing connection: %v", err)
	}
	return nil
}

// listens for new messages sent to the readChan and parses them before submitting them to the EventHandler
func (s *Session) listen() {
	for msg := range s.readChan {
		var payload gateway_structs.Payload
		if err := json.Unmarshal(msg, &payload); err != nil {
			s.errorChan <- fmt.Errorf("error parsing message: %v", err)
			continue
		}

		var err error
		payload.Data, err = gateway.NewReceiveEvent(payload)
		if err != nil {
			s.errorChan <- fmt.Errorf("error parsing event: %v", err)
			fmt.Println(payload.ToString())
			continue
		}

		// signal when HELLO is received
		if payload.OpCode == gateway_structs.Hello {
			close(s.helloReceived)
		}

		if err := s.EventHandler.HandleEvent(s, payload); err != nil {
			s.errorChan <- fmt.Errorf("error handling event: %v", err)
			fmt.Println(payload.ToString())
			continue
		}
	}
}

// writes messages as raw bytes to the writeChan
func (s *Session) Write(data []byte) {
	if len(s.writeChan) < cap(s.writeChan) {
		s.writeChan <- data
	} else {
		s.errorChan <- fmt.Errorf("failed to write data to write channel")
	}
}

func (s *Session) Reply(interactionOptions structs.InteractionResponseOptions, interaction *structs.Interaction) error {
	if s.replyAck == nil {
		s.replyAck = make(chan struct{})
	}

	done := make(chan error)

    go func() {
        err := s.sendReply(interactionOptions, interaction)
        if err != nil {
            done <- err
            return
        }

        // Wait for the replyAck channel to be closed
        select {
        case <-s.replyAck:
            done <- nil
        case <-time.After(10 * time.Second): // Optional timeout to prevent indefinite blocking
            done <- fmt.Errorf("timeout waiting for replyAck")
        }
    }()

    // Wait for the goroutine to finish
    return <-done
}

func (s *Session) sendReply(interactionOptions structs.InteractionResponseOptions, interaction *structs.Interaction) error {
	interactionID := interaction.ID.ToString()
	interactionToken := interaction.Token
	token := *s.GetToken()
	reqDto := dto.CreateInteractionResponseDto{
		WithResponse: util.ToPtr(true),
	}
	response := interactionOptions.InteractionResponse()
	_, err := requestutil.CreateInteractionResponse(interactionID, interactionToken, token, reqDto, *response)
	if err != nil {
		return err
	}

	// Signal that the reply is complete
	return nil
}

// RegisterCommands adds custom commands to the EventHandler
//
// The command key must match the name of the command that was registered using the Discord API
// You must ACK the command within 3 seconds or Discord will assume the command failed, to properly ACK a command you must
// call the session.Reply method
func (s *Session) RegisterCommands(commands map[string]func(*Session, gateway_structs.Payload) error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for name, command := range commands {
		s.EventHandler.AddCustomHandler(name, command)
	}
}

// reads frames from the gateway in increments of 1024 bytes
// dynamically resizes the buffer array to fit the full message and writes the message to the readChan
func (s *Session) handleRead() {
	defer close(s.readChan)

	var buffers [][]byte

	for {
		tempBuffer := make([]byte, 1024)
		n, err := s.Conn.Read(tempBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			s.errorChan <- err
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
				s.errorChan <- fmt.Errorf("error decoding raw message: %v", err)
				buffers = nil
				break
			}

			s.readChan <- msg

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

// reads from the writeChan and writes the message to the gateway
// has a retry mechanism with a delay of 2 seconds
// after 3 retries, give up and go home
func (s *Session) handleWrite() {
	defer close(s.writeChan)

	retryCount := 0
	maxRetries := 3
	retryDelay := time.Second * 2

	for data := range s.writeChan {
		for {
			if _, err := s.Conn.Write(data); err != nil {
				if retryCount < maxRetries {
					retryCount++
					log.Printf("write error: %v, retrying %d/%d", err, retryCount, maxRetries)
					time.Sleep(retryDelay)
					continue
				} else {
					s.errorChan <- fmt.Errorf("write error after %d retries: %v", maxRetries, err)
					if err := s.ReconnectSession(); err != nil {
						s.errorChan <- fmt.Errorf("error resuming session: %v", err)
						s.Exit()
						break
					}
					return
				}
			}
			retryCount = 0 // Reset retry count on successful write
			break
		}
	}
}

// reads from the errorChan and logs the error
func (s *Session) handleError() {
	defer close(s.errorChan)

	for err := range s.errorChan {
		log.Printf("error: %v\n", err)
	}
}

// when a session is disconnected but can be resumed for one of many reasons, use this
func (s *Session) ResumeSession() error {
	close(s.stopHeartbeat)

	// open a new connection using the cached url
	ws, err := s.dialer(nil, "/?v=10&encoding=json")
	if err != nil {
		return err
	}
	s.SetConn(ws)

	// Reinitialize the channels
	s.helloReceived = make(chan struct{})
	s.stopHeartbeat = make(chan struct{})
	s.replyAck = make(chan struct{})
	s.readChan = make(chan []byte)
	s.writeChan = make(chan []byte, 4096)
	s.errorChan = make(chan error)

	// Start the goroutines to listen, read, write, and handle errors
	go s.listen()
	go s.handleRead()
	go s.handleWrite()
	go s.handleError()

	var resumeData gateway_structs.Payload
	resumeData.OpCode = gateway_structs.Resume

	// let her rip tater chip
	if err := s.EventHandler.HandleEvent(s, resumeData); err != nil {
		return err
	}

	return nil
}

// when a session is disconnected and can not be resumed, use this
func (s *Session) ReconnectSession() error {
	close(s.stopHeartbeat)

	ws, err := s.dialer(nil, "/?v=10&encoding=json")
	if err != nil {
		return err
	}
	s.SetConn(ws)

	// reinitialize the channels
	s.helloReceived = make(chan struct{})
	s.stopHeartbeat = make(chan struct{})
	s.replyAck = make(chan struct{})
	s.readChan = make(chan []byte)
	s.writeChan = make(chan []byte, 4096)
	s.errorChan = make(chan error)

	// Start the goroutines to listen, read, write, and handle errors
	go s.listen()
	go s.handleRead()
	go s.handleWrite()
	go s.handleError()

	<-s.helloReceived

	var identifyData gateway_structs.Payload
	identifyData.OpCode = gateway_structs.Identify

	// let her rip tater chip
	if err := s.EventHandler.HandleEvent(s, identifyData); err != nil {
		return err
	}

	return nil
}

// it dial
func (s *Session) dialer(url *string, query string) (*websocket.Conn, error) {
	if url == nil {
		if s.GetResumeURL() != nil {
			url = s.GetResumeURL()
		} else {
			var err error
			if s.GetToken() == nil {
				return nil, fmt.Errorf("token not set for session")
			}
			gatewayUrl, err := requestutil.GetGatewayUrl(*s.GetToken())
			if err != nil {
				return nil, err
			}
			url = &gatewayUrl
		}
	}

	ws, err := websocket.Dial(*url+query, "", "http://localhost/")
	if err != nil {
		return nil, err
	}

	return ws, nil
}

// getters and setters because mutex
func (s *Session) GetConn() *websocket.Conn {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Conn
}

func (s *Session) SetConn(conn *websocket.Conn) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Conn = conn
}

func (s *Session) GetHeartbeatACK() *int {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.HeartbeatACK
}

func (s *Session) SetHeartbeatACK(ack *int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.HeartbeatACK = ack
}

func (s *Session) GetSequence() *int {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Sequence
}

func (s *Session) SetSequence(seq *int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Sequence = seq
}

func (s *Session) GetEventHandler() *EventHandler {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.EventHandler
}

func (s *Session) SetEventHandler(handler *EventHandler) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.EventHandler = handler
}

func (s *Session) GetID() *string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.ID
}

func (s *Session) SetID(id *string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.ID = id
}

func (s *Session) GetResumeURL() *string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.ResumeURL
}

func (s *Session) SetResumeURL(url *string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.ResumeURL = url
}

func (s *Session) GetToken() *string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Token
}

func (s *Session) SetToken(token *string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Token = token
}

func (s *Session) GetIntents() []structs.Intent {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Intents
}

func (s *Session) SetIntents(intents []structs.Intent) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Intents = intents
}

func (s *Session) GetServers() map[string]structs.Server {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Servers
}

func (s *Session) SetServers(servers map[string]structs.Server) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Servers = servers
}

func (s *Session) GetServerByName(name string) *structs.Server {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for _, guild := range s.Servers {
		if guild.Name == name {
			return &guild
		}
	}

	return nil
}

func (s *Session) AddServer(server structs.Server) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Servers[server.ID.ToString()] = server
}

func (s *Session) GetBotData() *structs.Bot {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.BotData
}

func (s *Session) SetBotData(bot *structs.Bot) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.BotData = bot
}

func (s *Session) SetVoiceSession(session *VoiceSession) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Voice = session
}

func (s *Session) GetVoiceSession() *VoiceSession {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Voice
}
