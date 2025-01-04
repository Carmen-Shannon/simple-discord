package dto

import (
	"encoding/json"
	"errors"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/util"
)

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

type MessageOptions interface {
	SetChannelID(structs.Snowflake)
	SetContent(string) error
	SetNonce(string) error
	SetTTS(bool)
	SetEmbeds([]structs.Embed) error
	SetAllowedMentions(structs.AllowedMentions)
	SetMessageReference(message structs.Message, guildID structs.Snowflake, refType *int) error
	SetComponents([]structs.MessageComponent)
	SetStickerIDs([]structs.Snowflake) error
	SetFiles(map[string][]byte)
	SetAttachments([]structs.Attachment)
	SetFlags(structs.Bitfield[structs.MessageFlag]) error
	SetEnforceNonce(bool)
	SetPoll(structs.Poll)
	Validate() error
	ConstructDtoFromOptions() (*CreateMessageDto, error)
}

var _ MessageOptions = (*CreateMessageDto)(nil)

func NewMessageOptions() MessageOptions {
	return &CreateMessageDto{}
}

func (c *CreateMessageDto) ConstructDtoFromOptions() (*CreateMessageDto, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *CreateMessageDto) Validate() error {
	if c.ChannelID.Equals(*structs.NewSnowflake(0)) {
		return errors.New("channel ID must be set")
	}
	if c.Content == nil &&
		len(c.Embeds) == 0 &&
		len(c.StickerIDs) == 0 &&
		len(c.Components) == 0 &&
		len(c.Files) == 0 &&
		c.Poll == nil {
		return errors.New("message must have content, embeds, stickers, components, files, or a poll")
	}
	if c.Nonce == nil && c.EnforceNonce != nil && *c.EnforceNonce {
		return errors.New("nonce must be set if enforce_nonce is true")
	}
	return nil
}

func (c *CreateMessageDto) SetChannelID(channelID structs.Snowflake) {
	c.ChannelID = channelID
}

func (c *CreateMessageDto) SetContent(content string) error {
	if len(content) > 2000 {
		return errors.New("content must be less than 2000 characters")
	}
	c.Content = &content
	return nil
}

func (c *CreateMessageDto) SetNonce(nonce string) error {
	if len(nonce) > 25 {
		return errors.New("nonce must be less than 25 characters")
	}
	c.Nonce = &nonce
	return nil
}

func (c *CreateMessageDto) SetTTS(tts bool) {
	c.TTS = &tts
}

func (c *CreateMessageDto) SetEmbeds(embeds []structs.Embed) error {
	if len(embeds) > 10 {
		return errors.New("embeds cannot exceed 10")
	}

	for _, embed := range embeds {
		bytes, err := json.Marshal(embed)
		if err != nil {
			return err
		}

		if len(string(bytes)) > 6000 {
			return errors.New("embeds cannot exceed 6000 characters")
		}
	}

	c.Embeds = embeds
	return nil
}

func (c *CreateMessageDto) SetAllowedMentions(allowedMentions structs.AllowedMentions) {
	c.AllowedMentions = &allowedMentions
}

func (c *CreateMessageDto) SetMessageReference(message structs.Message, guildID structs.Snowflake, refType *int) error {
	if refType != nil {
		if structs.MessageReferenceType(*refType) != structs.DefaultMessageReferenceType && structs.MessageReferenceType(*refType) != structs.ForwardMessageReferenceType {
			return errors.New("invalid message reference type")
		}
	} else {
		refType = util.ToPtr(int(structs.DefaultMessageReferenceType))
	}
	c.MessageReference = &structs.MessageReference{
		Type:            structs.MessageReferenceType(*refType),
		MessageID:       &message.ID,
		ChannelID:       &message.ChannelID,
		GuildID:         &guildID,
		FailIfNotExists: util.ToPtr(false),
	}
	return nil
}

func (c *CreateMessageDto) SetComponents(components []structs.MessageComponent) {
	c.Components = components
}

func (c *CreateMessageDto) SetStickerIDs(stickerIDs []structs.Snowflake) error {
	if len(stickerIDs) > 3 {
		return errors.New("sticker IDs cannot exceed 3")
	}
	c.StickerIDs = stickerIDs
	return nil
}

func (c *CreateMessageDto) SetFiles(files map[string][]byte) {
	c.Files = files
}

func (c *CreateMessageDto) SetAttachments(attachments []structs.Attachment) {
	c.Attachments = attachments
}

func (c *CreateMessageDto) SetFlags(flags structs.Bitfield[structs.MessageFlag]) error {
	for _, flag := range flags {
		if flag != structs.SurpressEmbedsMessageFlag && flag != structs.SurpressNotificationsMessageFlag {
			return errors.New("can only accept SUPRESS_EMBEDS and SUPRESS_NOTIFICATIONS flags")
		}
	}
	c.Flags = &flags
	return nil
}

func (c *CreateMessageDto) SetEnforceNonce(enforceNonce bool) {
	c.EnforceNonce = &enforceNonce
}

func (c *CreateMessageDto) SetPoll(poll structs.Poll) {
	c.Poll = &poll
}

type CreateMessageDto struct {
	ChannelID        structs.Snowflake                      `json:"-"`
	Content          *string                                `json:"content,omitempty"`
	Nonce            *string                                `json:"nonce,omitempty"`
	TTS              *bool                                  `json:"tts,omitempty"`
	Embeds           []structs.Embed                        `json:"embeds,omitempty"`
	AllowedMentions  *structs.AllowedMentions               `json:"allowed_mentions,omitempty"`
	MessageReference *structs.MessageReference              `json:"message_reference,omitempty"`
	Components       []structs.MessageComponent             `json:"components,omitempty"`
	StickerIDs       []structs.Snowflake                    `json:"sticker_ids,omitempty"`
	Files            map[string][]byte                      `json:"-"`
	Attachments      []structs.Attachment                   `json:"attachments,omitempty"`
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
	Messages  []structs.Snowflake `json:"messages"`
}
