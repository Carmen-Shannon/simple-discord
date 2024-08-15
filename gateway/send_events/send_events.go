package sendevents

import (
	"encoding/json"

	"github.com/Carmen-Shannon/simple-discord/structs"
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

type ResumeEvent struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
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
	GuildID   structs.Snowflake  `json:"guild_id"`
	ChannelID *structs.Snowflake `json:"channel_id,omitempty"`
	SelfMute  bool               `json:"self_mute"`
	SelfDeaf  bool               `json:"self_deaf"`
}

type PresenceUpdateEvent struct {
	Since      int                `json:"since"`
	Activities []structs.Activity `json:"activities"`
	Status     string             `json:"status"`
	Afk        bool               `json:"afk"`
}
