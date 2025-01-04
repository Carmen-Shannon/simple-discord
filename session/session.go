package session

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	gateway "github.com/Carmen-Shannon/simple-discord/gateway"
	requestutil "github.com/Carmen-Shannon/simple-discord/gateway/request_util"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	gateway_structs "github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"github.com/Carmen-Shannon/simple-discord/util"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type CommandFunc func(s Session, p gateway_structs.Payload) error

type Session interface {
	JoinVoice(guildID, channelID structs.Snowflake) error
	LeaveVoice(guildID structs.Snowflake) error
	Exit() error
	NewVoiceSession()
	Write(data []byte)
	SendMessage(messageOptions dto.MessageOptions) error
	InteractionReply(interactionOptions structs.InteractionResponseOptions, interaction *structs.Interaction) error
	RegisterCommands(commands map[string]CommandFunc)
	RegisterListeners(listeners map[Listener]CommandFunc)
	ResumeSession() error
	ReconnectSession() error
	GetConn() *websocket.Conn
	SetConn(conn *websocket.Conn)
	GetHeartbeatACK() *int
	SetHeartbeatACK(ack *int)
	GetSequence() *int
	SetSequence(seq *int)
	GetEventHandler() *EventHandler
	SetEventHandler(handler *EventHandler)
	GetID() *string
	SetID(id *string)
	GetResumeURL() *string
	SetResumeURL(url *string)
	GetToken() *string
	SetToken(token *string)
	GetIntents() []structs.Intent
	SetIntents(intents []structs.Intent)
	GetServers() map[string]structs.Server
	SetServers(servers map[string]structs.Server)
	GetServerByName(name string) *structs.Server
	GetServerByGuildID(guildID structs.Snowflake) *structs.Server
	AddServer(server structs.Server)
	GetBotData() *structs.BotData
	SetBotData(bot *structs.BotData)
	SetVoiceSession(session VoiceSession)
	GetVoiceSession() VoiceSession
	GetShard() *int
	SetShard(shard *int)
	GetShards() *int
	SetShards(shards *int)
	GetMaxConcurrency() *int
	SetMaxConcurrency(maxConcurrency *int)
	GetSession() *session
}

type session struct {
	// Mutex for thread safety
	Mu *sync.Mutex

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
	BotData *structs.BotData

	// Voice session
	Voice VoiceSession

	// Shard ID
	Shard *int

	// Shard count
	Shards *int

	// Max identify concurrency
	MaxConcurrency *int

	ctx    context.Context
	cancel context.CancelFunc

	// various channels, self explanatory what each one does
	helloReceived chan struct{}
	stopHeartbeat chan struct{}
	readChan      chan []byte
	writeChan     chan []byte
	errorChan     chan error
}

var _ Session = (*session)(nil)

// handles establishing a new session with the discord gateway
func NewSession(token string, intents []structs.Intent, shard *int) (Session, error) {
	var sess session
	sess.Mu = &sync.Mutex{}
	sess.Token = &token
	sess.EventHandler = NewEventHandler()
	sess.Intents = intents
	sess.Servers = make(map[string]structs.Server)
	sess.ctx, sess.cancel = context.WithCancel(context.Background())
	sess.helloReceived = make(chan struct{})
	sess.stopHeartbeat = make(chan struct{})
	sess.readChan = make(chan []byte)
	sess.writeChan = make(chan []byte, 4096)
	sess.errorChan = make(chan error)

	// set up a dummy voice session
	sess.NewVoiceSession()

	if shard != nil {
		sess.SetShard(shard)
	}

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

// RegisterCommands adds custom commands to the EventHandler
//
// The command key must match the name of the command that was registered using the Discord API
// You must ACK the command within 3 seconds or Discord will assume the command failed, to properly ACK a command you must
// call the session.Reply method
func (s *session) RegisterCommands(commands map[string]CommandFunc) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for name, command := range commands {
		s.EventHandler.AddCustomHandler(name, command)
	}
}

// RegisterListeners adds custom listeners to the EventHandler
//
// This would be used to interact with gateway events from Discord, like MESSAGE_CREATE
// The list of events you can listen to are defined in the Listener enum
func (s *session) RegisterListeners(listeners map[Listener]CommandFunc) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for listener, command := range listeners {
		s.EventHandler.AddListener(string(listener), command)
	}
}

// writes messages as raw bytes to the writeChan
func (s *session) Write(data []byte) {
	if len(s.writeChan) < cap(s.writeChan) {
		s.writeChan <- data
	} else {
		s.errorChan <- fmt.Errorf("failed to write data to write channel")
	}
}

// closes the hearbeat and websocket connection
func (s *session) Exit() error {
	if s.stopHeartbeat != nil {
		close(s.stopHeartbeat)
		s.stopHeartbeat = nil
	}

	// Close the voice connection if there is one still active
	if s.GetVoiceSession() != nil && s.GetVoiceSession().GetConnected() {
		if err := s.GetVoiceSession().Exit(); err != nil {
			return err
		}
	}
	// Close the connection
	if err := s.Conn.Close(websocket.StatusNormalClosure, "disconnect"); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			return fmt.Errorf("error closing connection: %v", err)
		}
	}

	s.cancel()
	time.Sleep(1 * time.Second) //arbitrary sleep to allow for cleanup
	// Close the channels
	close(s.readChan)
	close(s.writeChan)
	close(s.errorChan)

	return nil
}

