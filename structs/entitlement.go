package structs

import "time"

type EntitlementType int

const (
	PurchaseEntitlement                EntitlementType = 1
	PremiumSubscriptionEntitlement     EntitlementType = 2
	DeveloperGiftEntitlement           EntitlementType = 3
	TestModePurchaseEntitlement        EntitlementType = 4
	FreePurchaseEntitlement            EntitlementType = 5
	UserGiftEntitlement                EntitlementType = 6
	PremiumPurchaseEntitlement         EntitlementType = 7
	ApplicationSubscriptionEntitlement EntitlementType = 8
)

type Entitlement struct {
	ID            Snowflake     `json:"id"`
	SKUID         Snowflake     `json:"skuid"`
	ApplicationID Snowflake     `json:"application_id"`
	UserID        *Snowflake    `json:"user_id,omitempty"`
	Type          EntitlementType `json:"type"`
	Deleted       bool          `json:"deleted"`
	StartsAt      *time.Time    `json:"starts_at,omitempty"`
	EndsAt        *time.Time    `json:"ends_at,omitempty"`
	GuildID       *Snowflake    `json:"guild_id,omitempty"`
	Consumed      *bool         `json:"consumed,omitempty"`
}
