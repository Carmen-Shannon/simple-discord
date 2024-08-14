package structs

type StageInstance struct {
	ID                    Snowflake                 `json:"id"`
	GuildID               Snowflake                 `json:"guild_id"`
	ChannelID             Snowflake                 `json:"channel_id"`
	Topic                 string                    `json:"topic"`
	PrivacyLevel          StageInstancePrivacyLevel `json:"privacy_level"`
	DiscoverableDisabled  bool                      `json:"discoverable_disabled"`
	GuildScheduledEventID *Snowflake                `json:"guild_scheduled_event_id,omitempty"`
}

type StageInstancePrivacyLevel int

const (
	StageInstancePublic    StageInstancePrivacyLevel = 1
	StageInstanceGuildOnly StageInstancePrivacyLevel = 2
)
