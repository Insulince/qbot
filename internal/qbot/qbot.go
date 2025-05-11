package qbot

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/pkg/errors"

	"github.com/bwmarrin/discordgo"
)

const (
	sqliteDriver = "sqlite"
	databaseFile = "/var/lib/litefs/qbot.db"

	exclamationsRegexString = `^!+$`
)

// Define regex to check if the encountered command is only exclamation marks
var exclamationsRegex = regexp.MustCompile(exclamationsRegexString)

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
	session *discordgo.Session
	db      *sql.DB

	queue         []QueueItem
	queueMutex    sync.Mutex
	currentUser   *QueueItem
	enterTimeout  time.Duration
	fullTimeout   time.Duration
	warnThreshold time.Duration
}

// New initializes and returns a new QBot instance.
func New(token string) (*QBot, error) {
	q := &QBot{
		enterTimeout:  5 * time.Minute,
		fullTimeout:   30 * time.Minute,
		warnThreshold: 2 * time.Minute,
	}

	// Requisition a database connection for use in the command.
	db, err := sql.Open(sqliteDriver, databaseFile)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}
	q.db = db

	// Start a new session.
	if err := q.newSession(token); err != nil {
		return nil, errors.Wrap(err, "new session")
	}

	// Start the timeout checker.
	q.Go(q.timeoutChecker)

	// Start the scheduler.
	q.Go(q.startScheduler)

	defer func() {
		if v := recover(); v != nil {
			q.reportError(fmt.Errorf("panic: %v", v))
			panic(v)
		}
	}()

	fmt.Println("Q is ready")

	return q, nil
}

func (q *QBot) newSession(token string) error {
	// Create a new Discord session
	auth := fmt.Sprintf("Bot %s", token)
	session, err := discordgo.New(auth)
	if err != nil {
		return errors.Wrap(err, "error creating Discord session")
	}

	// Set intents
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	// Register handlers
	session.AddHandler(q.messageHandler)

	// Open connection
	if err := session.Open(); err != nil {
		return errors.Wrap(err, "opening session")
	}

	// Record session on q
	q.session = session

	return nil
}

// Close shuts down the bot.
func (q *QBot) Close() {
	jmust.MustClose(q.session)

	defer jmust.MustClose(q.db)
}

// post splits long messages into chunks and sends them serially.
func (q *QBot) post(channelId, msg string) error {
	const chunkSize = 2000 // Max Discord message length

	// Split message into chunks of at most 2000 characters
	for len(msg) > 0 {
		// Determine the length of the next chunk
		end := chunkSize
		if len(msg) < chunkSize {
			end = len(msg)
		}

		// Send the chunk
		if _, err := q.session.ChannelMessageSend(channelId, msg[:end]); err != nil {
			return errors.Wrapf(err, "failed to send message to channel %q", channelId)
		}

		// Move to the next chunk
		msg = msg[end:]
	}

	return nil
}

func (q *QBot) mustPost(channelId, msg string) {
	jmust.Must[any](q.post, channelId, msg)
}

func (q *QBot) postWithoutTags(channelId, msg string) error {
	const chunkSize = 2000 // Max Discord message length

	// Split message into chunks of at most 2000 characters
	for len(msg) > 0 {
		// Determine the length of the next chunk
		end := chunkSize
		if len(msg) < chunkSize {
			end = len(msg)
		}

		// Send the chunk
		_, err := q.session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
			Content: msg[:end],
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{}, // Prevents pinging
			},
		})
		if err != nil {
			return errors.Wrapf(err, "failed to send message to channel %q", channelId)
		}

		// Move to the next chunk
		msg = msg[end:]
	}

	return nil
}

func (q *QBot) mustPostWithoutTags(channelId, msg string) {
	jmust.Must[any](q.postWithoutTags, channelId, msg)
}

type Cmd struct {
	Message *discordgo.MessageCreate
	Command string
	Args    []string
}

var emptyCmd = Cmd{}

func (cmd Cmd) String() string {
	return fmt.Sprintf("%s %s", cmd.Command, strings.Join(cmd.Args, " "))
}

