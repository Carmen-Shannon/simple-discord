package structs

type Emoji struct {
	ID            *Snowflake  `json:"id,omitempty"`
	Name          *string     `json:"name,omitempty"`
	Roles         []Snowflake `json:"roles,omitempty"`
	User          *User       `json:"user,omitempty"`
	RequireColons *bool       `json:"require_colons,omitempty"`
	Managed       *bool       `json:"managed,omitempty"`
	Animated      *bool       `json:"animated,omitempty"`
	Available     *bool       `json:"available,omitempty"`
}

type EmojiAnimationType int

const (
	PremiumEmojiAnimationType EmojiAnimationType = 1
	BasicEmojiAnimationType   EmojiAnimationType = 2
)
