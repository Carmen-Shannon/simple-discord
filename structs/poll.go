package structs

import "time"

type Poll struct {
	Question         PollMedia
	Answers          []PollAnswer
	Expiry           *time.Time
	AllowMultiselect bool
	LayoutType       PollLayoutType
	Results          *PollResults
}

type PollCreate struct {
	Question         PollMedia
	Answers          []PollAnswer
	Duration         *int
	AllowMultiselect *bool
	LayoutType       *PollLayoutType
}

type PollMedia struct {
	Text  *string
	Emoji *Emoji
}

type PollAnswer struct {
	AnswerID  int
	PollMedia PollMedia
}

type PollLayoutType int

const (
	DefaultPollLayoutType PollLayoutType = 1
)

type PollResults struct {
	IsFinalized bool
	AnswerCount []PollAnswerCount
}

type PollAnswerCount struct {
	ID        int
	Count     int
	IsMeVoted bool
}