// when a session is disconnected but can be resumed for one of many reasons, use this
func (s *session) ResumeSession() error {
	if s.stopHeartbeat != nil {
		close(s.stopHeartbeat)
		s.stopHeartbeat = nil
	}

	if s.GetVoiceSession().GetConnected() {
		err := s.GetVoiceSession().ResumeSession()
		if err != nil {
			return err
		}
	}

	// open a new connection using the cached url
	ws, err := s.dialer(nil, "/?v=10&encoding=json")
	if err != nil {
		return err
	}
	s.SetConn(ws)

	s.cancel()
	// Reinitialize the channels
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.helloReceived = make(chan struct{})
	s.stopHeartbeat = make(chan struct{})
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
func (s *session) ReconnectSession() error {
	if s.stopHeartbeat != nil {
		close(s.stopHeartbeat)
		s.stopHeartbeat = nil
	}

	s.NewVoiceSession()

	ws, err := s.dialer(nil, "/?v=10&encoding=json")
	if err != nil {
		return err
	}
	s.SetConn(ws)

	// reinitialize the channels
	s.helloReceived = make(chan struct{})
	s.stopHeartbeat = make(chan struct{})
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

func (s *session) JoinVoice(guildID, channelID structs.Snowflake) error {
	if s.GetVoiceSession() == nil {
		return fmt.Errorf("voice session not initialized")
	}
	if s.GetVoiceSession().GetConnected() {
		if !channelID.Equals(*s.GetVoiceSession().GetChannelID()) &&
			guildID.Equals(*s.GetVoiceSession().GetGuildID()) {
			if err := s.GetVoiceSession().Exit(); err != nil {
				return err
			}
			s.NewVoiceSession()
		}
	}
	s.GetVoiceSession().SetGuildID(guildID)
	s.GetVoiceSession().SetChannelID(channelID)
	s.GetVoiceSession().SetBotData(s.GetBotData())

	// send the voice state update payload
	var voiceStateUpdate gateway_structs.Payload
	voiceStateUpdate.OpCode = gateway_structs.VoiceStateUpdate

	if err := s.EventHandler.HandleEvent(s, voiceStateUpdate); err != nil {
		return err
	}

	// wait for the voice gateway to be ready, timeout after 5 seconds
	select {
	case <-s.GetVoiceSession().GetSession().connectReady:
		if err := s.GetVoiceSession().Connect(); err != nil {
			return err
		}
	case <-time.After(5 * time.Second):
		return fmt.Errorf("voice gateway did not connect in time")
	}

	return nil
}

func (s *session) LeaveVoice(guildID structs.Snowflake) error {
	if s.GetVoiceSession() == nil {
		return fmt.Errorf("voice session not initialized")
	}
	s.GetVoiceSession().SetGuildID(guildID)
	s.GetVoiceSession().ClearChannelID()

	// send the voice state update payload
	var voiceStateUpdate gateway_structs.Payload
	voiceStateUpdate.OpCode = gateway_structs.VoiceStateUpdate

	if err := s.EventHandler.HandleEvent(s, voiceStateUpdate); err != nil {
		return err
	}

	if err := s.GetVoiceSession().Exit(); err != nil {
		s.SetVoiceSession(nil)
		s.NewVoiceSession()
		return err
	}

	s.NewVoiceSession()
	return nil
}

func (s *session) SendMessage(messageOptions dto.MessageOptions) error {
	done := make(chan error)

	go func() {
		_, err := s.sendMessage(messageOptions)
		if err != nil {
			done <- err
			return
		}

		done <- nil
	}()

	// Wait for the goroutine to finish
	return <-done
}

// InteractionReply is used to reply to interaction create events, this must be called within 3 seconds of receiving the event
func (s *session) InteractionReply(interactionOptions structs.InteractionResponseOptions, interaction *structs.Interaction) error {
	done := make(chan error)

	go func() {
		err := s.sendInteractionReply(interactionOptions, interaction)
		if err != nil {
			done <- err
			return
		}

		done <- nil
	}()

	// Wait for the goroutine to finish
	return <-done
}

// listens for new messages sent to the readChan and parses them before submitting them to the EventHandler
func (s *session) listen() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case msg := <-s.readChan:
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
				if s.helloReceived != nil {
					close(s.helloReceived)
					s.helloReceived = nil
				}
			}

			if err := s.EventHandler.HandleEvent(s, payload); err != nil {
				s.errorChan <- fmt.Errorf("error handling event: %v", err)
				fmt.Println(payload.ToString())
				continue
			}
		}
	}
}

