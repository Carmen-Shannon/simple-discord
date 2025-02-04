package structs

import (
	"errors"
)

type ResolvedData struct {
	Users       map[Snowflake]User        `json:"users,omitempty"`
	Members     map[Snowflake]GuildMember `json:"members,omitempty"`
	Roles       map[Snowflake]Role        `json:"roles,omitempty"`
	Channels    map[Snowflake]Channel     `json:"channels,omitempty"`
	Attachments map[Snowflake]Attachment  `json:"attachments,omitempty"`
}

type ApplicationCommandInteractionDataOption struct {
	Name    string                                    `json:"name"`
	Type    ApplicationCommandOptionType              `json:"type"`
	Value   any                                       `json:"value,omitempty"`
	Options []ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	Focused *bool                                     `json:"focused,omitempty"`
}

type InteractionData struct {
	ID       Snowflake                                 `json:"id"`
	Name     string                                    `json:"name"`
	Type     int                                       `json:"type"`
	Resolved *ResolvedData                             `json:"resolved,omitempty"`
	Options  []ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	GuildID  *Snowflake                                `json:"guild_id,omitempty"`
	TargetID *Snowflake                                `json:"target_id,omitempty"`
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
	LaunchActivityInteraction                       InteractionResponseType = 12
)

type IntegrationType int

const (
	GuildInstallType IntegrationType = 0
	UserInstallType  IntegrationType = 1
)

type IntegrationContextType int

const (
	GuildContextType          IntegrationContextType = 0
	BotDMContextType          IntegrationContextType = 1
	PrivateChannelContextType IntegrationContextType = 2
)

type InteractionResponseData struct {
	TTS             bool                  `json:"tts"`
	Content         string                `json:"content,omitempty"`
	Embeds          []Embed               `json:"embeds.omitempty"`
	AllowedMentions *AllowedMentions      `json:"allowed_mentions,omitempty"`
	Flags           Bitfield[MessageFlag] `json:"flags,omitempty"`
	Components      []MessageComponent    `json:"components,omitempty"`
	Attachments     []Attachment          `json:"attachments,omitempty"`
	Poll            *Poll                 `json:"poll,omitempty"`
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
	Context                      *IntegrationContextType    `json:"context,omitempty"`
}

type InteractionResponse struct {
	Type InteractionResponseType  `json:"type"`
	Data *InteractionResponseData `json:"data,omitempty"`
}

type InteractionCallbackResponse struct {
	Interaction InteractionCallbackObject    `json:"interaction"`
	Resource    *InteractionCallbackResource `json:"resource,omitempty"`
}

type InteractionCallbackObject struct {
	ID                       Snowflake       `json:"id"`
	Type                     InteractionType `json:"type"`
	ActivityInstanceID       *string         `json:"activity_instance_id,omitempty"`
	ResponseMessageID        *Snowflake      `json:"response_message_id,omitempty"`
	ResponseMessageLoading   *bool           `json:"response_message_loading,omitempty"`
	ResponseMessageEphemeral *bool           `json:"response_message_ephemeral,omitempty"`
}

type InteractionCallbackResource struct {
	Type             InteractionResponseType      `json:"type"`
	ActivityInstance *InteractionActivityInstance `json:"activity_instance,omitempty"`
	Message          *Message                     `json:"message,omitempty"`
}

type InteractionActivityInstance struct {
	ID string `json:"id"`
}

type InteractionResponseOptions interface {
	InteractionResponse() *InteractionResponse
	SetResponseType(InteractionResponseType)
	SetTTS(bool)
	SetContent(string)
	SetEmbeds([]Embed) error
	SetAllowedMentions(*AllowedMentions)
	SetFlags(Bitfield[MessageFlag]) error
	SetComponents([]MessageComponent)
	SetAttachments([]Attachment)
	SetPoll(*Poll)
}

var _ InteractionResponseOptions = (*InteractionResponse)(nil)

func (i *InteractionResponse) InteractionResponse() *InteractionResponse {
	return i
}

func (i *InteractionResponse) SetResponseType(responseType InteractionResponseType) {
	i.Type = responseType
}

func (i *InteractionResponse) SetTTS(tts bool) {
	i.Data.TTS = tts
}

func (i *InteractionResponse) SetContent(content string) {
	i.Data.Content = content
}

func (i *InteractionResponse) SetEmbeds(embeds []Embed) error {
	if len(embeds) > 10 {
		return errors.New("embeds cannot exceed 10")
	}
	i.Data.Embeds = embeds
	return nil
}

func (i *InteractionResponse) SetAllowedMentions(allowedMentions *AllowedMentions) {
	i.Data.AllowedMentions = allowedMentions
}

func (i *InteractionResponse) SetFlags(flags Bitfield[MessageFlag]) error {
	for _, flag := range flags {
		if flag != EphemeralMessageFlag && flag != SurpressEmbedsMessageFlag && flag != SurpressNotificationsMessageFlag {
			return errors.New("can only accept SUPRESS_EMBEDS, EPHEMERAL, and SUPRESS_NOTIFICATIONS flags")
		}
	}
	i.Data.Flags = flags
	return nil
}

func (i *InteractionResponse) SetComponents(components []MessageComponent) {
	i.Data.Components = components
}

func (i *InteractionResponse) SetAttachments(attachments []Attachment) {
	i.Data.Attachments = attachments
}

func (i *InteractionResponse) SetPoll(poll *Poll) {
	i.Data.Poll = poll
}

func NewInteractionResponseOptions() InteractionResponseOptions {
	return &InteractionResponse{
		Type: PongInteraction,
		Data: &InteractionResponseData{},
	}
}
