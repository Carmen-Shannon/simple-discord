package structs

import "time"

const None = 0

type VerificationLevel int

const (
	LowVerification      VerificationLevel = 1
	MediumVerification   VerificationLevel = 2
	HighVerification     VerificationLevel = 3
	VeryHighVerification VerificationLevel = 4
)

type DefaultMessageNotificationLevel int

const (
	AllMessagesLevel  DefaultMessageNotificationLevel = 0
	OnlyMentionsLevel DefaultMessageNotificationLevel = 1
)

type ExplicitContentFilterLevel int

const (
	DisabledLevel            ExplicitContentFilterLevel = 0
	MembersWithoutRolesLevel ExplicitContentFilterLevel = 1
	AllMembersLevel          ExplicitContentFilterLevel = 2
)

type GuildFeature string

const (
	AnimatedBanner                  GuildFeature = "ANIMATED_BANNER"
	AnimatedIcon                    GuildFeature = "ANIMATED_ICON"
	ApplicationCommandPermissionsV2 GuildFeature = "APPLICATION_COMMAND_PERMISSIONS_V2"
	AutoModeration                  GuildFeature = "AUTO_MODERATION"
	Banner                          GuildFeature = "BANNER"
	Community                       GuildFeature = "COMMUNITY"
	CreatorMonetizableProvisional   GuildFeature = "CREATOR_MONETIZABLE_PROVISIONAL"
	CreatorStorePage                GuildFeature = "CREATOR_STORE_PAGE"
	DeveloperSupportServer          GuildFeature = "DEVELOPER_SUPPORT_SERVER"
	Discoverable                    GuildFeature = "DISCOVERABLE"
	Featurable                      GuildFeature = "FEATURABLE"
	InvitesDisabled                 GuildFeature = "INVITES_DISABLED"
	InviteSplash                    GuildFeature = "INVITE_SPLASH"
	MemberVerificationGateEnabled   GuildFeature = "MEMBER_VERIFICATION_GATE_ENABLED"
	MoreStickers                    GuildFeature = "MORE_STICKERS"
	News                            GuildFeature = "NEWS"
	Partnered                       GuildFeature = "PARTNERED"
	PreviewEnabled                  GuildFeature = "PREVIEW_ENABLED"
	RaidAlertsDisabled              GuildFeature = "RAID_ALERTS_DISABLED"
	RoleIcons                       GuildFeature = "ROLE_ICONS"
	RoleSubscriptionsAvailable      GuildFeature = "ROLE_SUBSCRIPTIONS_AVAILABLE_FOR_PURCHASE"
	RoleSubscriptionsEnabled        GuildFeature = "ROLE_SUBSCRIPTIONS_ENABLED"
	TicketedEventsEnabled           GuildFeature = "TICKETED_EVENTS_ENABLED"
	VanityURL                       GuildFeature = "VANITY_URL"
	Verified                        GuildFeature = "VERIFIED"
	VIPRegions                      GuildFeature = "VIP_REGIONS"
	WelcomeScreenEnabled            GuildFeature = "WELCOME_SCREEN_ENABLED"
)

type MFALevel int

const (
	ElevatedMFA MFALevel = 1
)

type NSFWLevel int

const (
	DefaultNSFW       NSFWLevel = 0
	ExplicitNSFW      NSFWLevel = 1
	SafeNSFW          NSFWLevel = 2
	AgeRestrictedNSFW NSFWLevel = 3
)

type SystemChannelFlag int64

const (
	SurpressJoinNotificationsFlag                           SystemChannelFlag = 1 << 0
	SurpressPremiumSubscriptionsFlag                        SystemChannelFlag = 1 << 1
	SurpressGuildReminderNotificationsFlag                  SystemChannelFlag = 1 << 2
	SurpressJoinNotificationRepliesFlag                     SystemChannelFlag = 1 << 3
	SurpressRoleSubscriptionPurchaseNotificationsFlag       SystemChannelFlag = 1 << 4
	SurpressRoleSubscriptionPurchaseNotificationRepliesFlag SystemChannelFlag = 1 << 5
)

type PremiumTier int

const (
	Tier1 PremiumTier = 1
	Tier2 PremiumTier = 2
	Tier3 PremiumTier = 3
)

type WidgetStyle string

const (
	WidgetStyleShield  WidgetStyle = "shield"
	WidgetStyleBanner1 WidgetStyle = "banner1"
	WidgetStyleBanner2 WidgetStyle = "banner2"
	WidgetStyleBanner3 WidgetStyle = "banner3"
	WidgetStyleBanner4 WidgetStyle = "banner4"
)

