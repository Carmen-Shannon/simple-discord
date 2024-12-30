package structs

import "time"

type VoiceState struct {
	GuildID                 *Snowflake   `json:"guild_id,omitempty"`
	ChannelID               *Snowflake   `json:"channel_id,omitempty"`
	UserID                  Snowflake    `json:"user_id"`
	Member                  *GuildMember `json:"member,omitempty"`
	SessionID               string       `json:"session_id"`
	IsDeafened              bool         `json:"deaf"`
	IsMuted                 bool         `json:"mute"`
	IsSelfDeafened          bool         `json:"self_deaf"`
	IsSelfMuted             bool         `json:"self_mute"`
	IsStreaming             *bool        `json:"self_stream,omitempty"`
	IsVideo                 bool         `json:"self_video"`
	IsSurpressed            bool         `json:"suppress"`
	RequestToSpeakTimestamp *time.Time   `json:"request_to_speak_timestamp,omitempty"`
}

type VoiceRegion struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	IsOptimal     bool   `json:"optimal"`
	IsDeprecrated bool   `json:"deprecated"`
	IsCustom      bool   `json:"custom"`
}

type SpeakingEvent struct {
	Speaking Bitfield[SpeakingFlag] `json:"speaking"`
	Delay    int                    `json:"delay"`
	SSRC     int                    `json:"ssrc"`
}

type SpeakingFlag int64

const (
	SpeakingFlagMicrophone SpeakingFlag = 1 << 0
	SpeakingFlagSoundshare SpeakingFlag = 1 << 1
	SpeakingFlagPriority   SpeakingFlag = 1 << 2
)
