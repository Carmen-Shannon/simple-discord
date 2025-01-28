package session

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/Carmen-Shannon/simple-discord/structs/gateway/payload"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway/ws"
	"github.com/coder/websocket"
)

var errWriteLimit = errors.New("write limit exceeded")

type wReq struct {
	data   []byte
	binary bool
}

type session struct {
	mu     *sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	conn ws.Connection

	payloadDecoders map[string]payload.Payload
	eventDecoders   map[string]payload.Payload
	validCloseErrs  map[error]struct{}

	statusCodeHandlers map[websocket.StatusCode]func()
	errorHandlers      map[error]func()

	listenFunc func(p payload.Payload) (any, error)
	handleFunc func(p payload.Payload) error

	wLimit *int

	readChan  chan []byte
	writeChan chan wReq
	errorChan chan error
}

type Session interface {
	Write(data []byte, binary bool)
	Connect(gateway string, udp bool) error
	Exit(graceful bool) error
	Error(err error)
	SetPayloadDecoders(decoders ...payload.Payload)
	SetEventDecoders(decoder ...payload.Payload)
	SetValidCloseErrors(err ...error)
	SetStatusCodeHandlers(map[websocket.StatusCode]func())
	SetErrorHandlers(map[error]func())
	SetListenFunc(f func(p payload.Payload) (any, error))
	SetHandleFunc(f func(p payload.Payload) error)
	SetWriteLimit(limit int)
}

func NewSession() Session {
	s := &session{
		mu:        &sync.Mutex{},
		conn:      nil,
		readChan:  make(chan []byte),
		writeChan: make(chan wReq),
		errorChan: make(chan error),
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s
}

func (s *session) Write(data []byte, binary bool) {
	if s.canWrite(len(s.writeChan)) {
		writeReq := wReq{
			data:   data,
			binary: binary,
		}
		write(s.ctx, s.writeChan, writeReq)
	} else {
		write(s.ctx, s.errorChan, errWriteLimit)
	}
}

func (s *session) Connect(gateway string, udp bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if udp {
		addr, err := net.ResolveUDPAddr("udp", gateway)
		if err != nil {
			return err
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			return err
		}

		s.mu.Lock()
		s.conn = ws.NewUdpConn(conn)
		s.mu.Unlock()
	} else {
		conn, _, err := websocket.Dial(ctx, gateway, nil)
		if err != nil {
			return err
		}

		s.mu.Lock()
		s.conn = ws.NewWebSocketConn(conn)
		s.mu.Unlock()
	}

	go s.listen()
	go s.read()
	go s.write()
	go s.error()

	return nil
}

func (s *session) Exit(graceful bool) error {
	defer s.cancel()
	defer close(s.readChan)
	defer close(s.writeChan)
	defer close(s.errorChan)

	if graceful {
		if err := s.conn.Close(int(websocket.StatusNormalClosure), "disconnect"); err != nil {
			if !s.isValidCloseErr(err) {
				return err
			}
		}
		return nil
	}

	if err := s.conn.CloseNow(); err != nil {
		if !s.isValidCloseErr(err) {
			return err
		}
	}
	return nil
}

func (s *session) Error(err error) {
	write(s.ctx, s.errorChan, err)
}

func (s *session) listen() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case msg, ok := <-s.readChan:
			if !ok {
				return
			}

			event, err := s.decodeEvent(msg)
			if err != nil {
				write(s.ctx, s.errorChan, err)
				continue
			}
			if event == nil {
				write(s.ctx, s.errorChan, errors.New("no idea what went wrong here"))
				continue
			}
			if !s.isValidEvent(*event) {
				write(s.ctx, s.errorChan, errors.New("invalid event"))
				continue
			}
			if s.listenFunc != nil {
				if _, err = s.listenFunc(*event); err != nil {
					write(s.ctx, s.errorChan, err)
					continue
				}
			}
			if s.handleFunc != nil {
				if err = s.handleFunc(*event); err != nil {
					write(s.ctx, s.errorChan, err)
					continue
				}
			}
		}
	}
}

func (s *session) read() {
	var buffer bytes.Buffer
	if ws, ok := s.conn.(*ws.WsConn); ok {
		ws.SetReadLimit(-1)
	}

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			bytes, err := s.conn.Read(s.ctx)
			if err != nil {
				if status := websocket.CloseStatus(err); status >= 0 {
					s.handleStatusCode(status)
					return
				}
				s.handleError(err)
				return
			}

			buffer.Write(bytes)

			for {
				startOffset := buffer.Len()
				payload, err := s.decodePayload(buffer.Bytes())
				if err != nil {
					if err == io.EOF || err == io.ErrUnexpectedEOF {
						fmt.Println("INCOMPLETE MESSAGE?")
						if startOffset <= buffer.Len() {
							buffer.Truncate(startOffset)
						}
						break
					}
					write(s.ctx, s.errorChan, err)
					buffer.Reset()
					break
				}

				if payload == nil {
					write(s.ctx, s.errorChan, errors.New("nil payload received"))
					buffer.Reset()
					break
				}

				payloadBytes, _ := (*payload).Marshal()
				write(s.ctx, s.readChan, payloadBytes)

				remainingData := buffer.Bytes()[startOffset:]
				if len(remainingData) > 0 {
					buffer.Reset()
					buffer.Write(remainingData)
				} else {
					buffer.Reset()
					break
				}
			}
		}
	}
}

