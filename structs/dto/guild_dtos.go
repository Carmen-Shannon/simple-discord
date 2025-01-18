package dto

import (
	"time"

	"github.com/Carmen-Shannon/simple-discord/structs"
)

type CreateGuildDto struct {
	Name                        string                                       `json:"name"`
	Region                      *string                                      `json:"region,omitempty"`
	Icon                        *string                                      `json:"icon,omitempty"`
	VerificationLevel           *structs.VerificationLevel                   `json:"verification_level,omitempty"`
	DefaultMessageNotifications *structs.DefaultMessageNotificationLevel     `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *structs.ExplicitContentFilterLevel          `json:"explicit_content_filter,omitempty"`
	Roles                       []structs.Role                               `json:"roles,omitempty"`
	Channels                    []structs.Channel                            `json:"channels,omitempty"`
	AfkChannelID                *structs.Snowflake                           `json:"afk_channel_id,omitempty"`
	AfkTimeout                  *int                                         `json:"afk_timeout,omitempty"`
	SystemChannelID             *structs.Snowflake                           `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *structs.Bitfield[structs.SystemChannelFlag] `json:"system_channel_flags,omitempty"`
}

type GetGuildDto struct {
	GuildID    structs.Snowflake `json:"-"`
	WithCounts *bool             `json:"with_counts,omitempty"`
}

type GetGuildPreviewDto struct {
	GuildID structs.Snowflake `json:"-"`
}

type ModifyGuildDto struct {
	GuildID                     structs.Snowflake                            `json:"-"`
	Name                        string                                       `json:"name"`
	Region                      *string                                      `json:"region,omitempty"`
	VerificationLevel           *structs.VerificationLevel                   `json:"verification_level,omitempty"`
	DefaultMessageNotifications *structs.DefaultMessageNotificationLevel     `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *structs.ExplicitContentFilterLevel          `json:"explicit_content_filter,omitempty"`
	AfkChannelID                *structs.Snowflake                           `json:"afk_channel_id,omitempty"`
	AfkTimeout                  *int                                         `json:"afk_timeout,omitempty"`
	Icon                        *string                                      `json:"icon,omitempty"`
	OwnerID                     *structs.Snowflake                           `json:"owner_id,omitempty"`
	Splash                      *string                                      `json:"splash,omitempty"`
	DiscoverySplash             *string                                      `json:"discovery_splash,omitempty"`
	Banner                      *string                                      `json:"banner,omitempty"`
	SystemChannelID             *structs.Snowflake                           `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *structs.Bitfield[structs.SystemChannelFlag] `json:"system_channel_flags,omitempty"`
	RulesChannelID              *structs.Snowflake                           `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      *structs.Snowflake                           `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             *string                                      `json:"preferred_locale,omitempty"`
	Features                    []structs.GuildFeature                       `json:"features,omitempty"`
	Description                 *string                                      `json:"description,omitempty"`
	PremiumProgressBarEnabled   *bool                                        `json:"premium_progress_bar_enabled,omitempty"`
	SafetyAlertsChannelID       *structs.Snowflake                           `json:"safety_alerts_channel_id,omitempty"`
}

type CreateGuildChannelDto struct {
	GuildID                       structs.Snowflake         `json:"-"`
	Name                          string                    `json:"name"`
	Type                          structs.ChannelType       `json:"type"`
	Topic                         *string                   `json:"topic,omitempty"`
	Bitrate                       *int                      `json:"bitrate,omitempty"`
	UserLimit                     *int                      `json:"user_limit,omitempty"`
	RateLimitPerUser              *int                      `json:"rate_limit_per_user,omitempty"`
	Position                      *int                      `json:"position,omitempty"`
	PermissionOverwrites          []structs.Overwrite       `json:"permission_overwrites,omitempty"`
	ParentID                      *structs.Snowflake        `json:"parent_id,omitempty"`
	NSFW                          *bool                     `json:"nsfw,omitempty"`
	RtcRegion                     *string                   `json:"rtc_region,omitempty"`
	VideoQualityMode              *structs.VideoQualityMode `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration    *int                      `json:"default_auto_archive_duration,omitempty"`
	DefaultReactionEmoji          *structs.DefaultReaction  `json:"default_reaction_emoji,omitempty"`
	AvailableTags                 []structs.ForumTag        `json:"available_tags,omitempty"`
	DefaultSortOrder              *structs.SortOrderType    `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *structs.ForumLayoutType  `json:"default_forum_layout,omitempty"`
	DefaultThreadRateLimitPerUser *int                      `json:"default_thread_rate_limit_per_user,omitempty"`
}

type ModifyGuildChannelPositionsDto struct {
	GuildID         structs.Snowflake  `json:"-"`
	ID              structs.Snowflake  `json:"id"`
	Position        *int               `json:"position,omitempty"`
	LockPermissions *bool              `json:"lock_permissions,omitempty"`
	ParentID        *structs.Snowflake `json:"parent_id,omitempty"`
}

type GetGuildMemberDto struct {
	GuildID structs.Snowflake `json:"-"`
	UserID  structs.Snowflake `json:"-"`
}

type ListGuildMembersDto struct {
	GuildID structs.Snowflake  `json:"-"`
	Limit   *int               `json:"limit,omitempty"`
	After   *structs.Snowflake `json:"after,omitempty"`
}

type SearchGuildMembersDto struct {
	GuildID structs.Snowflake `json:"-"`
	Query   *string           `json:"query"`
	Limit   *int              `json:"limit,omitempty"`
}

type AddGuildMemberDto struct {
	GuildID     structs.Snowflake   `json:"-"`
	UserID      structs.Snowflake   `json:"-"`
	AccessToken string              `json:"access_token"`
	Nick        *string             `json:"nick,omitempty"`
	Roles       []structs.Snowflake `json:"roles,omitempty"`
	Mute        *bool               `json:"mute,omitempty"`
	Deaf        *bool               `json:"deaf,omitempty"`
}

