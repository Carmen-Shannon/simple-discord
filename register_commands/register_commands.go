package main

import (
	"fmt"
	"log"
	"os"

	requestutil "github.com/Carmen-Shannon/simple-discord/gateway/request_util"
	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	"github.com/Carmen-Shannon/simple-discord/util"
	"github.com/joho/godotenv"
)

const guildID = "219182306869379072"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalf("token not found")
	}

	applicationID := os.Getenv("APPLICATION_ID")
	if applicationID == "" {
		log.Fatalf("application id not found")
	}

	testCommand := dto.NewGuildApplicationCommandDto("ping", util.ToPtr(structs.ChatInputCommand))
	testCommand.SetDescription("This is a test command")

	command, err := requestutil.CreateGuildApplicationCommand(testCommand, applicationID, guildID, token)
	if err != nil {
		log.Fatalf("error creating command: %v", err)
	}

	fmt.Printf("Created command: %s", command)
}
