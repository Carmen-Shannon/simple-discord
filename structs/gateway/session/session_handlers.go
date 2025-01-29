package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway/payload"
	receiveevents "github.com/Carmen-Shannon/simple-discord/structs/gateway/receive_events"
	sendevents "github.com/Carmen-Shannon/simple-discord/structs/gateway/send_events"
	"github.com/Carmen-Shannon/simple-discord/util"
	requestutil "github.com/Carmen-Shannon/simple-discord/util/request_util"
)

func (e *eventHandler) handleInteractionCreateEvent(s ClientSession, p payload.SessionPayload) error {
	if interactionCreateEvent, ok := p.Data.(receiveevents.InteractionCreateEvent); ok {
		name := interactionCreateEvent.Data.Name
		if handler, ok := e.CustomHandlers[name]; ok && handler != nil {
			go func() {
				if err := handler(s, p); err != nil {
					s.Error(err)
				}
			}()
			return nil
		}
		return errors.New("no handler for interaction")
	}
	return errors.New("unexpected payload data type")
}

func handleSendRequestGuildMembersEvent(s ClientSession, p payload.SessionPayload) error {
	fmt.Println("HANDLING REQUEST GUILD MEMBERS EVENT")
	fmt.Println("REQUEST GUILD MEMBERS NOT IMPLEMENTED")
	return nil
}

func handleSendVoiceStateUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if voiceStateUpdateEvent, ok := p.Data.(sendevents.UpdateVoiceStateEvent); ok {
		if voiceStateUpdateEvent.GuildID == nil {
			return errors.New("guild ID not set")
		}

		voiceStateData, err := json.Marshal(p)
		if err != nil {
			return err
		}

		s.Write(voiceStateData, false)
		return nil
	} else {
		return errors.New("unexpected payload data type")
	}
}

func handleSendPresenceUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	fmt.Println("HANDLING PRESENCE UPDATE EVENT")
	fmt.Println("PRESENCE UPDATE NOT IMPLEMENTED")
	return nil
}

func handleSendResumeEvent(s ClientSession, p payload.SessionPayload) error {
	resumeEvent := sendevents.ResumeEvent{
		Token:     *s.GetToken(),
		SessionID: *s.GetSessionID(),
		Seq:       *s.GetSequence(),
	}
	resumePayload := payload.SessionPayload{
		OpCode: gateway.GatewayOpResume,
		Data:   resumeEvent,
	}

	resumeData, err := json.Marshal(resumePayload)
	if err != nil {
		return err
	}

	s.Write(resumeData, false)
	return nil
}

func handleSendIdentifyEvent(s ClientSession, p payload.SessionPayload) error {
	identifyEvent := sendevents.IdentifyEvent{
		Token: *s.GetToken(),
		Properties: sendevents.IdentifyProperties{
			Os:      runtime.GOOS,
			Browser: "discord",
			Device:  "discord",
		},
		Intents: structs.GetIntents(s.GetIntents()),
	}
	if s.GetShard() != nil && s.GetShards() != nil {
		identifyEvent.Shard = &[]int{*s.GetShard(), *s.GetShards()}
	}
	identifyPayload := payload.SessionPayload{
		OpCode: gateway.GatewayOpIdentify,
		Data:   identifyEvent,
	}

	identifyData, err := json.Marshal(identifyPayload)
	if err != nil {
		return err
	}

	s.Write(identifyData, false)
	return nil
}