type OnboardingMode int

const (
	OnboardingModeDefault  OnboardingMode = 0
	OnboardingModeAdvanced OnboardingMode = 1
)

type PromptType int

const (
	PromptTypeMultipleChoice PromptType = 0
	PromptTypeDropdown       PromptType = 1
)

type WelcomeScreenChannel struct {
	ChannelID   Snowflake  `json:"channel_id"`
	Description string     `json:"description"`
	EmojiID     *Snowflake `json:"emoji_id,omitempty"`
	EmojiName   *string    `json:"emoji_name,omitempty"`
}

type WelcomeScreen struct {
	Description     *string                `json:"description,omitempty"`
	WelcomeChannels []WelcomeScreenChannel `json:"welcome_channels"`
}

type Guild struct {
	ID                          Snowflake                       `json:"id"`
	Name                        string                          `json:"name"`
	Icon                        *string                         `json:"icon,omitempty"`
	IconHash                    *string                         `json:"icon_hash,omitempty"`
	Splash                      *string                         `json:"splash,omitempty"`
	DiscoverySplash             *string                         `json:"discovery_splash,omitempty"`
	Owner                       *bool                           `json:"owner,omitempty"`
	OwnerID                     Snowflake                       `json:"owner_id"`
	Permissions                 *string                         `json:"permissions,omitempty"`
	Region                      *string                         `json:"region,omitempty"` //DEPRECATED
	AFKChannelID                *Snowflake                      `json:"afk_channel_id,omitempty"`
	AFKTimeout                  int                             `json:"afk_timeout"`
	IsWidgetEnabled             *bool                           `json:"is_widget_enabled,omitempty"`
	WidgetChannelID             *Snowflake                      `json:"widget_channel_id,omitempty"`
	VerificationLevel           VerificationLevel               `json:"verification_level"`
	DefaultMessageNotifications DefaultMessageNotificationLevel `json:"default_message_notifications"`
	ExplicitContentFilter       ExplicitContentFilterLevel      `json:"explicit_content_filter"`
	Roles                       []Role                          `json:"roles"`
	Emojis                      []Emoji                         `json:"emojis"`
	Features                    []GuildFeature                  `json:"features"`
	MFALevel                    MFALevel                        `json:"mfa_level"`
	ApplicationID               *Snowflake                      `json:"application_id,omitempty"`
	SystemChannelID             *Snowflake                      `json:"system_channel_id,omitempty"`
	SystemChannelFlags          Bitfield[SystemChannelFlag]     `json:"system_channel_flags"`
	RulesChannelID              *Snowflake                      `json:"rules_channel_id,omitempty"`
	MaxPresences                *int                            `json:"max_presences,omitempty"`
	MaxMembers                  *int                            `json:"max_members,omitempty"`
	VanityURLCode               *string                         `json:"vanity_url_code,omitempty"`
	Description                 *string                         `json:"description,omitempty"`
	Banner                      *string                         `json:"banner,omitempty"`
	PremiumTier                 PremiumTier                     `json:"premium_tier"`
	PremiumSubscriptionCount    *int                            `json:"premium_subscription_count,omitempty"`
	PreferredLocale             string                          `json:"preferred_locale"`
	PublicUpdatesChannelID      *Snowflake                      `json:"public_updates_channel_id,omitempty"`
	MaxVideoChannelUsers        *int                            `json:"max_video_channel_users,omitempty"`
	MaxStageVideoChannelUsers   *int                            `json:"max_stage_video_channel_users,omitempty"`
	ApproximateMemberCount      *int                            `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    *int                            `json:"approximate_presence_count,omitempty"`
	WelcomeScreen               *WelcomeScreen                  `json:"welcome_screen,omitempty"`
	NSFWLevel                   NSFWLevel                       `json:"nsfw_level"`
	Stickers                    []Sticker                       `json:"stickers"`
	IsPremiumProgressBarEnabled bool                            `json:"is_premium_progress_bar_enabled"`
	SafetyAlertsChannelID       *Snowflake                      `json:"safety_alerts_channel_id,omitempty"`
}

type GuildPreview struct {
	ID                       Snowflake      `json:"id"`
	Name                     string         `json:"name"`
	Icon                     *string        `json:"icon,omitempty"`
	Splash                   *string        `json:"splash,omitempty"`
	DiscoverySplash          *string        `json:"discovery_splash,omitempty"`
	Emojis                   []Emoji        `json:"emojis"`
	Features                 []GuildFeature `json:"features"`
	ApproximateMemberCount   int            `json:"approximate_member_count"`
	ApproximatePresenceCount int            `json:"approximate_presence_count"`
	Description              *string        `json:"description,omitempty"`
	Stickers                 []Sticker      `json:"stickers"`
}

type UnavailableGuild struct {
	ID          Snowflake `json:"id"`
	Unavailable bool      `json:"unavailable"`
}

type GuildApplicationCommandPermissions struct {
	ID            Snowflake                       `json:"id"`
	ApplicationID Snowflake                       `json:"application_id"`
	GuildID       Snowflake                       `json:"guild_id"`
	Permissions   []ApplicationCommandPermissions `json:"permissions"`
}

type GuildIntegrationType string

const (
	Twitch            GuildIntegrationType = "twitch"
	Youtube           GuildIntegrationType = "youtube"
	Discord           GuildIntegrationType = "discord"
	GuildSubscription GuildIntegrationType = "guild_subscription"
)

type IntegrationExpireBehavior int

const (
	RemoveRoleBehavior IntegrationExpireBehavior = 0
	KickBehavior       IntegrationExpireBehavior = 1
)

type GuildIntegration struct {
	ID                Snowflake                  `json:"id"`
	Name              string                     `json:"name"`
	Type              GuildIntegrationType       `json:"type"`
	IsEnabled         bool                       `json:"enabled"`
	IsSyncing         *bool                      `json:"syncing,omitempty"`
	RoleID            *Snowflake                 `json:"role_id,omitempty"`
	EnabledEmoticons  *bool                      `json:"enabled_emoticons,omitempty"`
	ExpireBehavior    *IntegrationExpireBehavior `json:"expire_behavior,omitempty"`
	ExpireGracePeriod *int                       `json:"expire_grace_period,omitempty"`
	User              *User                      `json:"user,omitempty"`
	Account           IntegrationAccount         `json:"account"`
	SyncedAt          *time.Time                 `json:"synced_at,omitempty"`
	SubscriberCount   *int                       `json:"subscriber_count,omitempty"`
	IsRevoked         *bool                      `json:"revoked,omitempty"`
	Application       *IntegrationApplication    `json:"application,omitempty"`
	Scopes            []OAuth2Scope              `json:"scopes"`
}

type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type IntegrationApplication struct {
	ID          Snowflake `json:"id"`
	Name        string    `json:"name"`
	Icon        *string   `json:"icon,omitempty"`
	Description string    `json:"description"`
	Bot         *User     `json:"bot,omitempty"`
}

type Ban struct {
	Reason string `json:"reason,omitempty"`
	User   User   `json:"user"`
}

type GuildWidgetSettings struct {
	Enabled   bool       `json:"enabled"`
	ChannelID *Snowflake `json:"channel_id,omitempty"`
}

type GuildWidget struct {
	ID            Snowflake `json:"id"`
	Name          string    `json:"name"`
	InstantInvite *string   `json:"instant_invite"`
	Channels      []Channel `json:"channels"`
	Members       []User    `json:"members"`
	PresenceCount int       `json:"presence_count"`
}

type OnboardingPrompt struct {
	ID           Snowflake      `json:"id"`
	Type         PromptType     `json:"type"`
	Options      []PromptOption `json:"options"`
	Title        string         `json:"title"`
	SingleSelect bool           `json:"single_select"`
	Required     bool           `json:"required"`
	InOnboarding bool           `json:"in_onboarding"`
}

type PromptOption struct {
	ID            Snowflake   `json:"id"`
	ChannelIDs    []Snowflake `json:"channel_ids"`
	RoleIDs       []Snowflake `json:"role_ids"`
	Emoji         *Emoji      `json:"emoji,omitempty"`
	EmojiID       *Snowflake  `json:"emoji_id,omitempty"`
	EmojiName     *string     `json:"emoji_name,omitempty"`
	EmojiAnimated *bool       `json:"emoji_animated,omitempty"`
	Title         string      `json:"title"`
	Description   string      `json:"description,omitempty"`
}

type GuildOnboarding struct {
	GuildID           Snowflake          `json:"guild_id"`
	Prompts           []OnboardingPrompt `json:"prompts"`
	DefaultChannelIDs []Snowflake        `json:"default_channel_ids"`
	Enabled           bool               `json:"enabled"`
	Mode              OnboardingMode     `json:"mode"`
}