func (cmd Cmd) Empty() bool {
	return cmd.Command == ""
}

func interpretMessage(m *discordgo.MessageCreate) (Cmd, error) {
	content := m.Content

	// Remove spaces from the left of the message.
	content = strings.TrimLeftFunc(content, unicode.IsSpace)

	// Check that the message begins with an exclamation mark.
	if !strings.HasPrefix(content, "!") {
		return emptyCmd, nil
	}

	// The message appears to be a bot command, so lets split it up by space to get the command and args.
	segments := strings.Split(content, " ")
	if len(segments) == 0 {
		return Cmd{}, errors.Errorf("message %q does not contain any space-delimitted segments", content)
	}

	command := strings.ToLower(segments[0])
	args := segments[1:]

	if exclamationsRegex.MatchString(command) {
		// This command is composed of only exclamation marks which is not a bot command, do nothing.
		return emptyCmd, nil
	}

	var cmd Cmd

	cmd.Message = m
	cmd.Command = command
	cmd.Args = args

	return cmd, nil
}

// messageHandler routes commands to the proper handlers.
func (q *QBot) messageHandler(_ *discordgo.Session, m *discordgo.MessageCreate) {
	// NOTE(justin): We explicitly ignore the discord session because we already have one in q. It is required as part
	// of the contract to qualify this function as a MessageHandler to register with Discord's API. It is for this same
	// reason that we cannot make this function return an error and have to do this silly anonymous function stuff
	// in the following code.

	err := func() error {
		// Log every message for debugging.
		log.Printf("[%s] %s: %s\n", m.ChannelID, m.Author.Username, m.Content)

		// Ignore messages from the bot itself.
		if m.Author.ID == q.session.State.User.ID {
			return nil
		}

		// Check if the message directly mentions the bot.
		for _, user := range m.Mentions {
			if user.ID == q.session.State.User.ID {
				q.mustPost(m.ChannelID, fmt.Sprintf("I can't respond to direct mentions, use `!help` for usage details"))
				return nil
			}
		}

		// Extract command and args
		cmd, err := interpretMessage(m)
		if err != nil {
			return errors.Wrapf(err, "interpreting message")
		}

		if cmd.Empty() {
			// This message does not appear to be a bot command, do nothing.
			return nil
		}

		if strings.ToUpper(cmd.Message.Content) == cmd.Message.Content {
			q.mustPost(cmd.Message.ChannelID, "I heard you, no need to shout!")
		}

		// Command routing.
		switch cmd.Command {
		case `!queue`, `!q`, `!enqueue`, `!join`:
			return q.handleQueue(cmd)
		case `!enter`, `!enterbracket`:
			return q.handleEnter(cmd)
		case `!full`, `!bracketfull`:
			return q.handleFull(cmd)
		case `!view`, `!viewqueue`:
			return q.handleView(cmd)
		case `!leave`, `!leavequeue`:
			return q.handleLeave(cmd)
		case `!position`, `!currentposition`:
			return q.handlePosition(cmd)
		case `!help`:
			return q.handleHelp(cmd)
		case `!commands`:
			return q.handleCommands(cmd)
		case `!version`:
			return q.handleVersion(cmd)
		case `!skip`, `!skipcurrent`:
			return q.handleSkip(cmd)
		case `!reset`, `!resetqueue`:
			return q.handleReset(cmd)
		case `!remove`, `!removeplayer`:
			return q.handleRemove(cmd)
		case `!moretime`, `!extend`:
			return q.handleMoreTime(cmd)
		case `!submitwave`, `!submitwaves`, `!wave`, `!waves`:
			return q.handleSubmitWave(cmd)
		case `!leaderboard`, `!lb`:
			return q.handleLeaderboard(cmd, false)
		case `!history`:
			return q.handleHistory(cmd)
		case `!deverror`:
			return q.handleDevError(cmd)
		default:
			q.mustPost(m.ChannelID, fmt.Sprintf("unknown command (use `!help` for available commands): `%s`", cmd.Command))
			return nil
		}
	}()

	if err != nil {
		q.reportError(err)
	}
}
