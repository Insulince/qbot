package qbot

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// Handle !submitwave <wave>
func (q *QBot) handleSubmitWave(cmd Cmd) error {
	if len(cmd.Args) != 1 {
		q.mustPost(cmd.Message.ChannelID, "Usage: `!submitwave <wave>`")
		return nil
	}

	userId := cmd.Message.Author.ID
	username := cmd.Message.Author.Username
	id := cmd.Message.Author.ID
	wavesStr := cmd.Args[0]
	displayName, err := q.GetDisplayName(cmd.Message)
	if err != nil {
		return errors.Wrap(err, "getting display name")
	}

	waves, err := strconv.Atoi(wavesStr)
	if err != nil {
		return errors.Wrapf(err, "converting waves %q to int", wavesStr)
	}
	if waves < 1 || waves > 100000 {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s> Please enter a valid wave number between 1 and 100000.", id))
		return nil
	}

	tournament, err := q.store.GetLatestTournament()
	if err != nil {
		return errors.Wrap(err, "get latest tournament")
	}

	if err := q.store.InsertTournamentEntry(cmd.GuildId, tournament.Id, userId, username, displayName, waves); err != nil {
		return errors.Wrap(err, "inserting tournament entry")
	}

	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("✅ <@%s> set their waves to **%d**!", id, waves))

	return nil
}
