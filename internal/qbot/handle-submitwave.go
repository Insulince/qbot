package qbot

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
)

// Handle !submitwave <wave>
func handleSubmitWave(session *discordgo.Session, message *discordgo.MessageCreate, args []string, db *sql.DB) {
	if len(args) != 2 {
		session.ChannelMessageSend(message.ChannelID, "Usage: `!submitwave <wave>`")
		return
	}

	userID := message.Author.ID
	username := message.Author.Username
	id := message.Author.ID

	wave, err := strconv.Atoi(args[1])
	if err != nil || wave < 1 || wave > 10000 {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> Please enter a valid wave number between 1 and 10000.", id))
		return
	}

	const fetchLatestTournamentIdSql = `
SELECT
    MAX(id)
FROM tournaments;
`
	var tournamentId int
	if err := db.QueryRow(fetchLatestTournamentIdSql).Scan(&tournamentId); err != nil {
		log.Println("DB Error:", err)
		session.ChannelMessageSend(message.ChannelID, "Error retrieving leaderboard.")
		return
	}

	// Update or insert new high score
	insertWaveSql := `
INSERT INTO tournament_entries
    (tournament_id, user_id, username, waves)
VALUES
    (? ?, ?, ?)
ON CONFLICT(user_id) DO UPDATE SET wave = excluded.wave;
`
	_, err = db.Exec(insertWaveSql, tournamentId, userID, username, wave)
	if err != nil {
		log.Println("DB Error:", err)
		session.ChannelMessageSend(message.ChannelID, "Error saving your wave.")
		return
	}
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("✅ <@%s> set their wave to **%d**!", id, wave))
}
