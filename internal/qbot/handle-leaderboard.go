package qbot

import (
	"fmt"
	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"os"
	"strings"
)

const (
	ChannelIdTournamentQueue = "1343112046404833351"                      // #tournament-queue
	FoolMemeURL              = "https://your-image-url.com/fool-meme.jpg" // Replace with actual URL
	foolUsername             = "hengstler"
)

// Handle !leaderboard
func (q *QBot) handleLeaderboard(cmd Cmd, final bool) error {
	channelId := ChannelIdTournamentQueue
	if cmd.Message != nil {
		channelId = cmd.Message.ChannelID
	}

	const fetchTournamentEntriesSql = `
SELECT
    user_id,
    username,
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
	var lastPlaceUsername string
	for i := 1; rows.Next(); i++ {
		var userId string
		var username string
		var waves int
		if err := rows.Scan(&userId, &username, &waves); err != nil {
			return errors.Wrap(err, "scanning fields")
		}
		entry := fmt.Sprintf("%d. **<@%s>** - Wave %d", i, userId, waves)
		entries = append(entries, entry)

		lastPlaceUsername = username // Keep track of the last user
	}

	entriesMsg := "_No entries yet._"
	if len(entries) > 0 {
		entriesMsg = strings.Join(entries, "\n")
	}
	leaderboardMsg += entriesMsg

	// Send the leaderboard message first
	q.mustPostWithoutTags(channelId, leaderboardMsg)

	if final {
		// If "Fool" is in last place, post the meme
		if lastPlaceUsername == foolUsername {
			q.sendFoolMeme(channelId) // Send Fool meme as an attachment
		}
	}

	return nil
}

// Send a local image file
func (q *QBot) sendFoolMeme(channelID string) {
	const filePath = "/app/assets/wompwomp.png"
	file, err := os.Open(filePath)
	if err != nil {
		q.mustPost(channelID, "❌ Error: Could not load image.")
		return
	}
	defer file.Close()

	// Create a message with both text and an image
	message := &discordgo.MessageSend{
		Content: "😂 **Congrats Fool!** 😂",
		Files: []*discordgo.File{
			{
				Name:   "wompwomp.png",
				Reader: file,
			},
		},
	}

	// Send the message with the image
	_, err = q.session.ChannelMessageSendComplex(channelID, message)
	if err != nil {
		q.mustPost(channelID, "❌ Error: Failed to send image.")
	}
}
