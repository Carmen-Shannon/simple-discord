package integration_tests

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"time"

	receiveevents "github.com/Carmen-Shannon/simple-discord/gateway/receive_events"
	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"golang.org/x/net/websocket"
)

// for mocking the discord gateway, in case you want to manually test the Sessions response to certain payloads
func mockWebSocketServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		websocket.Handler(func(ws *websocket.Conn) {
			if err := mockHandleHello(ws); err != nil {
				return
			}

			for {
				var data []byte
				if err := websocket.Message.Receive(ws, &data); err != nil {
					return
				}

				var payload gateway.Payload
				if err := json.Unmarshal(data, &payload); err != nil {
					return
				}

				switch payload.OpCode {
				case gateway.Identify:
					if err := mockHandleIdentify(ws); err != nil {
						return
					}
				case gateway.Heartbeat:
					if err := mockHandleHeartbeatACK(ws); err != nil {
						return
					}
				default:
					continue
				}
			}
		}).ServeHTTP(w, r)
	})
	server := httptest.NewUnstartedServer(handler)
	server.Listener.Close()
	server.Listener, _ = net.Listen("tcp", "localhost:443")
	server.Start()
	return server
}

func mockWebSocketServerWithDisconnect(delay time.Duration) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		websocket.Handler(func(ws *websocket.Conn) {
			if err := mockHandleHello(ws); err != nil {
				return
			}

			go func() {
				time.Sleep(delay)
				if err := mockHandleDisconnect(ws); err != nil {
					return
				}
				ws.Close()
			}()

			for {
				var data []byte
				if err := websocket.Message.Receive(ws, &data); err != nil {
					return
				}

				var payload gateway.Payload
				if err := json.Unmarshal(data, &payload); err != nil {
					return
				}

				switch payload.OpCode {
				case gateway.Identify:
					if err := mockHandleIdentify(ws); err != nil {
						return
					}
				case gateway.Heartbeat:
					if err := mockHandleHeartbeatACK(ws); err != nil {
						return
					}
				default:
					continue
				}
			}
		}).ServeHTTP(w, r)
	})
	server := httptest.NewUnstartedServer(handler)
	server.Listener.Close()
	server.Listener, _ = net.Listen("tcp", "localhost:443")
	server.Start()
	return server
}

func mockHandleDisconnect(ws *websocket.Conn) error {
	disconnectPayload := gateway.Payload{
		OpCode: gateway.Reconnect,
		Data:   receiveevents.ReconnectEvent{},
	}

	rawPayload, _ := json.Marshal(disconnectPayload)

	if err := websocket.Message.Send(ws, rawPayload); err != nil {
		return err
	}
	return nil
}

func mockHandleHello(ws *websocket.Conn) error {
	helloPayload := gateway.Payload{
		OpCode: gateway.Hello,
		Data: receiveevents.HelloEvent{
			HeartbeatInterval: 41250,
		},
	}

	rawPayload, _ := json.Marshal(helloPayload)

	if err := websocket.Message.Send(ws, rawPayload); err != nil {
		return err
	}
	return nil
}

func mockHandleIdentify(ws *websocket.Conn) error {
	ready := "READY"
	readyPayload := gateway.Payload{
		OpCode: gateway.Dispatch,
		Data: receiveevents.ReadyEvent{
			Version: 69,
			User: structs.User{
				ID: structs.Snowflake{
					ID: 1234567890,
				},
				Username:             "test",
				Discriminator:        "test-1234",
				Flags:                structs.Bitfield[structs.UserFlag]{},
				AvatarDecorationData: structs.AvatarDecorationData{},
			},
			Guilds: []structs.UnavailableGuild{
				{
					ID: structs.Snowflake{
						ID: 1234567890,
					},
					Unavailable: false,
				},
			},
			SessionID:        "test-session-id",
			ResumeGatewayURL: "ws://localhost:443",
			Shard:            []int{0, 1},
			Application: structs.Application{
				ID: structs.Snowflake{
					ID: 1234567890,
				},
				Name:        "test-application",
				Description: "a test",
				Summary:     "a test summary",
				VerifyKey:   "a test verify key",
				Team: structs.Team{
					ID: structs.Snowflake{
						ID: 1234567890,
					},
					Name:    "test-team",
					Members: []structs.TeamMember{},
					OwnerUserID: structs.Snowflake{
						ID: 1234567890,
					},
				},
				Flags:                 structs.Bitfield[structs.ApplicationFlag]{},
				ApproximateGuildCount: 1,
				RedirectURIs:          []string{},
				Tags:                  []string{},
			},
		},
		EventName: &ready,
	}

	rawPayload, _ := json.Marshal(readyPayload)
	if err := websocket.Message.Send(ws, rawPayload); err != nil {
		return err
	}
	return nil
}

func mockHandleHeartbeatACK(ws *websocket.Conn) error {
	heartbeatACKPayload := gateway.Payload{
		OpCode: gateway.HeartbeatACK,
		Data:   nil,
	}

	rawPayload, _ := json.Marshal(heartbeatACKPayload)

	if err := websocket.Message.Send(ws, rawPayload); err != nil {
		return err
	}
	return nil
}
