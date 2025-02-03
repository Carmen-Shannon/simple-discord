package session

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway/payload"
)

type Listener string

const (
	HelloListener                      Listener = "HELLO"
	ReadyListener                      Listener = "READY"
	ResumedListener                    Listener = "RESUMED"
	ReconnectListener                  Listener = "RECONNECT"
	InvalidSessionListener             Listener = "INVALID_SESSION"
	ChannelCreateListener              Listener = "CHANNEL_CREATE"
	ChannelUpdateListener              Listener = "CHANNEL_UPDATE"
	ChannelDeleteListener              Listener = "CHANNEL_DELETE"
	GuildCreateListener                Listener = "GUILD_CREATE"
	GuildUpdateListener                Listener = "GUILD_UPDATE"
	GuildDeleteListener                Listener = "GUILD_DELETE"
	GuildBanAddListener                Listener = "GUILD_BAN_ADD"
	GuildBanRemoveListener             Listener = "GUILD_BAN_REMOVE"
	GuildEmojisUpdateListener          Listener = "GUILD_EMOJIS_UPDATE"
	GuildIntegrationsUpdateListener    Listener = "GUILD_INTEGRATIONS_UPDATE"
	GuildAuditLogEntryCreateListener   Listener = "GUILD_AUDIT_LOG_ENTRY_CREATE"
	GuildMemberAddListener             Listener = "GUILD_MEMBER_ADD"
	GuildMemberRemoveListener          Listener = "GUILD_MEMBER_REMOVE"
	GuildMemberUpdateListener          Listener = "GUILD_MEMBER_UPDATE"
	GuildMembersChunkListener          Listener = "GUILD_MEMBERS_CHUNK"
	GuildRoleCreateListener            Listener = "GUILD_ROLE_CREATE"
	GuildRoleUpdateListener            Listener = "GUILD_ROLE_UPDATE"
	GuildRoleDeleteListener            Listener = "GUILD_ROLE_DELETE"
	MessageCreateListener              Listener = "MESSAGE_CREATE"
	MessageUpdateListener              Listener = "MESSAGE_UPDATE"
	MessageDeleteListener              Listener = "MESSAGE_DELETE"
	MessageBulkDeleteListener          Listener = "MESSAGE_BULK_DELETE"
	MessageReactionAddListener         Listener = "MESSAGE_REACTION_ADD"
	MessageReactionRemoveListener      Listener = "MESSAGE_REACTION_REMOVE"
	MessageReactionRemoveAllListener   Listener = "MESSAGE_REACTION_REMOVE_ALL"
	MessageReactionRemoveEmojiListener Listener = "MESSAGE_REACTION_REMOVE_EMOJI"
	MessagePollVoteAddListener         Listener = "MESSAGE_POLL_VOTE_ADD"
	MessagePollVoteRemoveListener      Listener = "MESSAGE_POLL_VOTE_REMOVE"
	TypingStartListener                Listener = "TYPING_START"
	UserUpdateListener                 Listener = "USER_UPDATE"
	VoiceChannelEffectSendListener     Listener = "VOICE_CHANNEL_EFFECT_SEND"
	VoiceStateUpdateListener           Listener = "VOICE_STATE_UPDATE"
	VoiceServerUpdateListener          Listener = "VOICE_SERVER_UPDATE"
	VoiceChannelStatusUpdateListener   Listener = "VOICE_CHANNEL_STATUS_UPDATE"
	WebhooksUpdateListener             Listener = "WEBHOOKS_UPDATE"
	PresenceUpdateListener             Listener = "PRESENCE_UPDATE"
)

