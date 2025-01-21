package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"

	receiveevents "github.com/Carmen-Shannon/simple-discord/gateway/receive_events"
	requestutil "github.com/Carmen-Shannon/simple-discord/gateway/request_util"
	sendevents "github.com/Carmen-Shannon/simple-discord/gateway/send_events"
	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	gateway "github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"github.com/Carmen-Shannon/simple-discord/util"
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
	gateway.Heartbeat:           "Heartbeat",
	gateway.Identify:            "Identify",
	gateway.PresenceUpdate:      "PresenceUpdate",
	gateway.VoiceStateUpdate:    "VoiceStateUpdate",
	gateway.Resume:              "Resume",
	gateway.RequestGuildMembers: "RequestGuildMembers",
	gateway.Hello:               "Hello",
	gateway.HeartbeatACK:        "HeartbeatACK",
	gateway.Reconnect:           "Reconnect",
}

type EventFunc func(Session, gateway.Payload) error

type EventHandler struct {
	NamedHandlers    map[string]EventFunc
	OpCodeHandlers   map[gateway.GatewayOpCode]EventFunc
	CustomHandlers   map[string]EventFunc
	ListenerHandlers map[string]EventFunc
}

// sets up a new EventHandler with the default Discord handlers
func NewEventHandler() *EventHandler {
	e := &EventHandler{
		OpCodeHandlers: map[gateway.GatewayOpCode]EventFunc{
			gateway.Heartbeat:           handleHeartbeatEvent,
			gateway.Identify:            handleSendIdentifyEvent,
			gateway.PresenceUpdate:      handleSendPresenceUpdateEvent,
			gateway.VoiceStateUpdate:    handleSendVoiceStateUpdateEvent,
			gateway.Reconnect:           handleReconnectEvent,
			gateway.Resume:              handleSendResumeEvent,
			gateway.RequestGuildMembers: handleSendRequestGuildMembersEvent,
			gateway.Hello:               handleHelloEvent,
			gateway.HeartbeatACK:        handleHeartbeatACKEvent,
		},
		CustomHandlers:   map[string]EventFunc{},
		ListenerHandlers: map[string]EventFunc{},
	}

	e.NamedHandlers = map[string]EventFunc{
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

// HandleEvent handles events (duh)
// first we need to check if there is an EventName attached to the payload, so we can map it to the correct handler
// if there is no EventName then we use the OpCode handlers
// this function can handle sending events as well, just pass it the payload with the appropriate EventName or OpCode and let it fly
func (e *EventHandler) HandleEvent(s Session, payload gateway.Payload) error {
	// check first for the payload event name ("t" field in the raw payload) to see if it was omitted
	// if it's not there run with the OpCode
	if payload.EventName == nil {
		fmt.Printf("HANDLING OPCODE EVENT: %v, %s\n", payload.OpCode, opCodeNames[payload.OpCode])
		if handler, ok := e.OpCodeHandlers[payload.OpCode]; ok && handler != nil {
			// if the payload has a sequence number update the Session with the latest sequence
			if payload.Seq != nil {
				s.SetSequence(payload.Seq)
			}

			// let her rip tater chip
			go func() {
				if err := handler(s, payload); err != nil {
					s.GetSession().errorChan <- err
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
			s.SetSequence(payload.Seq)
		}

		// let her rip tater chip
		go func() {
			if err := handler(s, payload); err != nil {
				s.GetSession().errorChan <- err
			}

			// check if there are any listeners for this event
			if listener, ok := e.ListenerHandlers[*payload.EventName]; ok && listener != nil {
				if err := listener(s, payload); err != nil {
					s.GetSession().errorChan <- err
				}
			}
		}()
		return nil
	}
	return errors.New("no handler for event name")
}

func (e *EventHandler) AddCustomHandler(name string, handler func(Session, gateway.Payload) error) {
	e.CustomHandlers[name] = handler
}

func (e *EventHandler) AddListener(event string, handler func(Session, gateway.Payload) error) {
	e.ListenerHandlers[event] = handler
}

func (e *EventHandler) handleInteractionCreateEvent(s Session, p gateway.Payload) error {
	if interactionCreateEvent, ok := p.Data.(receiveevents.InteractionCreateEvent); ok {
		name := interactionCreateEvent.Data.Name
		if handler, ok := e.CustomHandlers[name]; ok && handler != nil {
			go func() {
				if err := handler(s, p); err != nil {
					s.GetSession().errorChan <- err
				}
			}()
			return nil
		}
		return errors.New("no handler for interaction")
	}
	return errors.New("unexpected payload data type")
}

func handleSendRequestGuildMembersEvent(s Session, p gateway.Payload) error {
	fmt.Println("HANDLING REQUEST GUILD MEMBERS EVENT")
	fmt.Println("REQUEST GUILD MEMBERS NOT IMPLEMENTED")
	return nil
}

func handleSendVoiceStateUpdateEvent(s Session, p gateway.Payload) error {
	if voiceStateUpdateEvent, ok := p.Data.(sendevents.UpdateVoiceStateEvent); ok {
		if voiceStateUpdateEvent.GuildID == nil {
			return errors.New("guild ID not set")
		}

		vs := s.GetVoiceSession(*voiceStateUpdateEvent.GuildID)
		if vs == nil {
			return errors.New("voice session not initialized")
		}

		p.Seq = vs.GetSequence()

		voiceStateData, err := json.Marshal(p)
		if err != nil {
			return err
		}

		s.Write(voiceStateData)
		return nil
	} else {
		return errors.New("unexpected payload data type")
	}
}

func handleSendPresenceUpdateEvent(s Session, p gateway.Payload) error {
	fmt.Println("HANDLING PRESENCE UPDATE EVENT")
	fmt.Println("PRESENCE UPDATE NOT IMPLEMENTED")
	return nil
}

func handleSendResumeEvent(s Session, p gateway.Payload) error {
	resumeEvent := sendevents.ResumeEvent{
		Token:     *s.GetToken(),
		SessionID: *s.GetID(),
		Seq:       *s.GetSequence(),
	}
	resumePayload := gateway.Payload{
		OpCode: gateway.Resume,
		Data:   resumeEvent,
	}

	resumeData, err := json.Marshal(resumePayload)
	if err != nil {
		return err
	}

	s.Write(resumeData)
	return nil
}

func handleSendIdentifyEvent(s Session, p gateway.Payload) error {
	identifyEvent := sendevents.IdentifyEvent{
		Token: *s.GetToken(),
		Properties: sendevents.IdentifyProperties{
			Os:      runtime.GOOS,
			Browser: "discord",
			Device:  "discord",
		},
		Intents: structs.GetIntents(s.GetIntents()),
	}
	if s.GetShard() != nil && s.GetShards() != nil {
		identifyEvent.Shard = &[]int{*s.GetShard(), *s.GetShards()}
	}
	identifyPayload := gateway.Payload{
		OpCode: gateway.Identify,
		Data:   identifyEvent,
	}

	identifyData, err := json.Marshal(identifyPayload)
	if err != nil {
		return err
	}

	s.Write(identifyData)
	return nil
}

func handleChannelDeleteEvent(s Session, p gateway.Payload) error {
	if channelDeleteEvent, ok := p.Data.(receiveevents.ChannelDeleteEvent); ok {
		servers := s.GetServers()
		server, exists := servers[channelDeleteEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}
		server.DeleteChannel(channelDeleteEvent.ID)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleChannelUpdateEvent(s Session, p gateway.Payload) error {
	if channelUpdateEvent, ok := p.Data.(receiveevents.ChannelUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[channelUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}
		channel := structs.Channel(*channelUpdateEvent.Channel)
		server.UpdateChannel(channel.ID, channel)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleChannelCreateEvent(s Session, p gateway.Payload) error {
	if channelCreateEvent, ok := p.Data.(receiveevents.ChannelCreateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[channelCreateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}
		channel := structs.Channel(*channelCreateEvent.Channel)
		channel.Typing = structs.NewTypingChannel()
		server.AddChannel(channel)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handlePresenceUpdateEvent(s Session, p gateway.Payload) error {
	if presenceUpdateEvent, ok := p.Data.(receiveevents.PresenceUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[presenceUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.UpdatePresence(presenceUpdateEvent.User.ID, *presenceUpdateEvent.PresenceUpdate)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleInvalidSessionEvent(s Session, p gateway.Payload) error {
	if invalidSessionEvent, ok := p.Data.(receiveevents.InvalidSessionEvent); ok {
		if invalidSessionEvent {
			if err := s.ReconnectSession(); err != nil {
				return err
			}
			fmt.Println("RECONNECTED SESSION")
		} else {
			s.Exit()
			fmt.Println("SESSION ENDED")
		}
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleReconnectEvent(s Session, p gateway.Payload) error {
	if _, ok := p.Data.(receiveevents.ReconnectEvent); ok {
		if err := s.ReconnectSession(); err != nil {
			return err
		}
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleResumedEvent(s Session, p gateway.Payload) error {
	return nil
}

func handleGuildCreateEvent(s Session, p gateway.Payload) error {
	if guildCreateEvent, ok := p.Data.(receiveevents.GuildCreateEvent); ok {
		if guildCreateEvent.Unavailable != nil && !*guildCreateEvent.Unavailable {
			server := structs.NewServer(guildCreateEvent.Guild)
			if err := util.UpdateFields(server, guildCreateEvent.Server); err != nil {
				return err
			}
			s.AddServer(*server)
		}
	} else if guildCreateUnavailableEvent, ok := p.Data.(receiveevents.GuildCreateUnavailableEvent); ok {
		unavailableGuild := &structs.Guild{}
		unavailableGuild.ID = guildCreateUnavailableEvent.ID
		server := structs.NewServer(unavailableGuild)
		server.Unavailable = &guildCreateUnavailableEvent.Unavailable
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleGuildUpdateEvent(s Session, p gateway.Payload) error {
	if guildUpdateEvent, ok := p.Data.(receiveevents.GuildUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildUpdateEvent.ID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.UpdateGuild(*guildUpdateEvent.Guild)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildDeleteEvent(s Session, p gateway.Payload) error {
	if guildDeleteEvent, ok := p.Data.(receiveevents.GuildDeleteEvent); ok {
		servers := s.GetServers()
		_, exists := servers[guildDeleteEvent.ID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		unavailableGuild := &structs.Guild{}
		unavailableGuild.ID = guildDeleteEvent.ID
		server := structs.NewServer(unavailableGuild)
		server.Unavailable = &guildDeleteEvent.Unavailable
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildBanAddEvent(s Session, p gateway.Payload) error {
	fmt.Println("HANDLING GUILD BAN ADD EVENT")
	fmt.Println("GUILD BAN ADD NOT IMPLEMENTED")
	return nil
}

func handleGuildBanRemoveEvent(s Session, p gateway.Payload) error {
	fmt.Println("HANDLING GUILD BAN REMOVE EVENT")
	fmt.Println("GUILD BAN REMOVE NOT IMPLEMENTED")
	return nil
}

func handleGuildEmojisUpdateEvent(s Session, p gateway.Payload) error {
	if guildEmojisUpdateEvent, ok := p.Data.(receiveevents.GuildEmojisUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildEmojisUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.Emojis = guildEmojisUpdateEvent.Emojis
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildIntegrationsUpdateEvent(s Session, p gateway.Payload) error {
	if _, ok := p.Data.(receiveevents.GuildIntegrationsUpdateEvent); ok {
		return nil
	}
	return errors.New("unexpected payload data type")
}

func handleGuildMemberAddEvent(s Session, p gateway.Payload) error {
	if guildMemberAddEvent, ok := p.Data.(receiveevents.GuildMemberAddEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildMemberAddEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.AddMember(*guildMemberAddEvent.GuildMember)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildMemberRemoveEvent(s Session, p gateway.Payload) error {
	if guildMemberRemoveEvent, ok := p.Data.(receiveevents.GuildMemberRemoveEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildMemberRemoveEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.DeleteMember(guildMemberRemoveEvent.User.ID)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildMemberUpdateEvent(s Session, p gateway.Payload) error {
	if guildMemberUpdateEvent, ok := p.Data.(receiveevents.GuildMemberUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildMemberUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMember := server.GetMember(guildMemberUpdateEvent.User.ID)
		if currentMember == nil {
			return errors.New("user not in server")
		}

		if err := util.UpdateFields(currentMember, &guildMemberUpdateEvent); err != nil {
			return err
		}

		server.UpdateMember(guildMemberUpdateEvent.User.ID, *currentMember)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildMembersChunkEvent(s Session, p gateway.Payload) error {
	if guildMembersChunkEvent, ok := p.Data.(receiveevents.GuildMembersChunk); ok {
		servers := s.GetServers()
		server, exists := servers[guildMembersChunkEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		for _, member := range guildMembersChunkEvent.Members {
			if server.HasMember(member.User.ID) {
				server.UpdateMember(member.User.ID, member)
			} else {
				server.AddMember(member)
			}
		}

		for _, presence := range guildMembersChunkEvent.Presences {
			if server.HasPresence(presence.User.ID) {
				server.UpdatePresence(presence.User.ID, *presence.PresenceUpdate)
			} else {
				server.AddPresence(*presence.PresenceUpdate)
			}
		}

		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildRoleCreateEvent(s Session, p gateway.Payload) error {
	if guildRoleCreateEvent, ok := p.Data.(receiveevents.GuildRoleCreateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildRoleCreateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.AddRole(guildRoleCreateEvent.Role)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildRoleUpdateEvent(s Session, p gateway.Payload) error {
	if guildRoleUpdateEvent, ok := p.Data.(receiveevents.GuildRoleUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildRoleUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.UpdateRole(guildRoleUpdateEvent.Role.ID, guildRoleUpdateEvent.Role)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildRoleDeleteEvent(s Session, p gateway.Payload) error {
	if guildRoleDeleteEvent, ok := p.Data.(receiveevents.GuildRoleDeleteEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildRoleDeleteEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.DeleteRole(guildRoleDeleteEvent.RoleID)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageCreateEvent(s Session, p gateway.Payload) error {
	if messageCreateEvent, ok := p.Data.(receiveevents.MessageCreateEvent); ok {
		if messageCreateEvent.GuildID != nil {
			servers := s.GetServers()
			server, exists := servers[messageCreateEvent.GuildID.ToString()]
			if !exists {
				return errors.New("server not found")
			}

			server.AddMessage(*messageCreateEvent.Message)
		} else {
			return nil
		}
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageUpdateEvent(s Session, p gateway.Payload) error {
	if messageUpdateEvent, ok := p.Data.(receiveevents.MessageUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[messageUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		if server.GetMessage(messageUpdateEvent.ChannelID, messageUpdateEvent.Message.ID) == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = messageUpdateEvent.ChannelID
			query.MessageID = messageUpdateEvent.Message.ID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}
			server.AddMessage(*message)
		}

		server.UpdateMessage(*messageUpdateEvent.Message)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageDeleteEvent(s Session, p gateway.Payload) error {
	if messageDeleteEvent, ok := p.Data.(receiveevents.MessageDeleteEvent); ok {
		servers := s.GetServers()
		server, exists := servers[messageDeleteEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.DeleteMessage(messageDeleteEvent.ChannelID, messageDeleteEvent.ID)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageBulkDeleteEvent(s Session, p gateway.Payload) error {
	if messageBulkDeleteEvent, ok := p.Data.(receiveevents.MessageDeleteBulkEvent); ok {
		servers := s.GetServers()
		server, exists := servers[messageBulkDeleteEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		for _, id := range messageBulkDeleteEvent.IDs {
			server.DeleteMessage(messageBulkDeleteEvent.ChannelID, id)
		}
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageReactionAddEvent(s Session, p gateway.Payload) error {
	if reactionAddEvent, ok := p.Data.(receiveevents.MessageReactionAddEvent); ok {
		servers := s.GetServers()
		server, exists := servers[reactionAddEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMessage := server.GetMessage(reactionAddEvent.ChannelID, reactionAddEvent.MessageID)
		if currentMessage == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = reactionAddEvent.ChannelID
			query.MessageID = reactionAddEvent.MessageID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}

			currentMessage = message
		}

		// check if there is already a "reaction" cached
		currentReaction := currentMessage.GetReaction(reactionAddEvent.Emoji)
		if currentReaction == nil {
			currentReaction = &structs.Reaction{
				Emoji:       reactionAddEvent.Emoji,
				Count:       0,
				BurstColors: reactionAddEvent.BurstColors,
			}
		}

		// check if the user reacting is the author of the message
		if reactionAddEvent.MessageAuthorID != nil && reactionAddEvent.MessageAuthorID.Equals(reactionAddEvent.UserID) {
			currentReaction.IsMe = true
		}

		// check if the reaction is a burst  or normal reaction
		if reactionAddEvent.Type == receiveevents.MessageReactionBurst {
			currentReaction.CountDetails.Burst++
		} else if reactionAddEvent.Type == receiveevents.MessageReactionNormal {
			currentReaction.CountDetails.Normal++
		}

		currentReaction.Count++
		currentMessage.UpdateReactions(*currentReaction)

		server.UpdateMessage(*currentMessage)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageReactionRemoveEvent(s Session, p gateway.Payload) error {
	if reactionRemoveEvent, ok := p.Data.(receiveevents.MessageReactionRemoveEvent); ok {
		servers := s.GetServers()
		server, exists := servers[reactionRemoveEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMessage := server.GetMessage(reactionRemoveEvent.ChannelID, reactionRemoveEvent.MessageID)
		if currentMessage == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = reactionRemoveEvent.ChannelID
			query.MessageID = reactionRemoveEvent.MessageID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}

			currentMessage = message
		}

		currentReaction := currentMessage.GetReaction(reactionRemoveEvent.Emoji)
		if currentReaction == nil {
			return nil
		}

		if currentReaction.Count--; currentReaction.Count == 0 {
			currentMessage.DeleteReaction(reactionRemoveEvent.Emoji)
		} else {
			currentMessage.UpdateReactions(*currentReaction)
		}

		server.UpdateMessage(*currentMessage)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageReactionRemoveAllEvent(s Session, p gateway.Payload) error {
	if reactionRemoveAllEvent, ok := p.Data.(receiveevents.MessageReactionRemoveAllEvent); ok {
		servers := s.GetServers()
		server, exists := servers[reactionRemoveAllEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMessage := server.GetMessage(reactionRemoveAllEvent.ChannelID, reactionRemoveAllEvent.MessageID)
		if currentMessage == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = reactionRemoveAllEvent.ChannelID
			query.MessageID = reactionRemoveAllEvent.MessageID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}

			currentMessage = message
		}

		currentMessage.Reactions = []structs.Reaction{}
		server.UpdateMessage(*currentMessage)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageReactionRemoveEmojiEvent(s Session, p gateway.Payload) error {
	if reactionRemoveEmojiEvent, ok := p.Data.(receiveevents.MessageReactionRemoveEmojiEvent); ok {
		servers := s.GetServers()
		server, exists := servers[reactionRemoveEmojiEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMessage := server.GetMessage(reactionRemoveEmojiEvent.ChannelID, reactionRemoveEmojiEvent.MessageID)
		if currentMessage == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = reactionRemoveEmojiEvent.ChannelID
			query.MessageID = reactionRemoveEmojiEvent.MessageID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}

			currentMessage = message
		}

		currentMessage.DeleteReaction(reactionRemoveEmojiEvent.Emoji)
		server.UpdateMessage(*currentMessage)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessagePollVoteAddEvent(s Session, p gateway.Payload) error {
	if messagePollVoteAddEvent, ok := p.Data.(receiveevents.MessagePollVoteAddEvent); ok {
		servers := s.GetServers()
		server, exists := servers[messagePollVoteAddEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMessage := server.GetMessage(messagePollVoteAddEvent.ChannelID, messagePollVoteAddEvent.MessageID)
		if currentMessage == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = messagePollVoteAddEvent.ChannelID
			query.MessageID = messagePollVoteAddEvent.MessageID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}

			currentMessage = message
		}
		if currentMessage.Poll == nil {
			return errors.New("no poll active")
		}

		answer := &structs.PollAnswer{
			AnswerID:  messagePollVoteAddEvent.AnswerID,
			PollMedia: structs.PollMedia{},
		}

		currentMessage.Poll.Answers = append(currentMessage.Poll.Answers, *answer)
		server.UpdateMessage(*currentMessage)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessagePollVoteRemoveEvent(s Session, p gateway.Payload) error {
	if messagePollVoteRemoveEvent, ok := p.Data.(receiveevents.MessagePollVoteRemoveEvent); ok {
		servers := s.GetServers()
		server, exists := servers[messagePollVoteRemoveEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMessage := server.GetMessage(messagePollVoteRemoveEvent.ChannelID, messagePollVoteRemoveEvent.MessageID)
		if currentMessage == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = messagePollVoteRemoveEvent.ChannelID
			query.MessageID = messagePollVoteRemoveEvent.MessageID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}

			currentMessage = message
		}
		if currentMessage.Poll == nil {
			return errors.New("no poll active")
		}

		for i, answer := range currentMessage.Poll.Answers {
			if answer.AnswerID == messagePollVoteRemoveEvent.AnswerID {
				currentMessage.Poll.Answers = append(currentMessage.Poll.Answers[:i], currentMessage.Poll.Answers[i+1:]...)
				break
			}
		}

		server.UpdateMessage(*currentMessage)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleTypingStartEvent(s Session, p gateway.Payload) error {
	if typingStartEvent, ok := p.Data.(receiveevents.TypingStartEvent); ok {
		servers := s.GetServers()
		server, exists := servers[typingStartEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentChannel := server.GetChannel(typingStartEvent.ChannelID)
		if currentChannel == nil {
			return errors.New("channel not found")
		}

		if currentChannel.Typing == nil {
			currentChannel.Typing = structs.NewTypingChannel()
		}
		currentChannel.Typing.AddUser(typingStartEvent.UserID)
		server.UpdateChannel(typingStartEvent.ChannelID, *currentChannel)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleUserUpdateEvent(s Session, p gateway.Payload) error {
	if userUpdateEvent, ok := p.Data.(receiveevents.UserUpdateEvent); ok {
		servers := s.GetServers()
		for _, server := range servers {
			members := server.GetMembers()
			for _, member := range members {
				if member.User.ID.Equals(userUpdateEvent.ID) {
					if err := util.UpdateFields(member.User, userUpdateEvent); err != nil {
						return err
					}
					server.UpdateMember(userUpdateEvent.ID, member)
					s.AddServer(server)
					break
				}
			}
		}
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceChannelEffectSendEvent(s Session, p gateway.Payload) error {
	if _, ok := p.Data.(receiveevents.VoiceChannelEffectSendEvent); ok {
		fmt.Println("VOICE CHANNEL EFFECT SEND NOT IMPLEMENTED IDK WHAT TO USE IT FOR")
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceStateUpdateEvent(s Session, p gateway.Payload) error {
	if voiceStateUpdateEvent, ok := p.Data.(receiveevents.VoiceStateUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[voiceStateUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.UpdateVoiceState(*voiceStateUpdateEvent.VoiceState)
		s.AddServer(server)

		// yuck!!!!
		if s.GetBotData().UserDetails.ID.Equals(voiceStateUpdateEvent.UserID) {
			vs := s.GetVoiceSession(*voiceStateUpdateEvent.GuildID)
			if vs != nil {
				vs.SetSessionID(voiceStateUpdateEvent.SessionID)
				if vs.GetGuildID() != nil && vs.GetToken() != nil {
					if vs.GetSessionID() != nil && !vs.IsConnectReady() && !vs.GetConnected() {
						vs.SetConnectReady(true)
					}
				}
			}
		}
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceServerUpdateEvent(s Session, p gateway.Payload) error {
	if voiceServerUpdateEvent, ok := p.Data.(receiveevents.VoiceServerUpdateEvent); ok {
		vs := s.GetVoiceSession(voiceServerUpdateEvent.GuildID)
		if vs == nil {
			return errors.New("voice session not initialized")
		}

		if voiceServerUpdateEvent.Endpoint != nil && !strings.Contains(*voiceServerUpdateEvent.Endpoint, "wss://") {
			*voiceServerUpdateEvent.Endpoint = "wss://" + *voiceServerUpdateEvent.Endpoint
		}
		if voiceServerUpdateEvent.Endpoint != nil {
			vs.SetResumeURL(*voiceServerUpdateEvent.Endpoint)
		}
		vs.SetGuildID(voiceServerUpdateEvent.GuildID)
		vs.SetToken(voiceServerUpdateEvent.Token)
		if vs.GetSessionID() != nil && !vs.IsConnectReady() && !vs.GetConnected() {
			vs.SetConnectReady(true)
		}
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleWebhooksUpdateEvent(s Session, p gateway.Payload) error {
	if _, ok := p.Data.(receiveevents.WebhooksUpdateEvent); ok {
		fmt.Println("WEBHOOKS UPDATE EVENT NOT IMPLEMENTED YET THIS IS USED FOR SERVER WEBHOOK UPDATE EVENTS")
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildAuditLogEntryCreateEvent(s Session, p gateway.Payload) error {
	if _, ok := p.Data.(receiveevents.GuildAuditLogEntryCreateEvent); ok {
		fmt.Println("GUILD AUDIT LOG ENTRY CREATE EVENT NOT IMPLEMENTED YET THIS IS USED FOR SERVER AUDIT LOG EVENTS")
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceChannelStatusUpdateEvent(s Session, p gateway.Payload) error {
	if _, ok := p.Data.(receiveevents.VoiceChannelStatusUpdateEvent); ok {
		fmt.Println("VOICE CHANNEL STATUS UPDATE EVENT NOT IMPLEMENTED YET THIS IS USED FOR VOICE CHANNEL STATUS UPDATES ??")
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleHeartbeatACKEvent(s Session, p gateway.Payload) error {
	return nil
}

func handleReadyEvent(s Session, p gateway.Payload) error {
	if readyEvent, ok := p.Data.(receiveevents.ReadyEvent); ok {
		s.SetID(&readyEvent.SessionID)
		s.SetResumeURL(&readyEvent.ResumeGatewayURL)
		s.SetBotData(structs.NewBotData(readyEvent.User, readyEvent.Application))
		s.SetVoiceSessionBotData(structs.NewBotData(readyEvent.User, readyEvent.Application))
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleHelloEvent(s Session, p gateway.Payload) error {
	if helloEvent, ok := p.Data.(receiveevents.HelloEvent); ok {
		heartbeatInterval := int(helloEvent.HeartbeatInterval)
		s.SetHeartbeatACK(&heartbeatInterval)
	} else {
		return errors.New("unexpected payload data type")
	}

	return startHeartbeatTimer(s)
}

func handleHeartbeatEvent(s Session, p gateway.Payload) error {
	if heartbeatEvent, ok := p.Data.(receiveevents.HeartbeatEvent); ok {
		if heartbeatEvent.LastSequence != nil {
			s.SetSequence(heartbeatEvent.LastSequence)
		}
		return sendHeartbeatEvent(s)
	}
	return errors.New("unexpected payload data type")
}

func sendHeartbeatEvent(s Session) error {
	if s.GetConn() == nil {
		return errors.New("connection unavailable")
	}

	heartbeatEvent := sendevents.HeartbeatEvent{
		LastSequence: s.GetSequence(),
	}
	ackPayload := gateway.Payload{
		OpCode: gateway.Heartbeat,
		Data:   heartbeatEvent,
	}

	heartbeatData, err := json.Marshal(ackPayload)
	if err != nil {
		return err
	}

	s.Write(heartbeatData)
	return nil
}

func heartbeatLoop(ticker *time.Ticker, s Session) {
	defer ticker.Stop()

	firstHeartbeat := true

	for {
		select {
		case <-ticker.C:
			if firstHeartbeat {
				jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
				time.Sleep(jitter)
				firstHeartbeat = false
			}

			if err := sendHeartbeatEvent(s); err != nil {
				return
			}
		case <-s.GetSession().stopHeartbeat:
			return
		}
	}
}

func startHeartbeatTimer(s Session) error {
	if s.GetHeartbeatACK() == nil {
		return errors.New("no heartbeat interval set")
	}

	ticker := time.NewTicker(time.Duration(*s.GetHeartbeatACK()) * time.Millisecond)
	go heartbeatLoop(ticker, s)
	return nil
}
