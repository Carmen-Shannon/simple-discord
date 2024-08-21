package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	gateway "github.com/Carmen-Shannon/simple-discord/gateway"
	receiveevents "github.com/Carmen-Shannon/simple-discord/gateway/receive_events"
	sendevents "github.com/Carmen-Shannon/simple-discord/gateway/send_events"
	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/util"
)

type EventHandler struct {
	NamedHandlers  map[string]func(*Session, gateway.Payload) error
	OpCodeHandlers map[gateway.GatewayOpCode]func(*Session, gateway.Payload) error
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		NamedHandlers: map[string]func(*Session, gateway.Payload) error{
			"HELLO":                     handleHelloEvent,
			"READY":                     handleReadyEvent,
			"RESUMED":                   handleResumedEvent,
			"RECONNECT":                 handleReconnectEvent,
			"INVALID_SESSION":           handleInvalidSessionEvent,
			"CHANNEL_CREATE":            handleChannelCreateEvent,
			"CHANNEL_UPDATE":            handleChannelUpdateEvent,
			"CHANNEL_DELETE":            handleChannelDeleteEvent,
			"GUILD_CREATE":              handleGuildCreateEvent,
			"GUILD_UPDATE":              handleGuildUpdateEvent,
			"GUILD_DELETE":              handleGuildDeleteEvent,
			"GUILD_BAN_ADD":             handleGuildBanAddEvent,
			"GUILD_BAN_REMOVE":          handleGuildBanRemoveEvent,
			"GUILD_EMOJIS_UPDATE":       handleGuildEmojisUpdateEvent,
			"GUILD_INTEGRATIONS_UPDATE": handleGuildIntegrationsUpdateEvent,
			"GUILD_MEMBER_ADD":          handleGuildMemberAddEvent,
			"GUILD_MEMBER_REMOVE":       handleGuildMemberRemoveEvent,
			"GUILD_MEMBER_UPDATE":       handleGuildMemberUpdateEvent,
			"GUILD_MEMBERS_CHUNK":       handleGuildMembersChunkEvent,
			"GUILD_ROLE_CREATE":         handleGuildRoleCreateEvent,
			"GUILD_ROLE_UPDATE":         handleGuildRoleUpdateEvent,
			"GUILD_ROLE_DELETE":         handleGuildRoleDeleteEvent,
			"MESSAGE_CREATE":            handleMessageCreateEvent,
			"MESSAGE_UPDATE":            nil, //placeholder
			"MESSAGE_DELETE":            nil, //placeholder
			"MESSAGE_BULK_DELETE":       nil, //placeholder
			"REACTION_ADD":              nil, //placeholder
			"REACTION_REMOVE":           nil, //placeholder
			"REACTION_REMOVE_ALL":       nil, //placeholder
			"TYPING_START":              nil, //placeholder
			"USER_UPDATE":               nil, //placeholder
			"VOICE_STATE_UPDATE":        nil, //placeholder
			"VOICE_SERVER_UPDATE":       nil, //placeholder
			"WEBHOOKS_UPDATE":           nil, //placeholder
			"PRESENCE_UPDATE":           handlePresenceUpdateEvent,
		},
		OpCodeHandlers: map[gateway.GatewayOpCode]func(*Session, gateway.Payload) error{
			gateway.Heartbeat:           handleHeartbeatEvent,
			gateway.Identify:            handleSendIdentifyEvent,
			gateway.PresenceUpdate:      handleSendPresenceUpdateEvent,
			gateway.VoiceStateUpdate:    handleSendVoiceStateUpdateEvent,
			gateway.Resume:              handleSendResumeEvent,
			gateway.RequestGuildMembers: handleSendRequestGuildMembersEvent,
			gateway.Hello:               handleHelloEvent,
			gateway.HeartbeatACK:        handleHeartbeatACKEvent,
		},
	}
}

// this is really just for helping me log more better, will remove
var opCodeNames = map[gateway.GatewayOpCode]string{
	gateway.Heartbeat:           "Heartbeat",
	gateway.Identify:            "Identify",
	gateway.PresenceUpdate:      "PresenceUpdate",
	gateway.VoiceStateUpdate:    "VoiceStateUpdate",
	gateway.Resume:              "Resume",
	gateway.RequestGuildMembers: "RequestGuildMembers",
	gateway.Hello:               "Hello",
	gateway.HeartbeatACK:        "HeartbeatACK",
}

