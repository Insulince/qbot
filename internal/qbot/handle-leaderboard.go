package qbot

import (
	"fmt"
	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/pkg/errors"
	"strings"
)

const (
	ChannelIdTournamentQueue = "1343112046404833351" // #tournament-queue
)

// Handle !leaderboard
func (q *QBot) handleLeaderboard(cmd Cmd) error {
	channelId := ChannelIdTournamentQueue
	if cmd.Message != nil {
		channelId = cmd.Message.ChannelID
	}

	const fetchTournamentEntriesSql = `
SELECT
    user_id,
    waves
FROM tournament_entries
WHERE tournament_id = (SELECT MAX(id) FROM tournaments)
ORDER BY waves DESC;
`
	rows, err := q.db.Query(fetchTournamentEntriesSql)
	if err != nil {
		q.mustPost(channelId, "Error retrieving leaderboard.")
		return errors.Wrap(err, "querying tournaments")
	}
	defer jmust.MustClose(rows)

	leaderboardMsg := fmt.Sprintf("🏆 **Latest Tournament Leaderboard** 🏆\n")
	var entries []string
	for i := 1; rows.Next(); i++ {
		var userId string
		var waves int
		if err := rows.Scan(&userId, &waves); err != nil {
			return errors.Wrap(err, "scanning fields")
		}
		entry := fmt.Sprintf("%d. **<@%s>** - Wave %d", i, userId, waves)
		entries = append(entries, entry)
	}

	entriesMsg := "_No entries yet._"
	if len(entries) > 0 {
		entriesMsg = strings.Join(entries, "\n")
	}
	leaderboardMsg += entriesMsg

	q.mustPostWithoutTags(channelId, leaderboardMsg)

	return nil
}
