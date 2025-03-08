package qbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"strconv"
)

// Handle !submitwave <wave>
func (q *QBot) handleSubmitWave(m *discordgo.MessageCreate, args []string) error {
	if len(args) != 2 {
		q.mustPost(m.ChannelID, "Usage: `!submitwave <wave>`")
		return nil
	}

	userID := m.Author.ID
	username := m.Author.Username
	id := m.Author.ID

	waves, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.Wrap(err, "converting wave to int")
	}
	if waves < 1 || waves > 10000 {
		q.mustPost(m.ChannelID, fmt.Sprintf("<@%s> Please enter a valid wave number between 1 and 10000.", id))
		return nil
	}

	const fetchLatestTournamentIdSql = `
SELECT
    MAX(id)
FROM tournaments;
`
	var tournamentId int
	if err := q.db.QueryRow(fetchLatestTournamentIdSql).Scan(&tournamentId); err != nil {
		q.mustPost(m.ChannelID, "Error retrieving leaderboard.")
		return errors.Wrap(err, "query row")
	}

	// Update or insert new high score
	insertWaveSql := `
INSERT INTO tournament_entries
    (tournament_id, user_id, username, waves)
VALUES
    (?, ?, ?, ?)
ON CONFLICT (tournament_id, user_id) DO UPDATE SET waves = excluded.waves;
`
	_, err = q.db.Exec(insertWaveSql, tournamentId, userID, username, waves)
	if err != nil {
		q.mustPost(m.ChannelID, "Error saving your waves.")
		return errors.Wrap(err, "exec query")
	}
	q.mustPost(m.ChannelID, fmt.Sprintf("✅ <@%s> set their waves to **%d**!", id, waves))

	return nil
}
