package simplesocket

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	GatewayURL = "wss://gateway.discord.gg/?v=10&encoding=json"
)

var (
	conn              net.Conn
	heartbeatInterval time.Duration
	lastSequence      int
)

func ConnectToGateway() error {
	u, err := url.Parse(GatewayURL); if err != nil { return err }

	hostPort := u.Host
	if u.Port() == "" {
		if u.Scheme == "wss" {
			hostPort += ":443"
		} else {
			hostPort += ":80"
		}
	}

	conn, err := net.Dial("tcp", hostPort); if err != nil { return err }
	key, err := GenerateWebSocketKey(); if err != nil { return err }

	req := &http.Request{
		Method: "GET",
		URL: u,
		Header: http.Header{
			"Connection": {"Upgrade"},
			"Upgrade": {"websocket"},
			"Sec-WebSocket-Key": {key},
			"Sec-WebSocket-Version": {"13"},
			"Host": {u.Host},
		},
	}

	if err := req.Write(conn); err != nil { return err }

	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, req); if err != nil { return err }
	if resp.StatusCode != http.StatusSwitchingProtocols {
		return errors.New("Failed to upgrade connection")
	}
	return nil
}

func GenerateWebSocketKey() (string, error) {
    key := make([]byte, 16)
    _, err := rand.Read(key)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(key), nil
}