package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	gateway "github.com/Carmen-Shannon/simple-discord/gateway"
	receiveevents "github.com/Carmen-Shannon/simple-discord/gateway/receive_events"
	sendevents "github.com/Carmen-Shannon/simple-discord/gateway/send_events"
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
			"RESUMED":                   nil, //placeholder
			"RECONNECT":                 nil, //placeholder
			"INVALID_SESSION":           nil, //placeholder
			"CHANNEL_CREATE":            nil, //placeholder
			"CHANNEL_UPDATE":            nil, //placeholder
			"CHANNEL_DELETE":            nil, //placeholder
			"GUILD_CREATE":              nil, //placeholder
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
			gateway.Heartbeat:           nil, //placeholder
			gateway.Hello:               handleHelloEvent,
			gateway.Identify:            nil, //placeholder
			gateway.PresenceUpdate:      nil, //placeholder
			gateway.VoiceStateUpdate:    nil, //placeholder
			gateway.Resume:              nil, //placeholder
			gateway.RequestGuildMembers: nil, //placeholder
			gateway.HeartbeatACK:        nil, //placeholder
		},
	}
}

func (e *EventHandler) HandleEvent(s *Session, payload gateway.Payload) error {
	if payload.EventName == nil {
		if handler, ok := e.OpCodeHandlers[payload.OpCode]; ok && handler != nil {
			if payload.Seq != nil {
				s.SetSequence(*payload.Seq)
			}
			return handler(s, payload)
		}
		return errors.New("no handler for opcode")
	}

	if handler, ok := e.NamedHandlers[*payload.EventName]; ok && handler != nil {
		if payload.Seq != nil {
			s.SetSequence(*payload.Seq)
		}
		return handler(s, payload)
	}
	return errors.New("no handler for event name")
}

func handleReadyEvent(s *Session, p gateway.Payload) error {
	if readyEvent, ok := p.Data.(receiveevents.ReadyEvent); ok {
		s.ID = &readyEvent.SessionID
		s.ResumeURL = &readyEvent.ResumeGatewayURL
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleHelloEvent(s *Session, p gateway.Payload) error {
	if helloEvent, ok := p.Data.(receiveevents.HelloEvent); ok {
		heartbeatInterval := int(helloEvent.HeartbeatInterval)
		s.HeartbeatACK = &heartbeatInterval
	} else {
		return errors.New("unexpected payload data type")
	}

	return startHeartbeatTimer(s)
}

func sendHeartbeatEvent(s *Session) error {
	if s.Conn == nil {
		return errors.New("connection unavailable")
	}

	heartbeatEvent := sendevents.HeartbeatEvent{
		LastSequence: s.Sequence,
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

	return s.Write(heartbeatData)
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
