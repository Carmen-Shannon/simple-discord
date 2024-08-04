package structs

type GatewayPresenceUpdate struct {
	Since      int        `json:"since"`
	Activities []Activity `json:"activities"`
	Status     string     `json:"status"`
	Afk        bool       `json:"afk"`
}
