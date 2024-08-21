package structs

type RoleFlag int64

const (
	InPrompt RoleFlag = 1 << 0
)

type Role struct {
	ID            Snowflake          `json:"id"`
	Name          string             `json:"name"`
	Color         int                `json:"color"`
	IsHoist       bool               `json:"hoist"`
	Icon          *string            `json:"icon,omitempty"`
	UnicodeEmoji  *string            `json:"unicode_emoji,omitempty"`
	Position      int                `json:"position"`
	Permissions   string             `json:"permissions"`
	IsManaged     bool               `json:"managed"`
	IsMentionable bool               `json:"mentionable"`
	Tags          RoleTags           `json:"tags"`
	Flags         Bitfield[RoleFlag] `json:"flags"`
}

type RoleTags struct {
	BotID                 *Snowflake `json:"bot_id,omitempty"`
	IntegrationID         *Snowflake `json:"integration_id,omitempty"`
	PremiumSubscriber     bool       `json:"premium_subscriber"`
	SubscriptionListingID *Snowflake `json:"subscription_listing_id,omitempty"`
	AvailableForPurchase  bool       `json:"available_for_purchase"`
	GuildConnections      bool       `json:"guild_connections"`
}
