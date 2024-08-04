package structs

type User struct {
	ID                   Snowflake
	Username             string
	Discriminator        string
	GlobalName           *string
	Avatar               *string
	IsBot                *bool
	IsSystem             *bool
	IsMFA                *bool
	Banner               *string
	AccentColor          *int
	Locale               *string
	IsVerified           *bool
	Email                *string
	Flags                *int
	PremiumType          *int
	PublicFlags          *int
	AvatarDecorationData AvatarDecorationData
}

type AvatarDecorationData struct {
	Asset string
	SKU   Snowflake
}
