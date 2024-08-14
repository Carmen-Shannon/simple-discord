package structs

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
	Flags                *int                 `json:"flags,omitempty"`
	PremiumType          *int                 `json:"premium_type,omitempty"`
	PublicFlags          *int                 `json:"public_flags,omitempty"`
	AvatarDecorationData AvatarDecorationData `json:"avatar_decoration_data"`
}

type AvatarDecorationData struct {
	Asset string    `json:"asset"`
	SKU   Snowflake `json:"sku"`
}
