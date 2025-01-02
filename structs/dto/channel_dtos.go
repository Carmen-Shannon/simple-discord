package dto

import (
	"time"

	"github.com/Carmen-Shannon/simple-discord/structs"
)

type GetChannelDto struct {
	ChannelID structs.Snowflake `json:"channel_id"`
}

type UpdateChannelDto struct {
	ChannelID                     structs.Snowflake                      `json:"-"`
	Name                          *string                                `json:"name,omitempty"`
	Icon                          *string                                `json:"icon,omitempty"`                               // only for group DMs
	Type                          *structs.ChannelType                   `json:"type,omitempty"`                               // only for guild channels
	Position                      *int                                   `json:"position,omitempty"`                           // only for guild channels
	Topic                         *string                                `json:"topic,omitempty"`                              // only for guild channels
	NSFW                          *bool                                  `json:"nsfw,omitempty"`                               // only for guild channels
	RateLimitPerUser              *int                                   `json:"rate_limit_per_user,omitempty"`                // for threads and guild channels
	Bitrate                       *int                                   `json:"bitrate,omitempty"`                            // only for guild channels
	UserLimit                     *int                                   `json:"user_limit,omitempty"`                         // only for guild channels
	PermissionOverwrites          *[]structs.Overwrite                   `json:"permission_overwrites,omitempty"`              // only for guild channels
	ParentID                      *structs.Snowflake                     `json:"parent_id,omitempty"`                          // only for guild channels
	RTCRegion                     *string                                `json:"rtc_region,omitempty"`                         // only for guild channels
	VideoQualityMode              *structs.VideoQualityMode              `json:"video_quality_mode,omitempty"`                 // only for guild channels
	DefaultAutoArchiveDuration    *int                                   `json:"default_auto_archive_duration,omitempty"`      // only for guild channels
	Flags                         *structs.Bitfield[structs.ChannelFlag] `json:"flags,omitempty"`                              // for threads and guild channels
	AvailableTags                 *[]structs.ForumTag                    `json:"available_tags,omitempty"`                     // only for guild channels
	DefaultReactionEmoji          *structs.DefaultReaction               `json:"default_reaction_emoji,omitempty"`             // only for guild channels
	DefaultThreadRateLimitPerUser *int                                   `json:"default_thread_rate_limit_per_user,omitempty"` // only for guild channels
	DefaultSortOrder              *structs.SortOrderType                 `json:"default_sort_order,omitempty"`                 // only for guild channels
	DefaultForumLayout            *structs.ForumLayoutType               `json:"default_forum_layout,omitempty"`               // only for guild channels
	Archived                      *bool                                  `json:"archived,omitempty"`                           // only for threads
	AutoArchiveDuration           *int                                   `json:"auto_archive_duration,omitempty"`              // only for threads
	Locked                        *bool                                  `json:"locked,omitempty"`                             // only for threads
	Invitable                     *bool                                  `json:"invitable,omitempty"`                          // only for threads
	AppliedTags                   *[]structs.Snowflake                   `json:"applied_tags,omitempty"`                       // only for threads
}

type EditChannelPermissionsDto struct {
	ChannelID   structs.Snowflake                     `json:"-"`
	OverwriteID structs.Snowflake                     `json:"-"`
	Allow       *structs.Bitfield[structs.Permission] `json:"allow"`
	Deny        *structs.Bitfield[structs.Permission] `json:"deny"`
	Type        *int                                  `json:"type"`
}

type CreateChannelInviteDto struct {
	ChannelID           structs.Snowflake  `json:"-"`
	MaxAge              *int               `json:"max_age,omitempty"`
	MaxUses             *int               `json:"max_uses,omitempty"`
	Temporary           *bool              `json:"temporary,omitempty"`
	Unique              *bool              `json:"unique,omitempty"`
	TargetType          *int               `json:"target_type,omitempty"`
	TargetUserId        *structs.Snowflake `json:"target_user_id,omitempty"`
	TargetApplicationID *structs.Snowflake `json:"target_application_id,omitempty"`
}

type DeleteChannelPermissionDto struct {
	ChannelID   structs.Snowflake `json:"channel_id"`
	OverwriteID structs.Snowflake `json:"overwrite_id"`
}

type FollowAnnouncementChannelDto struct {
	ChannelID        structs.Snowflake `json:"-"`
	WebhookChannelID structs.Snowflake `json:"webhook_channel_id"`
}

type TriggerTypingIndicatorDto struct {
	ChannelID structs.Snowflake `json:"-"`
}

type PinMessageDto struct {
	ChannelID structs.Snowflake `json:"-"`
	MessageID structs.Snowflake `json:"-"`
}

type GroupDMAddRecipientDto struct {
	ChannelID   structs.Snowflake `json:"-"`
	UserID      structs.Snowflake `json:"-"`
	AccessToken string            `json:"access_token"`
	Nick        string            `json:"nick"`
}

type GroupDMRemoveRecipientDto struct {
	ChannelID structs.Snowflake `json:"-"`
	UserID    structs.Snowflake `json:"-"`
}

type StartThreadFromMessageDto struct {
	ChannelID           structs.Snowflake `json:"-"`
	MessageID           structs.Snowflake `json:"-"`
	Name                string            `json:"name"`
	AutoArchiveDuration *int              `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    *int              `json:"rate_limit_per_user,omitempty"`
}

type StartThreadWithoutMessageDto struct {
	ChannelID           structs.Snowflake    `json:"-"`
	Name                string               `json:"name"`
	AutoArchiveDuration *int                 `json:"auto_archive_duration,omitempty"`
	Type                *structs.ChannelType `json:"type,omitempty"`
	Invitable           *bool                `json:"invitable,omitempty"`
	RateLimitPerUser    *int                 `json:"rate_limit_per_user,omitempty"`
}

type StartThreadInForumOrMediaChannelDto struct {
	ChannelID           structs.Snowflake   `json:"-"`
	Name                string              `json:"name"`
	AutoArchiveDuration *int                `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    *int                `json:"rate_limit_per_user,omitempty"`
	AppliedTags         []structs.Snowflake `json:"applied_tags,omitempty"`
	Files               map[string][]byte   `json:"-"`
	ForumAndMediaThreadMessage
}

type ForumAndMediaThreadMessage struct {
	Content         *string                                `json:"content.omitempty"`
	Embeds          []structs.Embed                        `json:"embeds,omitempty"`
	AllowedMentions *structs.AllowedMentions               `json:"allowed_mentions,omitempty"`
	Components      []structs.MessageComponent             `json:"components,omitempty"`
	StickerIDs      []structs.Snowflake                    `json:"sticker_ids,omitempty"`
	Attachments     []structs.Attachment                   `json:"attachments,omitempty"`
	Flags           *structs.Bitfield[structs.MessageFlag] `json:"flags,omitempty"`
}

type GetThreadMemberDto struct {
	ChannelID  structs.Snowflake `json:"-"`
	UserID     structs.Snowflake `json:"-"`
	WithMember *bool             `json:"with_member,omitempty"`
}

type ListThreadMembersDto struct {
	ChannelID  structs.Snowflake  `json:"-"`
	WithMember *bool              `json:"with_member,omitempty"`
	After      *structs.Snowflake `json:"after,omitempty"`
	Limit      *int               `json:"limit,omitempty"`
}

type ListPublicArchivedThreadsDto struct {
	ChannelID structs.Snowflake `json:"-"`
	Before    *time.Time        `json:"before,omitempty"`
	Limit     *int              `json:"limit,omitempty"`
}
