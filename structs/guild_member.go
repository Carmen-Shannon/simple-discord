package structs

import "time"

type GuildMember struct {
	User                 *User `json:"user,omitempty"`
	Nickname             *string `json:"nick,omitempty"`
	Avatar               *string `json:"avatar,omitempty"`
	Roles                []Snowflake `json:"roles"`
	Joined               time.Time `json:"joined_at"`
	PremiumSince         *time.Time `json:"premium_since,omitempty"`
	IsDeafened           bool `json:"deaf"`
	IsMute               bool `json:"mute"`
	Flags                GuildMemberFlag `json:"flags"`
	Pending              *bool `json:"pending,omitempty"`
	Permissions          *string `json:"permissions,omitempty"`
	TimeoutUntil         *time.Time `json:"timeout_until,omitempty"`
	AvatarDecorationData AvatarDecorationData `json:"avatar_decoration_data"`
}

type GuildMemberFlag int

const (
	GuildMemberFlagDidRejoin            GuildMemberFlag = 1 << 0
	GuildMemberFlagCompletedOnboarding  GuildMemberFlag = 1 << 1
	GuildMemberFlagBypassesVerification GuildMemberFlag = 1 << 2
	GuildMemberFlagStartedOnboarding    GuildMemberFlag = 1 << 3
)
