package structs

import "time"

type VoiceState struct {
	GuildID                 *Snowflake
	ChannelID               *Snowflake
	UserID                  Snowflake
	Member                  *GuildMember
	SessionID               string
	IsDeafened              bool
	IsMuted                 bool
	IsSelfDeafened          bool
	IsSelfMuted             bool
	IsStreaming             *bool
	IsVideo                 bool
	IsSurpressed            bool
	RequestToSpeakTimestamp *time.Time
}

type VoiceRegion struct {
	ID            string
	Name          string
	IsOptimal     bool
	IsDeprecrated bool
	IsCustom      bool
}
