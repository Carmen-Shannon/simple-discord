package structs

import "time"

type GuildMember struct {
	User                 *User
	Nickname             *string
	Avatar               *string
	Roles                []Snowflake
	Joined               time.Time
	PremiumSince         *time.Time
	IsDeafened           bool
	IsMute               bool
	Flags                int
	Pending              *bool
	Permissions          *string
	TimeoutUntil         *time.Time
	AvatarDecorationData AvatarDecorationData
}
