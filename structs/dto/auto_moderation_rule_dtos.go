package dto

import "github.com/Carmen-Shannon/simple-discord/structs"

type GetAutoModerationRuleDto struct {
	GuildID structs.Snowflake `json:"guild_id"`
	RuleID  structs.Snowflake `json:"rule_id"`
}

type CreateAutoModerationRuleDto struct {
	GuildID         structs.Snowflake                      `json:"-"`
	Name            string                                 `json:"name"`
	EventType       structs.AutoModerationEventType        `json:"event_type"`
	TriggerType     structs.AutoModerationTriggerType      `json:"trigger_type"`
	TriggerMetadata *structs.AutoModerationTriggerMetadata `json:"trigger_metadata,omitempty"`
	Actions         []structs.AutoModerationAction         `json:"actions"`
	Enabled         *bool                                  `json:"enabled,omitempty"`
	ExemptRoles     []structs.Snowflake                    `json:"exempt_roles"`
	ExemptChannels  []structs.Snowflake                    `json:"exempt_channels"`
}

type ModifyAutoModerationRuleDto struct {
	GuildID         structs.Snowflake                      `json:"-"`
	RuleID          structs.Snowflake                      `json:"-"`
	Name            *string                                `json:"name,omitempty"`
	EventType       *structs.AutoModerationEventType       `json:"event_type,omitempty"`
	TriggerMetaData *structs.AutoModerationTriggerMetadata `json:"trigger_metadata,omitempty"`
	Actions         *[]structs.AutoModerationAction        `json:"actions,omitempty"`
	Enabled         *bool                                  `json:"enabled,omitempty"`
	ExemptRoles     *[]structs.Snowflake                   `json:"exempt_roles,omitempty"`
	ExemptChannels  *[]structs.Snowflake                   `json:"exempt_channels,omitempty"`
}
