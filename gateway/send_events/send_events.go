package sendevents

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/voice"
)

type IdentifyProperties struct {
	Os      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

type IdentifyEvent struct {
	Token          string               `json:"token"`
	Properties     IdentifyProperties   `json:"properties"`
	Compress       *bool                `json:"compress,omitempty"`
	LargeThreshold *int                 `json:"large_threshold,omitempty"`
	Shard          *[]int               `json:"shard,omitempty"`
	Presence       *PresenceUpdateEvent `json:"presence,omitempty"`
	Intents        int                  `json:"intents"`
}

type VoiceIdentifyEvent struct {
	ServerID               structs.Snowflake `json:"server_id"`
	UserID                 structs.Snowflake `json:"user_id"`
	SessionID              string            `json:"session_id"`
	Token                  string            `json:"token"`
	MaxDaveProtocolVersion *int              `json:"max_dave_protocol_version,omitempty"`
}

type VoiceSelectProtocolEvent struct {
	Protocol string                  `json:"protocol"`
	Data     VoiceSelectProtocolData `json:"data"`
	Codecs   []voice.Codec           `json:"codecs"`
}

type VoiceSelectProtocolData struct {
	Address string                        `json:"address"`
	Port    int                           `json:"port"`
	Mode    voice.TransportEncryptionMode `json:"mode"`
}

type SpeakingEvent struct {
	structs.SpeakingEvent
}

type VoiceResumeEvent struct {
	ServerID  structs.Snowflake `json:"server_id"`
	SessionID string            `json:"session_id"`
	Token     string            `json:"token"`
	SeqAck    int               `json:"seq_ack"`
}

type VoiceDaveReadyForTransitionEvent struct {
	TransitionID int `json:"transition_id"`
}

type VoiceDaveMlsKeyPackageEvent struct {
	OpCode     uint8
	MLSMessage voice.KeyPackage
}

func (v *VoiceDaveMlsKeyPackageEvent) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Read OpCode
	if err := binary.Read(buf, binary.BigEndian, &v.OpCode); err != nil {
		return fmt.Errorf("failed to read opcode: %w", err)
	}

	// Read remaining bytes into KeyPackage
	remainingBytes := make([]byte, buf.Len())
	if _, err := buf.Read(remainingBytes); err != nil {
		return fmt.Errorf("failed to read remaining bytes: %w", err)
	}

	// Unmarshal remaining bytes into KeyPackage
	if err := v.MLSMessage.UnmarshalBinary(remainingBytes); err != nil {
		return fmt.Errorf("failed to unmarshal KeyPackage: %w", err)
	}

	return nil
}

func (v *VoiceDaveMlsKeyPackageEvent) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write OpCode
	if err := binary.Write(buf, binary.BigEndian, v.OpCode); err != nil {
		return nil, fmt.Errorf("failed to write opcode: %w", err)
	}

	// Marshal KeyPackage
	keyPackageBytes, err := v.MLSMessage.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal KeyPackage: %w", err)
	}

	// Write KeyPackage
	if _, err := buf.Write(keyPackageBytes); err != nil {
		return nil, fmt.Errorf("failed to write KeyPackage: %w", err)
	}

	return buf.Bytes(), nil
}

type VoiceDaveMlsCommitWelcomeEvent struct {
	OpCode     uint8
	MLSMessage voice.Commit
}

func (v *VoiceDaveMlsCommitWelcomeEvent) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Read OpCode
	if err := binary.Read(buf, binary.BigEndian, &v.OpCode); err != nil {
		return fmt.Errorf("failed to read opcode: %w", err)
	}

	// Read remaining bytes into Commit
	remainingBytes := make([]byte, buf.Len())
	if _, err := buf.Read(remainingBytes); err != nil {
		return fmt.Errorf("failed to read remaining bytes: %w", err)
	}

	// Unmarshal remaining bytes into Commit
	if err := v.MLSMessage.UnmarshalBinary(remainingBytes); err != nil {
		return fmt.Errorf("failed to unmarshal Commit: %w", err)
	}

	return nil
}

func (v *VoiceDaveMlsCommitWelcomeEvent) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write OpCode
	if err := binary.Write(buf, binary.BigEndian, v.OpCode); err != nil {
		return nil, fmt.Errorf("failed to write opcode: %w", err)
	}

	// Marshal Commit
	commitBytes, err := v.MLSMessage.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Commit: %w", err)
	}

	// Write Commit
	if _, err := buf.Write(commitBytes); err != nil {
		return nil, fmt.Errorf("failed to write Commit: %w", err)
	}

	return buf.Bytes(), nil
}

type ResumeEvent struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

type VoiceHeartbeatEvent struct {
	Timestamp int64 `json:"t"`
	SeqAck    *int  `json:"seq_ack,omitempty"`
}

func (v *VoiceHeartbeatEvent) UnmarshalJSON(data []byte) error {
	var sequence int
	if err := json.Unmarshal(data, &sequence); err == nil {
		v.SeqAck = &sequence
		return nil
	}

	return errors.New("unable to unmarshal VoiceHeartbeatEvent struct")
}

func (v *VoiceHeartbeatEvent) MarshalJSON() ([]byte, error) {
	if v.SeqAck == nil {
		return []byte("null"), nil
	}

	return json.Marshal(*v.SeqAck)
}

type HeartbeatEvent struct {
	LastSequence *int `json:"-"`
}

func (h *HeartbeatEvent) MarshalJSON() ([]byte, error) {
	if h.LastSequence == nil {
		return []byte("null"), nil
	}

	return json.Marshal(*h.LastSequence)
}

type RequestGuildMembersEvent struct {
	GuildID   structs.Snowflake   `json:"guild_id"`
	Query     string              `json:"query"`
	Limit     int                 `json:"limit"`
	Presences bool                `json:"presences"`
	UserIDs   []structs.Snowflake `json:"user_ids"`
	Nonce     *string             `json:"nonce,omitempty"`
}

type UpdateVoiceStateEvent struct {
	GuildID   *structs.Snowflake `json:"guild_id"`
	ChannelID *structs.Snowflake `json:"channel_id"`
	SelfMute  bool               `json:"self_mute"`
	SelfDeaf  bool               `json:"self_deaf"`
}

type PresenceUpdateEvent struct {
	Since      int                `json:"since"`
	Activities []structs.Activity `json:"activities"`
	Status     string             `json:"status"`
	Afk        bool               `json:"afk"`
}
