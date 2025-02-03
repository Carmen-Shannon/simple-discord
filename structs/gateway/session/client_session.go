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
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway/payload"
	receiveevents "github.com/Carmen-Shannon/simple-discord/structs/gateway/receive_events"
	sendevents "github.com/Carmen-Shannon/simple-discord/structs/gateway/send_events"
	"github.com/Carmen-Shannon/simple-discord/util"
	requestutil "github.com/Carmen-Shannon/simple-discord/util/request_util"
	"github.com/coder/websocket"
)

type CommandFunc func(s ClientSession, p payload.SessionPayload) error

type clientSession struct {
	Session
	mu     *sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	token          *string
	intents        []structs.Intent
	botData        *structs.BotData
	heartbeatAck   *int
	resumeUrl      *string
	sessionID      *string
	sequence       *int
	shard          *int
	shards         *int
	maxConcurrency *int
	version        string

	servers       map[string]*structs.Server
	voiceSessions map[string]VoiceSession

	eventHandler *eventHandler

	cb func(s ClientSession) error

	closeGroup     structs.SyncGroup
	helloReceived  chan struct{}
	readyReceived  chan struct{}
	resumeReceived chan struct{}
}

type ClientSession interface {
	Write(date []byte, binary bool)
	Exit(graceful bool) error
	Error(err error)
	Dial(init bool) error
	Resume(url string) error
	Reply(interactionOptions structs.InteractionResponseOptions, interaction *structs.Interaction) error
	Send(messageOptions dto.MessageOptions, response bool) (*structs.Message, error)
	JoinVoice(guildID, channelID structs.Snowflake) error
	DisconnectVoice(guildID structs.Snowflake) error
	Play(filepath string, guildID, channelID structs.Snowflake) error
	ReconnectSession() error
	ResumeSession() error
	RegisterCommands(commands map[string]CommandFunc)
	RegisterListeners(listeners map[Listener]CommandFunc)
	GetToken() *string
	SetToken(token string)
	GetIntents() []structs.Intent
	SetIntents(intents ...structs.Intent)
	GetBotData() *structs.BotData
	SetBotData(data structs.BotData)
	GetSessionID() *string
	SetSessionID(id string)
	GetSequence() *int
	SetSequence(seq int)
	GetShard() *int
	SetShard(shard int)
	GetShards() *int
	SetShards(shards int)
	GetMaxConcurrency() *int
	SetMaxConcurrency(concurrency int)
	AddServer(s structs.Server)
	GetServers() map[string]*structs.Server
	SetHeartbeatAck(ack int)
	GetHeartbeatAck() *int
	GetCtx() context.Context
	GetCancel() context.CancelFunc
	SetResumeUrl(url string)
	GetResumeUrl() *string
	GetServerByGuildID(guildID structs.Snowflake) *structs.Server
	SetCb(cb func(s ClientSession) error)
	SetEventHandler(handler *eventHandler)
	GetEventHandler() *eventHandler
	GetVoiceSession(guildID structs.Snowflake) VoiceSession
	AddVoiceSession(guildID structs.Snowflake, vs VoiceSession)
	CloseHelloReceived()
	CloseReadyReceived()
	CloseResumeReceived()
}

var _ ClientSession = (*clientSession)(nil)

