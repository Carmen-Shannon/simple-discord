package structs

type Bot struct {
	UserDetails        *User
	ApplicationDetails *Application
}

func NewBot(user User, application Application) *Bot {
	return &Bot{
		UserDetails:        &user,
		ApplicationDetails: &application,
	}
}