type ModifyGuildMemberDto struct {
	GuildID                    structs.Snowflake                          `json:"-"`
	UserID                     structs.Snowflake                          `json:"-"`
	Nick                       *string                                    `json:"nick,omitempty"`
	Roles                      []structs.Snowflake                        `json:"roles,omitempty"`
	Mute                       *bool                                      `json:"mute,omitempty"`
	Deaf                       *bool                                      `json:"deaf,omitempty"`
	ChannelID                  *structs.Snowflake                         `json:"channel_id,omitempty"`
	CommunicationDisabledUntil *time.Time                                 `json:"communication_disabled_until,omitempty"`
	Flags                      *structs.Bitfield[structs.GuildMemberFlag] `json:"flags,omitempty"`
}

type ModifyCurrentMemberDto struct {
	GuildID structs.Snowflake `json:"-"`
	Nick    *string           `json:"nick,omitempty"`
}

type AddGuildMemberRoleDto struct {
	GuildID structs.Snowflake `json:"-"`
	UserID  structs.Snowflake `json:"-"`
	RoleID  structs.Snowflake `json:"-"`
}

type GetGuildBansDto struct {
	GuildID structs.Snowflake  `json:"-"`
	Limit   *int               `json:"limit,omitempty"`
	Before  *structs.Snowflake `json:"before,omitempty"`
	After   *structs.Snowflake `json:"after,omitempty"`
}

type CreateGuildBanDto struct {
	GuildID              structs.Snowflake `json:"-"`
	UserID               structs.Snowflake `json:"-"`
	DeleteMessageDays    *int              `json:"delete_message_days,omitempty"`
	DeleteMessageSeconds *int              `json:"delete_message_seconds,omitempty"`
}

type BulkGuildBanDto struct {
	GuildID              structs.Snowflake   `json:"-"`
	UserIDs              []structs.Snowflake `json:"user_ids"`
	DeleteMessageSeconds *int                `json:"delete_message_seconds,omitempty"`
}

type GetGuildRoleDto struct {
	GuildID structs.Snowflake `json:"-"`
	RoleID  structs.Snowflake `json:"-"`
}

type CreateGuildRoleDto struct {
	GuildID      structs.Snowflake                     `json:"-"`
	Name         *string                               `json:"name,omitempty"`
	Permissions  *structs.Bitfield[structs.Permission] `json:"permissions,omitempty"`
	Color        *int                                  `json:"color,omitempty"`
	Hoist        *bool                                 `json:"hoist,omitempty"`
	Icon         *string                               `json:"icon,omitempty"`
	UnicodeEmoji *string                               `json:"unicode_emoji,omitempty"`
	Mentionable  *bool                                 `json:"mentionable,omitempty"`
}

type ModifyGuildRolePositionsDto struct {
	GuildID  structs.Snowflake `json:"-"`
	ID       structs.Snowflake `json:"id"`
	Position *int              `json:"position"`
}

type ModifyGuildRoleDto struct {
	GuildID      structs.Snowflake                     `json:"-"`
	RoleID       structs.Snowflake                     `json:"-"`
	Name         *string                               `json:"name,omitempty"`
	Permissions  *structs.Bitfield[structs.Permission] `json:"permissions,omitempty"`
	Color        *int                                  `json:"color,omitempty"`
	Hoist        *bool                                 `json:"hoist,omitempty"`
	Icon         *string                               `json:"icon,omitempty"`
	UnicodeEmoji *string                               `json:"unicode_emoji,omitempty"`
	Mentionable  *bool                                 `json:"mentionable,omitempty"`
}

type ModifyGuildMFALevelDto struct {
	GuildID structs.Snowflake `json:"-"`
	Level   *structs.MFALevel `json:"level"`
}

type GetGuildPruneCountDto struct {
	GuildID      structs.Snowflake   `json:"-"`
	Days         *int                `json:"days,omitempty"`
	IncludeRoles []structs.Snowflake `json:"include_roles,omitempty"`
}

type BeginGuildPruneDto struct {
	GuildID           structs.Snowflake   `json:"-"`
	Days              *int                `json:"days,omitempty"`
	ComputePruneCount *bool               `json:"compute_prune_count,omitempty"`
	IncludeRoles      []structs.Snowflake `json:"include_roles,omitempty"`
	Reason            *string             `json:"reason,omitempty"`
}

type DeleteGuildIntegrationDto struct {
	GuildID       structs.Snowflake `json:"-"`
	IntegrationID structs.Snowflake `json:"-"`
}

type ModifyGuildWidgetDto struct {
	GuildID structs.Snowflake `json:"-"`
	*structs.GuildWidgetSettings
}

type GetGuildWidgetImageDto struct {
	GuildID structs.Snowflake    `json:"-"`
	Style   *structs.WidgetStyle `json:"style,omitempty"`
}

type ModifyGuildWelcomeScreenDto struct {
	GuildID         structs.Snowflake              `json:"-"`
	Enabled         *bool                          `json:"enabled,omitempty"`
	WelcomeChannels []structs.WelcomeScreenChannel `json:"welcome_channels,omitempty"`
	Description     *string                        `json:"description,omitempty"`
}

type ModifyGuildOnboardingDto struct {
	GuildID           structs.Snowflake          `json:"-"`
	Prompts           []structs.OnboardingPrompt `json:"prompts"`
	DefaultChannelIDs []structs.Snowflake        `json:"default_channel_ids"`
	Enabled           bool                       `json:"enabled"`
	Mode              structs.OnboardingMode     `json:"mode"`
}
