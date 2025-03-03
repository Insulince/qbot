package qbot

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

// Handle !leaderboard
func handleLeaderboard(session *discordgo.Session, message *discordgo.MessageCreate, db *sql.DB) {
	channelId := "1343112046404833351"
	if message != nil {
		channelId = message.ChannelID
	}

	rows, err := db.Query("SELECT user_id, wave FROM leaderboard ORDER BY wave DESC LIMIT 30")
	if err != nil {
		log.Println("DB Error:", err)
		session.ChannelMessageSend(channelId, "Error retrieving leaderboard.")
		return
	}
	defer rows.Close()

	leaderboardMsg := "**🏆 Tournament Leaderboard 🏆**\n"
	i := 1
	for rows.Next() {
		var user_id string
		var wave int
		rows.Scan(&user_id, &wave)
		leaderboardMsg += fmt.Sprintf("%d. **<@%s>** - Wave %d\n", i, user_id, wave)
		i++
	}

	if i == 1 {
		leaderboardMsg += "_No records yet._"
	}

	session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Content: fmt.Sprintf(leaderboardMsg),
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{}, // Prevents pinging
		},
	})
}
