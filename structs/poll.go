package structs

import "time"

type Poll struct {
	Question         PollMedia      `json:"question"`
	Answers          []PollAnswer   `json:"answers"`
	Expiry           *time.Time     `json:"expiry,omitempty"`
	AllowMultiselect bool           `json:"allow_multiselect"`
	LayoutType       PollLayoutType `json:"layout_type"`
	Results          *PollResults   `json:"results,omitempty"`
}

type PollCreate struct {
	Question         PollMedia       `json:"question"`
	Answers          []PollAnswer    `json:"answers"`
	Duration         *int            `json:"duration,omitempty"`
	AllowMultiselect *bool           `json:"allow_multiselect,omitempty"`
	LayoutType       *PollLayoutType `json:"layout_type,omitempty"`
}

type PollMedia struct {
	Text  *string `json:"text,omitempty"`
	Emoji *Emoji  `json:"emoji,omitempty"`
}

type PollAnswer struct {
	AnswerID  int       `json:"answer_id"`
	PollMedia PollMedia `json:"poll_media"`
}

type PollLayoutType int

const (
	DefaultPollLayoutType PollLayoutType = 1
)

type PollResults struct {
	IsFinalized bool              `json:"is_finalized"`
	AnswerCount []PollAnswerCount `json:"answer_count"`
}

type PollAnswerCount struct {
	ID        int  `json:"id"`
	Count     int  `json:"count"`
	IsMeVoted bool `json:"is_me_voted"`
}
