package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Carmen-Shannon/simple-discord/session"
	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalf("token not found")
	}

	intents := []structs.Intent{structs.GuildsIntent, structs.MessageContentIntent, structs.GuildMessagesIntent}

	bot, err := session.NewSession(token, intents)
	if err != nil {
		log.Fatalf("error creating session: %v", err)
		bot.Exit()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	
	<-stop

	bot.Exit()
}