// this is really just for helping me log more better, will remove eventually
var opCodeNames = map[gateway.GatewayOpCode]string{
	gateway.GatewayOpHeartbeat:           "Heartbeat",
	gateway.GatewayOpIdentify:            "Identify",
	gateway.GatewayOpPresenceUpdate:      "PresenceUpdate",
	gateway.GatewayOpVoiceStateUpdate:    "VoiceStateUpdate",
	gateway.GatewayOpResume:              "Resume",
	gateway.GatewayOpRequestGuildMembers: "RequestGuildMembers",
	gateway.GatewayOpInvalidSession:      "InvalidSession",
	gateway.GatewayOpHello:               "Hello",
	gateway.GatewayOpHeartbeatACK:        "HeartbeatACK",
	gateway.GatewayOpReconnect:           "Reconnect",
}

var voiceOpCodeNames = map[gateway.VoiceOpCode]string{
	gateway.VoiceOpIdentify:                    "Identify",
	gateway.VoiceOpSelectProtocol:              "Select Protocol",
	gateway.VoiceOpReady:                       "Ready",
	gateway.VoiceOpHeartbeat:                   "Heartbeat",
	gateway.VoiceOpSessionDescription:          "Session Description",
	gateway.VoiceOpSpeaking:                    "Speaking",
	gateway.VoiceOpHeartbeatAck:                "Heartbeat Ack",
	gateway.VoiceOpResume:                      "Resume",
	gateway.VoiceOpHello:                       "Hello",
	gateway.VoiceOpResumed:                     "Resumed",
	gateway.VoiceOpClientsConnect:              "Clients Connect",
	gateway.VoiceOpClientDisconnect:            "Client Disconnect",
	gateway.VoiceOpPrepareTransition:           "DAVE Prepare Transition",
	gateway.VoiceOpExecuteTransition:           "DAVE Execute Transition",
	gateway.VoiceOpTransitionReady:             "DAVE Transition Ready",
	gateway.VoiceOpPrepareEpoch:                "DAVE Prepare Epoch",
	gateway.VoiceOpMLSExternalSender:           "DAVE MLS External Sender",
	gateway.VoiceOpMLSKeyPackage:               "DAVE MLS Key Package",
	gateway.VoiceOpMLSProposals:                "DAVE MLS Proposals",
	gateway.VoiceOpMLSCommitWelcome:            "DAVE MLS Commit Welcome",
	gateway.VoiceOpMLSAnnounceCommitTransition: "DAVE MLS Announce Commit Transition",
	gateway.VoiceOpMLSWelcome:                  "DAVE MLS Welcome",
	gateway.VoiceOpMLSInvalidCommitWelcome:     "DAVE MLS Invalid Commit Welcome",
}

type eventHandler struct {
	NamedHandlers    map[string]CommandFunc
	OpCodeHandlers   map[gateway.GatewayOpCode]CommandFunc
	CustomHandlers   map[string]CommandFunc
	ListenerHandlers map[string]CommandFunc
}

type voiceEventHandler struct {
	OpCodeHandlers map[gateway.VoiceOpCode]VoiceEventFunc
	BinaryHandlers map[gateway.VoiceOpCode]BinaryVoiceEventFunc
}

type udpEventHandler struct {
	VoicePacketHandlers map[string]UdpEventFunc
	DiscoveryHandlers   map[string]UdpDiscoveryEventFunc
}

type EventHandler interface {
	HandleEvent(s Session, p payload.Payload) error
}

func NewEventHandler[T any]() *T {
	var t T
	tType := reflect.TypeOf(t)

	switch tType {
	case reflect.TypeOf(eventHandler{}):
		eh := newEventHandler()
		return any(eh).(*T)
	case reflect.TypeOf(voiceEventHandler{}):
		eh := newVoiceEventHandler()
		return any(eh).(*T)
	case reflect.TypeOf(udpEventHandler{}):
		eh := newUdpEventHandler()
		return any(eh).(*T)
	default:
		return nil
	}
}