func NewClientSession(version string) ClientSession {
	cs := &clientSession{
		mu:             &sync.Mutex{},
		Session:        NewSession(),
		eventHandler:   NewEventHandler[eventHandler](),
		servers:        make(map[string]*structs.Server),
		voiceSessions:  make(map[string]VoiceSession),
		closeGroup:     *structs.NewSyncGroup(),
		helloReceived:  make(chan struct{}),
		readyReceived:  make(chan struct{}),
		resumeReceived: make(chan struct{}),
	}
	if version == "" {
		cs.version = "1"
	} else {
		cs.version = version
	}
	cs.ctx, cs.cancel = context.WithCancel(context.Background())

	cs.closeGroup.AddChannel("helloReceived")
	cs.closeGroup.AddChannel("readyReceived")
	cs.closeGroup.AddChannel("resumeReceived")

	cs.SetListenFunc(cs.validateEvent)
	cs.SetHandleFunc(cs.handleEvent)
	cs.SetWriteLimit(4096)
	cs.SetPayloadDecoders(&payload.RawMessagePayload{})
	cs.SetEventDecoders(&payload.SessionPayload{})
	cs.SetStatusCodeHandlers(map[websocket.StatusCode]func(){
		websocket.StatusNormalClosure: func() {},
		websocket.StatusGoingAway: func() {
			if err := cs.ResumeSession(); err != nil {
				if err := cs.ReconnectSession(); err != nil {
					cs.Exit(false)
				}
			}
		},
		4000: func() {
			cs.ReconnectSession()
		},
		4001: func() {
			cs.ReconnectSession()
		},
		4002: func() {
			cs.ReconnectSession()
		},
		4003: func() {
			cs.ReconnectSession()
		},
		4004: func() {
			fmt.Println("Session token is invalid, please check your token")
			cs.Exit(false)
		},
		4005: func() {
			cs.ReconnectSession()
		},
		4007: func() {
			cs.ReconnectSession()
		},
		4008: func() {
			cs.ReconnectSession()
		},
		4009: func() {
			cs.ReconnectSession()
		},
	})
	cs.SetErrorHandlers(map[error]func(){
		net.ErrClosed: func() {},
		errWsaRecv: func() {
			if err := cs.ResumeSession(); err != nil {
				cs.Exit(false)
			}
		},
		io.EOF: func() {
			if err := cs.ResumeSession(); err != nil {
				cs.Exit(false)
			}
		},
		io.ErrUnexpectedEOF: func() {
			if err := cs.ResumeSession(); err != nil {
				cs.Exit(false)
			}
		},
	})
	cs.SetValidCloseErrors(io.EOF, io.ErrUnexpectedEOF, net.ErrClosed, errWsaRecv, errWsaSend)

	return cs
}

func (s *clientSession) Exit(graceful bool) error {
	defer s.cancel()
	s.CloseHelloReceived()
	s.CloseReadyReceived()
	s.CloseResumeReceived()

	return s.Session.Exit(graceful)
}

func (s *clientSession) Error(err error) {
	wrapped := errors.Join(errors.New("client session error: "), err)
	s.Session.Error(wrapped)
}

func (s *clientSession) Dial(init bool) error {
	query := "/?v=10&encoding=json"

	if init {
		if err := s.initDialer(query); err != nil {
			return err
		}
	} else {
		if err := s.dialer(query); err != nil {
			return err
		}
	}

	select {
	case <-s.ctx.Done():
		return errors.New("context cancelled before hello received")
	case <-s.helloReceived:
	}

	if err := s.identify(); err != nil {
		return err
	}

	select {
	case <-s.ctx.Done():
		return errors.New("context cancelled before ready received")
	case <-s.readyReceived:
		return nil
	}
}

func (s *clientSession) Resume(url string) error {
	query := "/?v=10&encoding=json"

	if url == "" {
		if err := s.initDialer(query); err != nil {
			return err
		}
	}
	if err := s.dialer(query); err != nil {
		return err
	}

	select {
	case <-s.ctx.Done():
		return nil
	case <-s.helloReceived:
	}

	if err := s.resume(); err != nil {
		return err
	}

	select {
	case <-s.ctx.Done():
		return nil
	case <-s.resumeReceived:
		return nil
	case <-time.After(2 * time.Second):
		return errors.New("resume timeout")
	}
}

func (s *clientSession) Reply(interactionOptions structs.InteractionResponseOptions, interaction *structs.Interaction) error {
	done := make(chan error)

	go func() {
		defer close(done)
		if err := s.interactionReply(interactionOptions, interaction); err != nil {
			done <- err
			return
		}
	}()

	select {
	case <-s.ctx.Done():
		return nil
	case msg, ok := <-done:
		if !ok {
			return nil
		}
		return msg
	}
}