func handleChannelDeleteEvent(s ClientSession, p payload.SessionPayload) error {
	if channelDeleteEvent, ok := p.Data.(receiveevents.ChannelDeleteEvent); ok {
		servers := s.GetServers()
		server, exists := servers[channelDeleteEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}
		server.DeleteChannel(channelDeleteEvent.ID)
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleChannelUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if channelUpdateEvent, ok := p.Data.(receiveevents.ChannelUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[channelUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}
		channel := structs.Channel(*channelUpdateEvent.Channel)
		server.UpdateChannel(channel.ID, channel)
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleChannelCreateEvent(s ClientSession, p payload.SessionPayload) error {
	if channelCreateEvent, ok := p.Data.(receiveevents.ChannelCreateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[channelCreateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}
		channel := structs.Channel(*channelCreateEvent.Channel)
		channel.Typing = structs.NewTypingChannel()
		server.AddChannel(channel)
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handlePresenceUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if presenceUpdateEvent, ok := p.Data.(receiveevents.PresenceUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[presenceUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.UpdatePresence(presenceUpdateEvent.User.ID, *presenceUpdateEvent.PresenceUpdate)
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleInvalidSessionEvent(s ClientSession, p payload.SessionPayload) error {
	if invalidSessionEvent, ok := p.Data.(receiveevents.InvalidSessionEvent); ok {
		if invalidSessionEvent {
			if err := s.ResumeSession(); err != nil {
				return err
			}
		} else {
			s.CloseResumeReceived()
			if err := s.ReconnectSession(); err != nil {
				return err
			}
		}
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleReconnectEvent(s ClientSession, p payload.SessionPayload) error {
	if _, ok := p.Data.(receiveevents.ReconnectEvent); ok {
		if err := s.ResumeSession(); err != nil {
			return err
		}
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleResumedEvent(s ClientSession, p payload.SessionPayload) error {
	if _, ok := p.Data.(receiveevents.ResumedEvent); ok {
		s.CloseResumeReceived()
		return nil
	} else {
		return errors.New("unexpected payload data type")
	}
}

func handleGuildCreateEvent(s ClientSession, p payload.SessionPayload) error {
	if guildCreateEvent, ok := p.Data.(receiveevents.GuildCreateEvent); ok {
		if guildCreateEvent.Unavailable != nil && !*guildCreateEvent.Unavailable {
			server := structs.NewServer(guildCreateEvent.Guild)
			if err := util.UpdateFields(server, guildCreateEvent.Server); err != nil {
				return err
			}
			s.AddServer(*server)
		}
	} else if guildCreateUnavailableEvent, ok := p.Data.(receiveevents.GuildCreateUnavailableEvent); ok {
		unavailableGuild := &structs.Guild{}
		unavailableGuild.ID = guildCreateUnavailableEvent.ID
		server := structs.NewServer(unavailableGuild)
		server.Unavailable = &guildCreateUnavailableEvent.Unavailable
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleGuildUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if guildUpdateEvent, ok := p.Data.(receiveevents.GuildUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildUpdateEvent.ID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.UpdateGuild(*guildUpdateEvent.Guild)
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildDeleteEvent(s ClientSession, p payload.SessionPayload) error {
	if guildDeleteEvent, ok := p.Data.(receiveevents.GuildDeleteEvent); ok {
		servers := s.GetServers()
		_, exists := servers[guildDeleteEvent.ID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		unavailableGuild := &structs.Guild{}
		unavailableGuild.ID = guildDeleteEvent.ID
		server := structs.NewServer(unavailableGuild)
		server.Unavailable = &guildDeleteEvent.Unavailable
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildBanAddEvent(s ClientSession, p payload.SessionPayload) error {
	fmt.Println("HANDLING GUILD BAN ADD EVENT")
	fmt.Println("GUILD BAN ADD NOT IMPLEMENTED")
	return nil
}

func handleGuildBanRemoveEvent(s ClientSession, p payload.SessionPayload) error {
	fmt.Println("HANDLING GUILD BAN REMOVE EVENT")
	fmt.Println("GUILD BAN REMOVE NOT IMPLEMENTED")
	return nil
}

func handleGuildEmojisUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if guildEmojisUpdateEvent, ok := p.Data.(receiveevents.GuildEmojisUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildEmojisUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.Emojis = guildEmojisUpdateEvent.Emojis
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildIntegrationsUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if _, ok := p.Data.(receiveevents.GuildIntegrationsUpdateEvent); ok {
		return nil
	}
	return errors.New("unexpected payload data type")
}

func handleGuildMemberAddEvent(s ClientSession, p payload.SessionPayload) error {
	if guildMemberAddEvent, ok := p.Data.(receiveevents.GuildMemberAddEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildMemberAddEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.AddMember(*guildMemberAddEvent.GuildMember)
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildMemberRemoveEvent(s ClientSession, p payload.SessionPayload) error {
	if guildMemberRemoveEvent, ok := p.Data.(receiveevents.GuildMemberRemoveEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildMemberRemoveEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.DeleteMember(guildMemberRemoveEvent.User.ID)
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildMemberUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if guildMemberUpdateEvent, ok := p.Data.(receiveevents.GuildMemberUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildMemberUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMember := server.GetMember(guildMemberUpdateEvent.User.ID)
		if currentMember == nil {
			return errors.New("user not in server")
		}

		if err := util.UpdateFields(currentMember, &guildMemberUpdateEvent); err != nil {
			return err
		}

		server.UpdateMember(guildMemberUpdateEvent.User.ID, *currentMember)
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildMembersChunkEvent(s ClientSession, p payload.SessionPayload) error {
	if guildMembersChunkEvent, ok := p.Data.(receiveevents.GuildMembersChunk); ok {
		servers := s.GetServers()
		server, exists := servers[guildMembersChunkEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		for _, member := range guildMembersChunkEvent.Members {
			if server.HasMember(member.User.ID) {
				server.UpdateMember(member.User.ID, member)
			} else {
				server.AddMember(member)
			}
		}

		for _, presence := range guildMembersChunkEvent.Presences {
			if server.HasPresence(presence.User.ID) {
				server.UpdatePresence(presence.User.ID, *presence.PresenceUpdate)
			} else {
				server.AddPresence(*presence.PresenceUpdate)
			}
		}

		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildRoleCreateEvent(s ClientSession, p payload.SessionPayload) error {
	if guildRoleCreateEvent, ok := p.Data.(receiveevents.GuildRoleCreateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildRoleCreateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.AddRole(guildRoleCreateEvent.Role)
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildRoleUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if guildRoleUpdateEvent, ok := p.Data.(receiveevents.GuildRoleUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildRoleUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.UpdateRole(guildRoleUpdateEvent.Role.ID, guildRoleUpdateEvent.Role)
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildRoleDeleteEvent(s ClientSession, p payload.SessionPayload) error {
	if guildRoleDeleteEvent, ok := p.Data.(receiveevents.GuildRoleDeleteEvent); ok {
		servers := s.GetServers()
		server, exists := servers[guildRoleDeleteEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.DeleteRole(guildRoleDeleteEvent.RoleID)
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageCreateEvent(s ClientSession, p payload.SessionPayload) error {
	if messageCreateEvent, ok := p.Data.(receiveevents.MessageCreateEvent); ok {
		if messageCreateEvent.GuildID != nil {
			servers := s.GetServers()
			server, exists := servers[messageCreateEvent.GuildID.ToString()]
			if !exists {
				return errors.New("server not found")
			}

			server.AddMessage(*messageCreateEvent.Message)
		} else {
			return nil
		}
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if messageUpdateEvent, ok := p.Data.(receiveevents.MessageUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[messageUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		if server.GetMessage(messageUpdateEvent.ChannelID, messageUpdateEvent.Message.ID) == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = messageUpdateEvent.ChannelID
			query.MessageID = messageUpdateEvent.Message.ID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}
			server.AddMessage(*message)
		}

		server.UpdateMessage(*messageUpdateEvent.Message)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageDeleteEvent(s ClientSession, p payload.SessionPayload) error {
	if messageDeleteEvent, ok := p.Data.(receiveevents.MessageDeleteEvent); ok {
		servers := s.GetServers()
		server, exists := servers[messageDeleteEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.DeleteMessage(messageDeleteEvent.ChannelID, messageDeleteEvent.ID)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageBulkDeleteEvent(s ClientSession, p payload.SessionPayload) error {
	if messageBulkDeleteEvent, ok := p.Data.(receiveevents.MessageDeleteBulkEvent); ok {
		servers := s.GetServers()
		server, exists := servers[messageBulkDeleteEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		for _, id := range messageBulkDeleteEvent.IDs {
			server.DeleteMessage(messageBulkDeleteEvent.ChannelID, id)
		}
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageReactionAddEvent(s ClientSession, p payload.SessionPayload) error {
	if reactionAddEvent, ok := p.Data.(receiveevents.MessageReactionAddEvent); ok {
		servers := s.GetServers()
		server, exists := servers[reactionAddEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMessage := server.GetMessage(reactionAddEvent.ChannelID, reactionAddEvent.MessageID)
		if currentMessage == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = reactionAddEvent.ChannelID
			query.MessageID = reactionAddEvent.MessageID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}

			currentMessage = message
		}

		// check if there is already a "reaction" cached
		currentReaction := currentMessage.GetReaction(reactionAddEvent.Emoji)
		if currentReaction == nil {
			currentReaction = &structs.Reaction{
				Emoji:       reactionAddEvent.Emoji,
				Count:       0,
				BurstColors: reactionAddEvent.BurstColors,
			}
		}

		// check if the user reacting is the author of the message
		if reactionAddEvent.MessageAuthorID != nil && reactionAddEvent.MessageAuthorID.Equals(reactionAddEvent.UserID) {
			currentReaction.IsMe = true
		}

		// check if the reaction is a burst  or normal reaction
		if reactionAddEvent.Type == receiveevents.MessageReactionBurst {
			currentReaction.CountDetails.Burst++
		} else if reactionAddEvent.Type == receiveevents.MessageReactionNormal {
			currentReaction.CountDetails.Normal++
		}

		currentReaction.Count++
		currentMessage.UpdateReactions(*currentReaction)

		server.UpdateMessage(*currentMessage)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageReactionRemoveEvent(s ClientSession, p payload.SessionPayload) error {
	if reactionRemoveEvent, ok := p.Data.(receiveevents.MessageReactionRemoveEvent); ok {
		servers := s.GetServers()
		server, exists := servers[reactionRemoveEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMessage := server.GetMessage(reactionRemoveEvent.ChannelID, reactionRemoveEvent.MessageID)
		if currentMessage == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = reactionRemoveEvent.ChannelID
			query.MessageID = reactionRemoveEvent.MessageID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}

			currentMessage = message
		}

		currentReaction := currentMessage.GetReaction(reactionRemoveEvent.Emoji)
		if currentReaction == nil {
			return nil
		}

		if currentReaction.Count--; currentReaction.Count == 0 {
			currentMessage.DeleteReaction(reactionRemoveEvent.Emoji)
		} else {
			currentMessage.UpdateReactions(*currentReaction)
		}

		server.UpdateMessage(*currentMessage)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageReactionRemoveAllEvent(s ClientSession, p payload.SessionPayload) error {
	if reactionRemoveAllEvent, ok := p.Data.(receiveevents.MessageReactionRemoveAllEvent); ok {
		servers := s.GetServers()
		server, exists := servers[reactionRemoveAllEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMessage := server.GetMessage(reactionRemoveAllEvent.ChannelID, reactionRemoveAllEvent.MessageID)
		if currentMessage == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = reactionRemoveAllEvent.ChannelID
			query.MessageID = reactionRemoveAllEvent.MessageID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}

			currentMessage = message
		}

		currentMessage.Reactions = []structs.Reaction{}
		server.UpdateMessage(*currentMessage)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessageReactionRemoveEmojiEvent(s ClientSession, p payload.SessionPayload) error {
	if reactionRemoveEmojiEvent, ok := p.Data.(receiveevents.MessageReactionRemoveEmojiEvent); ok {
		servers := s.GetServers()
		server, exists := servers[reactionRemoveEmojiEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMessage := server.GetMessage(reactionRemoveEmojiEvent.ChannelID, reactionRemoveEmojiEvent.MessageID)
		if currentMessage == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = reactionRemoveEmojiEvent.ChannelID
			query.MessageID = reactionRemoveEmojiEvent.MessageID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}

			currentMessage = message
		}

		currentMessage.DeleteReaction(reactionRemoveEmojiEvent.Emoji)
		server.UpdateMessage(*currentMessage)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessagePollVoteAddEvent(s ClientSession, p payload.SessionPayload) error {
	if messagePollVoteAddEvent, ok := p.Data.(receiveevents.MessagePollVoteAddEvent); ok {
		servers := s.GetServers()
		server, exists := servers[messagePollVoteAddEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMessage := server.GetMessage(messagePollVoteAddEvent.ChannelID, messagePollVoteAddEvent.MessageID)
		if currentMessage == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = messagePollVoteAddEvent.ChannelID
			query.MessageID = messagePollVoteAddEvent.MessageID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}

			currentMessage = message
		}
		if currentMessage.Poll == nil {
			return errors.New("no poll active")
		}

		answer := &structs.PollAnswer{
			AnswerID:  messagePollVoteAddEvent.AnswerID,
			PollMedia: structs.PollMedia{},
		}

		currentMessage.Poll.Answers = append(currentMessage.Poll.Answers, *answer)
		server.UpdateMessage(*currentMessage)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleMessagePollVoteRemoveEvent(s ClientSession, p payload.SessionPayload) error {
	if messagePollVoteRemoveEvent, ok := p.Data.(receiveevents.MessagePollVoteRemoveEvent); ok {
		servers := s.GetServers()
		server, exists := servers[messagePollVoteRemoveEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentMessage := server.GetMessage(messagePollVoteRemoveEvent.ChannelID, messagePollVoteRemoveEvent.MessageID)
		if currentMessage == nil {
			var query dto.GetChannelMessageDto
			query.ChannelID = messagePollVoteRemoveEvent.ChannelID
			query.MessageID = messagePollVoteRemoveEvent.MessageID
			message, err := requestutil.GetChannelMessage(query, *s.GetToken())
			if err != nil {
				return err
			} else if message == nil {
				return errors.New("message not found")
			}

			currentMessage = message
		}
		if currentMessage.Poll == nil {
			return errors.New("no poll active")
		}

		for i, answer := range currentMessage.Poll.Answers {
			if answer.AnswerID == messagePollVoteRemoveEvent.AnswerID {
				currentMessage.Poll.Answers = append(currentMessage.Poll.Answers[:i], currentMessage.Poll.Answers[i+1:]...)
				break
			}
		}

		server.UpdateMessage(*currentMessage)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleTypingStartEvent(s ClientSession, p payload.SessionPayload) error {
	if typingStartEvent, ok := p.Data.(receiveevents.TypingStartEvent); ok {
		servers := s.GetServers()
		server, exists := servers[typingStartEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		currentChannel := server.GetChannel(typingStartEvent.ChannelID)
		if currentChannel == nil {
			return errors.New("channel not found")
		}

		if currentChannel.Typing == nil {
			currentChannel.Typing = structs.NewTypingChannel()
		}
		currentChannel.Typing.AddUser(typingStartEvent.UserID)
		server.UpdateChannel(typingStartEvent.ChannelID, *currentChannel)
		s.AddServer(*server)
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleUserUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if userUpdateEvent, ok := p.Data.(receiveevents.UserUpdateEvent); ok {
		servers := s.GetServers()
		for _, server := range servers {
			members := server.GetMembers()
			for _, member := range members {
				if member.User.ID.Equals(userUpdateEvent.ID) {
					if err := util.UpdateFields(member.User, userUpdateEvent); err != nil {
						return err
					}
					server.UpdateMember(userUpdateEvent.ID, member)
					s.AddServer(*server)
					break
				}
			}
		}
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceChannelEffectSendEvent(s ClientSession, p payload.SessionPayload) error {
	if _, ok := p.Data.(receiveevents.VoiceChannelEffectSendEvent); ok {
		fmt.Println("VOICE CHANNEL EFFECT SEND NOT IMPLEMENTED IDK WHAT TO USE IT FOR")
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceStateUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if voiceStateUpdateEvent, ok := p.Data.(receiveevents.VoiceStateUpdateEvent); ok {
		servers := s.GetServers()
		server, exists := servers[voiceStateUpdateEvent.GuildID.ToString()]
		if !exists {
			return errors.New("server not found")
		}

		server.UpdateVoiceState(*voiceStateUpdateEvent.VoiceState)
		s.AddServer(*server)

		// yuck!!!!
		if s.GetBotData().UserDetails.ID.Equals(voiceStateUpdateEvent.UserID) {
			vs := s.GetVoiceSession(*voiceStateUpdateEvent.GuildID)
			if vs != nil {
				vs.SetSessionID(voiceStateUpdateEvent.SessionID)
				vs.SignalVoiceStateReady()
			}
		}
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceServerUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if voiceServerUpdateEvent, ok := p.Data.(receiveevents.VoiceServerUpdateEvent); ok {
		vs := s.GetVoiceSession(voiceServerUpdateEvent.GuildID)
		if vs == nil {
			return errors.New("voice session not initialized")
		}

		if voiceServerUpdateEvent.Endpoint != nil && !strings.Contains(*voiceServerUpdateEvent.Endpoint, "wss://") {
			*voiceServerUpdateEvent.Endpoint = "wss://" + *voiceServerUpdateEvent.Endpoint
		}

		vs.SetConnectUrl(*voiceServerUpdateEvent.Endpoint)
		vs.SetGuildID(voiceServerUpdateEvent.GuildID)
		vs.SetToken(voiceServerUpdateEvent.Token)
		vs.SignalVoiceServerReady()
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleWebhooksUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if _, ok := p.Data.(receiveevents.WebhooksUpdateEvent); ok {
		fmt.Println("WEBHOOKS UPDATE EVENT NOT IMPLEMENTED YET THIS IS USED FOR SERVER WEBHOOK UPDATE EVENTS")
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleGuildAuditLogEntryCreateEvent(s ClientSession, p payload.SessionPayload) error {
	if _, ok := p.Data.(receiveevents.GuildAuditLogEntryCreateEvent); ok {
		fmt.Println("GUILD AUDIT LOG ENTRY CREATE EVENT NOT IMPLEMENTED YET THIS IS USED FOR SERVER AUDIT LOG EVENTS")
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleVoiceChannelStatusUpdateEvent(s ClientSession, p payload.SessionPayload) error {
	if _, ok := p.Data.(receiveevents.VoiceChannelStatusUpdateEvent); ok {
		fmt.Println("VOICE CHANNEL STATUS UPDATE EVENT NOT IMPLEMENTED YET THIS IS USED FOR VOICE CHANNEL STATUS UPDATES ??")
	} else {
		return errors.New("unexpected payload data type")
	}
	return nil
}

func handleHeartbeatACKEvent(s ClientSession, p payload.SessionPayload) error {
	return nil
}

func handleReadyEvent(s ClientSession, p payload.SessionPayload) error {
	if readyEvent, ok := p.Data.(receiveevents.ReadyEvent); ok {
		s.SetSessionID(readyEvent.SessionID)
		s.SetResumeUrl(readyEvent.ResumeGatewayURL)
		s.SetBotData(*structs.NewBotData(readyEvent.User, readyEvent.Application))
		s.CloseReadyReceived()
	} else {
		return errors.New("unexpected payload data type")
	}

	return nil
}

func handleHelloEvent(s ClientSession, p payload.SessionPayload) error {
	if helloEvent, ok := p.Data.(receiveevents.HelloEvent); ok {
		heartbeatInterval := int(helloEvent.HeartbeatInterval)
		s.SetHeartbeatAck(heartbeatInterval)
		s.CloseHelloReceived()
	} else {
		return errors.New("unexpected payload data type")
	}

	return startHeartbeatTimer(s)
}

func handleHeartbeatEvent(s ClientSession, p payload.SessionPayload) error {
	if heartbeatEvent, ok := p.Data.(receiveevents.HeartbeatEvent); ok {
		if heartbeatEvent.LastSequence != nil {
			s.SetSequence(*heartbeatEvent.LastSequence)
		}
		return sendHeartbeatEvent(s)
	}
	return errors.New("unexpected payload data type")
}

func sendHeartbeatEvent(s ClientSession) error {
	heartbeatEvent := sendevents.HeartbeatEvent{
		LastSequence: s.GetSequence(),
	}
	ackPayload := payload.SessionPayload{
		OpCode: gateway.GatewayOpHeartbeat,
		Data:   heartbeatEvent,
	}

	heartbeatData, err := json.Marshal(ackPayload)
	if err != nil {
		return err
	}

	s.Write(heartbeatData, false)
	return nil
}

func heartbeatLoop(ticker *time.Ticker, s ClientSession) {
	defer ticker.Stop()

	firstHeartbeat := true

	for {
		select {
		case <-s.GetCtx().Done():
			return
		case <-ticker.C:
			if firstHeartbeat {
				jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
				time.Sleep(jitter)
				firstHeartbeat = false
			}

			if err := sendHeartbeatEvent(s); err != nil {
				return
			}
		}
	}
}

func startHeartbeatTimer(s ClientSession) error {
	if s.GetHeartbeatAck() == nil {
		return errors.New("no heartbeat interval set")
	}

	ticker := time.NewTicker(time.Duration(*s.GetHeartbeatAck()) * time.Millisecond)
	go heartbeatLoop(ticker, s)
	return nil
}
