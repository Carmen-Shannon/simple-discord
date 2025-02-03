# simple-discord

simple-discord is designed to be a "simple" to use framework for interfacing with the Discord API. It takes advantage of the Gateway that Discord uses for requests back and forth between the simple-discord client and the actual API. Currently, the framework is limited in what it can do and is certainly not close to fully integrating with the Discord API.

## Features

- **Golang Performance**: Utilizes Go's performance by running the sending and receiving of data on different channels. Each event spawns a new Goroutine which handles the event asynchronously, allowing the session managers to handle multiple events concurrently.
- **Custom JSON Unmarshalling**: Populates an active session with easy-to-use Golang structs. Every event and response from the Discord API is mapped to a local Go struct, including binary data with custom marshalling.
- **Active Cache**: Maintains an active cache while the bot is running. (not fully implemented)
    - The bot will attempt to fetch requested data from its local cache that it actively builds while it receives gateway events
    - Each shard contains it's own cache, and the global cache of data can be accessed via the Bot instance
    - Most of the caching is abstracted to the user, so as you develop you shouldn't notice it
    - For now the caching is really non-existent, most of the regular session events are updating the cache appropriately, and it should still be possible to access the cache from within custom handlers but it's not fleshed out fully.
- **Auto Sharding**: Automatically takes advantage of discord sharding based on the needs of the Bot. Makes a bot gateway request to Discord and uses the recommended sharding details to automatically produce gateway sessions with the appropriate shard mapping.
- **Full Integration**: Fully integrated with the Discord API, as well as the gateway's used to manage every event. Most definitions match 1-1 with the Discord developer docs, allowing developers to use the mature documentation as a reference.
- **Audio Playback/Recording**: Capable of recording voice packets sent from the voice gateway/udp connection as well as encrypting voice packets to send to the channel to playback audio.
    - Currently only works with static files, should work with any PCM audio format ffmpeg supports.

## Installation
To run the project, simply install Go v1.23.5 or above, and run the following command to install the latest package into your project:

# NOTE
The ffmpeg binaries are included in this package, as well as the forked gopus package. This should allow devs to build bots in a platform agnostic way, without having to perform additional steps to set up their environment. The ffmpeg binaries WILL be loaded into memory if your bot has the proper intents to handle voice connections and audio playback, and the ffmpeg process is invoked. The binaries will be ignored and never loaded into memory if the ffmpeg process is never needed.

```sh
go get github.com/Carmen-Shannon/simple-discord@latest
```

## Getting Started
Make sure you install the simple-discord project:
```sh
go get github.com/Carmen-Shannon/simple-discord@latest
```

### Import the bot package into your main application, and set up a new `Bot` instance:
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


    version := "1.0.0" // you can just omit this, or pass an empty string as the first argument to `NewBot`
    // make sure your token is labelled TOKEN in your env
    token := os.GetEnv("TOKEN")
    if token == "" {
        log.Fatalf("token not found")
        return
    }
    intents := []structs.Intent{structs.MessageContentIntent}

    // init a new client instance
    client, stopChan, err := bot.NewBot(version, token, intents)
    if err != nil {
        log.Fatalf("error creating session: %v", err)
        return
    }

    // the second returned object is a stop channel, use this channel to block the main function closing until the stopChan is closed
    <-stopChan
}
```

### Setting up custom interaction handlers (for slash commands):
- custom commands need to accept two arguments; *session.Session, and gateway.Payload
- custom commands need to return an error
- with auto-sharding, you can use the *session.Session argument to access the shard receiving the event
- gateway.Payload is used to assert the event is an `Interaction Create Event` from the discord gateway. I haven't implemented a way to make this boilerplate code default behavior yet, but will do that eventually
```go
import (
    "github.com/Carmen-Shannon/simple-discord/bot"
	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/dto"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway/payload"
	receiveevents "github.com/Carmen-Shannon/simple-discord/structs/gateway/receive_events"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway/session"
	"github.com/Carmen-Shannon/simple-discord/util"
)

