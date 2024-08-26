package dto

import "github.com/Carmen-Shannon/simple-discord/structs"

type GetChannelMessagesDto struct {
	ChannelID structs.Snowflake  `json:"channel_id"`
	Around    *structs.Snowflake `json:"around,omitempty"`
	Before    *structs.Snowflake `json:"before,omitempty"`
	After     *structs.Snowflake `json:"after,omitempty"`
	Limit     *int               `json:"limit,omitempty"`
}

type GetChannelMessageDto struct {
	ChannelID structs.Snowflake `json:"channel_id"`
	MessageID structs.Snowflake `json:"message_id"`
}

// At least ONE of Content, Embeds, StickerIDs, Components, Files, or Poll is required to use this in a request
type CreateChannelMessageDto struct {
	ChannelID        structs.Snowflake                      `json:"channel_id"`
	Content          *string                                `json:"content,omitempty"`
	Nonce            *string                                `json:"nonce,omitempty"`
	TTS              *bool                                  `json:"tts,omitempty"`
	Embeds           *[]structs.Embed                       `json:"embeds,omitempty"`
	AllowedMentions  *structs.AllowedMentions               `json:"allowed_mentions,omitempty"`
	MessageReference *structs.MessageReference              `json:"message_reference,omitempty"`
	Components       *[]structs.MessageComponent            `json:"components,omitempty"`
	StickerIDs       *[]structs.Snowflake                   `json:"sticker_ids,omitempty"`
	Files            *string                                `json:"files,omitempty"` // TODO: this is for uploading files to messages: https://discord.com/developers/docs/reference#uploading-files DONT USE THIS YET
	PayloadJSON      *string                                `json:"payload_json,omitempty"`
	Attachments      *[]structs.Attachment                  `json:"attachments,omitempty"`
	Flags            *structs.Bitfield[structs.MessageFlag] `json:"flags,omitempty"`
	EnforceNonce     *bool                                  `json:"enforce_nonce,omitempty"`
	Poll             *structs.Poll                          `json:"poll,omitempty"`
}

type CreateReactionDto struct {
	ChannelID structs.Snowflake `json:"channel_id"`
	MessageID structs.Snowflake `json:"message_id"`
	Emoji     structs.Emoji     `json:"emoji"`
}

type DeleteUserReactionDto struct {
	CreateReactionDto
	UserID structs.Snowflake `json:"user_id"`
}

type GetReactionsDto struct {
	CreateReactionDto
	Type  *structs.ReactionType `json:"type,omitempty"`
	After *structs.Snowflake    `json:"after,omitempty"`
	Limit *int                  `json:"limit,omitempty"`
}

type EditMessageDto struct {
	GetChannelMessageDto
	Content         *string                                `json:"content,omitempty"`
	Embeds          *[]structs.Embed                       `json:"embeds,omitempty"`
	Flags           *structs.Bitfield[structs.MessageFlag] `json:"flags,omitempty"`
	AllowedMentions *structs.AllowedMentions               `json:"allowed_mentions,omitempty"`
	Components      *[]structs.MessageComponent            `json:"components,omitempty"`
	Files           *string                                `json:"files,omitempty"` // TODO: this is for uploading files to messages: https://discord.com/developers/docs/reference#uploading-files DONT USE THIS YET
	PayloadJson     *string                                `json:"payload_json,omitempty"`
	Attachments     *[]structs.Attachment                  `json:"attachments,omitempty"`
}

type BulkDeleteMessagesDto struct {
	ChannelID structs.Snowflake   `json:"channel_id"`
	Messages []structs.Snowflake `json:"messages"`
}
