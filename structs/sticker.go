package structs

type StickerType int

const (
	StandardSticker StickerType = 1
	GuildSticker    StickerType = 2
)

type StickerFormatType int

const (
	PNG    StickerFormatType = 1
	APNG   StickerFormatType = 2
	LOTTIE StickerFormatType = 3
	GIF    StickerFormatType = 4
)

type Sticker struct {
	ID          Snowflake         `json:"id"`
	PackID      *Snowflake        `json:"pack_id,omitempty"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Tags        string            `json:"tags"`
	Asset       string            `json:"asset"`
	Type        StickerType       `json:"type"`
	FormatType  StickerFormatType `json:"format_type"`
	Available   *bool             `json:"available,omitempty"`
	GuildID     *Snowflake        `json:"guild_id,omitempty"`
	User        *User             `json:"user,omitempty"`
	SortValue   *int              `json:"sort_value,omitempty"`
}

type StickerItem struct {
	ID         Snowflake         `json:"id"`
	Name       string            `json:"name"`
	FormatType StickerFormatType `json:"format_type"`
}
