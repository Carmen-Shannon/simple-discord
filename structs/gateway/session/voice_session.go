package session

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway/payload"
	receiveevents "github.com/Carmen-Shannon/simple-discord/structs/gateway/receive_events"
	"github.com/coder/websocket"
)

type VoiceEventFunc func(VoiceSession, payload.VoicePayload) error
type BinaryVoiceEventFunc func(VoiceSession, payload.BinaryVoicePayload) error

type voiceSession struct {
	Session
	mu     *sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	token        *string
	botData      *structs.BotData
	heartbeatAck *int
	connectUrl   *string
	sessionID    *string
	sequence     *int
	guildID      *structs.Snowflake
	channelID    *structs.Snowflake

	connected              bool
	voiceStateReadySignal  bool
	voiceServerReadySignal bool

	audioPlayer  AudioPlayer
	eventHandler *voiceEventHandler

	closeGroup    structs.SyncGroup
	connectReady  chan struct{}
	resumeReady   chan struct{}
	readyReceived chan struct{}
}

type VoiceSession interface {
	Write(data []byte, binary bool)
	Error(err error)
	Exit(graceful bool) error
	Connect() error
	IsConnected() bool
	SetToken(token string)
	GetToken() *string
	SetBotData(botData structs.BotData)
	GetBotData() *structs.BotData
	SetHeartbeatAck(heartbeatAck int)
	GetHeartbeatAck() *int
	SetConnectUrl(connectUrl string)
	SetSessionID(sessionID string)
	GetSessionID() *string
	SetSequence(sequence int)
	GetSequence() *int
	SetGuildID(guildID structs.Snowflake)
	GetGuildID() *structs.Snowflake
	SetChannelID(channelID structs.Snowflake)
	GetChannelID() *structs.Snowflake
	SetEventHandler(eventHandler *voiceEventHandler)
	GetCtx() context.Context
	GetConnectReady() <-chan struct{}
	GetAudioPlayer() AudioPlayer
	SignalVoiceStateReady()
	SignalVoiceServerReady()
	CloseConnectReady()
	CloseResumeReady()
	CloseReadyReceived()
}

var _ VoiceSession = (*voiceSession)(nil)

func NewVoiceSession() VoiceSession {
	vs := &voiceSession{
		mu:            &sync.Mutex{},
		Session:       NewSession(),
		eventHandler:  NewEventHandler[voiceEventHandler](),
		closeGroup:    *structs.NewSyncGroup(),
		connectReady:  make(chan struct{}),
		resumeReady:   make(chan struct{}),
		readyReceived: make(chan struct{}),
	}
	vs.ctx, vs.cancel = context.WithCancel(context.Background())

	vs.audioPlayer = NewAudioPlayer(vs.speaking, vs.selectProtocol)

	vs.closeGroup.AddChannel("connectReady")
	vs.closeGroup.AddChannel("resumeReady")
	vs.closeGroup.AddChannel("readyReceived")

	vs.SetListenFunc(vs.validateEvent)
	vs.SetHandleFunc(vs.handleEvent)
	vs.SetPayloadDecoders(&payload.RawMessagePayload{}, &payload.BinaryVoicePayload{})
	vs.SetEventDecoders(&payload.VoicePayload{}, &payload.BinaryVoicePayload{})
	vs.SetStatusCodeHandlers(map[websocket.StatusCode]func(){
		websocket.StatusNormalClosure: func() {},
		4006: func() {
			vs.Exit(false)
		},
		4014: func() {
			vs.Exit(false)
		},
	})
	vs.SetErrorHandlers(map[error]func(){
		io.EOF: func() {
			vs.Exit(false)
		},
		net.ErrClosed: func() {
			vs.Exit(false)
		},
	})
	vs.SetValidCloseErrors(io.EOF, io.ErrUnexpectedEOF, net.ErrClosed)

	return vs
}

func (v *voiceSession) Exit(graceful bool) error {
	defer v.cancel()
	v.CloseConnectReady()
	v.CloseResumeReady()

	if v.audioPlayer.IsConnected() {
		v.audioPlayer.Exit()
	}

	return v.Session.Exit(graceful)
}

func (v *voiceSession) Connect() error {
	query := "?v=8"
	url := *v.connectUrl

	if err := v.Session.Connect(url+query, false); err != nil {
		return err
	}
	v.mu.Lock()
	v.connected = true
	v.mu.Unlock()

	if err := v.identify(); err != nil {
		return err
	}

	select {
	case <-v.ctx.Done():
		return nil
	case <-v.readyReceived:
	}
	return nil
}

