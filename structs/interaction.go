package structs

type ResolvedData struct {
	Users       *map[Snowflake]User
	Members     *map[Snowflake]GuildMember
	Roles       *map[Snowflake]Role
	Channels    *map[Snowflake]Channel
	Attachments *map[Snowflake]Attachment
}

type InteractionData struct {
	ID       Snowflake
	Name     string
	Type     int
	Resolved ResolvedData
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
	TTS *bool
	Content *string
	Embeds []Embed
	AllowedMentions AllowedMentions
	Flags MessageFlag
	Components MessageComponent
}

type Interaction struct {
	ID                           Snowflake
	ApplicationID                Snowflake
	Type                         InteractionType
	Data                         *InteractionData
	Guild                        *Guild
	GuildID                      *Snowflake
	Channel                      *Channel
	ChannelID                    *Snowflake
	Member                       *GuildMember
	User                         *User
	Token                        string
	Version                      int
	Message                      *Message
	AppPermissions               string
	Locale                       *string
	GuildLocale                  *string
	Entitlements                 []Entitlement
	AuthorizingIntegrationOwners map[IntegrationType]string
	Context                      *ContextType
}

type InteractionResponse struct {
	Type InteractionResponseType
	Data InteractionResponseData
}
