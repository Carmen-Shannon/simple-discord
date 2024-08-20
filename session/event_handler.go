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
			"CHANNEL_CREATE":            nil,                //placeholder
			"CHANNEL_UPDATE":            nil,                //placeholder
			"CHANNEL_DELETE":            nil,                //placeholder
			"GUILD_CREATE":              handleGuildCreateEvent,
			"GUILD_UPDATE":              nil, //placeholder
			"GUILD_DELETE":              nil, //placeholder
			"GUILD_BAN_ADD":             nil, //placeholder
			"GUILD_BAN_REMOVE":          nil, //placeholder
			"GUILD_EMOJIS_UPDATE":       nil, //placeholder
			"GUILD_INTEGRATIONS_UPDATE": nil, //placeholder
			"GUILD_MEMBER_ADD":          nil, //placeholder
			"GUILD_MEMBER_REMOVE":       nil, //placeholder
			"GUILD_MEMBER_UPDATE":       nil, //placeholder
			"GUILD_MEMBERS_CHUNK":       nil, //placeholder
			"GUILD_ROLE_CREATE":         nil, //placeholder
			"GUILD_ROLE_UPDATE":         nil, //placeholder
			"GUILD_ROLE_DELETE":         nil, //placeholder
			"MESSAGE_CREATE":            nil, //placeholder
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

func (e *EventHandler) HandleEvent(s *Session, payload gateway.Payload) error {
	if payload.EventName == nil {
		fmt.Printf("HANDLING OPCODE EVENT: %v\n", payload.OpCode)
		if handler, ok := e.OpCodeHandlers[payload.OpCode]; ok && handler != nil {
			if payload.Seq != nil {
				s.SetSequence(payload.Seq)
			}
			return handler(s, payload)
		}
		return errors.New("no handler for opcode")
	}

	if handler, ok := e.NamedHandlers[*payload.EventName]; ok && handler != nil {
		fmt.Printf("HANDLING NAMED EVENT: %v\n", *payload.EventName)
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
	fmt.Println("HANDLING RESUME EVENT")
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
	fmt.Println("HANDLING IDENTIFY EVENT")
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

func handleInvalidSessionEvent(s *Session, p gateway.Payload) error {
	fmt.Println("HANDLING INVALID SESSION EVENT")
	if invalidSessionEvent, ok := p.Data.(receiveevents.InvalidSessionEvent); ok {
		if invalidSessionEvent {
			handleSendResumeEvent(s, p)
		} else {
			s.Exit(1000)
		}
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleReconnectEvent(s *Session, p gateway.Payload) error {
	fmt.Println("HANDLING RECONNECT EVENT")
	if _, ok := p.Data.(receiveevents.ReconnectEvent); ok {
		handleSendResumeEvent(s, p)
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
	fmt.Println("HANDLING GUILD CREATE EVENT")
	if guildCreateEvent, ok := p.Data.(receiveevents.GuildCreateEvent); ok {
		if guildCreateEvent.Unavailable != nil && *guildCreateEvent.Unavailable {
			server := structs.Server(*guildCreateEvent.Server)
			s.AddServer(server)
		}
	} else if guildCreateUnavailableEvent, ok := p.Data.(receiveevents.GuildCreateUnavailableEvent); ok {
		server := structs.Server{}
		server.ID = guildCreateUnavailableEvent.ID
		s.AddServer(server)
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleHeartbeatACKEvent(s *Session, p gateway.Payload) error {
	fmt.Println("HANDLING HEARTBEAT ACK EVENT")
	return nil
}

func handleReadyEvent(s *Session, p gateway.Payload) error {
	fmt.Println("HANDLING READY EVENT")
	if readyEvent, ok := p.Data.(receiveevents.ReadyEvent); ok {
		s.SetID(&readyEvent.SessionID)
		s.SetResumeURL(&readyEvent.ResumeGatewayURL)
		fmt.Printf("successfully connected to gateway Bot ID: %v\n", *s.GetID())
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleHelloEvent(s *Session, p gateway.Payload) error {
	fmt.Println("HANDLING HELLO EVENT")
	if helloEvent, ok := p.Data.(receiveevents.HelloEvent); ok {
		heartbeatInterval := int(helloEvent.HeartbeatInterval)
		s.SetHeartbeatACK(&heartbeatInterval)
	} else {
		return errors.New("unexpected payload data type")
	}

	return startHeartbeatTimer(s)
}

func handleHeartbeatEvent(s *Session, p gateway.Payload) error {
	fmt.Println("HANDLING HEARTBEAT EVENT")
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

	fmt.Println("SENDING HEARTBEAT")
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
