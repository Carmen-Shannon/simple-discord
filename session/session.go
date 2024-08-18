package session

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"

	gateway "github.com/Carmen-Shannon/simple-discord/gateway"
	sendevents "github.com/Carmen-Shannon/simple-discord/gateway/send_events"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"golang.org/x/net/websocket"
)

type Session struct {
	Mu           sync.Mutex
	Conn         *websocket.Conn
	HeartbeatACK *int
	Sequence     *int
	EventHandler *EventHandler
	ID           *string
	ResumeURL    *string
	Token        *string
}

func (s *Session) Exit() error {
	return s.Conn.Close()
}

func (s *Session) Listen() error {
	log.Println("Starting to listen for messages")
	for {
		msg, err := s.Read()
		if err != nil {
			fmt.Printf("error reading message: %v\n", err)
			return err
		}

		var payload gateway.Payload
		if err := json.Unmarshal(msg, &payload); err != nil {
			fmt.Printf("error parsing message: %v\n", err)
			fmt.Println(payload.ToString())
			continue
		}

		data, err := gateway.NewReceiveEvent(payload)
		if err != nil {
			fmt.Printf("error parsing event: %v\n", err)
			fmt.Println(payload.ToString())
			continue
		}

		payload.Data = data
		fmt.Printf("Received payload type: %T\n", payload.Data)

		if err := s.EventHandler.HandleEvent(s, payload); err != nil {
			fmt.Printf("error handling event: %v\n", err)
			fmt.Println(payload.ToString())
			continue
		}
	}
}

func (s *Session) SetSequence(seq int) {
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
	ws, err := dialer()
	if err != nil {
		return nil, err
	}

	sess := &Session{
		Conn:         ws,
		EventHandler: NewEventHandler(),
		Token:        &token,
	}

	if err := sess.Identify(token, intents); err != nil {
		return nil, err
	}

	go sess.Listen()

	return sess, nil
}

func dialer() (*websocket.Conn, error) {
	ws, err := websocket.Dial(gateway.GatewayURL, "", "http://localhost/")
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to gateway")
	return ws, nil
}
