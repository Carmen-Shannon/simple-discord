package structs

import "encoding/json"

type ActivityType int

const (
	GameActivity      ActivityType = 0
	StreamingActivity ActivityType = 1
	ListeningActivity ActivityType = 2
	WatchingActivity  ActivityType = 3
	CustomActivity    ActivityType = 4
	CompetingActivity ActivityType = 5
)

type ActivityFlag int64

const (
	ActivityInstance            ActivityFlag = 1 << 0
	ActivityJoin                ActivityFlag = 1 << 1
	ActivitySpectate            ActivityFlag = 1 << 2
	ActivityJoinRequest         ActivityFlag = 1 << 3
	ActivitySync                ActivityFlag = 1 << 4
	ActivityPlay                ActivityFlag = 1 << 5
	ActivityPartyPrivacyFriends ActivityFlag = 1 << 6
	ActivityPartyPrivacyVoice   ActivityFlag = 1 << 7
	ActivityEmbedded            ActivityFlag = 1 << 8
)

type Activity struct {
	Name          string                  `json:"name"`
	Type          ActivityType            `json:"type"`
	URL           *string                 `json:"url"`
	CreatedAt     int                     `json:"created_at"`
	Timestamps    *ActivityTimestamps     `json:"timestamps"`
	ApplicationId *Snowflake              `json:"application_id"`
	Details       *string                 `json:"details"`
	State         *string                 `json:"state"`
	Emoji         *ActivityEmoji          `json:"emoji"`
	Party         *ActivityParty          `json:"party"`
	Assets        *ActivityAssets         `json:"assets"`
	Secrets       *ActivitySecrets        `json:"secrets"`
	Instance      *bool                   `json:"instance"`
	Flags         *Bitfield[ActivityFlag] `json:"flags"`
	Buttons       []Button                `json:"buttons"`
}

type ActivityTimestamps struct {
	Start *int `json:"start"`
	End   *int `json:"end"`
}

type ActivityEmoji struct {
	Name     string     `json:"name"`
	ID       *Snowflake `json:"id"`
	Animated *bool      `json:"animated"`
}

type ActivityParty struct {
	ID   *string `json:"id"`
	Size *[]int  `json:"size"`
}

type ActivityAssets struct {
	LargeImage *string `json:"large_image"`
	LargeText  *string `json:"large_text"`
	SmallImage *string `json:"small_image"`
	SmallText  *string `json:"small_text"`
}

type ActivitySecrets struct {
	Join     *string `json:"join"`
	Spectate *string `json:"spectate"`
	Match    *string `json:"match"`
}

type Button struct {
	Label        string  `json:"label"`
	ReceiveLabel string  `json:"-"`
	URL          *string `json:"url,omitempty"`
}

func (b *Button) UnmarshalJSON(data []byte) error {
	// Try to unmarshal data as a string
	var label string
	if err := json.Unmarshal(data, &label); err == nil {
		b.ReceiveLabel = label
		b.Label = label
		return nil
	}

	// Try to unmarshal data as an object
	type Alias Button
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}
