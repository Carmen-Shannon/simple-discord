package receiveevents

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Carmen-Shannon/simple-discord/structs"
)

type HelloEvent struct {
	HeartbeatInterval int64 `json:"heartbeat_interval"`
}

type HeartbeatEvent struct {
	LastSequence *int `json:"-"`
}

func (h *HeartbeatEvent) UnmarshalJSON(data []byte) error {
	var sequence int
	if err := json.Unmarshal(data, &sequence); err == nil {
		h.LastSequence = &sequence
		return nil
	}

	return errors.New("unable to unmarshal HeartbeatEvent struct")
}

func (h *HeartbeatEvent) MarshalJSON() ([]byte, error) {
	if h.LastSequence == nil {
		return []byte("null"), nil
	}

	return json.Marshal(*h.LastSequence)
}

type HeartbeatACKEvent struct{}

type ReadyEvent struct {
	Version          int                        `json:"v"`
	User             structs.User               `json:"user"`
	Guilds           []structs.UnavailableGuild `json:"guilds"`
	SessionID        string                     `json:"session_id"`
	ResumeGatewayURL string                     `json:"resume_gateway_url"`
	Shard            []int                      `json:"shard,omitempty"`
	Application      structs.Application        `json:"application"`
}

type ResumedEvent struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

type ReconnectEvent struct{}

type InvalidSessionEvent bool

type ApplicationCommandPermissionsUpdateEvent struct {
	GuildApplicationCommand *structs.GuildApplicationCommandPermissions
	ApplicationCommand      *structs.ApplicationCommandPermissions
}

func (acp *ApplicationCommandPermissionsUpdateEvent) UnmarshalJSON(data []byte) error {
	var guildPermissions structs.GuildApplicationCommandPermissions
	if err := json.Unmarshal(data, &guildPermissions); err == nil {
		*acp.GuildApplicationCommand = guildPermissions
		return nil
	}

	var applicationPermissions structs.ApplicationCommandPermissions
	if err := json.Unmarshal(data, &applicationPermissions); err == nil {
		*acp.ApplicationCommand = applicationPermissions
		return nil
	}

	return errors.New("unable to marshal ApplicationCommand struct")
}

type AutoModerationRuleCreateEvent struct {
	*structs.AutoModerationRule
}

type AutoModerationRuleUpdateEvent struct {
	*structs.AutoModerationRule
}

type AutoModerationRuleDeleteEvent struct {
	*structs.AutoModerationRule
}

type AutoModerationActionExecutionEvent struct {
	GuildID              structs.Snowflake                 `json:"guild_id"`
	Action               structs.AutoModerationAction      `json:"action"`
	RuleID               structs.Snowflake                 `json:"rule_id"`
	RuleTriggerType      structs.AutoModerationTriggerType `json:"rule_trigger_type"`
	UserID               structs.Snowflake                 `json:"user_id"`
	ChannelID            *structs.Snowflake                `json:"channel_id,omitempty"`
	MessageID            *structs.Snowflake                `json:"message_id,omitempty"`
	AlertSystemMessageID *structs.Snowflake                `json:"alert_system_message_id,omitempty"`
	Content              string                            `json:"content"`
	MatchedKeyword       *string                           `json:"matched_keyword,omitempty"`
	MatchedContent       *string                           `json:"matched_content,omitempty"`
}

type ChannelCreateEvent struct {
	*structs.Channel
}

type ChannelUpdateEvent struct {
	*structs.Channel
}

type ChannelDeleteEvent struct {
	*structs.Channel
}

type ThreadCreateEvent struct {
	*structs.Channel
	IsNew bool `json:"newly_created"`
}

type ThreadUpdateEvent struct {
	*structs.Channel
}

type ThreadDeleteEvent struct {
	ID       structs.Snowflake   `json:"id"`
	GuildID  structs.Snowflake   `json:"guild_id"`
	ParentID structs.Snowflake   `json:"parent_id"`
	Type     structs.ChannelType `json:"type"`
}

type ThreadListSyncEvent struct {
	GuildID   structs.Snowflake      `json:"guild_id"`
	ChannelID []structs.Snowflake    `json:"channel_ids"`
	Threads   []structs.Channel      `json:"threads"`
	Members   []structs.ThreadMember `json:"members"`
}

type ThreadMemberUpdateEvent struct {
	*structs.ThreadMember
	GuildID structs.Snowflake `json:"guild_id"`
}

