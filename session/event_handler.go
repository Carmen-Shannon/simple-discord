package session

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"
	gateway "github.com/Carmen-Shannon/simple-discord/gateway"
	sendevents "github.com/Carmen-Shannon/simple-discord/gateway/send_events"
	receiveevents "github.com/Carmen-Shannon/simple-discord/gateway/receive_events"
)

type EventHandler struct {
	ReceiveHandlers map[string]func(*Session, interface{}) error
	SendHandlers map[gateway.GatewayOpCode]func(*Session, interface{}) error
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		ReceiveHandlers: map[string]func(*Session, interface{}) error{
			"HELLO":             handleHelloEvent,
			"READY":             nil, //placeholder
			"RESUMED":           nil, //placeholder
			"RECONNECT":         nil, //placeholder
			"INVALID_SESSION":   nil, //placeholder
			"CHANNEL_CREATE":    nil, //placeholder
			"CHANNEL_UPDATE":    nil, //placeholder
			"CHANNEL_DELETE":    nil, //placeholder
			"GUILD_CREATE":      nil, //placeholder
			"GUILD_UPDATE":      nil, //placeholder
			"GUILD_DELETE":      nil, //placeholder
			"GUILD_BAN_ADD":     nil, //placeholder
			"GUILD_BAN_REMOVE":  nil, //placeholder
			"GUILD_EMOJIS_UPDATE": nil, //placeholder
			"GUILD_INTEGRATIONS_UPDATE": nil, //placeholder
			"GUILD_MEMBER_ADD":  nil, //placeholder
			"GUILD_MEMBER_REMOVE": nil, //placeholder
			"GUILD_MEMBER_UPDATE": nil, //placeholder
			"GUILD_MEMBERS_CHUNK": nil, //placeholder
			"GUILD_ROLE_CREATE":  nil, //placeholder
			"GUILD_ROLE_UPDATE":  nil, //placeholder
			"GUILD_ROLE_DELETE":  nil, //placeholder
			"MESSAGE_CREATE":    nil, //placeholder
			"MESSAGE_UPDATE":    nil, //placeholder
			"MESSAGE_DELETE":    nil, //placeholder
			"MESSAGE_BULK_DELETE": nil, //placeholder
			"REACTION_ADD":      nil, //placeholder
			"REACTION_REMOVE":   nil, //placeholder
			"REACTION_REMOVE_ALL": nil, //placeholder
			"TYPING_START":      nil, //placeholder
			"USER_UPDATE":       nil, //placeholder
			"VOICE_STATE_UPDATE": nil, //placeholder
			"VOICE_SERVER_UPDATE": nil, //placeholder
			"WEBHOOKS_UPDATE":   nil, //placeholder
		},
		SendHandlers: map[gateway.GatewayOpCode]func(*Session, interface{}) error{
			gateway.Heartbeat: handleSendHeartbeatEvent,
			gateway.Identify: nil, //placeholder
			gateway.PresenceUpdate: nil, //placeholder
			gateway.VoiceStateUpdate: nil, //placeholder
			gateway.Resume: nil, //placeholder
			gateway.RequestGuildMembers: nil, //placeholder
		},
	}
}

func (e *EventHandler) HandleEvent(s *Session, payload gateway.Payload) error {
	if payload.EventName == nil {
		return nil
	}
	if handler, ok := e.ReceiveHandlers[*payload.EventName]; ok && handler != nil {
		if payload.Seq != nil {
			s.UpdateSequence(*payload.Seq)
		}
        return handler(s, payload.Data)
    }
    return errors.New("no handler for event")
}

func handleHelloEvent(s *Session, d interface{}) error {
	switch d := d.(type) {
	case receiveevents.HelloEvent:
		heartbeatInterval := int(d.HeartbeatInterval)
		*s.HeartbeatACK = heartbeatInterval
	default:
		return errors.New("unexpected payload data type")
	}

	return startHeartbeatTimer(s)
}

func handleSendHeartbeatEvent(s *Session, d interface{}) error {
	return sendHeartbeatEvent(s)
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
		Data: heartbeatEvent,
	}

	heartbeatData, err := json.Marshal(ackPayload)
	if err != nil {
		return err
	}

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