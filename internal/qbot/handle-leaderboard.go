package qbot

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// Handle !leaderboard
func handleLeaderboard(session *discordgo.Session, message *discordgo.MessageCreate, db *sql.DB) {
	channelId := "1343112046404833351"
	if message != nil {
		channelId = message.ChannelID
	}

	const fetchTournamentEntriesSql = `
SELECT
    user_id,
    waves
FROM tournament_entries
WHERE tournament_id = (SELECT MAX(id) FROM tournaments)
ORDER BY waves DESC;
`
	rows, err := db.Query(fetchTournamentEntriesSql)
	if err != nil {
		log.Println("DB Error:", err)
		session.ChannelMessageSend(channelId, "Error retrieving leaderboard.")
		return
	}
	defer rows.Close()

	leaderboardMsg := fmt.Sprintf("🏆 **Latest Tournament Leaderboard** 🏆\n")
	var entries []string
	for i := 1; rows.Next(); i++ {
		var userId string
		var waves int
		rows.Scan(&userId, &waves)
		entry := fmt.Sprintf("%d. **<@%s>** - Wave %d\n", i, userId, waves)
		entries = append(entries, entry)
	}

	entriesMsg := "_No entries yet._"
	if len(entries) > 0 {
		entriesMsg = strings.Join(entries, "\n")
	}
	leaderboardMsg += entriesMsg

	session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Content: leaderboardMsg,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{}, // Prevents pinging
		},
	})
}