type ThreadMembersUpdateEvent struct {
	ID               structs.Snowflake      `json:"id"`
	GuildID          structs.Snowflake      `json:"guild_id"`
	MemberCount      int                    `json:"member_count"`
	AddedMembers     []structs.ThreadMember `json:"added_members"`
	RemovedMemberIDs []structs.Snowflake    `json:"removed_member_ids"`
}

type ChannelPinsUpdateEvent struct {
	GuildID          *structs.Snowflake `json:"guild_id,omitempty"`
	ChannelID        structs.Snowflake  `json:"channel_id"`
	LastPinTimestamp time.Time          `json:"last_pin_timestamp"`
}

type EntitlementCreateEvent struct {
	*structs.Entitlement
}

type EntitlementUpdateEvent struct {
	*structs.Entitlement
}

type EntitlementDeleteEvent struct {
	*structs.Entitlement
}

type GuildCreateEvent struct {
	*structs.Guild
	JoinedAt             time.Time                     `json:"joined_at"`
	Large                bool                          `json:"large"`
	Unavailable          *bool                         `json:"unavailable,omitempty"`
	MemberCount          int                           `json:"member_count"`
	VoiceStates          []structs.VoiceState          `json:"voice_states"`
	Members              []structs.GuildMember         `json:"members"`
	Channels             []structs.Channel             `json:"channels"`
	Threads              []structs.Channel             `json:"threads"`
	Presences            []PresenceUpdateEvent         `json:"presences"`
	StageInstances       []structs.StageInstance       `json:"stage_instances"`
	GuildScheduledEvents []structs.GuildScheduledEvent `json:"guild_scheduled_events"`
}

type GuildCreateUnavailableEvent struct {
	*structs.UnavailableGuild
}

type GuildUpdateEvent struct {
	*structs.Guild
}

type GuildDeleteEvent struct {
	*structs.UnavailableGuild
}

type GuildAuditLogEntryCreateEvent struct {
	*structs.AuditLogEntry
	GuildID structs.Snowflake `json:"guild_id"`
}

type GuildBanAddEvent struct {
	GuildID structs.Snowflake `json:"guild_id"`
	User    structs.User      `json:"user"`
}

type GuildBanRemoveEvent struct {
	GuildID structs.Snowflake `json:"guild_id"`
	User    structs.User      `json:"user"`
}

type GuildEmojisUpdateEvent struct {
	GuildID structs.Snowflake `json:"guild_id"`
	Emojis  []structs.Emoji   `json:"emojis"`
}

type GuildStickersUpdateEvent struct {
	GuildID  structs.Snowflake `json:"guild_id"`
	Stickers []structs.Sticker `json:"stickers"`
}

type GuildIntegrationsUpdateEvent struct {
	GuildID structs.Snowflake `json:"guild_id"`
}

type GuildMemberAddEvent struct {
	*structs.GuildMember
	GuildID structs.Snowflake `json:"guild_id"`
}

type GuildMemberRemoveEvent struct {
	GuildID structs.Snowflake `json:"guild_id"`
	User    structs.User      `json:"user"`
}

type GuildMemberUpdateEvent struct {
	GuildID                    structs.Snowflake             `json:"guild_id"`
	Roles                      []structs.Snowflake           `json:"roles"`
	User                       structs.User                  `json:"user"`
	Nick                       *string                       `json:"nick,omitempty"`
	Avatar                     *string                       `json:"avatar,omitempty"`
	JoinedAt                   *time.Time                    `json:"joined_at,omitempty"`
	PremiumSince               *time.Time                    `json:"premium_since,omitempty"`
	IsDeafened                 *bool                         `json:"deaf,omitempty"`
	IsMuted                    *bool                         `json:"mute,omitempty"`
	IsPending                  *bool                         `json:"pending,omitempty"`
	CommunicationDisabledUntil *time.Time                    `json:"communication_disabled_until,omitempty"`
	Flags                      structs.GuildMemberFlag       `json:"flags"`
	AvatarDecorationData       *structs.AvatarDecorationData `json:"avatar_decoration_data,omitempty"`
}

type GuildRoleCreateEvent struct {
	GuildID structs.Snowflake `json:"guild_id"`
	Role    structs.Role      `json:"role"`
}

type GuildRoleUpdateEvent struct {
	GuildID structs.Snowflake `json:"guild_id"`
	RoleID  structs.Snowflake `json:"role_id"`
}

