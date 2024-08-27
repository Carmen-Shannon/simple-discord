package structs

type AuditLog struct {
	ApplicationCommands []ApplicationCommand `json:"application_commands"`
	AuditLogEntries []AuditLogEntry `json:"audit_log_entries"`
	AutoModerationRules []AutoModerationRule `json:"auto_moderation_rules"`
	GuildScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events"`
	Integrations []GuildIntegration `json:"integrations"`
	Threads []Channel `json:"threads"`
	Users []User `json:"users"`
	Webhooks []Webhook `json:"webhooks"`
}

type AuditLogEntry struct {
	TargetID   string                 `json:"target_id"`
	Changes    []AuditLogChange       `json:"changes"`
	UserID     *Snowflake             `json:"user_id,omitempty"`
	ID         Snowflake              `json:"id"`
	ActionType AuditLogEvent          `json:"action_type"`
	Options    OptionalAuditEntryInfo `json:"options"`
	Reason     string                 `json:"reason"`
}

type AuditLogChange struct {
	NewValue interface{} `json:"new_value"`
	OldValue interface{} `json:"old_value"`
	Key      string
}

type AuditLogEvent int

const (
	AuditGuildUpdate                             AuditLogEvent = 1
	AuditChannelCreate                           AuditLogEvent = 10
	AuditChannelUpdate                           AuditLogEvent = 11
	AuditChannelDelete                           AuditLogEvent = 12
	AuditChannelOverwriteCreate                  AuditLogEvent = 13
	AuditChannelOverwriteUpdate                  AuditLogEvent = 14
	AuditChannelOverwriteDelete                  AuditLogEvent = 15
	AuditMemberKick                              AuditLogEvent = 20
	AuditMemberPrune                             AuditLogEvent = 21
	AuditMemberBanAdd                            AuditLogEvent = 22
	AuditMemberBanRemove                         AuditLogEvent = 23
	AuditMemberUpdate                            AuditLogEvent = 24
	AuditMemberRoleUpdate                        AuditLogEvent = 25
	AuditMemberMove                              AuditLogEvent = 26
	AuditMemberDisconnect                        AuditLogEvent = 27
	AuditBotAdd                                  AuditLogEvent = 28
	AuditRoleCreate                              AuditLogEvent = 30
	AuditRoleUpdate                              AuditLogEvent = 31
	AuditRoleDelete                              AuditLogEvent = 32
	AuditInviteCreate                            AuditLogEvent = 40
	AuditInviteUpdate                            AuditLogEvent = 41
	AuditInviteDelete                            AuditLogEvent = 42
	AuditWebhookCreate                           AuditLogEvent = 50
	AuditWebhookUpdate                           AuditLogEvent = 51
	AuditWebhookDelete                           AuditLogEvent = 52
	AuditEmojiCreate                             AuditLogEvent = 60
	AuditEmojiUpdate                             AuditLogEvent = 61
	AuditEmojiDelete                             AuditLogEvent = 62
	AuditMessageDelete                           AuditLogEvent = 72
	AuditMessageBulkDelete                       AuditLogEvent = 73
	AuditMessagePin                              AuditLogEvent = 74
	AuditMessageUnpin                            AuditLogEvent = 75
	AuditIntegrationCreate                       AuditLogEvent = 80
	AuditIntegrationUpdate                       AuditLogEvent = 81
	AuditIntegrationDelete                       AuditLogEvent = 82
	AuditStageInstanceCreate                     AuditLogEvent = 83
	AuditStageInstanceUpdate                     AuditLogEvent = 84
	AuditStageInstanceDeleted                    AuditLogEvent = 85
	AuditStickerCreate                           AuditLogEvent = 90
	AuditStickerUpdate                           AuditLogEvent = 91
	AuditStickerDelete                           AuditLogEvent = 92
	AuditGuildScheduledEventCreate               AuditLogEvent = 100
	AuditGuildScheduledEventUpdate               AuditLogEvent = 101
	AuditGuildScheduledEventDelete               AuditLogEvent = 102
	AuditThreadCreate                            AuditLogEvent = 110
	AuditThreadUpdate                            AuditLogEvent = 111
	AuditThreadDelete                            AuditLogEvent = 112
	AuditApplicationCommandPermissionUpdate      AuditLogEvent = 121
	AuditAutoModerationRuleCreate                AuditLogEvent = 140
	AuditAutoModerationRuleUpdate                AuditLogEvent = 141
	AuditAutoModerationRuleDelete                AuditLogEvent = 142
	AuditAutoModerationBlockMessage              AuditLogEvent = 143
	AuditAutoModerationFlagToChannel             AuditLogEvent = 144
	AuditAutoModerationUserCommunicationDisabled AuditLogEvent = 145
	AuditCreatorMonetizationRequestCreated       AuditLogEvent = 150
	AuditCreatorMonetizationTermsAccepted        AuditLogEvent = 151
	AuditOnboardingPromptCreate                  AuditLogEvent = 163
	AuditOnboardingPrompUpdate                   AuditLogEvent = 164
	AuditOnboardingPrompDelete                   AuditLogEvent = 165
	AuditOnboardingCreate                        AuditLogEvent = 166
	AuditOnboardingUpdate                        AuditLogEvent = 167
	AuditHomeSettingsCreate                      AuditLogEvent = 190
	AuditHomeSettingsUpdate                      AuditLogEvent = 191
)

type OptionalAuditEntryInfo struct {
	ApplicationID                 Snowflake `json:"application_id"`
	AutoModerationRuleName        string    `json:"auto_moderation_rule_name"`
	AutoModerationRuleTriggerType string    `json:"auto_moderation_rule_trigger_type"`
	ChannelID                     Snowflake `json:"channel_id"`
	Count                         string    `json:"count"`
	DeleteMemberDays              string    `json:"delete_member_days"`
	ID                            Snowflake `json:"id"`
	MembersRemoved                string    `json:"members_removed"`
	MessageID                     Snowflake `json:"message_id"`
	RoleName                      string    `json:"role_name"`
	Type                          string    `json:"type"`
	IntegrationType               string    `json:"integration_type"`
}
