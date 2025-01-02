package structs

import (
	"encoding/json"
	"errors"
	"time"
)

type AttachmentFlag int64

const (
	IsRemix AttachmentFlag = 1 << 2
)

type AllowedMentionType string

const (
	RoleMentionsType     AllowedMentionType = "roles"
	UserMentionsType     AllowedMentionType = "users"
	EveryoneMentionsType AllowedMentionType = "everyone"
)

type MessageFlag int64

const (
	CrossPostedMessageFlag              MessageFlag = 1 << 0
	IsCrossPostedMessageFlag            MessageFlag = 1 << 1
	SurpressEmbedsMessageFlag           MessageFlag = 1 << 2
	SourceMessageDeletedMessageFlag     MessageFlag = 1 << 3
	UrgentMessageFlag                   MessageFlag = 1 << 4
	HasThreadMessageFlag                MessageFlag = 1 << 5
	EphemeralMessageFlag                MessageFlag = 1 << 6
	LoadingMessageFlag                  MessageFlag = 1 << 7
	FailedToMentionSomeRolesMessageFlag MessageFlag = 1 << 8
	SurpressNotificationsMessageFlag    MessageFlag = 1 << 12
	IsVoiceMessageMessageFlag           MessageFlag = 1 << 13
)

type MessageActivityType int

const (
	JoinMessageActivityType        MessageActivityType = 1
	SpectateMessageActivityType    MessageActivityType = 2
	ListenMessageActivityType      MessageActivityType = 3
	JoinRequestMessageActivityType MessageActivityType = 5
)

type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyID string              `json:"party_id"`
}

type AllowedMentions struct {
	Parse         []AllowedMentionType `json:"parse"`
	Roles         []Snowflake          `json:"roles"`
	Users         []Snowflake          `json:"users"`
	IsRepliedUser bool                 `json:"replied_user"`
}

type ReactionCountDetails struct {
	Burst  ReactionType `json:"burst"`
	Normal ReactionType `json:"normal"`
}

type ReactionType int

const (
	NormalReactionType ReactionType = 0
	BurstReactionType  ReactionType = 1
)

type Reaction struct {
	Count        int                  `json:"count"`
	CountDetails ReactionCountDetails `json:"count_details"`
	IsMe         bool                 `json:"me"`
	IsMeBurst    bool                 `json:"me_burst"`
	Emoji        Emoji                `json:"emoji"`
	BurstColors  []string             `json:"burst_colors"`
}

type MessageType struct {
	Type      string `json:"type"`
	Value     int    `json:"value"`
	Deletable bool   `json:"deletable"`
}

func (m *MessageType) UnmarshalJSON(data []byte) error {
	var messageType int
	if err := json.Unmarshal(data, &messageType); err != nil {
		return err
	}

	m = GetMessageType(messageType)
	if m == nil {
		return errors.New("invalid message type")
	}
	return nil
}

