package dto

import "github.com/Carmen-Shannon/simple-discord/structs"

type ListGuildEmojisDto struct {
	GuildID structs.Snowflake `json:"-"`
}

type GetGuildEmojiDto struct {
	GuildID structs.Snowflake `json:"-"`
	EmojiID structs.Snowflake `json:"-"`
}

type CreateGuildEmojiDto struct {
	GuildID structs.Snowflake   `json:"-"`
	Name    string              `json:"name"`
	Image   string              `json:"image"`
	Roles   []structs.Snowflake `json:"roles"`
}

type ModifyGuildEmojiDto struct {
	GuildID structs.Snowflake   `json:"-"`
	EmojiID structs.Snowflake   `json:"-"`
	Name    string              `json:"name"`
	Roles   []structs.Snowflake `json:"roles,omitempty"`
}

type DeleteGuildEmojiDto struct {
	GuildID structs.Snowflake `json:"-"`
	EmojiID structs.Snowflake `json:"-"`
}

type ListApplicationEmojisDto struct {
	ApplicationID structs.Snowflake `json:"-"`
}

type GetApplicationEmojiDto struct {
	ApplicationID structs.Snowflake `json:"-"`
	EmojiID       structs.Snowflake `json:"-"`
}

type CreateApplicationEmojiDto struct {
	ApplicationID structs.Snowflake `json:"-"`
	Name          string            `json:"name"`
	Image         string            `json:"image"`
}

type ModifyApplicationEmojiDto struct {
	ApplicationID structs.Snowflake `json:"-"`
	EmojiID       structs.Snowflake `json:"-"`
	Name          string            `json:"name"`
}

type DeleteApplicationEmojiDto struct {
	ApplicationID structs.Snowflake `json:"-"`
	EmojiID       structs.Snowflake `json:"-"`
}
