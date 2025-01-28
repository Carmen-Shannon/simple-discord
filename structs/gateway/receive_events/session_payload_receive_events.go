package receiveevents

import (
	"encoding/json"
	"errors"

	"github.com/Carmen-Shannon/simple-discord/structs/gateway/payload"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
)

// this function will take an un-typed payload and return the appropriate type based on the OpCode, this will handle dispatch events
func NewReceiveEvent(eventData payload.SessionPayload) (any, error) {
	jsonData, err := json.Marshal(eventData.Data)
	if err != nil {
		return nil, err
	}

	switch eventData.OpCode {
	case gateway.GatewayOpDispatch:
		event, err := handleDispatchEvent(jsonData, eventData)
		if err != nil {
			return nil, err
		}
		eventData.Data = event
		return event, nil
	case gateway.GatewayOpHeartbeat:
		var event HelloEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}
		eventData.Data = event
		return event, nil
	case gateway.GatewayOpReconnect:
		var event ReconnectEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}
		eventData.Data = event
		return event, nil
	case gateway.GatewayOpInvalidSession:
		var event InvalidSessionEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}
		eventData.Data = event
		return event, nil
	case gateway.GatewayOpHello:
		var event HelloEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}
		eventData.Data = event
		return event, nil
	case gateway.GatewayOpHeartbeatACK:
		var event HeartbeatACKEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}
		eventData.Data = event
		return event, nil
	default:
		return nil, errors.New("gateway event assignment failed")
	}
}

// this is where un-typed payloads with an EventName will be assigned to the appropriate struct
func handleDispatchEvent(data []byte, payload payload.SessionPayload) (any, error) {
	if payload.EventName == nil {
		return nil, errors.New("event name is nil")
	}
	switch *payload.EventName {
	case "HELLO":
		var event HelloEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "READY":
		var event ReadyEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "RESUMED":
		var event ResumedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "RECONNECT":
		var event ReconnectEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INVALID_SESSION":
		var event InvalidSessionEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "APPLICATION_COMMAND_PERMISSIONS_UPDATE":
		var event ApplicationCommandPermissionsUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "AUTO_MODERATION_RULE_CREATE":
		var event AutoModerationRuleCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "AUTO_MODERATION_RULE_UPDATE":
		var event AutoModerationRuleUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "AUTO_MODERATION_RULE_DELETE":
		var event AutoModerationRuleDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "AUTO_MODERATION_ACTION_EXECUTION":
		var event AutoModerationActionExecutionEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "CHANNEL_CREATE":
		var event ChannelCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "CHANNEL_UPDATE":
		var event ChannelUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "CHANNEL_DELETE":
		var event ChannelDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "THREAD_CREATE":
		var event ThreadCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "THREAD_UPDATE":
		var event ThreadUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "THREAD_DELETE":
		var event ThreadDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "THREAD_LIST_SYNC":
		var event ThreadListSyncEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "THREAD_MEMBER_UPDATE":
		var event ThreadMemberUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "THREAD_MEMBERS_UPDATE":
		var event ThreadMembersUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "CHANNEL_PINS_UPDATE":
		var event ChannelPinsUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_CREATE":
		var event GuildCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_UPDATE":
		var event GuildUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_DELETE":
		var event GuildDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_BAN_ADD":
		var event GuildBanAddEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_BAN_REMOVE":
		var event GuildBanRemoveEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_EMOJIS_UPDATE":
		var event GuildEmojisUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_STICKERS_UPDATE":
		var event GuildStickersUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_INTEGRATIONS_UPDATE":
		var event GuildIntegrationsUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_MEMBER_ADD":
		var event GuildMemberAddEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_MEMBER_REMOVE":
		var event GuildMemberRemoveEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_MEMBER_UPDATE":
		var event GuildMemberUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_ROLE_CREATE":
		var event GuildRoleCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_ROLE_UPDATE":
		var event GuildRoleUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_ROLE_DELETE":
		var event GuildRoleDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_SCHEDULED_EVENT_CREATE":
		var event GuildScheduledEventCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_SCHEDULED_EVENT_UPDATE":
		var event GuildScheduledEventUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_SCHEDULED_EVENT_DELETE":
		var event GuildScheduledEventDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_SCHEDULED_EVENT_USER_ADD":
		var event GuildScheduledEventUserAddEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_SCHEDULED_EVENT_USER_REMOVE":
		var event GuildScheduledEventUserRemoveEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INTEGRATION_CREATE":
		var event IntegrationCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INTEGRATION_UPDATE":
		var event IntegrationUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INTEGRATION_DELETE":
		var event IntegrationDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INTERACTION_CREATE":
		var event InteractionCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INVITE_CREATE":
		var event InviteCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INVITE_DELETE":
		var event InviteDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_CREATE":
		var event MessageCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_UPDATE":
		var event MessageUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_DELETE":
		var event MessageDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_DELETE_BULK":
		var event MessageDeleteBulkEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_REACTION_ADD":
		var event MessageReactionAddEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_REACTION_REMOVE":
		var event MessageReactionRemoveEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_REACTION_REMOVE_ALL":
		var event MessageReactionRemoveAllEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_REACTION_REMOVE_EMOJI":
		var event MessageReactionRemoveEmojiEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "PRESENCE_UPDATE":
		var event PresenceUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "STAGE_INSTANCE_CREATE":
		var event StageInstanceCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "STAGE_INSTANCE_UPDATE":
		var event StageInstanceUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "STAGE_INSTANCE_DELETE":
		var event StageInstanceDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "TYPING_START":
		var event TypingStartEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "USER_UPDATE":
		var event UserUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "VOICE_STATE_UPDATE":
		var event VoiceStateUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "VOICE_SERVER_UPDATE":
		var event VoiceServerUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "VOICE_CHANNEL_STATUS_UPDATE":
		var event VoiceChannelStatusUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "WEBHOOKS_UPDATE":
		var event WebhooksUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_AUDIT_LOG_ENTRY_CREATE":
		var event GuildAuditLogEntryCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	default:
		return nil, errors.New("dispatch event assignment failed")
	}
}