func (s *clientSession) Send(messageOptions dto.MessageOptions, response bool) (*structs.Message, error) {
	token := *s.GetToken()
	reqDto, err := messageOptions.ConstructDtoFromOptions()
	if err != nil {
		return nil, err
	}

	msg, err := requestutil.CreateMessage(*reqDto, token)
	if err != nil {
		return nil, err
	}

	if response {
		return msg, nil
	}

	return nil, nil
}

func (s *clientSession) JoinVoice(guildID, channelID structs.Snowflake) error {
	vs := s.GetVoiceSession(guildID)
	if vs == nil {
		vs = NewVoiceSession()
		vs.SetCleanupFunc(s.vsCleanup(guildID))
		vs.SetResumeFunc(s.vsResume(guildID, vs))
		vs.SetReconnectFunc(s.vsReconnect(guildID, channelID))
	} else if !vs.IsConnected() {
		vs = NewVoiceSession()
		vs.SetCleanupFunc(s.vsCleanup(guildID))
		vs.SetResumeFunc(s.vsResume(guildID, vs))
		vs.SetReconnectFunc(s.vsReconnect(guildID, channelID))
	} else if vs.IsConnected() {
		if vs.GetChannelID().Equals(channelID) {
			return errors.New("already connected to the voice channel")
		} else {
			if err := vs.Exit(false); err != nil {
				return err
			}
			vs = NewVoiceSession()
			vs.SetCleanupFunc(s.vsCleanup(guildID))
			vs.SetResumeFunc(s.vsResume(guildID, vs))
			vs.SetReconnectFunc(s.vsReconnect(guildID, channelID))
		}
	}

	vs.SetBotData(*s.GetBotData())
	vs.SetGuildID(guildID)
	vs.SetChannelID(channelID)
	s.AddVoiceSession(guildID, vs)

	if err := s.voiceStateUpdate(&guildID, &channelID); err != nil {
		return err
	}

	select {
	case <-s.ctx.Done():
		return errors.New("context cancelled before voice ready received")
	case <-vs.GetCtx().Done():
		return errors.New("voice context cancelled before voice ready received")
	case <-vs.GetConnectReady():
		if err := vs.Connect(); err != nil {
			return err
		}
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("voice connect timeout")
	}
}

func (s *clientSession) DisconnectVoice(guildID structs.Snowflake) error {
	vs := s.GetVoiceSession(guildID)
	if vs == nil {
		return nil
	}

	if err := vs.Exit(true); err != nil {
		return err
	}

	if err := s.voiceStateUpdate(&guildID, nil); err != nil {
		return err
	}

	delete(s.voiceSessions, guildID.ToString())

	return nil
}

func (s *clientSession) Play(filepath string, guildID, channelID structs.Snowflake) error {
	vs := s.GetVoiceSession(guildID)
	if vs == nil {
		if err := s.JoinVoice(guildID, channelID); err != nil {
			return err
		}

		vs = s.GetVoiceSession(guildID)
	}

	ap := vs.GetAudioPlayer()
	if err := ap.Play(filepath); err != nil {
		return err
	}

	return nil
}

func (s *clientSession) ReconnectSession() error {
	if err := s.Exit(true); err != nil {
		return err
	}

	sess := NewClientSession(s.version)
	sess.SetToken(*s.GetToken())
	sess.SetIntents(s.GetIntents()...)
	sess.SetShard(*s.GetShard())
	sess.SetShards(*s.GetShards())
	sess.SetMaxConcurrency(*s.GetMaxConcurrency())
	sess.SetEventHandler(s.eventHandler)
	sess.SetCb(s.cb)
	if err := sess.Dial(false); err != nil {
		return err
	}

	s.mu.Lock()
	if err := s.cb(sess); err != nil {
		return err
	}
	s.mu.Unlock()

	return nil
}

