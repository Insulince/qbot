package qbot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// QueueItem represents an entry in the queue.
type QueueItem struct {
	UserID    string
	AddedAt   time.Time // timestamp for the current phase (enter or full)
	ChannelID string    // channel the user joined from
	Entered   bool      // whether the user has signaled they've joined a bracket
	Warned    bool      // whether a warning has been sent for the current phase
}

// QBot struct holds the Discord session and state.
type QBot struct {
	session       *discordgo.Session
	queue         []QueueItem
	queueMutex    sync.Mutex
	currentUser   *QueueItem
	enterTimeout  time.Duration
	fullTimeout   time.Duration
	warnThreshold time.Duration
}

// NewQBot initializes and returns a new bot instance.
func NewQBot() (*QBot, error) {
	token, found := os.LookupEnv("DISCORD_BOT_TOKEN")
	if !found {
		return nil, fmt.Errorf("DISCORD_BOT_TOKEN not set")
	}

	// Create a new Discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	q := &QBot{
		session:       dg,
		enterTimeout:  5 * time.Minute,
		fullTimeout:   15 * time.Minute,
		warnThreshold: 2 * time.Minute,
	}

	// Set intents
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	// Register handlers
	dg.AddHandler(q.messageHandler)

	// Open connection
	if err := dg.Open(); err != nil {
		return nil, fmt.Errorf("error opening connection: %w", err)
	}

	fmt.Println("QBot Q is now running.")
	go q.timeoutChecker(dg)

	return q, nil
}

// Close shuts down the bot gracefully.
func (q *QBot) Close() {
	q.session.Close()
}

// messageHandler routes commands to the proper handlers.
func (q *QBot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Log every message for debugging.
	log.Printf("[%s] %s: %s\n", m.ChannelID, m.Author.Username, m.Content)

	// Ignore messages from the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the message directly mentions the bot.
	for _, user := range m.Mentions {
		if user.ID == s.State.User.ID {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I can't respond to direct mentions, below is my `!help` message"))
			q.handleHelp(s, m)
			return
		}
	}

	content := m.Content

	// Command routing.
	switch {
	case content == "!queue":
		q.handleQueue(s, m)
	case content == "!enter":
		q.handleEnter(s, m)
	case content == "!full":
		q.handleFull(s, m)
	case content == "!view":
		q.handleView(s, m)
	case content == "!leave":
		q.handleLeave(s, m)
	case content == "!position":
		q.handlePosition(s, m)
	case content == "!help":
		q.handleHelp(s, m)
	case content == "!commands":
		q.handleCommands(s, m)
	case content == "!skip":
		q.handleSkip(s, m)
	case content == "!reset":
		q.handleReset(s, m)
	case strings.HasPrefix(content, "!remove"):
		q.handleRemove(s, m)
	case content == "!moretime":
		q.handleMoreTime(s, m)
	case content == "!version":
		q.handleVersion(s, m)
	}
}
