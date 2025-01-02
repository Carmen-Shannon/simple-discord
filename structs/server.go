package structs

import (
	"sync"
	"time"
)

type Server struct {
	*Guild
	JoinedAt             time.Time             `json:"joined_at"`
	Large                bool                  `json:"large"`
	Unavailable          *bool                 `json:"unavailable,omitempty"`
	MemberCount          int                   `json:"member_count"`
	VoiceStates          []VoiceState          `json:"voice_states"`
	Members              []GuildMember         `json:"members"`
	Channels             []Channel             `json:"channels"`
	Threads              []Channel             `json:"threads"`
	Presences            []PresenceUpdate      `json:"presences"`
	StageInstances       []StageInstance       `json:"stage_instances"`
	GuildScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events"`
	mu                   *sync.RWMutex         `json:"-"`
}

func NewServer(guild *Guild) *Server {
	return &Server{
		Guild: guild,
		mu:    &sync.RWMutex{},
	}
}

func (s *Server) UpdateGuild(newGuild Guild) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Guild = &newGuild
}

func (s *Server) AddMessage(message Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, channel := range s.Channels {
		if channel.ID.Equals(message.ChannelID) {
			s.Channels[i].Messages = append(s.Channels[i].Messages, message)
			return
		}
	}
}

func (s *Server) GetMessage(channelId, messageId Snowflake) *Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, channel := range s.Channels {
		if channel.ID.Equals(channelId) {
			for _, message := range channel.Messages {
				if message.ID.Equals(messageId) {
					return &message
				}
			}
		}
	}

	return nil
}

func (s *Server) UpdateMessage(message Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, channel := range s.Channels {
		if channel.ID.Equals(message.ChannelID) {
			// try to get out early if the message doesnt exist, just add it
			if channel.GetMessage(message.ID) == nil {
				s.Channels[i].Messages = append(s.Channels[i].Messages, message)
				return
			}

			// scan babyyyy
			for j, m := range channel.Messages {
				if m.ID.Equals(message.ID) {
					s.Channels[i].Messages[j] = message
					return
				}
			}
		}
	}
}

func (s *Server) DeleteMessage(channelId, messageId Snowflake) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, channel := range s.Channels {
		if channel.ID.Equals(channelId) {
			for j, message := range channel.Messages {
				if message.ID.Equals(messageId) {
					s.Channels[i].Messages = append(channel.Messages[:j], channel.Messages[j+1:]...)
					return
				}
			}
		}
	}
}

func (s *Server) AddRole(role Role) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Roles = append(s.Roles, role)
}

func (s *Server) UpdateRole(roleId Snowflake, newRole Role) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, role := range s.Roles {
		if role.ID.Equals(roleId) {
			s.Roles[i] = newRole
			return
		}
	}
}

func (s *Server) DeleteRole(roleId Snowflake) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, role := range s.Roles {
		if role.ID.Equals(roleId) {
			s.Roles = append(s.Roles[:i], s.Roles[i+1:]...)
			return
		}
	}
}

func (s *Server) AddMember(member GuildMember) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Members = append(s.Members, member)
}

func (s *Server) GetMember(memberId Snowflake) *GuildMember {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, member := range s.Members {
		if member.User.ID.Equals(memberId) {
			return &member
		}
	}
	return nil
}

func (s *Server) GetMembers() []GuildMember {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Members
}

func (s *Server) UpdateMember(memberId Snowflake, newMember GuildMember) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, member := range s.Members {
		if member.User.ID.Equals(memberId) {
			s.Members[i] = newMember
			return
		}
	}
}

func (s *Server) DeleteMember(memberId Snowflake) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, member := range s.Members {
		if member.User.ID.Equals(memberId) {
			s.Members = append(s.Members[:i], s.Members[i+1:]...)
			return
		}
	}
}

func (s *Server) HasMember(memberId Snowflake) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, member := range s.Members {
		if member.User.ID.Equals(memberId) {
			return true
		}
	}
	return false
}

func (s *Server) AddChannel(channel Channel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Channels = append(s.Channels, channel)
}

func (s *Server) GetChannel(channelId Snowflake) *Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, channel := range s.Channels {
		if channel.ID.Equals(channelId) {
			return &channel
		}
	}
	return nil
}

func (s *Server) UpdateChannel(channelId Snowflake, newChannel Channel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, channel := range s.Channels {
		if channel.ID.Equals(channelId) {
			s.Channels[i] = newChannel
			return
		}
	}
}

func (s *Server) DeleteChannel(channelId Snowflake) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, channel := range s.Channels {
		if channel.ID.Equals(channelId) {
			s.Channels = append(s.Channels[:i], s.Channels[i+1:]...)
			return
		}
	}
}

func (s *Server) AddPresence(presence PresenceUpdate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Presences = append(s.Presences, presence)
}

func (s *Server) UpdatePresence(userId Snowflake, newPresence PresenceUpdate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, presence := range s.Presences {
		if presence.User.ID.Equals(userId) {
			s.Presences[i] = newPresence
			return
		}
	}
}

func (s *Server) HasPresence(userId Snowflake) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, presence := range s.Presences {
		if presence.User.ID.Equals(userId) {
			return true
		}
	}
	return false
}

func (s *Server) GetVoiceState(userId Snowflake) *VoiceState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, voiceState := range s.VoiceStates {
		if voiceState.UserID.Equals(userId) {
			return &voiceState
		}
	}
	return nil
}

func (s *Server) AddVoiceState(voiceState VoiceState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.VoiceStates = append(s.VoiceStates, voiceState)
}

func (s *Server) UpdateVoiceState(voiceState VoiceState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, vs := range s.VoiceStates {
		if vs.UserID.Equals(voiceState.UserID) {
			s.VoiceStates[i] = voiceState
			return
		}
	}
	s.VoiceStates = append(s.VoiceStates, voiceState)
}

type PresenceUpdate struct {
	User         User           `json:"user"`
	GuildID      Snowflake      `json:"guild_id"`
	Status       UserStatusType `json:"status"`
	Activities   []Activity     `json:"activities"`
	ClientStatus ClientStatus   `json:"client_status"`
	Nonce        *string        `json:"nonce,omitempty"`
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
