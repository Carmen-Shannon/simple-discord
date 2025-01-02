package bot

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Carmen-Shannon/simple-discord/session"
	"github.com/Carmen-Shannon/simple-discord/structs"
	"github.com/Carmen-Shannon/simple-discord/structs/gateway"
)

type Bot interface {
	GetSession(shardID int) (*session.Session, error)
	GetSessionByGuildID(guildID structs.Snowflake) (*session.Session, error)
	RegisterCommands(commands map[string]session.CommandFunc)
}

type bot struct {
	sessions    map[int]*session.Session
	activeShard int
	mu          sync.RWMutex
}

var _ Bot = (*bot)(nil)

func NewBot(token string, intents []structs.Intent) (Bot, <-chan struct{}, error) {
	initialSession, err := session.NewSession(token, intents, nil)
	if err != nil {
		return nil, nil, err
	}

	b := &bot{
		sessions: make(map[int]*session.Session),
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

func (b *bot) GetSession(shardID int) (*session.Session, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	session, ok := b.sessions[shardID]
	if !ok {
		return nil, fmt.Errorf("no session found for shard ID %d", shardID)
	}

	return session, nil
}

func (b *bot) GetSessionByGuildID(guildID structs.Snowflake) (*session.Session, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, session := range b.sessions {
		if session.GetServerByGuildID(guildID) != nil {
			return session, nil
		}
	}

	return nil, fmt.Errorf("no session found for guild ID %s", guildID.ToString())
}

func (b *bot) RegisterCommands(commands map[string]session.CommandFunc) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Assert the type of commands to match the expected type
	convertedCommands := make(map[string]func(*session.Session, gateway.Payload) error)
	for k, v := range commands {
		convertedCommands[k] = v
	}

	for _, session := range b.sessions {
		session.RegisterCommands(commands)
	}
}

func (b *bot) run(stopChan chan struct{}) error {
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    <-stop
    close(stopChan)
    return nil
}
