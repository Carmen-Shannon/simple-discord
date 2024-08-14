package structs

import "time"

type GuildScheduledEvent struct {
	ID                 Snowflake                          `json:"id"`
	GuildID            Snowflake                          `json:"guild_id"`
	ChannelID          *Snowflake                         `json:"channel_id,omitempty"`
	CreatorID          *Snowflake                         `json:"creator_id,omitempty"`
	Name               string                             `json:"name"`
	Description        string                             `json:"description"`
	ScheduledStartTime time.Time                          `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time                         `json:"scheduled_end_time,omitempty"`
	PrivacyLevel       GuildScheduledEventPrivacyLevel    `json:"privacy_level"`
	Status             GuildScheduledEventStatus          `json:"status"`
	EntityType         GuildScheduledEventEntityType      `json:"entity_type"`
	EntityID           *Snowflake                         `json:"entity_id,omitempty"`
	EntityMetadata     *GuildScheduledEventMetadata       `json:"entity_metadata,omitempty"`
	Creator            *User                              `json:"creator,omitempty"`
	UserCount          *int                               `json:"user_count,omitempty"`
	Image              *string                            `json:"image,omitempty"`
	RecurrenceRule     *GuildScheduledEventRecurrenceRule `json:"recurrence_rule,omitempty"`
}

type GuildScheduledEventPrivacyLevel int

const (
	GuildOnly GuildScheduledEventPrivacyLevel = 2
)

type GuildScheduledEventStatus int

const (
	GuildScheduledEventScheduled GuildScheduledEventStatus = 1
	GuildScheduledEventActive    GuildScheduledEventStatus = 2
	GuildScheduledEventCompleted GuildScheduledEventStatus = 3
	GuildScheduledEventCanceled  GuildScheduledEventStatus = 5
)

type GuildScheduledEventEntityType int

const (
	GuildScheduledStageInstance GuildScheduledEventEntityType = 1
	GuildScheduledVoice         GuildScheduledEventEntityType = 2
	GuildScheduledExternal      GuildScheduledEventEntityType = 3
)

type GuildScheduledEventMetadata struct {
	Location string `json:"location"`
}

type GuildScheduledEventRecurrenceRule struct {
	Start      time.Time                `json:"start"`
	End        *time.Time               `json:"end,omitempty"`
	Frequency  RecurrenceRuleFrequency  `json:"frequency"`
	Interval   int                      `json:"interval"`
	ByWeekday  []RecurrenceRuleWeekday  `json:"by_weekday"`
	ByNWeekday []RecurrenceRuleNWeekday `json:"by_n_weekday"`
	ByMonth    []RecurrenceRuleMonth    `json:"by_month"`
	ByMonthDay []int                    `json:"by_month_day"`
	ByYearDay  []int                    `json:"by_year_day"`
	Count      *int                     `json:"count,omitempty"`
}

type RecurrenceRuleFrequency int

const (
	Yearly  RecurrenceRuleFrequency = 0
	Monthly RecurrenceRuleFrequency = 1
	Weekly  RecurrenceRuleFrequency = 2
	Daily   RecurrenceRuleFrequency = 3
)

type RecurrenceRuleWeekday int

const (
	Monday    RecurrenceRuleWeekday = 0
	Tuesday   RecurrenceRuleWeekday = 1
	Wednesday RecurrenceRuleWeekday = 2
	Thursday  RecurrenceRuleWeekday = 3
	Friday    RecurrenceRuleWeekday = 4
	Saturday  RecurrenceRuleWeekday = 5
	Sunday    RecurrenceRuleWeekday = 6
)

type RecurrenceRuleNWeekday struct {
	N   int                   `json:"n"`
	Day RecurrenceRuleWeekday `json:"day"`
}

type RecurrenceRuleMonth int

const (
	January   RecurrenceRuleMonth = 1
	February  RecurrenceRuleMonth = 2
	March     RecurrenceRuleMonth = 3
	April     RecurrenceRuleMonth = 4
	May       RecurrenceRuleMonth = 5
	June      RecurrenceRuleMonth = 6
	July      RecurrenceRuleMonth = 7
	August    RecurrenceRuleMonth = 8
	September RecurrenceRuleMonth = 9
	October   RecurrenceRuleMonth = 10
	November  RecurrenceRuleMonth = 11
	December  RecurrenceRuleMonth = 12
)
