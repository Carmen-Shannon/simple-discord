package gateway

import "github.com/Carmen-Shannon/simple-discord/structs"

type GatewayOpCode int

const (
	GatewayOpDispatch            GatewayOpCode = 0
	GatewayOpHeartbeat           GatewayOpCode = 1
	GatewayOpIdentify            GatewayOpCode = 2
	GatewayOpPresenceUpdate      GatewayOpCode = 3
	GatewayOpVoiceStateUpdate    GatewayOpCode = 4
	GatewayOpResume              GatewayOpCode = 6
	GatewayOpReconnect           GatewayOpCode = 7
	GatewayOpRequestGuildMembers GatewayOpCode = 8
	GatewayOpInvalidSession      GatewayOpCode = 9
	GatewayOpHello               GatewayOpCode = 10
	GatewayOpHeartbeatACK        GatewayOpCode = 11
)

type VoiceOpCode int

const (
	VoiceOpIdentify VoiceOpCode = iota
	VoiceOpSelectProtocol
	VoiceOpReady
	VoiceOpHeartbeat
	VoiceOpSessionDescription
	VoiceOpSpeaking
	VoiceOpHeartbeatAck
	VoiceOpResume
	VoiceOpHello
	VoiceOpResumed
	VoiceOpClientsConnect              VoiceOpCode = 11
	VoiceOpClientDisconnect            VoiceOpCode = 13
	VoiceOpPrepareTransition           VoiceOpCode = 21
	VoiceOpExecuteTransition           VoiceOpCode = 22
	VoiceOpTransitionReady             VoiceOpCode = 23
	VoiceOpPrepareEpoch                VoiceOpCode = 24
	VoiceOpMLSExternalSender           VoiceOpCode = 25
	VoiceOpMLSKeyPackage               VoiceOpCode = 26
	VoiceOpMLSProposals                VoiceOpCode = 27
	VoiceOpMLSCommitWelcome            VoiceOpCode = 28
	VoiceOpMLSAnnounceCommitTransition VoiceOpCode = 29
	VoiceOpMLSWelcome                  VoiceOpCode = 30
	VoiceOpMLSInvalidCommitWelcome     VoiceOpCode = 31
)

var (
	Opus = structs.Codec{
		Name:        "opus",
		Type:        "audio",
		Priority:    1000,
		PayloadType: 120,
	}
)

type TransportEncryptionMode string

const (
	AEAD_AES256_GCM         TransportEncryptionMode = "aead_aes256_gcm_rtpsize"
	AEAD_XCHACHA20_POLY1305 TransportEncryptionMode = "aead_xchacha20_poly1305_rtpsize"
)

type UdpData struct {
	SSRC    int                     `json:"ssrc"`
	Address string                  `json:"address"`
	Port    int                     `json:"port"`
	Mode    TransportEncryptionMode `json:"mode"`
}
