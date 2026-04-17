package qbot

import (
	"bytes"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

// handleShame posts the shame GIF, optionally targeting a mentioned user.
// Usage: !shame [@user]
func (q *QBot) handleShame(cmd Cmd) error {
	var targetUserId string
	if len(cmd.Args) > 0 {
		if id, ok := parseMention(cmd.Args[0]); ok {
			targetUserId = id
		}
	}

	var content string
	if targetUserId != "" {
		content = fmt.Sprintf("<@%s> 🔔 **SHAME! SHAME! SHAME!** 🔔", targetUserId)
	} else {
		content = "🔔 **SHAME! SHAME! SHAME!** 🔔"
	}

	data, err := os.ReadFile(assetShame)
	if err != nil {
		// Graceful fallback if the file hasn't been dropped in yet.
		q.mustPost(cmd.Message.ChannelID, content)
		return nil
	}

	msg := &discordgo.MessageSend{
		Content: content,
		Files: []*discordgo.File{
			{
				Name:   "shame.gif",
				Reader: bytes.NewReader(data),
			},
		},
	}

	if targetUserId != "" {
		msg.AllowedMentions = &discordgo.MessageAllowedMentions{
			Users: []string{targetUserId},
			Roles: []string{},
			Parse: []discordgo.AllowedMentionType{},
		}
	}

	if _, err = q.session.ChannelMessageSendComplex(cmd.Message.ChannelID, msg); err != nil {
		q.mustPost(cmd.Message.ChannelID, "❌ Error: Failed to send shame image.")
	}

	return nil
}
