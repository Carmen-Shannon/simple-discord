package session

import (
	gateway "github.com/Carmen-Shannon/simple-discord/gateway"
	receiveevents "github.com/Carmen-Shannon/simple-discord/gateway/receive_events"
	sendevents "github.com/Carmen-Shannon/simple-discord/gateway/send_events"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"golang.org/x/net/websocket"
)

type Session struct {
	Mu           sync.Mutex
	Conn         *websocket.Conn
	HeartbeatACK *int
	Sequence     *int
	EventHandler *EventHandler
}

func (s *Session) Exit() error {
	return s.Conn.Close()
}

func (s *Session) Listen() {
	for {
		msg, err := s.Read()
		if err != nil {
			fmt.Printf("error reading message: %v\n", err)
			continue
		}

		var payload gateway.Payload
		if err := json.Unmarshal(msg, &payload); err != nil {
			fmt.Printf("error parsing message: %v\n", err)
			continue
		} else if err := gateway.NewReceiveEvent(payload); err != nil {
			fmt.Printf("error parsing data: %v\n", err)
			continue
		}

		if err := s.Eventhandler.HandleEvent(payload); err != nil {
			fmt.Printf("error handling event: %v\n", err)
			continue
		}
	}
}

func (s *Session) StartHeartbeat() error {
	if s.HeartbeatACK == nil {
		return errors.New("no heartbeat interval set")
	}

	ticker := time.NewTicker(time.Duration(*s.HeartbeatACK) * time.Millisecond)
	go s.heartbeatLoop(ticker)
	return nil
}

func (s *Session) heartbeatLoop(ticker *time.Ticker) {
	if ticker == nil {
		return
	}

	firstHeartbeat := true

	for range ticker.C {
		if firstHeartbeat {
			jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
			time.Sleep(jitter)
			firstHeartbeat = false
		}

		if err := s.Ack(); err != nil {
			ticker.Stop()
			return
		}
	}
}

func (s *Session) Ack() error {
	if s.Conn == nil {
		return errors.New("connection unavailable")
	}

	heartbeatEvent := sendevents.HeartbeatEvent{
		LastSequence: s.HeartbeatACK,
	}
	ackPayload := gateway.Payload{
		OpCode: gateway.Heartbeat,
		Data:   heartbeatEvent,
	}

	heartbeatData, err := json.Marshal(ackPayload)
	if err != nil {
		return err
	}

	return s.Write(heartbeatData)
}

func (s *Session) UpdateSequence(seq int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Sequence = &seq
}

func (s *Session) Read() ([]byte, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var msg []byte
	buffer := make([]byte, 512)
	for {
		n, err := s.Conn.Read(buffer)
		if err != nil {
			return nil, err
		}

		msg = append(msg, buffer[:n]...)
		if n < len(buffer) {
			break
		}
	}

	return msg, nil
}

func (s *Session) Write(data []byte) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if _, err := s.Conn.Write(data); err != nil {
		return err
	}

	return nil
}

func (s *Session) Identify(token string, intents []structs.Intent) error {
	identifyEvent := sendevents.IdentifyEvent{
		Token: token,
		Properties: sendevents.IdentifyProperties{
			Os:      runtime.GOOS,
			Browser: "discord",
			Device:  "discord",
		},
		Intents: structs.GetIntents(intents),
	}
	identifyPayload := gateway.Payload{
		OpCode: gateway.Identify,
		Data:   identifyEvent,
	}

	identifyData, err := json.Marshal(identifyPayload)
	if err != nil {
		return err
	}

	return s.Write(identifyData)
}

func NewSession(token string, intents []structs.Intent) (*Session, error) {
	ws, err := websocket.Dial(gateway.GatewayURL, "", "http://localhost/")
	if err != nil {
		return nil, err
	}

	sess := &Session{
		Conn: ws,
	}

	if err := sess.Identify(token, intents); err != nil {
		return nil, err
	}

	rawResponse, err := sess.Read()
	if err != nil {
		return nil, err
	}

	var helloPayload gateway.Payload
	if err := json.Unmarshal(rawResponse, &helloPayload); err != nil {
		return nil, err
	} else if err := gateway.NewReceiveEvent(helloPayload); err != nil {
		return nil, err
	}

	switch helloPayload.Data.(type) {
	case receiveevents.HelloEvent:
		heartbeatInterval := int(helloPayload.Data.(receiveevents.HelloEvent).HeartbeatInterval)
		sess.HeartbeatACK = &heartbeatInterval
		sess.UpdateSequence(*helloPayload.Seq)
	default:
		return nil, errors.New("unexpected payload data type")
	}

	if err := sess.StartHeartbeat(); err != nil {
		return nil, err
	}

	return sess, nil
}
