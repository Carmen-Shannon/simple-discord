package gateway

import "encoding/json"

type GatewayOpCode int

const (
	Dispatch            GatewayOpCode = 0
	Heartbeat           GatewayOpCode = 1
	Identify            GatewayOpCode = 2
	PresenceUpdate      GatewayOpCode = 3
	VoiceStateUpdate    GatewayOpCode = 4
	Resume              GatewayOpCode = 6
	Reconnect           GatewayOpCode = 7
	RequestGuildMembers GatewayOpCode = 8
	InvalidSession      GatewayOpCode = 9
	Hello               GatewayOpCode = 10
	HeartbeatACK        GatewayOpCode = 11
)

type Payload struct {
	OpCode    GatewayOpCode `json:"op"`
	Data      any           `json:"d"`
	Seq       *int          `json:"s,omitempty"`
	EventName *string       `json:"t,omitempty"`
}

func (p *Payload) ToString() string {
	jsonData, _ := json.Marshal(p)
	return string(jsonData)
}