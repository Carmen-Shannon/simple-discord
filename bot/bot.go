package bot

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway/session"
)

// Bot is the interface for the bot package. It allows users to access the bot's sessions, register commands, and register listeners.
// You should use `NewBot` to create a new bot instance, and then use the session to interact with the Discord API.
// The bot will potentially have multiple sessions, all self-managed. The reason for this is auto-sharding, and how the bot handles sharding sessions upon initialization.
// Auto-sharding is controlled by the get gateway bot response from the Discord API, and if your bot is in a large number of guilds you will have more shards.
// Sharding is generally something you don't need to think about, since when you create a `CommandFunc` or `Listener`, the `ClientSession` used in the function will always be the one handling the interaction.
type Bot interface {
	GetSession(shardID int) (session.ClientSession, error)
	GetSessionByGuildID(guildID structs.Snowflake) session.ClientSession
	RegisterCommands(commands map[string]session.CommandFunc)
	RegisterListeners(listeners map[session.Listener]session.CommandFunc)
}

type bot struct {
	mu *sync.Mutex

	sessions map[int]session.ClientSession
}

var _ Bot = (*bot)(nil)

// NewBot creates a new auto-sharded bot with the given token and intents.
// Auto-sharding happens locally, and the Discord API will provide a recommended number of shards
// based on how many guilds the bot is in. Each shard will be attached to a session, and only that
// session will respond to that particular shard's events. For massive bots in many guilds, you might
// require more CPU to handle the load.
//
// Parameters:
//   - token: The bot token used for authentication with the Discord API.
//   - intents: A slice of intents specifying the intents the bot needs to have in order to function.
//
// Returns:
//   - newBot: The created bot instance.
//   - stopChan: A channel that will be closed when the bot is stopped.
//   - err: An error if the bot could not be created.
//
// Example:
//
//	bot, stopChan, err := bot.NewBot("your-bot-token", []structs.Intent{structs.GuildsIntent, structs.MessagesIntent})
//	if err != nil {
//	    log.Fatalf("error creating bot: %v", err)
//	}
//	<-stopChan
func NewBot(token string, intents []structs.Intent) (newBot Bot, stopChan chan struct{}, err error) {
	b := &bot{
		mu:       &sync.Mutex{},
		sessions: make(map[int]session.ClientSession),
	}

	initialSession := session.NewClientSession()
	initialSession.SetToken(token)
	initialSession.SetIntents(intents...)
	initialSession.SetShard(0)
	initialSession.SetCb(b.reconnectCb)
	if err := initialSession.Dial(true); err != nil {
		return nil, nil, err
	}

	shards := *initialSession.GetShards()
	maxConcurrency := *initialSession.GetMaxConcurrency()

	// Add the initial session to the sessions map
	b.sessions[*initialSession.GetShard()] = initialSession

	for i := 1; i < shards; i++ {
		if i%maxConcurrency == 0 {
			// Cooldown period of 5 seconds after creating maxConcurrency sessions
			time.Sleep(5 * time.Second)
		}

		shardID := i
		sess := session.NewClientSession()
		sess.SetToken(token)
		sess.SetIntents(intents...)
		sess.SetShard(shardID)
		sess.SetShards(shards)
		sess.SetCb(b.reconnectCb)
		if err := sess.Dial(false); err != nil {
			return nil, nil, err
		}

		b.sessions[shardID] = sess
	}

	stopChan = make(chan struct{})
	go func() {
		if err := b.run(stopChan); err != nil {
			fmt.Printf("error running bot: %v\n", err)
			os.Exit(1)
		}
	}()

	return b, stopChan, nil
}

func (b *bot) exit() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, session := range b.sessions {
		if err := session.Exit(true); err != nil {
			return err
		}
		fmt.Println("Shutting down session with shard ID:", *session.GetShard())
	}

	return nil
}

