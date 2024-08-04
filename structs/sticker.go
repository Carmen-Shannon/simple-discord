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
	ID          Snowflake
	PackID      *Snowflake
	Name        string
	Description *string
	Tags        string
	Asset       string
	Type        StickerType
	FormatType  StickerFormatType
	Available   *bool
	GuildID     *Snowflake
	User        *User
	SortValue   *int
}

type StickerItem struct {
	ID Snowflake
	Name string
	FormatType StickerFormatType
}
