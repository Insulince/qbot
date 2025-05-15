package qbot

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// handleDev is the developer function for testing various dev things.
func (q *QBot) handleDev(cmd Cmd) error {
	if !q.isModerator(cmd.Message) {
		return nil
	}

	if len(cmd.Args) < 1 {
		q.mustPost(cmd.Message.ChannelID, "❌ Usage: `!dev <cmd> [<args...>]`")
		return nil
	}

	subCmd := strings.ToLower(cmd.Args[0])
	subArgs := cmd.Args[1:]

	cmd.Command = subCmd
	cmd.Args = subArgs

	switch cmd.Command {
	case `err`, `error`:
		return q.handleDevError(cmd)
	case `displayname`:
		return q.handleDevDisplayName(cmd)
	case `parsetourneyname`:
		return q.handleDevParseTourneyName(cmd)
	default:
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("❌ unknown dev sub-command: `%s`", cmd.Command))
		return nil
	}
}

func (q *QBot) handleDevError(cmd Cmd) error {
	msg := strings.Join(cmd.Args, " ")
	err := errors.Errorf("intentional dev error triggered with payload %q", msg)

	return err
}

func (q *QBot) handleDevDisplayName(cmd Cmd) error {
	displayName, err := q.GetDisplayName(cmd.Message)
	if err != nil {
		q.mustPost(cmd.Message.ChannelID, "❌ Error getting display name")
		return errors.Wrap(err, "get display name")
	}

	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("✅ %s's display name is %q!", cmd.Message.Author.Username, displayName))

	return nil
}

func (q *QBot) handleDevParseTourneyName(cmd Cmd) error {
	if len(cmd.Args) != 1 {
		q.mustPost(cmd.Message.ChannelID, "❌ must provide exactly one argument of the tourney short name")
	}
	tourneyShortName := cmd.Args[0]

	tourneyName, err := parseTournamentShortName(tourneyShortName)
	if err != nil {
		q.mustPost(cmd.Message.ChannelID, "❌ Error parsing tourney short name")
		return errors.Wrap(err, "parse tourney short name")
	}

	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("✅ tourney short name %q parsed into %q!", tourneyShortName, tourneyName))

	return nil
}