// GetSession returns the session with the given shard ID.
// This is used once you already have the desired shard ID.
// If you need to find the shard ID that contains a specific guild, use GetSessionByGuildID.
// If the session does not exist, an error will be returned.
// This is useful for sending messages, creating channels, and dedicating that work to a specific session based on the shard ID.
//
// Parameters:
//   - shardID: The shard ID of the session you want to retrieve.
//
// Returns:
//   - *session.Session: The session with the given shard ID.
//   - error: An error if the session does not exist.
//
// Example:
//
//	session, err := bot.GetSession(0)
//	if err != nil {
//	    log.Fatalf("error getting session: %v", err)
//	}
//	// Use the session to send messages, create channels, etc...
//	_, err := session.SendMessage(msgOptions, "guild-id", nil)
//	if err != nil {
//	    log.Fatalf("error sending message: %v", err)
//	}
func (b *bot) GetSession(shardID int) (session.ClientSession, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	session, ok := b.sessions[shardID]
	if !ok {
		return nil, fmt.Errorf("no session found for shard ID %d", shardID)
	}

	return session, nil
}

// GetSessionByGuildID returns the session that contains the guild with the given guild ID.
// This is useful for finding which shard is managing which guild.
// If the session does not exist, nil will be returned.
//
// Parameters:
//   - guildID: The guild ID of the guild you want to find the session for.
//
// Returns:
//   - *session.Session: The session that contains the guild with the given guild ID.
//
// Example:
//
//	session := bot.GetSessionByGuildID("guild-id")
//	if session == nil {
//	    log.Fatalf("session not found for guild ID")
//	}
func (b *bot) GetSessionByGuildID(guildID structs.Snowflake) session.ClientSession {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, session := range b.sessions {

		if session.GetServerByGuildID(guildID) != nil {
			return session
		}
	}

	return nil
}

// RegisterCommands allows you to register any number of custom named handlers.
// Use this to handle interaction events, or slash commands.
// The commands map should be in the following format:
//   - map[string]session.CommandFunc
//
// The key is the name of the command, and the value is the function that will be called when the command is triggered.
//
// Parameters:
//   - commands: A map of command names and their corresponding functions.
//
// Example:
//
//	testCommand := func(sess *session.Session, payload gateway.Payload) error {
//		// make the session do stuff with the payload data from the interaction
//		if interactionEvent, ok := payload.Data.(receiveevents.InteractionCreateEvent); ok {
//			response := structs.NewInteractionResponseOptions()
//			response.SetResponseType(structs.ChannelMessageWithSourceInteraction)
//			response.SetFlags(structs.Bitfield[structs.MessageFlag]{structs.SurpressNotificationsMessageFlag})
//			response.SetContent("Pong!")
//			if err := sess.InteractionReply(response, interactionEvent.Interaction); err != nil {
//				return fmt.Errorf("could not reply to interaction: %v", err)
//			}
//		}
//		return nil
//	}
//	commands := map[string]session.CommandFunc{
//	    "ping": testCommand,
//	}
//	bot.RegisterCommands(commands)
func (b *bot) RegisterCommands(commands map[string]session.CommandFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, session := range b.sessions {
		session.RegisterCommands(commands)
	}
}

