package structs

type Intent int64

const (
	GuildsIntent                      Intent = 1 << 0
	GuildMembersIntent                Intent = 1 << 1
	GuildModerationIntent             Intent = 1 << 2
	GuildEmojisAndStickersIntens      Intent = 1 << 3
	GuildIntegrationsIntent           Intent = 1 << 4
	GuildWebhooksIntent               Intent = 1 << 5
	GuildInvitesIntent                Intent = 1 << 6
	GuildVoiceStatesIntent            Intent = 1 << 7
	GuildPresencesIntent              Intent = 1 << 8
	GuildMessagesIntent               Intent = 1 << 9
	GuildMessageReactionsIntent       Intent = 1 << 10
	GuildMessageTypingIntent          Intent = 1 << 11
	DirectMessagesIntent              Intent = 1 << 12
	DirectMessageReactionsIntent      Intent = 1 << 13
	DirectMessageTypingIntent         Intent = 1 << 14
	MessageContentIntent              Intent = 1 << 15
	GuildScheduledEventsIntent        Intent = 1 << 16
	AutoModerationConfigurationIntent Intent = 1 << 20
	AutoModerationExecutionIntent     Intent = 1 << 21
	GuildMessagePollsIntent           Intent = 1 << 24
	DirectMessagePollsIntent          Intent = 1 << 25
)

func GetIntents(intents []Intent) int {
	var intent int
	for _, i := range intents {
		intent |= int(i)
	}
	return intent
}