func (s *clientSession) ResumeSession() error {
	if err := s.Exit(false); err != nil {
		return err
	}

	sess := NewClientSession(s.version)
	sess.SetToken(*s.GetToken())
	sess.SetIntents(s.GetIntents()...)
	sess.SetShard(*s.GetShard())
	sess.SetShards(*s.GetShards())
	sess.SetMaxConcurrency(*s.GetMaxConcurrency())
	sess.SetResumeUrl(*s.GetResumeUrl())
	sess.SetEventHandler(s.eventHandler)
	sess.SetSessionID(*s.GetSessionID())
	sess.SetBotData(*s.GetBotData())
	sess.SetSequence(*s.GetSequence())
	sess.SetCb(s.cb)
	for _, server := range s.GetServers() {
		sess.AddServer(*server)
	}
	for _, vs := range s.voiceSessions {
		sess.AddVoiceSession(*vs.GetGuildID(), vs)
	}
	if err := sess.Resume("1"); err != nil {
		return err
	}

	s.mu.Lock()
	if err := s.cb(sess); err != nil {
		return err
	}
	s.mu.Unlock()

	return nil
}

func (s *clientSession) RegisterCommands(commands map[string]CommandFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for name, cmd := range commands {
		s.eventHandler.AddCustomHandler(name, cmd)
	}
}

func (s *clientSession) RegisterListeners(listeners map[Listener]CommandFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for listener, cmd := range listeners {
		s.eventHandler.AddListener(string(listener), cmd)
	}
}

func (s *clientSession) GetToken() *string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.token
}

func (s *clientSession) SetToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = &token
}

func (s *clientSession) GetIntents() []structs.Intent {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.intents
}

func (s *clientSession) SetIntents(intents ...structs.Intent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(intents) == 0 {
		s.intents = nil
		return
	}

	s.intents = intents
}

func (s *clientSession) GetBotData() *structs.BotData {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.botData
}

func (s *clientSession) SetBotData(data structs.BotData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.botData = &data
}

func (s *clientSession) GetSessionID() *string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sessionID
}

func (s *clientSession) SetSessionID(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessionID = &id
}

func (s *clientSession) GetSequence() *int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sequence
}

func (s *clientSession) SetSequence(seq int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sequence = &seq
}

func (s *clientSession) GetShard() *int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.shard
}

func (s *clientSession) SetShard(shard int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shard = &shard
}

func (s *clientSession) GetShards() *int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.shards
}

func (s *clientSession) SetShards(shards int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shards = &shards
}

func (s *clientSession) GetMaxConcurrency() *int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.maxConcurrency
}

func (s *clientSession) SetMaxConcurrency(concurrency int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.maxConcurrency = &concurrency
}

func (s *clientSession) AddServer(server structs.Server) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.servers[server.ID.ToString()] = &server
}

func (s *clientSession) GetServers() map[string]*structs.Server {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.servers
}

func (s *clientSession) SetHeartbeatAck(ack int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ack <= 0 {
		s.heartbeatAck = nil
		return
	}

	s.heartbeatAck = &ack
}

func (s *clientSession) GetHeartbeatAck() *int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heartbeatAck
}

func (s *clientSession) CloseHelloReceived() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeGroup.CloseChannels["helloReceived"].Do(func() {
		close(s.helloReceived)
	})
}

func (s *clientSession) CloseReadyReceived() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeGroup.CloseChannels["readyReceived"].Do(func() {
		close(s.readyReceived)
	})
}

func (s *clientSession) CloseResumeReceived() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeGroup.CloseChannels["resumeReceived"].Do(func() {
		close(s.resumeReceived)
	})
}

func (s *clientSession) GetCtx() context.Context {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ctx
}

func (s *clientSession) GetCancel() context.CancelFunc {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cancel
}

func (s *clientSession) SetResumeUrl(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resumeUrl = &url
}

func (s *clientSession) GetResumeUrl() *string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.resumeUrl
}

func (s *clientSession) GetServerByGuildID(guildID structs.Snowflake) *structs.Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, guild := range s.servers {
		if guild.ID.Equals(guildID) {
			return guild
		}
	}

	return nil
}

func (s *clientSession) SetCb(cb func(s ClientSession) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cb = cb
}

func (s *clientSession) SetEventHandler(handler *eventHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eventHandler = handler
}

func (s *clientSession) GetEventHandler() *eventHandler {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.eventHandler
}

func (s *clientSession) GetVoiceSession(guildID structs.Snowflake) VoiceSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.voiceSessions[guildID.ToString()]
}