// sets up a new EventHandler with the default Discord handlers
func newEventHandler() *eventHandler {
	e := &eventHandler{
		OpCodeHandlers: map[gateway.GatewayOpCode]CommandFunc{
			gateway.GatewayOpHeartbeat:           handleHeartbeatEvent,
			gateway.GatewayOpIdentify:            handleSendIdentifyEvent,
			gateway.GatewayOpPresenceUpdate:      handleSendPresenceUpdateEvent,
			gateway.GatewayOpVoiceStateUpdate:    handleSendVoiceStateUpdateEvent,
			gateway.GatewayOpReconnect:           handleReconnectEvent,
			gateway.GatewayOpResume:              handleSendResumeEvent,
			gateway.GatewayOpRequestGuildMembers: handleSendRequestGuildMembersEvent,
			gateway.GatewayOpInvalidSession:      handleInvalidSessionEvent,
			gateway.GatewayOpHello:               handleHelloEvent,
			gateway.GatewayOpHeartbeatACK:        handleHeartbeatACKEvent,
		},
		CustomHandlers:   map[string]CommandFunc{},
		ListenerHandlers: map[string]CommandFunc{},
	}

	e.NamedHandlers = map[string]CommandFunc{
		"HELLO":                         handleHelloEvent,
		"READY":                         handleReadyEvent,
		"RESUMED":                       handleResumedEvent,
		"RECONNECT":                     handleReconnectEvent,
		"INVALID_SESSION":               handleInvalidSessionEvent,
		"CHANNEL_CREATE":                handleChannelCreateEvent,
		"CHANNEL_UPDATE":                handleChannelUpdateEvent,
		"CHANNEL_DELETE":                handleChannelDeleteEvent,
		"GUILD_CREATE":                  handleGuildCreateEvent,
		"GUILD_UPDATE":                  handleGuildUpdateEvent,
		"GUILD_DELETE":                  handleGuildDeleteEvent,
		"GUILD_BAN_ADD":                 handleGuildBanAddEvent,
		"GUILD_BAN_REMOVE":              handleGuildBanRemoveEvent,
		"GUILD_EMOJIS_UPDATE":           handleGuildEmojisUpdateEvent,
		"GUILD_INTEGRATIONS_UPDATE":     handleGuildIntegrationsUpdateEvent,
		"GUILD_AUDIT_LOG_ENTRY_CREATE":  handleGuildAuditLogEntryCreateEvent,
		"GUILD_MEMBER_ADD":              handleGuildMemberAddEvent,
		"GUILD_MEMBER_REMOVE":           handleGuildMemberRemoveEvent,
		"GUILD_MEMBER_UPDATE":           handleGuildMemberUpdateEvent,
		"GUILD_MEMBERS_CHUNK":           handleGuildMembersChunkEvent,
		"GUILD_ROLE_CREATE":             handleGuildRoleCreateEvent,
		"GUILD_ROLE_UPDATE":             handleGuildRoleUpdateEvent,
		"GUILD_ROLE_DELETE":             handleGuildRoleDeleteEvent,
		"MESSAGE_CREATE":                handleMessageCreateEvent,
		"MESSAGE_UPDATE":                handleMessageUpdateEvent,
		"MESSAGE_DELETE":                handleMessageDeleteEvent,
		"MESSAGE_BULK_DELETE":           handleMessageBulkDeleteEvent,
		"MESSAGE_REACTION_ADD":          handleMessageReactionAddEvent,
		"MESSAGE_REACTION_REMOVE":       handleMessageReactionRemoveEvent,
		"MESSAGE_REACTION_REMOVE_ALL":   handleMessageReactionRemoveAllEvent,
		"MESSAGE_REACTION_REMOVE_EMOJI": handleMessageReactionRemoveEmojiEvent,
		"MESSAGE_POLL_VOTE_ADD":         handleMessagePollVoteAddEvent,
		"MESSAGE_POLL_VOTE_REMOVE":      handleMessagePollVoteRemoveEvent,
		"TYPING_START":                  handleTypingStartEvent,
		"USER_UPDATE":                   handleUserUpdateEvent,
		"VOICE_CHANNEL_EFFECT_SEND":     handleVoiceChannelEffectSendEvent,
		"VOICE_STATE_UPDATE":            handleVoiceStateUpdateEvent,
		"VOICE_SERVER_UPDATE":           handleVoiceServerUpdateEvent,
		"VOICE_CHANNEL_STATUS_UPDATE":   handleVoiceChannelStatusUpdateEvent,
		"WEBHOOKS_UPDATE":               handleWebhooksUpdateEvent,
		"PRESENCE_UPDATE":               handlePresenceUpdateEvent,
		"INTERACTION_CREATE":            e.handleInteractionCreateEvent,
	}
	return e
}