func (s *session) write() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case data, ok := <-s.writeChan:
			if !ok {
				return
			}
			if err := s.conn.Write(s.ctx, data.data, data.binary); err != nil {
				if !s.isValidCloseErr(err) {
					write(s.ctx, s.errorChan, err)
				}
				s.handleError(err)
				return
			}
		}
	}
}

func (s *session) error() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case err, ok := <-s.errorChan:
			if !ok {
				return
			}
			log.Printf("session error: %v", err)
		}
	}
}

func (s *session) canWrite(len int) bool {
	if s.wLimit == nil {
		return true
	}
	return len < *s.wLimit
}

func (s *session) SetPayloadDecoders(decoders ...payload.Payload) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(decoders) == 0 {
		return
	}

	if s.payloadDecoders == nil {
		s.payloadDecoders = make(map[string]payload.Payload)
	}

	for _, d := range decoders {
		s.payloadDecoders[d.Type()] = d
	}
}

func (s *session) decodePayload(data []byte) (*payload.Payload, error) {
	var err error
	var msg payload.Payload
	for _, d := range s.payloadDecoders {
		dType := reflect.TypeOf(d)

		if dType.Kind() == reflect.Ptr {
			dType = dType.Elem()
		}

		msg = reflect.New(dType).Interface().(payload.Payload)
		err = msg.Unmarshal(data)
		if err == nil {
			return &msg, nil
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return nil, errors.New("no payload decoder found")
}

func (s *session) SetEventDecoders(decoder ...payload.Payload) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(decoder) == 0 {
		return
	}

	if s.eventDecoders == nil {
		s.eventDecoders = make(map[string]payload.Payload)
	}

	for _, d := range decoder {
		s.eventDecoders[d.Type()] = d
	}
}

func (s *session) decodeEvent(data []byte) (*payload.Payload, error) {
	var err error
	var msg payload.Payload
	for _, d := range s.eventDecoders {
		dType := reflect.TypeOf(d)

		if dType.Kind() == reflect.Ptr {
			dType = dType.Elem()
		}

		msg = reflect.New(dType).Interface().(payload.Payload)
		err = msg.Unmarshal(data)
		if err == nil {
			return &msg, nil
		}
	}
	if err != nil {
		return nil, err
	}
	return nil, errors.New("no event decoder found")
}

func (s *session) isValidEvent(payload payload.Payload) bool {
	if s.listenFunc != nil {

		return true
	}

	_, ok := s.eventDecoders[payload.Type()]
	return ok
}

func (s *session) SetValidCloseErrors(err ...error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(err) == 0 {
		return
	}

	if s.validCloseErrs == nil {
		s.validCloseErrs = make(map[error]struct{})
	}

	for _, e := range err {
		s.validCloseErrs[e] = struct{}{}
	}
}

func (s *session) isValidCloseErr(err error) bool {
	_, ok := s.validCloseErrs[err]
	return ok
}

func (s *session) SetStatusCodeHandlers(handlers map[websocket.StatusCode]func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(handlers) == 0 {
		return
	}

	if s.statusCodeHandlers == nil {
		s.statusCodeHandlers = make(map[websocket.StatusCode]func())
	}

	for code, handler := range handlers {
		s.statusCodeHandlers[code] = handler
	}
}

func (s *session) handleStatusCode(code websocket.StatusCode) {
	if handler, ok := s.statusCodeHandlers[code]; ok {
		handler()
	}
}

func (s *session) SetErrorHandlers(handlers map[error]func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(handlers) == 0 {
		return
	}

	if s.errorHandlers == nil {
		s.errorHandlers = make(map[error]func())
	}

	for err, handler := range handlers {
		s.errorHandlers[err] = handler
	}
}

func (s *session) handleError(err error) {
	if handler, ok := s.errorHandlers[err]; ok {
		fmt.Println("HANDLING ERROR:", err)
		handler()
	}
}

func (s *session) SetListenFunc(f func(p payload.Payload) (any, error)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.listenFunc = f
}

func (s *session) SetHandleFunc(f func(p payload.Payload) error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.handleFunc = f
}

func (s *session) SetWriteLimit(limit int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limit == 0 {
		return
	}

	s.wLimit = &limit
}

func write[T any](ctx context.Context, channel chan T, data T) {
	select {
	case channel <- data:
	case <-ctx.Done():
		return
	}
}
