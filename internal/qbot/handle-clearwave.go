package qbot

import (
	"fmt"

	"github.com/pkg/errors"
)

// handleClearWave removes the invoking user's entry from the current tournament leaderboard.
func (q *QBot) handleClearWave(cmd Cmd) error {
	tournament, err := q.store.GetLatestTournament()
	if err != nil {
		return errors.Wrap(err, "get latest tournament")
	}

	userId := cmd.Message.Author.ID

	deleted, err := q.store.DeleteTournamentEntry(tournament.Id, userId)
	if err != nil {
		return errors.Wrap(err, "delete tournament entry")
	}

	if !deleted {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s> You don't have a wave submitted for the current tournament.", userId))
		return nil
	}

	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("✅ <@%s> Your wave entry has been removed from the current tournament.", userId))
	return nil
}
