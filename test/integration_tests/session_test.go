package integration_tests

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/Carmen-Shannon/simple-discord/gateway"
	requestutil "github.com/Carmen-Shannon/simple-discord/gateway/request_util"
	"github.com/Carmen-Shannon/simple-discord/session"
	"github.com/Carmen-Shannon/simple-discord/structs"
)

// example of a possible test structure for testing the actual session logic
func TestNewSession(t *testing.T) {
	mockGateway := mockWebSocketServer()
	defer mockGateway.Close()

	monkey.Patch(requestutil.GetGatewayUrl, func(token string) (string, error) {
		return "ws://localhost:443", nil
	})
	defer monkey.Unpatch(requestutil.GetGatewayUrl)

	monkey.PatchInstanceMethod(reflect.TypeOf(&session.EventHandler{}), "HandleEvent", func(_ *session.EventHandler, _ *session.Session, _ gateway.Payload) error {
		return nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(&session.EventHandler{}), "HandleEvent")

	mockSession, err := session.NewSession("test", []structs.Intent{})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if mockSession == nil {
		t.Fatalf("Expected session, got nil")
	}
}

// func TestSessionReconnect(t *testing.T) {
// 	mockGateway := mockWebSocketServer()
// 	defer mockGateway.Close()

// 	monkey.Patch(session.GetGatewayUrl, func(token string) (string, error) {
// 		return "ws://localhost:443", nil
// 	})
// 	defer monkey.Unpatch(session.GetGatewayUrl)


// }
