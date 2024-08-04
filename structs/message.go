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
	Type    MessageActivityType
	PartyID string
}

type AllowedMentions struct {
	Parse         []AllowedMentionType
	Roles         []Snowflake
	Users         []Snowflake
	IsRepliedUser bool
}

type ReactionCountDetails struct {
	Burst  int
	Normal int
}

type Reaction struct {
	Count        int
	CountDetails ReactionCountDetails
	IsMe         bool
	IsMeBurst    bool
	Emoji        Emoji
	BurstColors  []string
}

type MessageType struct {
	Type      string
	Value     int
	Deletable bool
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
	Type            MessageReferenceType
	MessageID       *Snowflake
	ChannelID       *Snowflake
	GuildID         *Snowflake
	FailIfNotExists *bool
}

type MessageSnapshot struct {
	Message *Message
}

type MessageInteractionMetadata struct {
	ID                            Snowflake
	Type                          InteractionType
	User                          User
	AuthorizingIntegrationOwners  map[ApplicationIntegrationType]string
	OriginalMessageID             *Snowflake
	InteractedMessageID           *Snowflake
	TriggeringInteractionMetadata *MessageInteractionMetadata
}

type MessageInteraction struct {
	ID     Snowflake
	Type   InteractionType
	Name   string
	User   User
	Member *GuildMember
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
	ID                   Snowflake
	ChannelID            Snowflake
	Author               User
	Content              string
	Timestamp            time.Time
	EditedTimestamp      *time.Time
	TTS                  bool
	MentionEveryone      bool
	Mentions             []User
	MentionRoles         []Role
	MentionChannels      []ChannelMention
	Attachments          []Attachment
	Embeds               []Embed
	Reactions            []Reaction
	Nonce                *string
	Pinned               bool
	WebhookID            *Snowflake
	Type                 MessageType
	Activity             MessageActivity
	Application          *Application
	ApplicationID        *Snowflake
	Flags                *MessageFlag
	MessageReference     *MessageReference
	MessageSnapshots     []MessageSnapshot
	ReferencedMessage    *Message
	InteractionMetadata  MessageInteractionMetadata
	Interaction          MessageInteraction
	Thread               *Channel
	Components           []MessageComponent
	StickerItems         []StickerItem
	Stickers             []Sticker
	Position             *int
	RoleSubscriptionData RoleSubscriptionData
	Resolved             ResolvedData
	Poll                 Poll
	Call                 MessageCall
}

type MessageCall struct {
	Participants   []Snowflake
	EndedTimestamp *time.Time
}

type RoleSubscriptionData struct {
	RoleSubscriptionListingID Snowflake
	TierName                  string
	TotalMonthsSubscribed     int
	IsRenewal                 bool
}

type ActionRow struct {
	Type       MessageComponentType
	Components []MessageComponent
}

type ChannelMention struct {
	ID      Snowflake
	GuildID Snowflake
	Type    ChannelType
	Name    string
}

type Attachment struct {
	ID              Snowflake
	FileName        string
	Title           *string
	Description     *string
	ContentType     *string
	Size            int
	URL             string
	ProxyURL        string
	Height          *int
	Width           *int
	Ephemeral       *bool
	DurationSeconds *float64
	Waveform        *string
	Flags           *AttachmentFlag
}

type Embed struct {
	Title       *string
	Type        *EmbedType
	Description *string
	URL         *string
	Timestamp   *time.Time
	Color       *int
	Footer      *EmbedFooter
	Image       *EmbedImage
	Thumbnail   *EmbedThumbnail
	Video       *EmbedVideo
	Provider    *EmbedProvider
	Author      *EmbedAuthor
	Fields      *[]EmbedField
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
	URL      string
	ProxyURL *string
	Height   *int
	Width    *int
}

type EmbedVideo struct {
	URl      *string
	ProxyURL *string
	Height   *int
	Width    *int
}

type EmbedImage struct {
	URL      string
	ProxyURL *string
	Height   *int
	Width    *int
}

type EmbedProvider struct {
	Name *string
	URL  *string
}

type EmbedAuthor struct {
	Name         string
	URL          *string
	IconURL      *string
	ProxyIconURL *string
}

type EmbedFooter struct {
	Text         string
	IconURL      *string
	ProxyIconURL *string
}

type EmbedField struct {
	Name     string
	Value    string
	IsInline *bool
}
