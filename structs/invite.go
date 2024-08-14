package structs

import "time"

type Invite struct {
	Type                     InviteType           `json:"type"`
	Code                     string               `json:"code"`
	Guild                    *Guild               `json:"guild,omitempty"`
	Channel                  *Channel             `json:"channel,omitempty"`
	Inviter                  *User                `json:"inviter,omitempty"`
	TargetType               *InviteTargetType    `json:"target_type,omitempty"`
	TargetUser               *User                `json:"target_user,omitempty"`
	TargetApplication        *Application         `json:"target_application,omitempty"`
	ApproximatePresenceCount *int                 `json:"approximate_presence_count,omitempty"`
	ApproximateMemberCount   *int                 `json:"approximate_member_count,omitempty"`
	ExpiresAt                *time.Time           `json:"expires_at,omitempty"`
	StageInstance            *InviteStageInstance `json:"stage_instance,omitempty"`
	GuildScheduledEvent      *GuildScheduledEvent `json:"guild_scheduled_event,omitempty"`
}

type InviteType int

const (
	GuildInviteType   InviteType = 0
	GroupDMInviteType InviteType = 1
	FriendInviteType  InviteType = 2
)

type InviteTargetType int

const (
	StreamTargetType              InviteTargetType = 1
	EmbeddedApplicationTargetType InviteTargetType = 2
)

// DEPRECATED: This struct is deprecated and will be removed in the future.
type InviteStageInstance struct {
	Members          []GuildMember `json:"members"`
	ParticipantCount int           `json:"participant_count"`
	SpeakerCount     int           `json:"speaker_count"`
	Topic            string        `json:"topic"`
}
