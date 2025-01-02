package structs

import (
	"sync"
	"time"
)

type ChannelType int

const (
	GuildTextChannel          ChannelType = 0
	DMChannel                 ChannelType = 1
	GuildVoiceChannel         ChannelType = 2
	GroupDMChannel            ChannelType = 3
	GuildCategoryChannel      ChannelType = 4
	GuildAnnouncementChannel  ChannelType = 5
	AnnouncementThreadChannel ChannelType = 10
	PublicThreadChannel       ChannelType = 11
	PrivateThreadChannel      ChannelType = 12
	GuildStageVoiceChannel    ChannelType = 13
	GuildDirectoryChannel     ChannelType = 14
	GuildForumChannel         ChannelType = 15
	GuildMediaChannel         ChannelType = 16
)

type ChannelFlag int64

const (
	PinnedFlag                   ChannelFlag = 1 << 1
	RequireTagFlag               ChannelFlag = 1 << 4
	HideMediaDownloadOptionsFlag ChannelFlag = 1 << 15
)

type SortOrderType int

const (
	LatestActivitySort SortOrderType = 0
	CreationDateSort   SortOrderType = 1
)

type ForumLayoutType int

const (
	NotSetLayout      ForumLayoutType = 0
	ListViewLayout    ForumLayoutType = 1
	GalleryViewLayout ForumLayoutType = 2
)

type Channel struct {
	ID                     Snowflake              `json:"id"`
	Type                   ChannelType            `json:"type"`
	GuildID                *Snowflake             `json:"guild_id,omitempty"`
	Position               *int                   `json:"position,omitempty"`
	PermissionOverwrites   []Overwrite            `json:"permission_overwrites,omitempty"`
	Name                   *string                `json:"name,omitempty"`
	Topic                  *string                `json:"topic,omitempty"`
	NSFW                   *bool                  `json:"nsfw,omitempty"`
	LastMessageID          *Snowflake             `json:"last_message_id,omitempty"`
	Bitrate                *int                   `json:"bitrate,omitempty"`
	UserLimit              *int                   `json:"user_limit,omitempty"`
	RateLimitPerUser       *int                   `json:"rate_limit_per_user,omitempty"`
	Recipients             []User                 `json:"recipients,omitempty"`
	Icon                   *string                `json:"icon,omitempty"`
	OwnerID                *Snowflake             `json:"owner_id,omitempty"`
	ApplicationID          *Snowflake             `json:"application_id,omitempty"`
	Managed                *bool                  `json:"managed,omitempty"`
	ParentID               *Snowflake             `json:"parent_id,omitempty"`
	LastPinTimestamp       *time.Time             `json:"last_pin_timestamp,omitempty"`
	RtcRegion              *string                `json:"rtc_region,omitempty"`
	VideoQualityMode       *int                   `json:"video_quality_mode,omitempty"`
	MessageCount           *int                   `json:"message_count,omitempty"`
	MemberCount            *int                   `json:"member_count,omitempty"`
	ThreadMetadata         *ThreadMetaData        `json:"thread_metadata,omitempty"`
	ThreadMember           *ThreadMember          `json:"thread_member,omitempty"`
	AutoArchiveDuration    *int                   `json:"auto_archive_duration,omitempty"`
	Permissions            *string                `json:"permissions,omitempty"`
	Flags                  *Bitfield[ChannelFlag] `json:"flags,omitempty"`
	TotalMessageSent       *int                   `json:"total_message_sent,omitempty"`
	AvailableTags          []ForumTag             `json:"available_tags,omitempty"`
	AppliedTags            []Snowflake            `json:"applied_tags,omitempty"`
	DefaultReactionEmoji   *DefaultReaction       `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimit *int                   `json:"default_thread_rate_limit,omitempty"`
	DefaultSortOrder       *SortOrderType         `json:"default_sort_order,omitempty"`
	DefaultForumLayout     *ForumLayoutType       `json:"default_forum_layout,omitempty"`
	Messages               []Message              `json:"-"`
	Typing                 *TypingChannel         `json:"-"`
}

func (c *Channel) GetMessage(messageId Snowflake) *Message {
	for _, message := range c.Messages {
		if message.ID.Equals(messageId) {
			return &message
		}
	}
	return nil
}

type TypingChannel struct {
	Users  map[string]*Snowflake  `json:"-"`
	mu     sync.RWMutex           `json:"-"`
	timers map[string]*time.Timer `json:"-"`
}

func NewTypingChannel() *TypingChannel {
	return &TypingChannel{
		Users:  make(map[string]*Snowflake),
		timers: make(map[string]*time.Timer),
	}
}

func (t *TypingChannel) GetUserId(userId Snowflake) *Snowflake {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, user := range t.Users {
		if user.Equals(userId) {
			return user
		}
	}

	return nil
}

func (t *TypingChannel) AddUser(user Snowflake) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Users[user.ToString()] = &user
	timer := t.startTimeout(user)
	t.timers[user.ToString()] = timer
}

func (t *TypingChannel) DeleteUser(user Snowflake) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.Users, user.ToString())
	t.timers[user.ToString()].Stop()
	delete(t.timers, user.ToString())
}

func (t *TypingChannel) startTimeout(user Snowflake) *time.Timer {
	return time.AfterFunc(15*time.Second, func() {
		t.DeleteUser(user)
	})
}

type Overwrite struct {
	ID    Snowflake `json:"id"`
	Type  int       `json:"type"`
	Allow string    `json:"allow"`
	Deny  string    `json:"deny"`
}

type ThreadMetaData struct {
	IsArchived          bool       `json:"is_archived"`
	AutoArchiveDuration int        `json:"auto_archive_duration"`
	ArchiveTimestamp    time.Time  `json:"archive_timestamp"`
	IsLocked            bool       `json:"is_locked"`
	IsInvitable         *bool      `json:"is_invitable,omitempty"`
	CreatedAt           *time.Time `json:"created_at,omitempty"`
}

type ThreadMember struct {
	ID       *Snowflake      `json:"id,omitempty"`
	UserID   *Snowflake      `json:"user_id,omitempty"`
	JoinedAt time.Time       `json:"joined_at"`
	Flags    Bitfield[int64] `json:"flags"`
	Member   GuildMember     `json:"member"`
}

type ForumTag struct {
	ID          Snowflake  `json:"id"`
	Name        string     `json:"name"`
	IsModerated bool       `json:"is_moderated"`
	EmojiID     *Snowflake `json:"emoji_id,omitempty"`
	EmojiName   *string    `json:"emoji_name,omitempty"`
}

type DefaultReaction struct {
	EmojiID   *Snowflake `json:"emoji_id,omitempty"`
	EmojiName *string    `json:"emoji_name,omitempty"`
}

type VideoQualityMode int

const (
	AutoVideoQualityMode VideoQualityMode = 1
	FullVideoQualityMode VideoQualityMode = 2
)

type FollowedChannel struct {
	ChannelID Snowflake `json:"channel_id"`
	WebhookID Snowflake `json:"webhook_id"`
}
