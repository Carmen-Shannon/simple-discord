package structs

type Webhook struct {
	ID            Snowflake   `json:"id"`
	Type          WebhookType `json:"type"`
	GuildID       *Snowflake  `json:"guild_id,omitempty"`
	ChannelID     Snowflake   `json:"channel_id"`
	User          *User       `json:"user"`
	Name          string      `json:"name"`
	Avatar        string      `json:"avatar"`
	Token         string      `json:"token"`
	ApplicationID Snowflake   `json:"application_id,omitempty"`
	SourceGuild   *Guild      `json:"source_guild,omitempty"`
	SourceChannel *Channel    `json:"source_channel,omitempty"`
	URL           string      `json:"url"`
}

type WebhookType int

const (
	IncomingWebhook        WebhookType = 1
	ChannelFollowerWebhook WebhookType = 2
	ApplicationWebhook     WebhookType = 3
)