func (s *session) sendMessage(messageOptions dto.MessageOptions) (*structs.Message, error) {
	token := *s.GetToken()
	reqDto, err := messageOptions.ConstructDtoFromOptions()
	if err != nil {
		return nil, err
	}

	message, err := requestutil.CreateMessage(*reqDto, token)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *session) sendInteractionReply(interactionOptions structs.InteractionResponseOptions, interaction *structs.Interaction) error {
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

// reads frames from the gateway in increments of 1024 bytes
// dynamically resizes the buffer array to fit the full message and writes the message to the readChan
func (s *session) handleRead() {
	var buffer bytes.Buffer
	s.Conn.SetReadLimit(-1)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			_, bytes, err := s.Conn.Read(s.ctx)
			if err != nil {
				if !errors.Is(err, net.ErrClosed) && !errors.As(err, &websocket.CloseError{}) {
					s.errorChan <- fmt.Errorf("error reading from session websocket: %v", err)
				}
				return
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
					s.errorChan <- fmt.Errorf("error decoding raw message: %v", err)
					buffer.Reset()
					break
				}

				s.readChan <- msg

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
}

// reads from the writeChan and writes the message to the gateway
// has a retry mechanism with a delay of 2 seconds
// after 3 retries, give up and go home
func (s *session) handleWrite() {
	retryCount := 0
	maxRetries := 3
	retryDelay := time.Second * 2

	for {
		select {
		case <-s.ctx.Done():
			return
		case data := <-s.writeChan:
			for {
				var msg json.RawMessage
				if err := json.Unmarshal(data, &msg); err != nil {
					s.errorChan <- fmt.Errorf("error unmarshalling data: %v", err)
					break
				}
				if err := wsjson.Write(s.ctx, s.GetConn(), msg); err != nil {
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
}

// reads from the errorChan and logs the error
func (s *session) handleError() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case err := <-s.errorChan:
			log.Printf("session error: %v\n", err)
		}
	}
}

// it dial
func (s *session) dialer(url *string, query string) (*websocket.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if url == nil {
		if s.GetResumeURL() != nil {
			url = s.GetResumeURL()
		} else {
			var err error
			if s.GetToken() == nil {
				return nil, fmt.Errorf("token not set for session")
			}
			gatewayBot, err := requestutil.GetGatewayBot(*s.GetToken())
			if err != nil {
				return nil, err
			}
			url = &gatewayBot.URL
			s.SetShards(&gatewayBot.Shards)
			if s.GetShard() == nil {
				s.SetShard(util.ToPtr(0))
			}
			if s.GetMaxConcurrency() == nil {
				s.SetMaxConcurrency(util.ToPtr(gatewayBot.SessionStartLimit.MaxConcurrency))
			}
		}
	}

	ws, _, err := websocket.Dial(ctx, *url+query, nil)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

// getters and setters because mutex
func (s *session) GetConn() *websocket.Conn {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Conn
}

func (s *session) SetConn(conn *websocket.Conn) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Conn = conn
}

func (s *session) GetHeartbeatACK() *int {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.HeartbeatACK
}

func (s *session) SetHeartbeatACK(ack *int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.HeartbeatACK = ack
}

func (s *session) GetSequence() *int {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Sequence
}

func (s *session) SetSequence(seq *int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Sequence = seq
}

func (s *session) GetEventHandler() *EventHandler {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.EventHandler
}

func (s *session) SetEventHandler(handler *EventHandler) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.EventHandler = handler
}

func (s *session) GetID() *string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.ID
}

func (s *session) SetID(id *string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.ID = id
}

func (s *session) GetResumeURL() *string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.ResumeURL
}

func (s *session) SetResumeURL(url *string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.ResumeURL = url
}

func (s *session) GetToken() *string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Token
}

func (s *session) SetToken(token *string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Token = token
}

func (s *session) GetIntents() []structs.Intent {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Intents
}

func (s *session) SetIntents(intents []structs.Intent) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Intents = intents
}

func (s *session) GetServers() map[string]structs.Server {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Servers
}

func (s *session) SetServers(servers map[string]structs.Server) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Servers = servers
}

func (s *session) GetServerByName(name string) *structs.Server {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for _, guild := range s.Servers {
		if guild.Name == name {
			return &guild
		}
	}

	return nil
}

func (s *session) GetServerByGuildID(guildID structs.Snowflake) *structs.Server {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for _, guild := range s.Servers {
		if guild.ID.Equals(guildID) {
			return &guild
		}
	}

	return nil
}

func (s *session) AddServer(server structs.Server) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Servers[server.ID.ToString()] = server
}

func (s *session) GetBotData() *structs.BotData {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.BotData
}

func (s *session) SetBotData(bot *structs.BotData) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.BotData = bot
}

func (s *session) SetVoiceSession(session VoiceSession) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Voice = session
}

func (s *session) GetVoiceSession() VoiceSession {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Voice
}

func (s *session) GetShard() *int {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Shard
}

func (s *session) SetShard(shard *int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Shard = shard
}

func (s *session) GetShards() *int {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Shards
}

func (s *session) SetShards(shards *int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Shards = shards
}

func (s *session) GetMaxConcurrency() *int {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.MaxConcurrency
}

func (s *session) SetMaxConcurrency(maxConcurrency *int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.MaxConcurrency = maxConcurrency
}

func (s *session) GetSession() *session {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s
}
