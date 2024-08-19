package session

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	gateway "github.com/Carmen-Shannon/simple-discord/gateway"
	requestutil "github.com/Carmen-Shannon/simple-discord/gateway/request_util"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"golang.org/x/net/websocket"
)

type Session struct {
	Mu            sync.Mutex
	Conn          *websocket.Conn
	HeartbeatACK  *int
	Sequence      *int
	EventHandler  *EventHandler
	ID            *string
	ResumeURL     *string
	Token         *string
	Intents       []structs.Intent
	helloReceived chan struct{}
	readChan      chan []byte
	writeChan     chan []byte
	errorChan     chan error
}

func (s *Session) Exit() error {
	return s.Conn.Close()
}

func (s *Session) Listen() {
	log.Println("Starting to listen for messages")
	for msg := range s.readChan {
		var payload gateway.Payload
		fmt.Printf("New message received: %s\n", string(msg))
		if err := json.Unmarshal(msg, &payload); err != nil {
			s.errorChan <- fmt.Errorf("error parsing message: %v", err)
			fmt.Println(payload.ToString())
			continue
		}

		data, err := gateway.NewReceiveEvent(payload)
		if err != nil {
			s.errorChan <- fmt.Errorf("error parsing event: %v", err)
			fmt.Println(payload.ToString())
			continue
		}

		payload.Data = data
		fmt.Printf("Received payload type: %T\n", payload.Data)

		// signal when HELLO is received
		if payload.OpCode == gateway.Hello {
			close(s.helloReceived)
		}

		if err := s.EventHandler.HandleEvent(s, payload); err != nil {
			s.errorChan <- fmt.Errorf("error handling event: %v", err)
			fmt.Println(payload.ToString())
			continue
		}
	}
}

func (s *Session) handleRead() {
	var buffers [][]byte

	for {
		tempBuffer := make([]byte, 1024)
		n, err := s.Conn.Read(tempBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			s.errorChan <- err
			continue
		}

		buffers = append(buffers, tempBuffer[:n])

		for {
			combinedBuffer := bytes.Join(buffers, nil)
			decoder := json.NewDecoder(bytes.NewReader(combinedBuffer))
			var msg json.RawMessage
			startOffset := len(combinedBuffer)

			if err := decoder.Decode(&msg); err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					// incomplete message
					break
				}
				s.errorChan <- fmt.Errorf("error decoding raw message: %v", err)
				buffers = nil
				break
			}

			s.readChan <- msg

			offset := startOffset - int(decoder.InputOffset())
			if offset > 0 && offset <= len(combinedBuffer) {
				remainingData := combinedBuffer[decoder.InputOffset():]
				buffers = [][]byte{remainingData}
			} else {
				buffers = nil
			}
		}
	}
}

func (s *Session) Write(data []byte) {
	if len(s.writeChan) < cap(s.writeChan) {
		s.writeChan <- data
	} else {
		s.errorChan <- fmt.Errorf("failed to write data to write channel")
	}
}

func (s *Session) handleWrite() {
	for data := range s.writeChan {
		fmt.Printf("WRITING DATA: %s\n", string(data))
		if _, err := s.Conn.Write(data); err != nil {
			s.errorChan <- err
		}
	}
}

func (s *Session) handleError() {
	for err := range s.errorChan {
		log.Printf("error: %v\n", err)
	}
}

func NewSession(token string, intents []structs.Intent) (*Session, error) {
	ws, err := dialer(token)
	if err != nil {
		return nil, err
	}

	var sess Session
	sess.SetConn(ws)
	sess.SetEventHandler(NewEventHandler())
	sess.SetToken(&token)
	sess.SetIntents(intents)
	sess.helloReceived = make(chan struct{})
	sess.readChan = make(chan []byte)
	sess.writeChan = make(chan []byte, 4096)
	sess.errorChan = make(chan error)

	go sess.Listen()
	go sess.handleRead()
	go sess.handleWrite()
	go sess.handleError()

	<-sess.helloReceived

	var identifyData gateway.Payload
	identifyData.OpCode = gateway.Identify

	if err := sess.EventHandler.HandleEvent(&sess, identifyData); err != nil {
		return nil, err
	}

	return &sess, nil
}

func getGatewayUrl(token string) (string, error) {
	botUrl, err := getBotUrl()
	if err != nil {
		return "", err
	}

	botVersion, err := getBotVersion()
	if err != nil {
		return "", err
	}
	headers := map[string]string{
		"Authorization": "Bot " + token,
		"User-Agent":    fmt.Sprintf("DiscordBot (%s, %s)", botUrl, botVersion),
	}
	resp, err := requestutil.HttpRequest("GET", "/gateway", headers, nil)
	if err != nil {
		return "", err
	}

	var gatewayResponse structs.GetGatewayResponse
	if err := json.Unmarshal(resp, &gatewayResponse); err != nil {
		return "", err
	}

	fmt.Println(gatewayResponse.URL)

	return gatewayResponse.URL, nil
}

func getBotUrl() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

func getBotVersion() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func dialer(token string) (*websocket.Conn, error) {
	url, err := getGatewayUrl(token)
	if err != nil {
		return nil, err
	}

	ws, err := websocket.Dial(url+"/?v=10&encoding=json", "", "http://localhost/")
	if err != nil {
		return nil, err
	}
	return ws, nil
}

// getters and setters because mutex
func (s *Session) GetConn() *websocket.Conn {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Conn
}

func (s *Session) SetConn(conn *websocket.Conn) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Conn = conn
}

func (s *Session) GetHeartbeatACK() *int {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.HeartbeatACK
}

func (s *Session) SetHeartbeatACK(ack *int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.HeartbeatACK = ack
}

func (s *Session) GetSequence() *int {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Sequence
}

func (s *Session) SetSequence(seq *int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Sequence = seq
}

func (s *Session) GetEventHandler() *EventHandler {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.EventHandler
}

func (s *Session) SetEventHandler(handler *EventHandler) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.EventHandler = handler
}

func (s *Session) GetID() *string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.ID
}

func (s *Session) SetID(id *string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.ID = id
}

func (s *Session) GetResumeURL() *string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.ResumeURL
}

func (s *Session) SetResumeURL(url *string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.ResumeURL = url
}

func (s *Session) GetToken() *string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Token
}

func (s *Session) SetToken(token *string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Token = token
}

func (s *Session) GetIntents() []structs.Intent {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Intents
}

func (s *Session) SetIntents(intents []structs.Intent) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Intents = intents
}
