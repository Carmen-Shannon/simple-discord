package structs

import (
	"time"
)

type Server struct {
	*Guild
	JoinedAt             time.Time             `json:"joined_at"`
	Large                bool                  `json:"large"`
	Unavailable          *bool                 `json:"unavailable,omitempty"`
	MemberCount          int                   `json:"member_count"`
	VoiceStates          []VoiceState          `json:"voice_states"`
	Members              []GuildMember         `json:"members"`
	Channels             []Channel             `json:"channels"`
	Threads              []Channel             `json:"threads"`
	Presences            []PresenceUpdate      `json:"presences"`
	StageInstances       []StageInstance       `json:"stage_instances"`
	GuildScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events"`
}

type PresenceUpdate struct {
	User         User           `json:"user"`
	GuildID      Snowflake      `json:"guild_id"`
	Status       UserStatusType `json:"status"`
	Activities   []Activity     `json:"activities"`
	ClientStatus ClientStatus   `json:"client_status"`
	Nonce        *string        `json:"nonce,omitempty"`
}

type UserStatusType string

const (
	UserOnline    UserStatusType = "online"
	UserDND       UserStatusType = "dnd"
	UserIdle      UserStatusType = "idle"
	UserInvisible UserStatusType = "invisible"
	UserOffline   UserStatusType = "offline"
)

type ClientStatus struct {
	Desktop *string `json:"desktop,omitempty"`
	Mobile  *string `json:"mobile,omitempty"`
	Web     *string `json:"web,omitempty"`
}