func main() {

    ...

    // set up a new handler function with the proper arguments and error response
    testCommand := func(sess session.ClientSession, p payload.SessionPayload) error {
        // I set up this ValidateEvent function to be able to decode the payload data into an actual interaction event
        // feel free to use it or use `interactionEvent, ok := p.Data.(receiveevents.InteractionCreateEvent)
        interactionEvent, ok := payload.ValidateEvent[receiveevents.InteractionCreateEvent](p.Data)
        if !ok {
            return fmt.Errorf("could not assert payload.Data to InteractionCreateEvent")
        }

        // for now using the NewInteractionResponseOptions function to create an interaction reply
        response := structs.NewInteractionResponseOptions()
        // setting the response type
        response.SetResponseType(structs.ChannelMessageWithSourceInteraction)
        // setting the flags, use structs.MessageFlag with the built-in Bitfield struct.
        response.SetFlags(structs.Bitfield[structs.MessageFlag]{structs.SurpressNotificationsMessageFlag})

        // using the ClientSession GetServerByGuildID function to grab the `Server` the bot was invoked in
        server := sess.GetServerByGuildID(*interactionEvent.GuildID)
        // using the Server GetVoiceState function to grab the voice channel of the invoking user
        voiceChannel := server.GetVoiceState(interactionEvent.Member.User.ID)
        // making sure the user invoking the bot is in a voice channel before triggering the command
        if voiceChannel == nil || voiceChannel.ChannelID == nil {
            response.SetContent("You must be in a voice channel to use this command")
            // always call Reply
            err := sess.Reply(response, interactionEvent.Interaction)
            if err != nil {
                return fmt.Errorf("could not reply to interaction: %v", err)
            }
            return nil
        }

        // setting the response message content
        response.SetContent("Pong!")
        // always call Reply
        err = sess.Reply(response, interactionEvent.Interaction)
        if err != nil {
            return fmt.Errorf("could not reply to interaction: %v", err)
        }

        // using the ClientSession JoinVoice function to join a voice channel, this is how you init a voice gateway connection
        err = sess.JoinVoice(*interactionEvent.GuildID, *voiceChannel.ChannelID)
        if err != nil {
            return fmt.Errorf("could not join voice channel: %v", err)
        }

        // if you want to track which shard handled the command
        fmt.Println("Command was triggered by shard: ", *sess.GetShard())
        return nil
    }

    // creating a map we can use to pass to the RegisterCommands function
    commands := map[string]session.CommandFunc{
        "hello": testCommand
    }

    // make sure to register the commands with the bot:
    client.RegisterCommands(commands)
}
```

### Setting up custom handlers (for discord gateway events)
- these handlers use pre-set Listeners to listen to the discord gateway events
- you can register custom handlers so the bot will behave a certain way when receiving one of these events
- see session/eventhandler.go for the full list of Listeners, or see the docs for the RegisterListeners function
```go

...

messageListener := func(sess session.Session, payload gateway.Payload) error {
    messageEvent, ok := payload.Data.(receiveevents.MessageCreateEvent)
    if !ok {
        return fmt.Errorf("could not assert payload.Data to MessageCreateEvent")
    }

    // don't do anything if the message comes from the bot
    if messageEvent.Author.ID.Equals(sess.GetBotData().UserDetails.ID) {
        return nil
    }

    // if it's not a message we can interact with just return early
    if messageEvent.Type.Value != 0 {
        return nil
    }

    // make the bot respond to a message with specific starting content, i.e a message that starts with ?
    if messageEvent.Content[0] == '?' {
        msg := dto.NewMessageOptions()
        msg.SetChannelID(messageEvent.ChannelID)
        if err := msg.SetContent("testing a response"); err != nil {
            return fmt.Errorf("could not set content for message: %v", err)
        }
        if err := msg.SetMessageReference(*messageEvent.Message, *messageEvent.GuildID, nil); err != nil {
            return fmt.Errorf("could not set message reference for message: %v", err)
        }

        err := sess.SendMessage(msg)
        if err != nil {
            return fmt.Errorf("could not send message: %v", err)
        }
    }
    return nil
}

...

client.RegisterListeners(map[session.Listener]session.CommandFunc{
    session.MessageCreateListener: messageListener,
})
```

### Registering commands with the Discord API:
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
Latest stable release is `v0.5.1`

## In-Progress

This list will change as I want to add things
- [x] Struct definitions
- [ ] Gateway management
    - [x] Standard gateway management
    - [ ] Voice gateway management
        - [x] Voice Gateway Connection/Upkeep
        - [x] Voice Gateway UDP Connection/Upkeep
            - [ ] Voice Encoding (playing audio)
                - [x] Playing from file (any PCM compatible data such as mp3)
                - [ ] Controlling audio playback state (pausing, resuming)
            - [ ] Voice Decoding (recording audio)
        - [ ] DAVE Voice support
- [x] Event handler
- [x] Shard management
- [ ] HTTP requests
    - [ ] Rate Limiting
        - [ ] Global Rate Limiting
        - [ ] Local Rate Limiting
    - [x] Application requests
    - [x] Application Role Connection Metadata requests
    - [x] Audit Log requests
    - [x] Auto Moderation requests
    - [x] Channel requests
    - [x] Emoji requests
    - [x] Entitlement requests
    - [x] Guild requests
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
    - [x] Registering Custom Handlers
- [ ] Tests
    - [ ] Integration tests
        - [ ] Gateway integration tests
        - [ ] HTTP integration tests
    - [ ] Unit tests