func (s *clientSession) AddVoiceSession(guildID structs.Snowflake, vs VoiceSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.voiceSessions[guildID.ToString()] = vs
}

func (s *clientSession) dialer(query string) error {
	url := ""
	if s.GetResumeUrl() != nil {
		url = *s.GetResumeUrl()
	} else {
		gateway, err := requestutil.GetGatewayUrl(s.version)
		if err != nil {
			return err
		}
		url = gateway
		s.SetResumeUrl(url)
	}

	if err := s.Session.Connect(url+query, false); err != nil {
		return err
	}
	return nil
}

func (s *clientSession) initDialer(query string) error {
	if s.GetToken() == nil {
		return errors.New("token is required for bot init")
	}
	url := ""
	gateway, err := requestutil.GetGatewayBot(*s.GetToken(), s.version)
	if err != nil {
		return err
	}
	url = gateway.URL
	s.SetShards(gateway.Shards)
	s.SetMaxConcurrency(gateway.SessionStartLimit.MaxConcurrency)
	s.SetResumeUrl(url)

	if err := s.Session.Connect(url+query, false); err != nil {
		return err
	}
	return nil
}

func (s *clientSession) identify() error {
	idPayload := payload.SessionPayload{
		OpCode: gateway.GatewayOpIdentify,
	}

	if err := s.eventHandler.HandleEvent(s, idPayload); err != nil {
		return err
	}

	return nil
}

func (s *clientSession) resume() error {
	resumePayload := payload.SessionPayload{
		OpCode: gateway.GatewayOpResume,
	}

	if err := s.eventHandler.HandleEvent(s, resumePayload); err != nil {
		return err
	}

	return nil
}

func (s *clientSession) interactionReply(interactionOptions structs.InteractionResponseOptions, interaction *structs.Interaction) error {
	interactionID := interaction.ID.ToString()
	interactionToken := interaction.Token
	token := *s.GetToken()
	reqDto := dto.CreateInteractionResponseDto{
		WithResponse: util.ToPtr(true),
	}
	response := interactionOptions.InteractionResponse()
	if _, err := requestutil.CreateInteractionResponse(interactionID, interactionToken, token, reqDto, *response); err != nil {
		return err
	}
	return nil
}

func (s *clientSession) voiceStateUpdate(guildID, channelID *structs.Snowflake) error {
	vsuPayload := payload.SessionPayload{
		OpCode: gateway.GatewayOpVoiceStateUpdate,
		Data: sendevents.UpdateVoiceStateEvent{
			GuildID:   guildID,
			ChannelID: channelID,
			SelfMute:  false,
			SelfDeaf:  false,
		},
	}

	if err := s.eventHandler.HandleEvent(s, vsuPayload); err != nil {
		return err
	}

	return nil
}

func (s *clientSession) handleEvent(p payload.Payload) error {
	sp, ok := p.(*payload.SessionPayload)
	if !ok {
		return errors.New("invalid payload type - handle error")
	}

	return s.eventHandler.HandleEvent(s, *sp)
}

func (s *clientSession) validateEvent(p payload.Payload) (any, error) {
	var err error
	sp, ok := p.(*payload.SessionPayload)
	if !ok {
		return nil, errors.New("invalid session payload type - validate error: " + p.ToString())
	}
	sp.Data, err = receiveevents.NewReceiveEvent(*sp)
	if err != nil {
		return nil, err
	}

	return &sp, nil
}

func (s *clientSession) vsCleanup(guildID structs.Snowflake) func() {
	return func() {
		if err := s.voiceStateUpdate(&guildID, nil); err != nil {
			s.Error(err)
		}

		delete(s.voiceSessions, guildID.ToString())
	}
}

func (s *clientSession) vsResume(guildID structs.Snowflake, vs VoiceSession) func() {
	return func() {
		s.AddVoiceSession(guildID, vs)
	}
}

func (s *clientSession) vsReconnect(guildID, channelID structs.Snowflake) func() {
	return func() {
		if err := s.JoinVoice(guildID, channelID); err != nil {
			s.Error(err)
			return
		}
	}
}