// RegisterListeners registers listeners for all sessions.
// Use this to handle events such as message creation, message deletion, etc...
// The listeners map should be in the following format:
//   - map[session.Listener]session.CommandFunc
//
// The key is the listener you want to listen for, and the value is the function that will be called when the listener is triggered.
//
// Parameters:
//   - listeners: A map of listeners and their corresponding functions.
//
// Example:
//
//	listeners := map[session.Listener]session.CommandFunc{
//	    session.MessageCreateListener: func(sess *session.Session, payload gateway.Payload) error {
//	        // make the session do stuff with the payload data from the event
//	        return nil
//	    },
//	}
//	bot.RegisterListeners(listeners)
//
// The available listeners are in the following format:
//   - HelloListener = "HELLO"
//   - ReadyListener = "READY"
//   - ResumedListener = "RESUMED"
//   - ReconnectListener = "RECONNECT"
//   - InvalidSessionListener = "INVALID_SESSION"
//   - ChannelCreateListener = "CHANNEL_CREATE"
//   - ChannelUpdateListener = "CHANNEL_UPDATE"
//   - ChannelDeleteListener = "CHANNEL_DELETE"
//   - GuildCreateListener = "GUILD_CREATE"
//   - GuildUpdateListener = "GUILD_UPDATE"
//   - GuildDeleteListener = "GUILD_DELETE"
//   - GuildBanAddListener = "GUILD_BAN_ADD"
//   - GuildBanRemoveListener = "GUILD_BAN_REMOVE"
//   - GuildEmojisUpdateListener = "GUILD_EMOJIS_UPDATE"
//   - GuildIntegrationsUpdateListener = "GUILD_INTEGRATIONS_UPDATE"
//   - GuildAuditLogEntryCreateListener = "GUILD_AUDIT_LOG_ENTRY_CREATE"
//   - GuildMemberAddListener = "GUILD_MEMBER_ADD"
//   - GuildMemberRemoveListener = "GUILD_MEMBER_REMOVE"
//   - GuildMemberUpdateListener = "GUILD_MEMBER_UPDATE"
//   - GuildMembersChunkListener = "GUILD_MEMBERS_CHUNK"
//   - GuildRoleCreateListener = "GUILD_ROLE_CREATE"
//   - GuildRoleUpdateListener = "GUILD_ROLE_UPDATE"
//   - GuildRoleDeleteListener = "GUILD_ROLE_DELETE"
//   - MessageCreateListener = "MESSAGE_CREATE"
//   - MessageUpdateListener = "MESSAGE_UPDATE"
//   - MessageDeleteListener = "MESSAGE_DELETE"
//   - MessageBulkDeleteListener = "MESSAGE_BULK_DELETE"
//   - MessageReactionAddListener = "MESSAGE_REACTION_ADD"
//   - MessageReactionRemoveListener = "MESSAGE_REACTION_REMOVE"
//   - MessageReactionRemoveAllListener = "MESSAGE_REACTION_REMOVE_ALL"
//   - MessageReactionRemoveEmojiListener = "MESSAGE_REACTION_REMOVE_EMOJI"
//   - MessagePollVoteAddListener = "MESSAGE_POLL_VOTE_ADD"
//   - MessagePollVoteRemoveListener = "MESSAGE_POLL_VOTE_REMOVE"
//   - TypingStartListener = "TYPING_START"
//   - UserUpdateListener = "USER_UPDATE"
//   - VoiceChannelEffectSendListener = "VOICE_CHANNEL_EFFECT_SEND"
//   - VoiceStateUpdateListener = "VOICE_STATE_UPDATE"
//   - VoiceServerUpdateListener = "VOICE_SERVER_UPDATE"
//   - VoiceChannelStatusUpdateListener = "VOICE_CHANNEL_STATUS_UPDATE"
//   - WebhooksUpdateListener = "WEBHOOKS_UPDATE"
//   - PresenceUpdateListener = "PRESENCE_UPDATE"
func (b *bot) RegisterListeners(listeners map[session.Listener]session.CommandFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, session := range b.sessions {
		session.RegisterListeners(listeners)
	}
}

func (b *bot) run(stopChan chan struct{}) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Capture panics and close stopChan
	go func() {
		if r := recover(); r != nil {
			fmt.Printf("panic occurred: %v\n", r)
			close(stopChan)
		}
	}()

	<-stop

	// Wait for processes to terminate
	time.Sleep(2 * time.Second)

	// Cleanup temporary ffmpeg binaries
	cleanupFFmpegBinaries()

	if err := b.exit(); err != nil {
		close(stopChan)
		return err
	}

	close(stopChan)
	return nil
}

func (b *bot) reconnectCb(sess session.ClientSession) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if sess.GetShard() == nil {
		return fmt.Errorf("session shard ID is nil")
	}

	b.sessions[*sess.GetShard()] = sess
	return nil
}

func cleanupFFmpegBinaries() error {
	tmpDir := os.TempDir()
	pattern := filepath.Join(tmpDir, "ffmpeg*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		os.Remove(match)
	}

	return nil
}