func (e *EventHandler) HandleEvent(s *Session, payload gateway.Payload) error {
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
			return handler(s, payload)
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
		return handler(s, payload)
	}
	return errors.New("no handler for event name")
}

func handleSendRequestGuildMembersEvent(s *Session, p gateway.Payload) error {
	fmt.Println("HANDLING REQUEST GUILD MEMBERS EVENT")
	fmt.Println("REQUEST GUILD MEMBERS NOT IMPLEMENTED")
	return nil
}

func handleSendVoiceStateUpdateEvent(s *Session, p gateway.Payload) error {
	fmt.Println("HANDLING VOICE STATE UPDATE EVENT")
	fmt.Println("VOICE STATE UPDATE NOT IMPLEMENTED")
	return nil
}

func handleSendPresenceUpdateEvent(s *Session, p gateway.Payload) error {
	fmt.Println("HANDLING PRESENCE UPDATE EVENT")
	fmt.Println("PRESENCE UPDATE NOT IMPLEMENTED")
	// presenceUpdateEvent := sendevents.PresenceUpdateEvent{
	// 	Activities: []structs.Activity{
	// 		{
	// 			Name: "discord",
	// 			Type: 0,
	// 		},
	// 	},
	// 	Status: "online",
	// }
	// presencePayload := gateway.Payload{
	// 	OpCode: gateway.PresenceUpdate,
	// 	Data:   presenceUpdateEvent,
	// }

	// presenceData, err := json.Marshal(presencePayload)
	// if err != nil {
	// 	return err
	// }

	// s.Write(presenceData)
	return nil
}

