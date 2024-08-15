package session

type EventHandler struct {
	Handlers map[string]func(*Session, interface{}) error
}