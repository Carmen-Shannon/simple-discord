package bot

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/Carmen-Shannon/simple-discord/session"
	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/util/ffmpeg"
)

type Bot interface {
	GetSession(shardID int) (session.Session, error)
	GetSessionByGuildID(guildID structs.Snowflake) session.Session
	RegisterCommands(commands map[string]session.CommandFunc)
	RegisterListeners(listeners map[session.Listener]session.CommandFunc)
}

type bot struct {
	sessions map[int]session.Session
	mu       sync.RWMutex
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
//   - Bot: The created bot instance.
//   - <-chan struct{}: A channel that will be closed when the bot is stopped.
//   - error: An error if the bot could not be created.
//
// Example:
//
//	bot, stopChan, err := bot.NewBot("your-bot-token", []structs.Intent{structs.GuildsIntent, structs.MessagesIntent})
//	if err != nil {
//	    log.Fatalf("error creating bot: %v", err)
//	}
//	<-stopChan
func NewBot(token string, intents []structs.Intent) (Bot, <-chan struct{}, error) {
	initialSession, err := session.NewSession(token, intents, nil)
	if err != nil {
		return nil, nil, err
	}

	b := &bot{
		sessions: make(map[int]session.Session),
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
		sess, err := session.NewSession(token, intents, &shardID)
		if err != nil {
			return nil, nil, err
		}
		sess.SetShard(&shardID)
		sess.SetShards(&shards)
		b.sessions[shardID] = sess
	}

	stopChan := make(chan struct{})
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
		if err := session.Exit(); err != nil {
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
func (b *bot) GetSession(shardID int) (session.Session, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

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
func (b *bot) GetSessionByGuildID(guildID structs.Snowflake) session.Session {
	b.mu.RLock()
	defer b.mu.RUnlock()

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
	b.mu.RLock()
	defer b.mu.RUnlock()

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
	b.mu.RLock()
	defer b.mu.RUnlock()

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

	// this just doesn't work tbh need to find a solution, seems to be intermittently holding ffmpeg in-memory after it finishes using it
	if ffmpeg.FfprobeCmd != nil && ffmpeg.FfprobeCmd.Process != nil {
		if err := ffmpeg.FfprobeCmd.Process.Signal(syscall.Signal(0)); err == nil {
			fmt.Println("ffprobe process is running")
			if err := ffmpeg.FfprobeCmd.Process.Kill(); err != nil {
				fmt.Printf("error killing ffprobe process: %v\n", err)
			} else {
				fmt.Println("ffprobe process terminated successfully")
			}
		}
	}
	if ffmpeg.FfmpegCmd != nil && ffmpeg.FfmpegCmd.Process != nil {
		if err := ffmpeg.FfmpegCmd.Process.Signal(syscall.Signal(0)); err == nil {
			fmt.Println("ffmpeg process is running")
			if err := ffmpeg.FfmpegCmd.Process.Kill(); err != nil {
				fmt.Printf("error killing ffmpeg process: %v\n", err)
			} else {
				fmt.Println("ffmpeg process terminated successfully")
			}
		}
	}

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

func cleanupFFmpegBinaries() error {
	tmpDir := os.TempDir()
	binaries := []string{"ffmpeg", "ffprobe"}

	for _, binary := range binaries {
		binaryPath := filepath.Join(tmpDir, binary)
		if runtime.GOOS == "windows" {
			binaryPath += ".exe"
		}
		if err := os.Remove(binaryPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove binary %s: %w", binaryPath, err)
		}
	}

	return nil
}