func handleSendResumeEvent(s *Session, p gateway.Payload) error {
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

func handleSendIdentifyEvent(s *Session, p gateway.Payload) error {
	identifyEvent := sendevents.IdentifyEvent{
		Token: *s.GetToken(),
		Properties: sendevents.IdentifyProperties{
			Os:      runtime.GOOS,
			Browser: "discord",
			Device:  "discord",
		},
		Intents: structs.GetIntents(s.GetIntents()),
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

func handleChannelDeleteEvent(s *Session, p gateway.Payload) error {
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

func handleChannelUpdateEvent(s *Session, p gateway.Payload) error {
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

func handleChannelCreateEvent(s *Session, p gateway.Payload) error {
	if channelCreateEvent, ok := p.Data.(receiveevents.ChannelCreateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[channelCreateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}
		channel := structs.Channel(*channelCreateEvent.Channel)
		server.AddChannel(channel)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handlePresenceUpdateEvent(s *Session, p gateway.Payload) error {
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

func handleInvalidSessionEvent(s *Session, p gateway.Payload) error {
	if invalidSessionEvent, ok := p.Data.(receiveevents.InvalidSessionEvent); ok {
		if invalidSessionEvent {
			if err := s.ResumeSession(); err != nil {
				return err
			}
			fmt.Println("RESUMED SESSION")
		} else {
			s.Exit(1001)
			var err error
			var newSess *Session
			newSess, err = NewSession(*s.GetToken(), s.GetIntents())
			if err != nil {
				return err
			}
			s.RegenerateSession(newSess)
			fmt.Println("REGENERATED SESSION")
		}
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleReconnectEvent(s *Session, p gateway.Payload) error {
	if _, ok := p.Data.(receiveevents.ReconnectEvent); ok {
		if err := s.ResumeSession(); err != nil {
			return err
		}
		fmt.Println("RESUMED SESSION")
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleResumedEvent(s *Session, p gateway.Payload) error {
	fmt.Println("HANDLING RESUMED EVENT")
	fmt.Println(p.ToString())
	return nil
}

func handleGuildCreateEvent(s *Session, p gateway.Payload) error {
	if guildCreateEvent, ok := p.Data.(receiveevents.GuildCreateEvent); ok {
		if guildCreateEvent.Unavailable != nil && !*guildCreateEvent.Unavailable {
			server := structs.Server(*guildCreateEvent.Server)
			s.AddServer(server)
		}
	} else if guildCreateUnavailableEvent, ok := p.Data.(receiveevents.GuildCreateUnavailableEvent); ok {
		server := structs.Server{}
		server.ID = guildCreateUnavailableEvent.ID
		server.Unavailable = &guildCreateUnavailableEvent.Unavailable
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleGuildUpdateEvent(s *Session, p gateway.Payload) error {
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

func handleGuildDeleteEvent(s *Session, p gateway.Payload) error {
	if guildDeleteEvent, ok := p.Data.(receiveevents.GuildDeleteEvent); ok {
		servers := s.GetServers()
		_, exists := servers[guildDeleteEvent.ID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		unavailableServer := structs.Server{}
		unavailableServer.ID = guildDeleteEvent.ID
		unavailableServer.Unavailable = &guildDeleteEvent.Unavailable
		s.AddServer(unavailableServer)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildBanAddEvent(s *Session, p gateway.Payload) error {
	fmt.Println(p.ToString())
	return nil
}

func handleGuildBanRemoveEvent(s *Session, p gateway.Payload) error {
	fmt.Println(p.ToString())
	return nil
}

func handleGuildEmojisUpdateEvent(s *Session, p gateway.Payload) error {
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

func handleGuildIntegrationsUpdateEvent(s *Session, p gateway.Payload) error {
	if _, ok := p.Data.(receiveevents.GuildIntegrationsUpdateEvent); ok {
		return nil
	}
	return errors.New("unexpected payload data type")
}

func handleGuildMemberAddEvent(s *Session, p gateway.Payload) error {
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

func handleGuildMemberRemoveEvent(s *Session, p gateway.Payload) error {
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

func handleGuildMemberUpdateEvent(s *Session, p gateway.Payload) error {
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

		if err := util.UpdateFields(currentMember, guildMemberUpdateEvent); err != nil {
			return err
		}

		server.UpdateMember(guildMemberUpdateEvent.User.ID, *currentMember)
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildMembersChunkEvent(s *Session, p gateway.Payload) error {
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

func handleGuildRoleCreateEvent(s *Session, p gateway.Payload) error {
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

func handleGuildRoleUpdateEvent(s *Session, p gateway.Payload) error {
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

func handleGuildRoleDeleteEvent(s *Session, p gateway.Payload) error {
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

func handleMessageCreateEvent(s *Session, p gateway.Payload) error {
	if _, ok := p.Data.(receiveevents.MessageCreateEvent); ok {

		fmt.Println("MESSAGE CREATE EVENT NOT IMPLEMENTED")
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleHeartbeatACKEvent(s *Session, p gateway.Payload) error {
	return nil
}

func handleReadyEvent(s *Session, p gateway.Payload) error {
	if readyEvent, ok := p.Data.(receiveevents.ReadyEvent); ok {
		s.SetID(&readyEvent.SessionID)
		s.SetResumeURL(&readyEvent.ResumeGatewayURL)
		fmt.Printf("successfully connected to gateway\n---------- %s ----------\n", readyEvent.User.Username)
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleHelloEvent(s *Session, p gateway.Payload) error {
	if helloEvent, ok := p.Data.(receiveevents.HelloEvent); ok {
		heartbeatInterval := int(helloEvent.HeartbeatInterval)
		s.SetHeartbeatACK(&heartbeatInterval)
	} else {
		return errors.New("unexpected payload data type")
	}

	return startHeartbeatTimer(s)
}

func handleHeartbeatEvent(s *Session, p gateway.Payload) error {
	if heartbeatEvent, ok := p.Data.(receiveevents.HeartbeatEvent); ok {
		if heartbeatEvent.LastSequence != nil {
			s.SetSequence(heartbeatEvent.LastSequence)
		}
		return sendHeartbeatEvent(s)
	}
	return errors.New("unexpected payload data type")
}

func sendHeartbeatEvent(s *Session) error {
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

func heartbeatLoop(ticker *time.Ticker, s *Session) {
	if ticker == nil {
		return
	} else if s.HeartbeatACK == nil {
		ticker.Stop()
		return
	}

	firstHeartbeat := true

	for range ticker.C {
		if firstHeartbeat {
			jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
			time.Sleep(jitter)
			firstHeartbeat = false
		}

		if err := sendHeartbeatEvent(s); err != nil {
			ticker.Stop()
			return
		}
	}
}

func startHeartbeatTimer(s *Session) error {
	if s.HeartbeatACK == nil {
		return errors.New("no heartbeat interval set")
	}

	ticker := time.NewTicker(time.Duration(*s.HeartbeatACK) * time.Millisecond)
	go heartbeatLoop(ticker, s)
	return nil
}

func (e *EventHandler) AddEvent() error {
	return nil
}
