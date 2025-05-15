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

	userID := cmd.Message.Author.ID
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
	if waves < 1 || waves > 10000 {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s> Please enter a valid wave number between 1 and 10000.", id))
		return nil
	}

	const fetchLatestTournamentIdSql = `
SELECT
    MAX(id)
FROM tournaments;
`
	var tournamentId int
	if err := q.db.QueryRow(fetchLatestTournamentIdSql).Scan(&tournamentId); err != nil {
		q.mustPost(cmd.Message.ChannelID, "Error retrieving leaderboard.")
		return errors.Wrap(err, "query row")
	}

	// Update or insert new high score
	insertWaveSql := `
INSERT INTO tournament_entries
    (tournament_id, user_id, username, waves, display_name)
VALUES
    (?, ?, ?, ?, ?)
ON CONFLICT (tournament_id, user_id) DO UPDATE SET waves = excluded.waves;
`
	_, err = q.db.Exec(insertWaveSql, tournamentId, userID, username, waves, displayName)
	if err != nil {
		q.mustPost(cmd.Message.ChannelID, "Error saving your waves.")
		return errors.Wrap(err, "exec query")
	}
	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("✅ <@%s> set their waves to **%d**!", id, waves))

	return nil
}
