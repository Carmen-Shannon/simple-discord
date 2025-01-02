package structs

type BotData struct {
	UserDetails        *User
	ApplicationDetails *Application
}

func NewBotData(user User, application Application) *BotData {
	return &BotData{
		UserDetails:        &user,
		ApplicationDetails: &application,
	}
}
