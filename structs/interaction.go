package structs

type ResolvedData struct {
	Users       *map[Snowflake]User        `json:"users,omitempty"`
	Members     *map[Snowflake]GuildMember `json:"members,omitempty"`
	Roles       *map[Snowflake]Role        `json:"roles,omitempty"`
	Channels    *map[Snowflake]Channel     `json:"channels,omitempty"`
	Attachments *map[Snowflake]Attachment  `json:"attachments,omitempty"`
}

type InteractionData struct {
	ID       Snowflake    `json:"id"`
	Name     string       `json:"name"`
	Type     int          `json:"type"`
	Resolved ResolvedData `json:"resolved,omitempty"`
}

type InteractionType int

const (
	PingInteraction                           InteractionType = 1
	ApplicationCommandInteraction             InteractionType = 2
	MessageComponentInteraction               InteractionType = 3
	ApplicationCommandAutocompleteInteraction InteractionType = 4
	ModalSubmitInteraction                    InteractionType = 5
)

type InteractionResponseType int

const (
	PongInteraction                                 InteractionResponseType = 1
	ChannelMessageWithSourceInteraction             InteractionResponseType = 4
	DeferredChannelMessageWithSourceInteraction     InteractionResponseType = 5
	DeferredUpdatedMessageInteraction               InteractionResponseType = 6
	UpdateMessageInteraction                        InteractionResponseType = 7
	ApplicationCommandAutocompleteResultInteraction InteractionResponseType = 8
	ModalInteraction                                InteractionResponseType = 9
	PremiumRequiredInteraction                      InteractionResponseType = 10
)

type IntegrationType int

const (
	GuildInstallType IntegrationType = 0
	UserInstallType  IntegrationType = 1
)

type ContextType int

const (
	GuildContextType          ContextType = 0
	BotDMContextType          ContextType = 1
	PrivateChannelContextType ContextType = 2
)

type InteractionResponseData struct {
	TTS             *bool                 `json:"tts,omitempty"`
	Content         *string               `json:"content,omitempty"`
	Embeds          []Embed               `json:"embeds"`
	AllowedMentions AllowedMentions       `json:"allowed_mentions"`
	Flags           Bitfield[MessageFlag] `json:"flags"`
	Components      MessageComponent      `json:"components"`
}

type Interaction struct {
	ID                           Snowflake                  `json:"id"`
	ApplicationID                Snowflake                  `json:"application_id"`
	Type                         InteractionType            `json:"type"`
	Data                         *InteractionData           `json:"data,omitempty"`
	Guild                        *Guild                     `json:"guild,omitempty"`
	GuildID                      *Snowflake                 `json:"guild_id,omitempty"`
	Channel                      *Channel                   `json:"channel,omitempty"`
	ChannelID                    *Snowflake                 `json:"channel_id,omitempty"`
	Member                       *GuildMember               `json:"member,omitempty"`
	User                         *User                      `json:"user,omitempty"`
	Token                        string                     `json:"token"`
	Version                      int                        `json:"version"`
	Message                      *Message                   `json:"message,omitempty"`
	AppPermissions               string                     `json:"app_permissions"`
	Locale                       *string                    `json:"locale,omitempty"`
	GuildLocale                  *string                    `json:"guild_locale,omitempty"`
	Entitlements                 []Entitlement              `json:"entitlements"`
	AuthorizingIntegrationOwners map[IntegrationType]string `json:"authorizing_integration_owners"`
	Context                      *ContextType               `json:"context,omitempty"`
}

type InteractionResponse struct {
	Type InteractionResponseType `json:"type"`
	Data InteractionResponseData `json:"data"`
}
