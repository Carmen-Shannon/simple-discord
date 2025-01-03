package gateway

import (
	"encoding/json"
	"errors"

	receiveevents "github.com/Carmen-Shannon/simple-discord/gateway/receive_events"
	sendevents "github.com/Carmen-Shannon/simple-discord/gateway/send_events"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
)

// this function will take an un-typed Payload and return the appropriate type based on the OpCode
func NewSendEvent(eventData gateway.Payload) (any, error) {
	jsonData, err := json.Marshal(eventData.Data)
	if err != nil {
		return nil, err
	}

	switch eventData.OpCode {

	case gateway.Heartbeat:
		var event sendevents.HeartbeatEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.Identify:
		var event sendevents.IdentifyEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.PresenceUpdate:
		var event sendevents.PresenceUpdateEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.VoiceStateUpdate:
		var event sendevents.UpdateVoiceStateEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.Resume:
		var event sendevents.ResumeEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	case gateway.RequestGuildMembers:
		var event sendevents.RequestGuildMembersEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}

		eventData.Data = event
		return event, nil
	default:
		return nil, errors.New("gateway event assignment failed")
	}
}

// this function will take an un-typed payload and return the appropriate type based on the OpCode, this will handle dispatch events
func NewReceiveEvent(eventData gateway.Payload) (any, error) {
	jsonData, err := json.Marshal(eventData.Data)
	if err != nil {
		return nil, err
	}

	switch eventData.OpCode {
	case gateway.Dispatch:
		event, err := handleDispatchEvent(jsonData, eventData)
		if err != nil {
			return nil, err
		}
		eventData.Data = event
		return event, nil
	case gateway.Heartbeat:
		var event receiveevents.HelloEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}
		eventData.Data = event
		return event, nil
	case gateway.Reconnect:
		var event receiveevents.ReconnectEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}
		eventData.Data = event
		return event, nil
	case gateway.InvalidSession:
		var event receiveevents.InvalidSessionEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}
		eventData.Data = event
		return event, nil
	case gateway.Hello:
		var event receiveevents.HelloEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return nil, err
		}
		eventData.Data = event
		return event, nil
	case gateway.HeartbeatACK:
		var event receiveevents.HeartbeatACKEvent
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
func handleDispatchEvent(data []byte, payload gateway.Payload) (any, error) {
	if payload.EventName == nil {
		return nil, errors.New("event name is nil")
	}
	switch *payload.EventName {
	case "HELLO":
		var event receiveevents.HelloEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "READY":
		var event receiveevents.ReadyEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "RESUMED":
		var event receiveevents.ResumedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "RECONNECT":
		var event receiveevents.ReconnectEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INVALID_SESSION":
		var event receiveevents.InvalidSessionEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "APPLICATION_COMMAND_PERMISSIONS_UPDATE":
		var event receiveevents.ApplicationCommandPermissionsUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "AUTO_MODERATION_RULE_CREATE":
		var event receiveevents.AutoModerationRuleCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "AUTO_MODERATION_RULE_UPDATE":
		var event receiveevents.AutoModerationRuleUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "AUTO_MODERATION_RULE_DELETE":
		var event receiveevents.AutoModerationRuleDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "AUTO_MODERATION_ACTION_EXECUTION":
		var event receiveevents.AutoModerationActionExecutionEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "CHANNEL_CREATE":
		var event receiveevents.ChannelCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "CHANNEL_UPDATE":
		var event receiveevents.ChannelUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "CHANNEL_DELETE":
		var event receiveevents.ChannelDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "THREAD_CREATE":
		var event receiveevents.ThreadCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "THREAD_UPDATE":
		var event receiveevents.ThreadUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "THREAD_DELETE":
		var event receiveevents.ThreadDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "THREAD_LIST_SYNC":
		var event receiveevents.ThreadListSyncEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "THREAD_MEMBER_UPDATE":
		var event receiveevents.ThreadMemberUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "THREAD_MEMBERS_UPDATE":
		var event receiveevents.ThreadMembersUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "CHANNEL_PINS_UPDATE":
		var event receiveevents.ChannelPinsUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_CREATE":
		var event receiveevents.GuildCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_UPDATE":
		var event receiveevents.GuildUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_DELETE":
		var event receiveevents.GuildDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_BAN_ADD":
		var event receiveevents.GuildBanAddEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_BAN_REMOVE":
		var event receiveevents.GuildBanRemoveEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_EMOJIS_UPDATE":
		var event receiveevents.GuildEmojisUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_STICKERS_UPDATE":
		var event receiveevents.GuildStickersUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_INTEGRATIONS_UPDATE":
		var event receiveevents.GuildIntegrationsUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_MEMBER_ADD":
		var event receiveevents.GuildMemberAddEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_MEMBER_REMOVE":
		var event receiveevents.GuildMemberRemoveEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_MEMBER_UPDATE":
		var event receiveevents.GuildMemberUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_ROLE_CREATE":
		var event receiveevents.GuildRoleCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_ROLE_UPDATE":
		var event receiveevents.GuildRoleUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_ROLE_DELETE":
		var event receiveevents.GuildRoleDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_SCHEDULED_EVENT_CREATE":
		var event receiveevents.GuildScheduledEventCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_SCHEDULED_EVENT_UPDATE":
		var event receiveevents.GuildScheduledEventUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_SCHEDULED_EVENT_DELETE":
		var event receiveevents.GuildScheduledEventDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_SCHEDULED_EVENT_USER_ADD":
		var event receiveevents.GuildScheduledEventUserAddEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_SCHEDULED_EVENT_USER_REMOVE":
		var event receiveevents.GuildScheduledEventUserRemoveEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INTEGRATION_CREATE":
		var event receiveevents.IntegrationCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INTEGRATION_UPDATE":
		var event receiveevents.IntegrationUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INTEGRATION_DELETE":
		var event receiveevents.IntegrationDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INTERACTION_CREATE":
		var event receiveevents.InteractionCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INVITE_CREATE":
		var event receiveevents.InviteCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "INVITE_DELETE":
		var event receiveevents.InviteDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_CREATE":
		var event receiveevents.MessageCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_UPDATE":
		var event receiveevents.MessageUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_DELETE":
		var event receiveevents.MessageDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_DELETE_BULK":
		var event receiveevents.MessageDeleteBulkEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_REACTION_ADD":
		var event receiveevents.MessageReactionAddEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_REACTION_REMOVE":
		var event receiveevents.MessageReactionRemoveEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_REACTION_REMOVE_ALL":
		var event receiveevents.MessageReactionRemoveAllEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "MESSAGE_REACTION_REMOVE_EMOJI":
		var event receiveevents.MessageReactionRemoveEmojiEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "PRESENCE_UPDATE":
		var event receiveevents.PresenceUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "STAGE_INSTANCE_CREATE":
		var event receiveevents.StageInstanceCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "STAGE_INSTANCE_UPDATE":
		var event receiveevents.StageInstanceUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "STAGE_INSTANCE_DELETE":
		var event receiveevents.StageInstanceDeleteEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "TYPING_START":
		var event receiveevents.TypingStartEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "USER_UPDATE":
		var event receiveevents.UserUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "VOICE_STATE_UPDATE":
		var event receiveevents.VoiceStateUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "VOICE_SERVER_UPDATE":
		var event receiveevents.VoiceServerUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "VOICE_CHANNEL_STATUS_UPDATE":
		var event receiveevents.VoiceChannelStatusUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "WEBHOOKS_UPDATE":
		var event receiveevents.WebhooksUpdateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	case "GUILD_AUDIT_LOG_ENTRY_CREATE":
		var event receiveevents.GuildAuditLogEntryCreateEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		payload.Data = event
		return event, nil
	default:
		return nil, errors.New("dispatch event assignment failed")
	}
}
