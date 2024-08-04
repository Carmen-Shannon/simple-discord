package structs

import "time"

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

type ChannelFlag int

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
	ID                     Snowflake
	Type                   ChannelType
	GuildID                *Snowflake
	Position               *int
	PermissionOverwrites   *Overwrite
	Name                   *string
	Topic                  *string
	NSFW                   *bool
	LastMessageID          *Snowflake
	Bitrate                *int
	UserLimit              *int
	RateLimitPerUser       *int
	Recipients             []User
	Icon                   *string
	OwnerID                *Snowflake
	ApplicationID          *Snowflake
	Managed                *bool
	ParentID               *Snowflake
	LastPinTimestamp       *time.Time
	RtcRegion              *string
	VideoQualityMode       *int
	MessageCount           *int
	MemberCount            *int
	ThreadMetadata         *ThreadMetaData
	ThreadMember           *ThreadMember
	AutoArchiveDuration    *int
	Permissions            *string
	Flags                  *ChannelFlag
	TotalMessageSent       *int
	AvailableTags          []ForumTag
	AppliedTags            []Snowflake
	DefaultReactionEmoji   *DefaultReaction
	DefaultThreadRateLimit *int
	DefaultSortOrder       *SortOrderType
	DefaultForumLayout     *ForumLayoutType
}

type Overwrite struct {
	ID    Snowflake
	Type  int
	Allow string
	Deny  string
}

type ThreadMetaData struct {
	IsArchived          bool
	AutoArchiveDuration int
	ArchiveTimestamp    time.Time
	IsLocked            bool
	IsInvitable         *bool
	CreatedAt           *time.Time
}

type ThreadMember struct {
	ID       *Snowflake
	UserID   *Snowflake
	JoinedAt time.Time
	Flags    int
	Member   GuildMember
}

type ForumTag struct {
	ID          Snowflake
	Name        string
	IsModerated bool
	EmojiID     *Snowflake
	EmojiName   *string
}

type DefaultReaction struct {
	EmojiID   *Snowflake
	EmojiName *string
}
