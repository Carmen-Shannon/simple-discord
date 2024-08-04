package structs

type Role struct {
	ID            Snowflake
	Name          string
	Color         int
	IsHoist       bool
	Icon          *string
	UnicodeEmoji  *string
	Position      int
	Permissions   string
	IsManaged     bool
	IsMentionable bool
	Tags          RoleTags
	Flags         int
}

type RoleTags struct {
	BotID                 *Snowflake
	Integrationid         *Snowflake
	PremiumSubscriber     bool
	SubscriptionListingID *Snowflake
	AvailableForPurchase  bool
	GuildConnections      bool
}
