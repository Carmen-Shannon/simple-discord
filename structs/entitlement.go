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
	ID            Snowflake
	SKUID         Snowflake
	ApplicationID Snowflake
	UserID        *Snowflake
	Type          EntitlementType
	Deleted       bool
	StartsAt      *time.Time
	EndsAt        *time.Time
	GuildID       *Snowflake
	Consumed      *bool
}
