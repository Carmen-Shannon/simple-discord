package structs

import (
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
}

func (s *Server) AddRole(role Role) {
	s.Roles = append(s.Roles, role)
}

func (s *Server) UpdateRole(roleId Snowflake, newRole Role) {
	for i, role := range s.Roles {
		if role.ID.Equals(roleId) {
			s.Roles[i] = newRole
			return
		}
	}
}

func (s *Server) DeleteRole(roleId Snowflake) {
	for i, role := range s.Roles {
		if role.ID.Equals(roleId) {
			s.Roles = append(s.Roles[:i], s.Roles[i+1:]...)
			return
		}
	}
}

func (s *Server) UpdateGuild(newGuild Guild) {
	s.Guild = &newGuild
}

func (s *Server) AddMember(member GuildMember) {
	s.Members = append(s.Members, member)
}

func (s *Server) UpdateMember(memberId Snowflake, newMember GuildMember) {
	for i, member := range s.Members {
		if member.User.ID.Equals(memberId) {
			s.Members[i] = newMember
			return
		}
	}
}

func (s *Server) GetMember(memberId Snowflake) *GuildMember {
	for _, member := range s.Members {
		if member.User.ID.Equals(memberId) {
			return &member
		}
	}
	return nil
}

func (s *Server) HasMember(memberId Snowflake) bool {
	for _, member := range s.Members {
		if member.User.ID.Equals(memberId) {
			return true
		}
	}
	return false
}

func (s *Server) DeleteMember(memberId Snowflake) {
	for i, member := range s.Members {
		if member.User.ID.Equals(memberId) {
			s.Members = append(s.Members[:i], s.Members[i+1:]...)
			return
		}
	}
}

func (s *Server) AddChannel(channel Channel) {
	s.Channels = append(s.Channels, channel)
}

func (s *Server) UpdateChannel(channelId Snowflake, newChannel Channel) {
	for i, channel := range s.Channels {
		if channel.ID.Equals(channelId) {
			s.Channels[i] = newChannel
			return
		}
	}
}

func (s *Server) DeleteChannel(channelId Snowflake) {
	for i, channel := range s.Channels {
		if channel.ID.Equals(channelId) {
			s.Channels = append(s.Channels[:i], s.Channels[i+1:]...)
			return
		}
	}
}

func (s *Server) HasPresence(userId Snowflake) bool {
	for _, presence := range s.Presences {
		if presence.User.ID.Equals(userId) {
			return true
		}
	}
	return false
}

func (s *Server) UpdatePresence(userId Snowflake, newPresence PresenceUpdate) {
	for i, presence := range s.Presences {
		if presence.User.ID.Equals(userId) {
			s.Presences[i] = newPresence
			return
		}
	}
}

func (s *Server) AddPresence(presence PresenceUpdate) {
	s.Presences = append(s.Presences, presence)
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
