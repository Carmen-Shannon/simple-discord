# simple-discord

simple-discord is designed to be a "simple" to use framework for interfacing with the Discord API. It takes advantage of the Gateway that Discord uses for requests back and forth between the simple-discord client and the actual API. Currently, the framework is limited in what it can do and is certainly not close to fully integrating with the Discord API.

## Features

- **Golang Performance**: Utilizes Go's performance by running the sending and receiving of data on different channels.
- **Custom JSON Unmarshalling**: Populates an active session with easy-to-use Golang structs.
- **Active Cache**: Maintains an active cache while the bot is running.
- **Auto Sharding**: Automatically takes advantage of discord sharding based on the needs of the Bot

## Installation

To run the project, simply install Go v1.22.3 or above, and run the following command to install the latest package into your project:

```sh
go get github.com/Carmen-Shannon/simple-discord@latest
```

## Getting Started
Make sure you install the simple-discord project:
```go get github.com/Carmen-Shannon/simple-discord@latest```

Import the bot package into your main application, and set up a new `Bot` instance:
```go
import (
    "github.com/Carmen-Shannon/simple-discord/bot"
)

func main() {
    // load your token as an env variable
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("error loading .env")
        return
    }

    // make sure your token is labelled TOKEN in your env
    token := os.GetEnv("TOKEN")
    if token == "" {
        log.Fatalf("token not found")
        return
    }

    // init a new client instance
    client, stopChan, err := bot.NewBot(token, intents)
    if err != nil {
        log.Fatalf("error creating session: %v", err)
        return
    }

    // the second returned object is a stop channel, use this channel to block the main function closing until the stopChan is closed
    <-stopChan
}
```

Setting up custom interaction handlers (for slash commands):
- custom commands need to accept two arguments; *session.Session, and gateway.Payload
- custom commands need to return an error
- with auto-sharding, you can use the *session.Session argument to access the shard receiving the event
- gateway.Payload is used to assert the event is an `Interaction Create Event` from the discord gateway. I haven't implemented a way to make this boilerplate code default behavior yet, but will do that eventually
```go
import (
    "github.com/Carmen-Shannon/simple-discord/bot"
    "github.com/Carmen-Shannon/simple-discord/session"
    "github.com/Carmen-Shannon/simple-discord/structs"
    "github.com/Carmen-Shannon/simple-discord/structs/gateway"
    receiveevents "github.com/Carmen-Shannon/simple-discord/gateway/receive_events"
)

func main() {
    ...
    // set up a new handler function with the proper arguments and error response
    testCommand := func(sess *session.Session, payload gateway.Payload) error {
        // this is boilerplate, you don't need to include it but it does allow you to use the interactionEvent directly and access all of the associated properties
        interactionEvent, ok := payload.Data.(receieveevents.InteractionCreateEvent)
        if !ok {
            return fmt.Errorf("could not assert payload.Data to InteractionCreateEvent")
        }

        // the bot needs to reply and ACK an interaction, so you need to call sess.Reply() at some point before returning.
        // discord requires an ACK within 3 seconds, but you are able to edit the message by sending follow-up HTTP requests
        // use NewInteractionResponseOptions to create the response that Reply needs
        response := structs.NewInteractionResponseOptions()
        // see strucs.InteractionResponseType for available types
        response.SetResponseType(structs.ChannelMessageWithSourceInteraction)
        // see structs.MessageFlag for available message flags
        response.SetFlags(structs.Bitfield[structs.MessageFlag]{structs.SurpressNotificationsMessageFlag})

        response.SetContent("Hello world!")
        if err := sess.Reply(response, interactionEvent.Interaction); err != nil {
            return fmt.Errorf("could not reply to interaction: %v", err)
        }

        // return once the command has finished
        return nil
    }

    commands := map[string]session.CommandFunc{
        "hello": testCommand
    }

    // make sure to register the commands with the bot:
    client.RegisterCommands(commands)
}
```

Registering commands with the Discord API:
- registering commands requires at least the Application ID of your bot which can be found in your developer portal
- you can register global commands with just an Application ID but there is a delay of about 15 minutes from registering the command until the bot will have access to it
- to test commands in a local server, you can register Guild commands using the Application ID of the bot and the Guild ID of the server you are using to test
- see dto.CreateGuildApplicationCommandDto for properties of the command
```go
import(
    requestutil "github.com/Carmen-Shannon/simple-discord/gateway/request_util"
    "github.com/Carmen-Shannon/simple-discord/structs"
    "github.com/Carmen-Shannon/simple-discord/structs/dto"
    "github.com/Carmen-Shannon/simple-discord/util"
)

func main() {
    token := "12345"
    applicationID := "12345"
    guildID := "12345"

    // call NewGuildApplicationCommandDto with 2 arguments, first is the name of the command second is the ApplicationCommandType of the command
    testCommand := dto.NewGuildApplicationCommandDto("hello", util.ToPtr(structs.ChatInputCommand))
    testCommand.SetDescription("this is a test command - hello")

    _, err := requestutil.CreateGuildApplicationCommand(testCommand, applicationID, guildID, token)
    if err != nil {
        log.Fatalf("error creating command: %v", err)
    }
}
```

## Version
N/A no release yet, v0.1.0 will be the first release

## In-Progress

This list will change as I want to add things
- [x] Struct definitions
- [ ] Gateway management
    - [x] Standard gateway management
    - [ ] Voice gateway management
        - [x] Voice Gateway Connection/Upkeep
        - [ ] Voice Gateway UDP Connection/Upkeep
        - [ ] DAVE Voice support
- [x] Event handler
- [ ] Shard management
- [ ] HTTP requests
    - [x] Application requests
    - [x] Application Role Connection Metadata requests
    - [x] Audit Log requests
    - [x] Auto Moderation requests
    - [x] Channel requests
    - [ ] Emoji requests
    - [ ] Guild requests
    - [ ] Guild Scheduled Event requests
    - [ ] Guild Template requests
    - [x] Interaction Requests
    - [ ] Invite requests
    - [x] Message requests
    - [ ] Poll requests
    - [ ] Stage Instance requests
    - [ ] Sticker requests
    - [ ] User requests
    - [ ] Voice requests
    - [ ] Webhook requests
- [X] Registering Custom Commands
    - [x] Registering Global Commands
    - [x] Registering Guild Commands
- [ ] Tests
    - [ ] Integration tests
        - [ ] Gateway integration tests
        - [ ] HTTP integration tests
    - [ ] Unit tests