func (v *voiceSession) IsConnected() bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.connected
}

func (v *voiceSession) SetToken(token string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.token = &token
}

func (v *voiceSession) GetToken() *string {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.token
}

func (v *voiceSession) SetBotData(botData structs.BotData) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.botData = &botData
}

func (v *voiceSession) GetBotData() *structs.BotData {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.botData
}

func (v *voiceSession) SetHeartbeatAck(heartbeatAck int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.heartbeatAck = &heartbeatAck
}

func (v *voiceSession) GetHeartbeatAck() *int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.heartbeatAck
}

func (v *voiceSession) SetConnectUrl(connectUrl string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.connectUrl = &connectUrl
}

func (v *voiceSession) SetSessionID(sessionID string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.sessionID = &sessionID
}

func (v *voiceSession) GetSessionID() *string {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.sessionID
}

func (v *voiceSession) SetSequence(sequence int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.sequence = &sequence
}

func (v *voiceSession) GetSequence() *int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.sequence
}

func (v *voiceSession) SetGuildID(guildID structs.Snowflake) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.guildID = &guildID
}

func (v *voiceSession) GetGuildID() *structs.Snowflake {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.guildID
}

func (v *voiceSession) SetChannelID(channelID structs.Snowflake) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.channelID = &channelID
}

func (v *voiceSession) GetChannelID() *structs.Snowflake {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.channelID
}

func (v *voiceSession) SetEventHandler(eventHandler *voiceEventHandler) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.eventHandler = eventHandler
}

func (v *voiceSession) GetCtx() context.Context {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.ctx
}

func (v *voiceSession) GetConnectReady() <-chan struct{} {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.connectReady
}

func (v *voiceSession) GetAudioPlayer() AudioPlayer {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.audioPlayer
}

func (v *voiceSession) SignalVoiceStateReady() {
	v.mu.Lock()
	v.voiceStateReadySignal = true
	v.mu.Unlock()

	if v.voiceServerReadySignal {
		v.CloseConnectReady()
	}
}

func (v *voiceSession) SignalVoiceServerReady() {
	v.mu.Lock()
	v.voiceServerReadySignal = true
	v.mu.Unlock()

	if v.voiceStateReadySignal {
		v.CloseConnectReady()
	}
}

func (v *voiceSession) CloseConnectReady() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.closeGroup.CloseChannels["connectReady"].Do(func() {
		close(v.connectReady)
	})
}

func (v *voiceSession) CloseResumeReady() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.closeGroup.CloseChannels["resumeReady"].Do(func() {
		close(v.resumeReady)
	})
}

func (v *voiceSession) CloseReadyReceived() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.closeGroup.CloseChannels["readyReceived"].Do(func() {
		close(v.readyReceived)
	})
}

func (v *voiceSession) identify() error {
	idp := &payload.VoicePayload{
		OpCode: gateway.VoiceOpIdentify,
	}

	if err := v.eventHandler.HandleEvent(v, idp); err != nil {
		return err
	}
	return nil
}

func (v *voiceSession) speaking() error {
	sp := &payload.VoicePayload{
		OpCode: gateway.VoiceOpSpeaking,
		Data:   nil,
	}

	if err := v.eventHandler.HandleEvent(v, sp); err != nil {
		return err
	}
	return nil
}

func (v *voiceSession) selectProtocol() error {
	sp := &payload.VoicePayload{
		OpCode: gateway.VoiceOpSelectProtocol,
	}

	if err := v.eventHandler.HandleEvent(v, sp); err != nil {
		return err
	}
	return nil
}

func (v *voiceSession) handleEvent(p payload.Payload) error {
	return v.eventHandler.HandleEvent(v, p)
}

func (v *voiceSession) validateEvent(p payload.Payload) (any, error) {
	var err error
	vp, ok := p.(*payload.VoicePayload)
	if !ok {
		bp, ok := p.(*payload.BinaryVoicePayload)
		if !ok {
			return nil, errors.New("invalid voice payload type - validate error: " + p.ToString())
		}
		v.Error(errors.New("binary voice events not implemented"))
		// bp.Data, err = receiveevents.NewBinaryVoiceReceiveEvent(*bp)
		// if err != nil {
		// 	return nil, err
		// }

		return &bp, nil
	}
	vp.Data, err = receiveevents.NewVoiceReceiveEvent(*vp)
	if err != nil {
		return nil, err
	}

	return &vp, nil
}
