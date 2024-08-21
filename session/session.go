package session

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
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
	Servers       map[string]structs.Server
	helloReceived chan struct{}
	readChan      chan []byte
	writeChan     chan []byte
	errorChan     chan error
}

func (s *Session) Exit(closeCode int) error {
	// Construct the close frame
	closeFrame := make([]byte, 2)
	binary.BigEndian.PutUint16(closeFrame[:2], uint16(closeCode))

	// Send the close frame
	if _, err := s.Conn.Write(closeFrame); err != nil {
		return fmt.Errorf("error sending close frame: %v", err)
	}

	// Close the connection
	if err := s.Conn.Close(); err != nil {
		return fmt.Errorf("error closing connection: %v", err)
	}

	// Close the channels
	close(s.readChan)
	close(s.writeChan)
	close(s.errorChan)

	return nil
}

func (s *Session) Listen() {
	for msg := range s.readChan {
		var payload gateway.Payload
		if err := json.Unmarshal(msg, &payload); err != nil {
			s.errorChan <- fmt.Errorf("error parsing message: %v", err)
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

func (s *Session) RegenerateSession(newSession *Session) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Conn = newSession.Conn
	s.HeartbeatACK = newSession.HeartbeatACK
	s.Sequence = newSession.Sequence
	s.EventHandler = newSession.EventHandler
	s.ID = newSession.ID
	s.ResumeURL = newSession.ResumeURL
	s.Token = newSession.Token
	s.Intents = newSession.Intents
	s.Servers = newSession.Servers
	s.helloReceived = newSession.helloReceived
	s.readChan = newSession.readChan
	s.writeChan = newSession.writeChan
	s.errorChan = newSession.errorChan

	return nil
}

func (s *Session) ResumeSession() error {
	// ensure the connection is really closed
	if err := s.Exit(1000); err != nil {
		return err
	}

	// open a new connection using the cached url
	ws, err := s.dialer()
	if err != nil {
		return err
	}
	s.SetConn(ws)

	// clean out the cached Guilds from the previous session
	s.SetServers(make(map[string]structs.Server))

	// Exit() already closed the channels, so we are re-opening them here
	s.helloReceived = make(chan struct{})
	s.readChan = make(chan []byte)
	s.writeChan = make(chan []byte, 4096)
	s.errorChan = make(chan error)

	// set up the goroutines to listen, read, write, and handle errors
	go s.Listen()
	go s.handleRead()
	go s.handleWrite()
	go s.handleError()

	var resumeData gateway.Payload
	resumeData.OpCode = gateway.Resume

	// let her rip tater chip
	if err := s.EventHandler.HandleEvent(s, resumeData); err != nil {
		return err
	}

	return nil
}

func NewSession(token string, intents []structs.Intent) (*Session, error) {
	var sess Session
	sess.SetToken(&token)
	sess.SetEventHandler(NewEventHandler())
	sess.SetIntents(intents)
	sess.SetServers(make(map[string]structs.Server))
	sess.helloReceived = make(chan struct{})
	sess.readChan = make(chan []byte)
	sess.writeChan = make(chan []byte, 4096)
	sess.errorChan = make(chan error)

	ws, err := sess.dialer()
	if err != nil {
		return nil, err
	}
	sess.SetConn(ws)

	// set up the goroutines to listen, read, write, and handle errors
	go sess.Listen()
	go sess.handleRead()
	go sess.handleWrite()
	go sess.handleError()

	// stop here until the HELLO event is receieved
	<-sess.helloReceived

	var identifyData gateway.Payload
	identifyData.OpCode = gateway.Identify

	// let her rip tater chip
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

func (s *Session) dialer() (*websocket.Conn, error) {
	var url string
	if s.GetResumeURL() != nil {
		url = *s.GetResumeURL()
	} else {
		var err error
		if s.GetToken() == nil {
			return nil, fmt.Errorf("token not set for session")
		}
		url, err = getGatewayUrl(*s.GetToken())
		if err != nil {
			return nil, err
		}
	}

	ws, err := websocket.Dial(url+"/?v=10&encoding=json", "", "http://localhost/")
	if err != nil {
		return nil, err
	}
	return ws, nil
}

// http request functions
func (s *Session) GetMessageRequest(channelId, messageId string) (*structs.Message, error) {
	if s.GetToken() == nil {
		return nil, errors.New("token not set for session")
	}
	path := "/channels/" + channelId + "/messages/" + messageId
	headers := map[string]string{
		"Authorization": "Bot " + *s.Token,
	}

	resp, err := requestutil.HttpRequest("GET", path, headers, nil)
	if err != nil {
		return nil, err
	}

	var message structs.Message

	err = json.Unmarshal(resp, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
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

func (s *Session) GetServers() map[string]structs.Server {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.Servers
}

func (s *Session) SetServers(servers map[string]structs.Server) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Servers = servers
}

func (s *Session) GetServerByName(name string) *structs.Server {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for _, guild := range s.Servers {
		if guild.Name == name {
			return &guild
		}
	}

	return nil
}

func (s *Session) AddServer(server structs.Server) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Servers[server.ID.ToString()] = server
}
