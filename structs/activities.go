package structs

type ActivityType int

const (
	GameActivity      ActivityType = 0
	StreamingActivity ActivityType = 1
	ListeningActivity ActivityType = 2
	WatchingActivity  ActivityType = 3
	CustomActivity    ActivityType = 4
	CompetingActivity ActivityType = 5
)

type ActivityFlag int

const ()

type Activity struct {
	Name          string       `json:"name"`
	Type          ActivityType `json:"type"`
	URL           *string      `json:"url"`
	CreatedAt     int          `json:"created_at"`
	Timestamps    *Timestamps  `json:"timestamps"`
	ApplicationId *Snowflake   `json:"application_id"` // TODO this should actually be a Snowflake, which is just a string representation of a bigint
	Details       *string      `json:"details"`
	State         *string      `json:"state"`
	Emoji         *Emoji       `json:"emoji"`
	Party         *Party       `json:"party"`
	Assets        *Assets      `json:"assets"`
	Secrets       *Secrets     `json:"secrets"`
	Instance      *bool        `json:"instance"`
	Flags         *int         `json:"flags"` // TODO this is something different than just an integer
	Buttons       *[]Button    `json:"buttons"`
}

type Timestamps struct {
	Start *int `json:"start"`
	End   *int `json:"end"`
}

type Emoji struct {
	Name     string     `json:"name"`
	ID       *Snowflake `json:"id"` // TODO this should actually be a Snowflake, which is just a string representation of a bigint
	Animated *bool      `json:"animated"`
}

type Party struct {
	ID   *string `json:"id"`
	Size *[]int  `json:"size"`
}

type Assets struct {
	LargeImage *string `json:"large_image"`
	LargeText  *string `json:"large_text"`
	SmallImage *string `json:"small_image"`
	SmallText  *string `json:"small_text"`
}

type Secrets struct {
	Join     *string `json:"join"`
	Spectate *string `json:"spectate"`
	Match    *string `json:"match"`
}

type Button struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}
