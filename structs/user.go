package structs

type UserFlag int64

const (
	StaffUserFlag                 UserFlag = 1 << 0
	PartnerUserFlag               UserFlag = 1 << 1
	HypesquadUserFlag             UserFlag = 1 << 2
	BugHunterLevel1UserFlag       UserFlag = 1 << 3
	HypesquadOnlineHouse1UserFlag UserFlag = 1 << 6
	HypesquadOnlineHouse2UserFlag UserFlag = 1 << 7
	HypesquadOnlineHouse3UserFlag UserFlag = 1 << 8
	PremiumEarlySupporterUserFlag UserFlag = 1 << 9
	TeamPseudoUserUserFlag        UserFlag = 1 << 10
	BugHunterLevel2UserFlag       UserFlag = 1 << 14
	VerifiedBotUserFlag           UserFlag = 1 << 16
	VerifiedDeveloperUserFlag     UserFlag = 1 << 17
	CertifiedModeratorUserFlag    UserFlag = 1 << 18
	BotHttpInteractionsUserFlag   UserFlag = 1 << 19
	ActiveDeveloperUserFlag       UserFlag = 1 << 20
)

type User struct {
	ID                   Snowflake            `json:"id"`
	Username             string               `json:"username"`
	Discriminator        string               `json:"discriminator"`
	GlobalName           *string              `json:"global_name,omitempty"`
	Avatar               *string              `json:"avatar,omitempty"`
	IsBot                *bool                `json:"is_bot,omitempty"`
	IsSystem             *bool                `json:"is_system,omitempty"`
	IsMFA                *bool                `json:"is_mfa,omitempty"`
	Banner               *string              `json:"banner,omitempty"`
	AccentColor          *int                 `json:"accent_color,omitempty"`
	Locale               *string              `json:"locale,omitempty"`
	IsVerified           *bool                `json:"is_verified,omitempty"`
	Email                *string              `json:"email,omitempty"`
	Flags                Bitfield[UserFlag]   `json:"flags,omitempty"`
	PremiumType          *int                 `json:"premium_type,omitempty"`
	PublicFlags          *int                 `json:"public_flags,omitempty"`
	AvatarDecorationData AvatarDecorationData `json:"avatar_decoration_data"`
}

type AvatarDecorationData struct {
	Asset string    `json:"asset"`
	SKU   Snowflake `json:"sku"`
}