func newVoiceEventHandler() *voiceEventHandler {
	e := &voiceEventHandler{
		OpCodeHandlers: map[gateway.VoiceOpCode]VoiceEventFunc{
			gateway.VoiceOpIdentify:                    handleSendVoiceIdentifyEvent,
			gateway.VoiceOpSelectProtocol:              handleSendVoiceSelectProtocolEvent,
			gateway.VoiceOpReady:                       handleVoiceReadyEvent,
			gateway.VoiceOpHeartbeat:                   handleVoiceSendHeartbeatEvent,
			gateway.VoiceOpSessionDescription:          handleVoiceSessionDescriptionEvent,
			gateway.VoiceOpSpeaking:                    handleVoiceSpeakingEvent,
			gateway.VoiceOpHeartbeatAck:                handleVoiceHeartbeatAckEvent,
			gateway.VoiceOpResume:                      handleSendVoiceResumeEvent,
			gateway.VoiceOpHello:                       handleVoiceHelloEvent,
			gateway.VoiceOpResumed:                     handleVoiceResumedEvent,
			gateway.VoiceOpClientsConnect:              handleVoiceClientsConnectEvent,
			gateway.VoiceOpClientDisconnect:            handleVoiceClientDisconnectEvent,
			gateway.VoiceOpPrepareTransition:           handleVoicePrepareTransitionEvent,
			gateway.VoiceOpExecuteTransition:           handleVoiceExecuteTransitionEvent,
			gateway.VoiceOpTransitionReady:             handleSendVoiceTransitionReadyEvent,
			gateway.VoiceOpPrepareEpoch:                handleVoicePrepareEpochEvent,
			gateway.VoiceOpMLSAnnounceCommitTransition: handleVoiceMLSAnnounceCommitTransitionEvent,
			gateway.VoiceOpMLSInvalidCommitWelcome:     handleSendVoiceMLSInvalidCommitWelcomeEvent,
		},
		BinaryHandlers: map[gateway.VoiceOpCode]BinaryVoiceEventFunc{
			gateway.VoiceOpMLSExternalSender: handleVoiceMLSExternalSenderEvent,
			gateway.VoiceOpMLSKeyPackage:     handleSendVoiceMLSKeyPackageEvent,
			gateway.VoiceOpMLSProposals:      handleVoiceMLSProposalsEvent,
			gateway.VoiceOpMLSCommitWelcome:  handleSendVoiceMLSCommitWelcomeEvent,
			gateway.VoiceOpMLSWelcome:        handleVoiceMLSWelcomeEvent,
		},
	}
	return e
}

func newUdpEventHandler() *udpEventHandler {
	e := &udpEventHandler{}
	return e
}

