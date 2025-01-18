package dto

import "github.com/Carmen-Shannon/simple-discord/structs"

type ListEntitlementsDto struct {
	ApplicationID  structs.Snowflake   `json:"-"`
	UserID         *structs.Snowflake  `json:"user_id,omitempty"`
	SkuIDs         []structs.Snowflake `json:"sku_ids,omitempty"`
	Before         *structs.Snowflake  `json:"before,omitempty"`
	After          *structs.Snowflake  `json:"after,omitempty"`
	Limit          *int                `json:"limit,omitempty"`
	GuildID        *structs.Snowflake  `json:"guild_id,omitempty"`
	ExcludeEnded   *bool               `json:"exclude_ended,omitempty"`
	ExcludeDeleted *bool               `json:"exclude_deleted,omitempty"`
}

type GetEntitlementDto struct {
	ApplicationID structs.Snowflake `json:"-"`
	EntitlementID structs.Snowflake `json:"-"`
}

type CreateTestEntitlementDto struct {
	ApplicationID structs.Snowflake  `json:"-"`
	SkuID         *structs.Snowflake `json:"sku_id,omitempty"`
	OwnerID       *structs.Snowflake `json:"owner_id,omitempty"`
	OwnerType     *int               `json:"owner_type,omitempty"`
}
