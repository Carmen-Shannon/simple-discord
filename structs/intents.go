package structs

type Intents int64

const (
	GuildsIntents                      Intents = 1 << 0
	GuildMembersIntents                Intents = 1 << 1
	GuildModerationIntents             Intents = 1 << 2
	GuildEmojisAndStickersIntens       Intents = 1 << 3
	GuildIntegrationsIntents           Intents = 1 << 4
	GuildWebhooksIntents               Intents = 1 << 5
	GuildInvitesIntents                Intents = 1 << 6
	GuildVoiceStatesIntents            Intents = 1 << 7
	GuildPresencesIntents              Intents = 1 << 8
	GuildMessagesIntents               Intents = 1 << 9
	GuildMessageReactionsIntents       Intents = 1 << 10
	GuildMessageTypingIntents          Intents = 1 << 11
	DirectMessagesIntents              Intents = 1 << 12
	DirectMessageReactionsIntents      Intents = 1 << 13
	DirectMessageTypingIntents         Intents = 1 << 14
	MessageContentIntents              Intents = 1 << 15
	GuildScheduledEventsIntents        Intents = 1 << 16
	AutoModerationConfigurationIntents Intents = 1 << 20
	AutoModerationExecutionIntents     Intents = 1 << 21
	GuildMessagePollsIntents           Intents = 1 << 24
	DirectMessagePollsIntents          Intents = 1 << 25
)
