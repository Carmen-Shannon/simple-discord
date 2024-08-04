package structs

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

type SystemChannelFlag int

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

type WelcomeScreenChannel struct {
	ChannelID   Snowflake
	Description string
	EmojiID     *Snowflake
	EmojiName   *string
}

type WelcomeScreen struct {
	Description     *string
	WelcomeChannels []WelcomeScreenChannel
}

type Guild struct {
	ID                          Snowflake
	Name                        string
	Icon                        *string
	IconHash                    *string
	Splash                      *string
	DiscoverySplash             *string
	Owner                       *bool
	OwnerID                     Snowflake
	Permissions                 *string
	Region                      *VoiceRegion
	AFKChannelID                *Snowflake
	AFKTimeout                  int
	IsWidgetEnabled             *bool
	WidgetChannelID             *Snowflake
	VerificationLevel           VerificationLevel
	DefaultMessageNotifications DefaultMessageNotificationLevel
	ExplicitContentFilter       ExplicitContentFilterLevel
	Roles                       []Role
	Emojis                      []Emoji
	Features                    []GuildFeature
	MFALevel                    MFALevel
	ApplicationID               *Snowflake
	SystemChannelID             *Snowflake
	SystemChannelFlags          SystemChannelFlag
	RulesChannelID              *Snowflake
	MaxPresences                *int
	MaxMembers                  *int
	VanityURLCode               *string
	Description                 *string
	Banner                      *string
	PremiumTier                 PremiumTier
	PremiumSubscriptionCount    *int
	PreferredLocale             string
	PublicUpdatesChannelID      *Snowflake
	MaxVideoChannelUsers        *int
	MaxStageVideoChannelUsers   *int
	ApproximateMemberCount      *int
	ApproximatePresenceCount    *int
	WelcomeScreen               *WelcomeScreen
	NSFWLevel                   NSFWLevel
	Stickers                    []Sticker
	IsPremiumProgressBarEnabled bool
	SafetyAlertsChannelID       *Snowflake
}