// HandleEvent handles events (duh)
// first we need to check if there is an EventName attached to the payload, so we can map it to the correct handler
// if there is no EventName then we use the OpCode handlers
// this function can handle sending events as well, just pass it the payload with the appropriate EventName or OpCode and let it fly
func (e *eventHandler) HandleEvent(s ClientSession, payload payload.SessionPayload) error {
	// check first for the payload event name ("t" field in the raw payload) to see if it was omitted
	// if it's not there run with the OpCode
	if payload.EventName == nil {
		fmt.Printf("HANDLING OPCODE EVENT: %v, %s\n", payload.OpCode, opCodeNames[payload.OpCode])
		if handler, ok := e.OpCodeHandlers[payload.OpCode]; ok && handler != nil {
			// if the payload has a sequence number update the Session with the latest sequence
			if payload.Seq != nil {
				s.SetSequence(*payload.Seq)
			}

			// let her rip tater chip
			go func() {
				if err := handler(s, payload); err != nil {
					fmt.Printf("ERROR HANDLING OPCODE EVENT: %v, %s, %v\n", payload.OpCode, opCodeNames[payload.OpCode], err)
					s.Error(err)
				}
			}()
			return nil
		}
		return errors.New("no handler for opcode")
	}

	// if we haven't returned from the above if-else, check the actual event name
	if handler, ok := e.NamedHandlers[*payload.EventName]; ok && handler != nil {
		fmt.Printf("HANDLING NAMED EVENT: %v\n", *payload.EventName)
		// if the payload has a sequence number update the Session with the latest sequence
		if payload.Seq != nil {
			s.SetSequence(*payload.Seq)
		}

		// let her rip tater chip
		go func() {
			if err := handler(s, payload); err != nil {
				s.Error(err)
			}

			// check if there are any listeners for this event
			if listener, ok := e.ListenerHandlers[*payload.EventName]; ok && listener != nil {
				if err := listener(s, payload); err != nil {
					s.Error(err)
				}
			}
		}()
		return nil
	}
	return errors.New("no handler for event name")
}

func (e *eventHandler) AddCustomHandler(name string, handler func(ClientSession, payload.SessionPayload) error) {
	e.CustomHandlers[name] = handler
}

func (e *eventHandler) AddListener(event string, handler func(ClientSession, payload.SessionPayload) error) {
	e.ListenerHandlers[event] = handler
}

func (e *voiceEventHandler) HandleEvent(s VoiceSession, p payload.Payload) error {
	if vp, ok := p.(*payload.VoicePayload); ok {
		return e.handleStandardEvent(s, *vp)
	}
	if bp, ok := p.(*payload.BinaryVoicePayload); ok {
		return e.handleBinaryEvent(s, *bp)
	}
	return errors.New("invalid payload type")
}

// HandleEvent handles voice events the same way as HandleEvent for normal events
func (e *voiceEventHandler) handleStandardEvent(s VoiceSession, p payload.VoicePayload) error {
	fmt.Printf("HANDLING VOICE EVENT: %v, %s\n", p.OpCode, voiceOpCodeNames[p.OpCode])
	if handler, ok := e.OpCodeHandlers[p.OpCode]; ok && handler != nil {
		if p.Seq != nil {
			s.SetSequence(*p.Seq)
		}
		go func() {
			if err := handler(s, p); err != nil {
				s.Error(fmt.Errorf("error handling voice event: %v\n%s", err, p.ToString()))
			}
		}()
		return nil
	}
	return errors.New("no handler for voice opcode")
}

func (e *voiceEventHandler) handleBinaryEvent(s VoiceSession, p payload.BinaryVoicePayload) error {
	if handler, ok := e.BinaryHandlers[gateway.VoiceOpCode(p.OpCode)]; ok && handler != nil {
		if p.SequenceNumber != nil {
			s.SetSequence(int(*p.SequenceNumber))
		}

		go func() {
			if err := handler(s, p); err != nil {
				s.Error(err)
			}
		}()
		return nil
	}
	return nil
}

func (e *udpEventHandler) HandleEvent(s UdpSession, p payload.Payload) error {
	go func() {
		if dp, ok := p.(*payload.DiscoveryPacket); ok {
			if err := handleDiscoveryEvent(s, *dp); err != nil {
				s.Error(err)
			}
		}
		if vp, ok := p.(*payload.VoicePacket); ok {
			if err := handleVoicePacketEvent(s, *vp); err != nil {
				s.Error(err)
			}
		}
		if _, ok := p.(*payload.SenderReportPacket); ok {
			// sender report packet, currently has no use when receiving
			return
		}
	}()
	return nil
}
