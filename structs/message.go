package structs

import "time"

type AttachmentFlag int

const (
	IsRemix AttachmentFlag = 1 << 2
)

type AllowedMentionType string

const (
	RoleMentionsType     AllowedMentionType = "roles"
	UserMentionsType     AllowedMentionType = "users"
	EveryoneMentionsType AllowedMentionType = "everyone"
)

type MessageFlag int

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
	Burst  int `json:"burst"`
	Normal int `json:"normal"`
}

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

func GetMessageType(messageType string) *MessageType {
	switch messageType {
	case "DEFAULT":
		return &MessageType{Type: messageType, Value: 0, Deletable: true}
	case "RECIPIENT_ADD":
		return &MessageType{Type: messageType, Value: 1, Deletable: false}
	case "RECIPIENT_REMOVE":
		return &MessageType{Type: messageType, Value: 2, Deletable: false}
	case "CALL":
		return &MessageType{Type: messageType, Value: 3, Deletable: false}
	case "CHANNEL_NAME_CHANGE":
		return &MessageType{Type: messageType, Value: 4, Deletable: false}
	case "CHANNEL_ICON_CHANGE":
		return &MessageType{Type: messageType, Value: 5, Deletable: false}
	case "CHANNEL_PINNED_MESSAGE":
		return &MessageType{Type: messageType, Value: 6, Deletable: true}
	case "USER_JOIN":
		return &MessageType{Type: messageType, Value: 7, Deletable: true}
	case "GUILD_BOOST":
		return &MessageType{Type: messageType, Value: 8, Deletable: true}
	case "GUILD_BOOST_TIER_1":
		return &MessageType{Type: messageType, Value: 9, Deletable: true}
	case "GUILD_BOOST_TIER_2":
		return &MessageType{Type: messageType, Value: 10, Deletable: true}
	case "GUILD_BOOST_TIER_3":
		return &MessageType{Type: messageType, Value: 11, Deletable: true}
	case "CHANNEL_FOLLOW_ADD":
		return &MessageType{Type: messageType, Value: 12, Deletable: true}
	case "GUILD_DISCOVERY_DISQUALIFIED":
		return &MessageType{Type: messageType, Value: 14, Deletable: true}
	case "GUILD_DISCOVERY_REQUALIFIED":
		return &MessageType{Type: messageType, Value: 15, Deletable: true}
	case "GUILD_DISCOVERY_GRACE_PERIOD_INITIAL_WARNING":
		return &MessageType{Type: messageType, Value: 16, Deletable: true}
	case "GUILD_DISCOVERY_GRACE_PERIOD_FINAL_WARNING":
		return &MessageType{Type: messageType, Value: 17, Deletable: true}
	case "THREAD_CREATED":
		return &MessageType{Type: messageType, Value: 18, Deletable: true}
	case "REPLY":
		return &MessageType{Type: messageType, Value: 19, Deletable: true}
	case "CHAT_INPUT_COMMAND":
		return &MessageType{Type: messageType, Value: 20, Deletable: true}
	case "THREAD_STARTER_MESSAGE":
		return &MessageType{Type: messageType, Value: 21, Deletable: false}
	case "GUILD_INVITE_REMINDER":
		return &MessageType{Type: messageType, Value: 22, Deletable: true}
	case "CONTEXT_MENU_COMMAND":
		return &MessageType{Type: messageType, Value: 23, Deletable: true}
	case "AUTO_MODERATION_ACTION":
		return &MessageType{Type: messageType, Value: 24, Deletable: true}
	case "ROLE_SUBSCRIPTION_PURCHASE":
		return &MessageType{Type: messageType, Value: 25, Deletable: true}
	case "INTERACTION_PREMIUM_UPSELL":
		return &MessageType{Type: messageType, Value: 26, Deletable: true}
	case "STAGE_START":
		return &MessageType{Type: messageType, Value: 27, Deletable: true}
	case "STAGE_END":
		return &MessageType{Type: messageType, Value: 28, Deletable: true}
	case "STAGE_SPEAKER":
		return &MessageType{Type: messageType, Value: 29, Deletable: true}
	case "STAGE_TOPIC":
		return &MessageType{Type: messageType, Value: 31, Deletable: true}
	case "GUILD_APPLICATION_PREMIUM_SUBSCRIPTION":
		return &MessageType{Type: messageType, Value: 32, Deletable: true}
	case "GUILD_INCIDENT_ALERT_MODE_ENABLED":
		return &MessageType{Type: messageType, Value: 36, Deletable: true}
	case "GUILD_INCIDENT_ALERT_MODE_DISABLED":
		return &MessageType{Type: messageType, Value: 37, Deletable: true}
	case "GUILD_INCIDENT_REPORT_RAID":
		return &MessageType{Type: messageType, Value: 38, Deletable: true}
	case "GUILD_INCIDENT_REPORT_FALSE_ALARM":
		return &MessageType{Type: messageType, Value: 39, Deletable: true}
	case "PURCHASE_NOTIFICATION":
		return &MessageType{Type: messageType, Value: 44, Deletable: true}
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
	Flags                *MessageFlag                `json:"flags,omitempty"`
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
	ID              Snowflake       `json:"id"`
	FileName        string          `json:"file_name"`
	Title           *string         `json:"title,omitempty"`
	Description     *string         `json:"description,omitempty"`
	ContentType     *string         `json:"content_type,omitempty"`
	Size            int             `json:"size"`
	URL             string          `json:"url"`
	ProxyURL        string          `json:"proxy_url"`
	Height          *int            `json:"height,omitempty"`
	Width           *int            `json:"width,omitempty"`
	Ephemeral       *bool           `json:"ephemeral,omitempty"`
	DurationSeconds *float64        `json:"duration_seconds,omitempty"`
	Waveform        *string         `json:"waveform,omitempty"`
	Flags           *AttachmentFlag `json:"flags,omitempty"`
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