type GuildRoleDeleteEvent struct {
	GuildID structs.Snowflake `json:"guild_id"`
	RoleID  structs.Snowflake `json:"role_id"`
}

type GuildScheduledEventCreateEvent struct {
	*structs.GuildScheduledEvent
}

type GuildScheduledEventUpdateEvent struct {
	*structs.GuildScheduledEvent
}

type GuildScheduledEventDeleteEvent struct {
	*structs.GuildScheduledEvent
}

type GuildScheduledEventUserAddEvent struct {
	GuildScheduledEventID structs.Snowflake `json:"guild_scheduled_event_id"`
	UserID                structs.Snowflake `json:"user_id"`
	GuildID               structs.Snowflake `json:"guild_id"`
}

type GuildScheduledEventUserRemoveEvent struct {
	GuildScheduledEventID structs.Snowflake `json:"guild_scheduled_event_id"`
	UserID                structs.Snowflake `json:"user_id"`
	GuildID               structs.Snowflake `json:"guild_id"`
}

type IntegrationCreateEvent struct {
	*structs.GuildIntegration
	GuildID structs.Snowflake `json:"guild_id"`
}

type IntegrationUpdateEvent struct {
	*structs.GuildIntegration
	GuildID structs.Snowflake `json:"guild_id"`
}

type IntegrationDeleteEvent struct {
	ID            structs.Snowflake  `json:"id"`
	GuildID       structs.Snowflake  `json:"guild_id"`
	ApplicationID *structs.Snowflake `json:"application_id,omitempty"`
}

type InteractionCreateEvent struct {
	*structs.Interaction
}

type InviteCreateEvent struct {
	ChannelID         structs.Snowflake         `json:"channel_id"`
	Code              string                    `json:"code"`
	CreatedAt         time.Time                 `json:"created_at"`
	GuildID           *structs.Snowflake        `json:"guild_id,omitempty"`
	Inviter           *structs.User             `json:"inviter,omitempty"`
	MaxAge            int                       `json:"max_age"`
	MaxUses           int                       `json:"max_uses"`
	TargetType        *structs.InviteTargetType `json:"target_type,omitempty"`
	TargetUser        *structs.User             `json:"target_user,omitempty"`
	TargetApplication *structs.Application      `json:"target_application,omitempty"`
	IsTemporary       bool                      `json:"temporary"`
	Uses              int                       `json:"uses"`
}

type InviteDeleteEvent struct {
	ChannelID structs.Snowflake  `json:"channel_id"`
	GuildID   *structs.Snowflake `json:"guild_id,omitempty"`
	Code      string             `json:"code"`
}

type MessageCreateEvent struct {
	*structs.Message
	GuildID  *structs.Snowflake   `json:"guild_id,omitempty"`
	Member   *structs.GuildMember `json:"member,omitempty"`
	Mentions []MessageCreateUser  `json:"mentions"`
}

type MessageUpdateEvent struct {
	*structs.Message
	GuildID  *structs.Snowflake   `json:"guild_id,omitempty"`
	Member   *structs.GuildMember `json:"member,omitempty"`
	Mentions []MessageCreateUser  `json:"mentions"`
}

type MessageDeleteEvent struct {
	ID        structs.Snowflake  `json:"id"`
	ChannelID structs.Snowflake  `json:"channel_id"`
	GuildID   *structs.Snowflake `json:"guild_id,omitempty"`
}

type MessageDeleteBulkEvent struct {
	IDs       []structs.Snowflake `json:"ids"`
	ChannelID structs.Snowflake   `json:"channel_id"`
	GuildID   *structs.Snowflake  `json:"guild_id,omitempty"`
}

type MessageReactionAddEvent struct {
	UserID          structs.Snowflake    `json:"user_id"`
	ChannelID       structs.Snowflake    `json:"channel_id"`
	MessageID       structs.Snowflake    `json:"message_id"`
	GuildID         *structs.Snowflake   `json:"guild_id,omitempty"`
	Member          *structs.GuildMember `json:"member,omitempty"`
	Emoji           structs.Emoji        `json:"emoji"`
	MessageAuthorID *structs.Snowflake   `json:"message_author_id,omitempty"`
	Burst           bool                 `json:"burst"`
	BustColors      []string             `json:"bust_colors"`
	Type            MessageReactionType  `json:"type"`
}

