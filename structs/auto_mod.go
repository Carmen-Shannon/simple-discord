package structs

type AutoModerationRule struct {
	ID              Snowflake                     `json:"id"`
	GuildID         Snowflake                     `json:"guild_id"`
	Name            string                        `json:"name"`
	CreatorID       Snowflake                     `json:"creator_id"`
	EventType       AutoModerationEventType       `json:"event_type"`
	TriggerType     AutoModerationTriggerType     `json:"trigger_type"`
	TriggerMetadata AutoModerationTriggerMetadata `json:"trigger_metadata"`
	Actions         []AutoModerationAction        `json:"actions"`
	Enabled         bool                          `json:"enabled"`
	ExemptRoles     []Snowflake                   `json:"exempt_roles"`
	ExemptChannles  []Snowflake                   `json:"exempt_channels"`
}

type AutoModerationEventType int

const (
	AutoModMessageSend  AutoModerationEventType = 1
	AutoModMemberUpdate AutoModerationEventType = 2
)

type AutoModerationTriggerType int

const (
	AutoModKeywordTrigger       AutoModerationTriggerType = 1
	AutoModSpamTrigger          AutoModerationTriggerType = 3
	AutoModKeywordPresetTrigger AutoModerationTriggerType = 4
	AutoModMentionSpamTrigger   AutoModerationTriggerType = 5
	AutoModMemberProfile        AutoModerationTriggerType = 6
)

type KeywordPresetType int

const (
	ProfanityKeywordPreset     KeywordPresetType = 1
	SexualContentKeywordPreset KeywordPresetType = 2
	SlursKeywordPreset         KeywordPresetType = 3
)

type AutoModerationTriggerMetadata struct {
	KeywordFilter                []string            `json:"keyword_filter"`
	RegexPatterns                []string            `json:"regex_patterns"`
	Presets                      []KeywordPresetType `json:"presets"`
	AllowList                    []string            `json:"allow_list"`
	MentionTotalLimit            int                 `json:"mention_total_limit"`
	MentionRaidProtectionEnabled bool                `json:"mention_raid_protection_enabled"`
}

type AutoModerationActionType int

const (
	AutoModBlockMessageAction           AutoModerationActionType = 1
	AutoModSendAlertMessageAction       AutoModerationActionType = 2
	AutoModTimeoutAction                AutoModerationActionType = 3
	AutoModBlockMemberInteractionAction AutoModerationActionType = 4
)

type AutoModerationActionMetadata struct {
	ChannelID       Snowflake `json:"channel_id"`
	DurationSeconds int       `json:"duration_seconds"`
	CustomMessage   *string   `json:"custom_message"`
}

type AutoModerationAction struct {
	Type     AutoModerationActionType     `json:"type"`
	Metadata AutoModerationActionMetadata `json:"metadata"`
}
