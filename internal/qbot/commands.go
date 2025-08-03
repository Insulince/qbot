package qbot

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

const (
	exclamationsRegexString = `^!+$`
)

// Define regex to check if the encountered command is only exclamation marks
var exclamationsRegex = regexp.MustCompile(exclamationsRegexString)

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
		cmd, err := q.interpretMessage(m)
		if err != nil {
			return errors.Wrapf(err, "interpreting message")
		}

		if cmd.IsEmpty() {
			// This message does not appear to be a bot command, do nothing.
			return nil
		}

		if messageIsShouting(cmd.Message.Content) {
			q.mustPost(cmd.Message.ChannelID, "I heard you, no need to shout!")
		}

		// Command routing.
		switch cmd.Command {
		case `!queue`, `!q`, `!enqueue`, `!join`:
			return q.handleQueue(cmd)
		case `!enter`, `!enterbracket`, `!e`:
			return q.handleEnter(cmd)
		case `!full`, `!bracketfull`, `!f`:
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
		case `!moretime`, `!extend`, `!moar`:
			return q.handleMoreTime(cmd)
		case `!submitwave`, `!submitwaves`, `!wave`, `!waves`, `!w`:
			return q.handleSubmitWave(cmd)
		case `!forcewave`:
			return q.handleForceWave(cmd)
		case `!leaderboard`, `!lb`:
			return q.handleLeaderboard(cmd, false)
		case `!owned`:
			return q.handleOwned(cmd)
		case `!history`:
			return q.handleHistory(cmd)
		case `!dev`:
			return q.handleDev(cmd)
		default:
			if isNumber(cmd.Command) {
				cmd.Command = cmd.Command[1:]                         // Remove the "!" from the number represented in the command.
				cmd.Args = append([]string{cmd.Command}, cmd.Args...) // Move the number into the args, preserving any other provided args after.
				cmd.Command = "submitwave"                            // Re-set the actual command to be "submitwave".
				return q.handleSubmitWave(cmd)
			}
			q.mustPost(m.ChannelID, fmt.Sprintf("unknown command (use `!help` for available commands): `%s`", cmd.Command))
			return nil
		}
	}()

	if err != nil {
		q.mustPost(m.ChannelID, "Could not process command, an error occurred and has been logged, contact a Q-Dev for assistance.")
		q.reportError(err)
	}
}

// A shouted message is one which both contains at least 1 letter and for each letter contained each one is uppercase.
func messageIsShouting(s string) bool {
	hasLetter := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetter = true
			if unicode.IsLower(r) {
				return false
			}
		}
	}
	return hasLetter
}

// If the string is "!" followed by a string that represents an integer, then return true, else return false.
func isNumber(s string) bool {
	s = s[1:]                 // Ignore the prefixed "!"
	_, err := strconv.Atoi(s) // Check if the rest is an integer.
	if err != nil {
		return false
	}
	return true
}

func (q *QBot) interpretMessage(m *discordgo.MessageCreate) (Cmd, error) {
	content := m.Content

	// Remove spaces from the left of the message.
	content = strings.TrimLeftFunc(content, unicode.IsSpace)

	if strings.HasPrefix(content, "❗") || strings.HasPrefix(content, "❕") {
		content = "!" + content[3:]
	}

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

	guildId := m.GuildID
	if _, err := q.guilds.Get(guildId); err != nil {
		return emptyCmd, errors.Wrap(err, "check guild id")
	}

	var cmd Cmd

	cmd.Message = m
	cmd.Command = command
	cmd.Args = args
	cmd.GuildId = guildId

	return cmd, nil
}
