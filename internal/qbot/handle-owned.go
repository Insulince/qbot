package qbot

import (
	"os"

	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// Handle !owned
func (q *QBot) handleOwned(cmd Cmd) error {
	const ownedPath = "/app/assets/owned.png"

	file, err := os.Open(ownedPath)
	if err != nil {
		q.mustPost(cmd.Message.ChannelID, "❌ Error: Could not open owned image.")
		return errors.Wrapf(err, "open owned image %q", ownedPath)
	}
	defer jmust.MustClose(file)

	// Create a message with the modified image
	message := &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   "owned.png",
				Reader: file,
			},
		},
	}

	// Send the message with the image
	if _, err = q.session.ChannelMessageSendComplex(cmd.Message.ChannelID, message); err != nil {
		q.mustPost(cmd.Message.ChannelID, "❌ Error: Failed to send image.")
		return errors.Wrapf(err, "send owned image")
	}

	return nil
}
