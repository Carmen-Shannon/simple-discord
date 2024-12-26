package dto

import "github.com/Carmen-Shannon/simple-discord/structs"

type GetChannelDto struct {
	ChannelID structs.Snowflake `json:"channel_id"`
}

type UpdateChannelDto struct {
	ChannelID                     structs.Snowflake                      `json:"channel_id"`
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
	ChannelID   structs.Snowflake                     `json:"channel_id"`
	OverwriteID structs.Snowflake                     `json:"overwrite_id"`
	Allow       *structs.Bitfield[structs.Permission] `json:"allow"`
	Deny        *structs.Bitfield[structs.Permission] `json:"deny"`
	Type        *int                                  `json:"type"`
}
