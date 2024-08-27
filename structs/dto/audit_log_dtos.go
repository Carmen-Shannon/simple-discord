package dto

import "github.com/Carmen-Shannon/simple-discord/structs"

type GetGuildAuditLogDto struct {
	GuildID    structs.Snowflake  `json:"guild_id"`
	UserID     *structs.Snowflake `json:"user_id,omitempty"`
	ActionType *int               `json:"action_type,omitempty"`
	Before     *structs.Snowflake `json:"before,omitempty"`
	After      *structs.Snowflake `json:"after,omitempty"`
	Limit      *int               `json:"limit,omitempty"`
}