type MessageReactionRemoveEvent struct {
	UserID    structs.Snowflake   `json:"user_id"`
	ChannelID structs.Snowflake   `json:"channel_id"`
	MessageID structs.Snowflake   `json:"message_id"`
	GuildID   *structs.Snowflake  `json:"guild_id,omitempty"`
	Emoji     structs.Emoji       `json:"emoji"`
	Burst     bool                `json:"burst"`
	Type      MessageReactionType `json:"type"`
}

type MessageReactionRemoveAllEvent struct {
	ChannelID structs.Snowflake  `json:"channel_id"`
	MessageID structs.Snowflake  `json:"message_id"`
	GuildID   *structs.Snowflake `json:"guild_id,omitempty"`
}

type MessageReactionRemoveEmojiEvent struct {
	ChannelID structs.Snowflake  `json:"channel_id"`
	GuildID   *structs.Snowflake `json:"guild_id,omitempty"`
	MessageID structs.Snowflake  `json:"message_id"`
	Emoji     structs.Emoji      `json:"emoji"`
}

type PresenceUpdateEvent struct {
	User         structs.User       `json:"user"`
	GuildID      structs.Snowflake  `json:"guild_id"`
	Status       UserStatusType     `json:"status"`
	Activities   []structs.Activity `json:"activities"`
	ClientStatus ClientStatus       `json:"client_status"`
	Nonce        *string            `json:"nonce,omitempty"`
}

type StageInstanceCreateEvent struct {
	*structs.StageInstance
}

type StageInstanceUpdateEvent struct {
	*structs.StageInstance
}

type StageInstanceDeleteEvent struct {
	*structs.StageInstance
}

type TypingStartEvent struct {
	ChannelID structs.Snowflake    `json:"channel_id"`
	GuildID   *structs.Snowflake   `json:"guild_id,omitempty"`
	UserID    structs.Snowflake    `json:"user_id"`
	Timestamp time.Time            `json:"timestamp"`
	Member    *structs.GuildMember `json:"member,omitempty"`
}

type UserUpdateEvent struct {
	*structs.User
}

type VoiceStateUpdateEvent struct {
	*structs.VoiceState
}

type VoiceServerUpdateEvent struct {
	Token    string            `json:"token"`
	GuildID  structs.Snowflake `json:"guild_id"`
	Endpoint *string           `json:"endpoint,omitempty"`
}

type WebhooksUpdateEvent struct {
	GuildID   structs.Snowflake `json:"guild_id"`
	ChannelID structs.Snowflake `json:"channel_id"`
}

type MessagePollVoteAddEvent struct {
	UserID    structs.Snowflake  `json:"user_id"`
	ChannelID structs.Snowflake  `json:"channel_id"`
	MessageID structs.Snowflake  `json:"message_id"`
	GuildID   *structs.Snowflake `json:"guild_id,omitempty"`
	AnswerID  int                `json:"answer_id"`
}

type MessagePollVoteRemoveEvent struct {
	UserID    structs.Snowflake  `json:"user_id"`
	ChannelID structs.Snowflake  `json:"channel_id"`
	MessageID structs.Snowflake  `json:"message_id"`
	GuildID   *structs.Snowflake `json:"guild_id,omitempty"`
	AnswerID  int                `json:"answer_id"`
}

type MessageReactionType int

const (
	MessageReactionNormal MessageReactionType = 0
	MessageReactionBurst  MessageReactionType = 1
)

type MessageCreateUser struct {
	*structs.User
	Member structs.GuildMember `json:"member"`
}

type GuildMembersChunk struct {
	GuildID    structs.Snowflake     `json:"guild_id"`
	Members    []structs.GuildMember `json:"members"`
	ChunkIndex int                   `json:"chunk_index"`
	ChunkCount int                   `json:"chunk_count"`
	NotFound   []structs.Snowflake   `json:"not_found"`
	Presences  []PresenceUpdateEvent
}

type UserStatusType string

const (
	UserOnline    UserStatusType = "online"
	UserDND       UserStatusType = "dnd"
	UserIdle      UserStatusType = "idle"
	UserInvisible UserStatusType = "invisible"
	UserOffline   UserStatusType = "offline"
)

type ClientStatus struct {
	Desktop *string `json:"desktop,omitempty"`
	Mobile  *string `json:"mobile,omitempty"`
	Web     *string `json:"web,omitempty"`
}