func (m *MessageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

func GetMessageType(messageValue int) *MessageType {
	switch messageValue {
	case 0:
		return &MessageType{Type: "DEFAULT", Value: messageValue, Deletable: true}
	case 1:
		return &MessageType{Type: "RECIPIENT_ADD", Value: messageValue, Deletable: false}
	case 2:
		return &MessageType{Type: "RECIPIENT_REMOVE", Value: messageValue, Deletable: false}
	case 3:
		return &MessageType{Type: "CALL", Value: messageValue, Deletable: false}
	case 4:
		return &MessageType{Type: "CHANNEL_NAME_CHANGE", Value: messageValue, Deletable: false}
	case 5:
		return &MessageType{Type: "CHANNEL_ICON_CHANGE", Value: messageValue, Deletable: false}
	case 6:
		return &MessageType{Type: "CHANNEL_PINNED_MESSAGE", Value: messageValue, Deletable: true}
	case 7:
		return &MessageType{Type: "USER_JOIN", Value: messageValue, Deletable: true}
	case 8:
		return &MessageType{Type: "GUILD_BOOST", Value: messageValue, Deletable: true}
	case 9:
		return &MessageType{Type: "GUILD_BOOST_TIER_1", Value: messageValue, Deletable: true}
	case 10:
		return &MessageType{Type: "GUILD_BOOST_TIER_2", Value: messageValue, Deletable: true}
	case 11:
		return &MessageType{Type: "GUILD_BOOST_TIER_3", Value: messageValue, Deletable: true}
	case 12:
		return &MessageType{Type: "CHANNEL_FOLLOW_ADD", Value: messageValue, Deletable: true}
	case 14:
		return &MessageType{Type: "GUILD_DISCOVERY_DISQUALIFIED", Value: messageValue, Deletable: true}
	case 15:
		return &MessageType{Type: "GUILD_DISCOVERY_REQUALIFIED", Value: messageValue, Deletable: true}
	case 16:
		return &MessageType{Type: "GUILD_DISCOVERY_GRACE_PERIOD_INITIAL_WARNING", Value: messageValue, Deletable: true}
	case 17:
		return &MessageType{Type: "GUILD_DISCOVERY_GRACE_PERIOD_FINAL_WARNING", Value: messageValue, Deletable: true}
	case 18:
		return &MessageType{Type: "THREAD_CREATED", Value: messageValue, Deletable: true}
	case 19:
		return &MessageType{Type: "REPLY", Value: messageValue, Deletable: true}
	case 20:
		return &MessageType{Type: "CHAT_INPUT_COMMAND", Value: messageValue, Deletable: true}
	case 21:
		return &MessageType{Type: "THREAD_STARTER_MESSAGE", Value: messageValue, Deletable: false}
	case 22:
		return &MessageType{Type: "GUILD_INVITE_REMINDER", Value: messageValue, Deletable: true}
	case 23:
		return &MessageType{Type: "CONTEXT_MENU_COMMAND", Value: messageValue, Deletable: true}
	case 24:
		return &MessageType{Type: "AUTO_MODERATION_ACTION", Value: messageValue, Deletable: true}
	case 25:
		return &MessageType{Type: "ROLE_SUBSCRIPTION_PURCHASE", Value: messageValue, Deletable: true}
	case 26:
		return &MessageType{Type: "INTERACTION_PREMIUM_UPSELL", Value: messageValue, Deletable: true}
	case 27:
		return &MessageType{Type: "STAGE_START", Value: messageValue, Deletable: true}
	case 28:
		return &MessageType{Type: "STAGE_END", Value: messageValue, Deletable: true}
	case 29:
		return &MessageType{Type: "STAGE_SPEAKER", Value: messageValue, Deletable: true}
	case 31:
		return &MessageType{Type: "STAGE_TOPIC", Value: messageValue, Deletable: true}
	case 32:
		return &MessageType{Type: "GUILD_APPLICATION_PREMIUM_SUBSCRIPTION", Value: messageValue, Deletable: true}
	case 36:
		return &MessageType{Type: "GUILD_INCIDENT_ALERT_MODE_ENABLED", Value: messageValue, Deletable: true}
	case 37:
		return &MessageType{Type: "GUILD_INCIDENT_ALERT_MODE_DISABLED", Value: messageValue, Deletable: true}
	case 38:
		return &MessageType{Type: "GUILD_INCIDENT_REPORT_RAID", Value: messageValue, Deletable: true}
	case 39:
		return &MessageType{Type: "GUILD_INCIDENT_REPORT_FALSE_ALARM", Value: messageValue, Deletable: true}
	case 44:
		return &MessageType{Type: "PURCHASE_NOTIFICATION", Value: messageValue, Deletable: true}
	default:
		return nil
	}
}

type MessageReferenceType int

const (
	DefaultMessageReferenceType MessageReferenceType = 0
	ForwardMessageReferenceType MessageReferenceType = 1
)

type MessageReference struct {
	Type            MessageReferenceType `json:"type"`
	MessageID       *Snowflake           `json:"message_id,omitempty"`
	ChannelID       *Snowflake           `json:"channel_id,omitempty"`
	GuildID         *Snowflake           `json:"guild_id,omitempty"`
	FailIfNotExists *bool                `json:"fail_if_not_exists,omitempty"`
}

type MessageSnapshot struct {
	Message *Message `json:"message"`
}

type MessageInteractionMetadata struct {
	ID                            Snowflake                             `json:"id"`
	Type                          InteractionType                       `json:"type"`
	User                          User                                  `json:"user"`
	AuthorizingIntegrationOwners  map[ApplicationIntegrationType]string `json:"authorizing_integration_owners"`
	OriginalMessageID             *Snowflake                            `json:"original_message_id,omitempty"`
	InteractedMessageID           *Snowflake                            `json:"interacted_message_id,omitempty"`
	TriggeringInteractionMetadata *MessageInteractionMetadata           `json:"triggering_interaction_metadata,omitempty"`
}

type MessageInteraction struct {
	ID     Snowflake       `json:"id"`
	Type   InteractionType `json:"type"`
	Name   string          `json:"name"`
	User   User            `json:"user"`
	Member *GuildMember    `json:"member,omitempty"`
}

type MessageComponentType int

const (
	ActionRowMessageComponent         MessageComponentType = 1
	ButtonMessageComponent            MessageComponentType = 2
	StringSelectMessageComponent      MessageComponentType = 3
	TextInputMessageComponent         MessageComponentType = 4
	UserSelectMessageComponent        MessageComponentType = 5
	RoleSelectMessageComponent        MessageComponentType = 6
	MentionableSelectMessageComponent MessageComponentType = 7
	ChannelSelectMessageComponent     MessageComponentType = 8
)

type MessageComponent struct {
	Type       MessageComponentType
	Components []ActionRow
}

type Message struct {
	ID                   Snowflake                   `json:"id"`
	ChannelID            Snowflake                   `json:"channel_id"`
	Author               User                        `json:"author"`
	Content              string                      `json:"content"`
	Timestamp            time.Time                   `json:"timestamp"`
	EditedTimestamp      *time.Time                  `json:"edited_timestamp,omitempty"`
	TTS                  bool                        `json:"tts"`
	MentionEveryone      bool                        `json:"mention_everyone"`
	Mentions             []User                      `json:"mentions"`
	MentionRoles         []Role                      `json:"mention_roles"`
	MentionChannels      []ChannelMention            `json:"mention_channels"`
	Attachments          []Attachment                `json:"attachments"`
	Embeds               []Embed                     `json:"embeds"`
	Reactions            []Reaction                  `json:"reactions"`
	Nonce                *string                     `json:"nonce,omitempty"`
	Pinned               bool                        `json:"pinned"`
	WebhookID            *Snowflake                  `json:"webhook_id,omitempty"`
	Type                 MessageType                 `json:"type"`
	Activity             *MessageActivity            `json:"activity,omitempty"`
	Application          *Application                `json:"application,omitempty"`
	ApplicationID        *Snowflake                  `json:"application_id,omitempty"`
	Flags                Bitfield[MessageFlag]       `json:"flags,omitempty"`
	MessageReference     *MessageReference           `json:"message_reference,omitempty"`
	MessageSnapshots     []MessageSnapshot           `json:"message_snapshots"`
	ReferencedMessage    *Message                    `json:"referenced_message,omitempty"`
	InteractionMetadata  *MessageInteractionMetadata `json:"interaction_metadata,omitempty"`
	Interaction          *MessageInteraction         `json:"interaction,omitempty"`
	Thread               *Channel                    `json:"thread,omitempty"`
	Components           []MessageComponent          `json:"components"`
	StickerItems         []StickerItem               `json:"sticker_items"`
	Stickers             []Sticker                   `json:"stickers"` //DEPRECATED: remove in the future
	Position             *int                        `json:"position,omitempty"`
	RoleSubscriptionData *RoleSubscriptionData       `json:"role_subscription_data,omitempty"`
	Resolved             *ResolvedData               `json:"resolved,omitempty"`
	Poll                 *Poll                       `json:"poll,omitempty"`
	Call                 *MessageCall                `json:"call,omitempty"`
}

func (m *Message) GetReaction(emoji Emoji) *Reaction {
	for _, r := range m.Reactions {
		// first check if it has an ID and the query Emoji does as well
		if r.Emoji.ID != nil && emoji.ID != nil {
			if r.Emoji.ID.Equals(*emoji.ID) {
				return &r
			}
		}

		// then check if it has a name and the query Emoji does as well
		if r.Emoji.Name != nil && emoji.Name != nil {
			if *r.Emoji.Name == *emoji.Name {
				return &r
			}
		}
	}
	return nil
}

func (m *Message) UpdateReactions(reaction Reaction) error {
	if m.GetReaction(reaction.Emoji) == nil {
		m.Reactions = append(m.Reactions, reaction)
		return nil
	}

	for i, r := range m.Reactions {
		if r.Emoji.ID != nil && reaction.Emoji.ID != nil {
			if r.Emoji.ID.Equals(*reaction.Emoji.ID) {
				m.Reactions[i] = reaction
				return nil
			}
		}

		if r.Emoji.Name != nil && reaction.Emoji.Name != nil {
			if *r.Emoji.Name == *reaction.Emoji.Name {
				m.Reactions[i] = reaction
				return nil
			}
		}
	}

	return errors.New("error updating reaction")
}

func (m *Message) DeleteReaction(emoji Emoji) {
	for i, r := range m.Reactions {
		if r.Emoji.ID != nil && emoji.ID != nil {
			if r.Emoji.ID.Equals(*emoji.ID) {
				m.Reactions = append(m.Reactions[:i], m.Reactions[i+1:]...)
				return
			}
		}

		if r.Emoji.Name != nil && emoji.Name != nil {
			if *r.Emoji.Name == *emoji.Name {
				m.Reactions = append(m.Reactions[:i], m.Reactions[i+1:]...)
				return
			}
		}
	}
}

type MessageCall struct {
	Participants   []Snowflake `json:"participants"`
	EndedTimestamp *time.Time  `json:"ended_timestamp,omitempty"`
}

type RoleSubscriptionData struct {
	RoleSubscriptionListingID Snowflake `json:"role_subscription_listing_id"`
	TierName                  string    `json:"tier_name"`
	TotalMonthsSubscribed     int       `json:"total_months_subscribed"`
	IsRenewal                 bool      `json:"is_renewal"`
}

type ActionRow struct {
	Type       MessageComponentType `json:"type"`
	Components []MessageComponent   `json:"components"`
}

type ChannelMention struct {
	ID      Snowflake   `json:"id"`
	GuildID Snowflake   `json:"guild_id"`
	Type    ChannelType `json:"type"`
	Name    string      `json:"name"`
}

type Attachment struct {
	ID              Snowflake                 `json:"id"`
	FileName        string                    `json:"filename"`
	Title           *string                   `json:"title,omitempty"`
	Description     *string                   `json:"description,omitempty"`
	ContentType     *string                   `json:"content_type,omitempty"`
	Size            int                       `json:"size"`
	URL             string                    `json:"url"`
	ProxyURL        string                    `json:"proxy_url"`
	Height          *int                      `json:"height,omitempty"`
	Width           *int                      `json:"width,omitempty"`
	Ephemeral       *bool                     `json:"ephemeral,omitempty"`
	DurationSeconds *float64                  `json:"duration_secs,omitempty"`
	Waveform        *string                   `json:"waveform,omitempty"`
	Flags           *Bitfield[AttachmentFlag] `json:"flags,omitempty"`
}

type Embed struct {
	Title       *string         `json:"title,omitempty"`
	Type        *EmbedType      `json:"type,omitempty"`
	Description *string         `json:"description,omitempty"`
	URL         *string         `json:"url,omitempty"`
	Timestamp   *time.Time      `json:"timestamp,omitempty"`
	Color       *int            `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      *[]EmbedField   `json:"fields,omitempty"`
}

type EmbedType string

const (
	Rich    EmbedType = "rich"
	Image   EmbedType = "image"
	Video   EmbedType = "video"
	Gifv    EmbedType = "gifv"
	Article EmbedType = "article"
	Link    EmbedType = "link"
)

type EmbedThumbnail struct {
	URL      string  `json:"url"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
}

type EmbedVideo struct {
	URl      *string `json:"url,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
}

type EmbedImage struct {
	URL      string  `json:"url"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         string  `json:"name"`
	URL          *string `json:"url,omitempty"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

type EmbedFooter struct {
	Text         string  `json:"text"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

type EmbedField struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	IsInline *bool  `json:"inline,omitempty"`
}
