package qbot

import (
	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// post splits long messages into chunks and sends them serially.
func (q *QBot) post(channelId, msg string) error {
	const chunkSize = 2000 // Max Discord message length

	// Split message into chunks of at most 2000 characters
	for len(msg) > 0 {
		// Determine the length of the next chunk
		end := chunkSize
		if len(msg) < chunkSize {
			end = len(msg)
		}

		// Send the chunk
		if _, err := q.session.ChannelMessageSend(channelId, msg[:end]); err != nil {
			return errors.Wrapf(err, "failed to send message to channel %q", channelId)
		}

		// Move to the next chunk
		msg = msg[end:]
	}

	return nil
}

func (q *QBot) mustPost(channelId, msg string) {
	jmust.Must[any](q.post, channelId, msg)
}

func (q *QBot) postWithoutTags(channelId, msg string) error {
	const chunkSize = 2000 // Max Discord message length

	// Split message into chunks of at most 2000 characters
	for len(msg) > 0 {
		// Determine the length of the next chunk
		end := chunkSize
		if len(msg) < chunkSize {
			end = len(msg)
		}

		// Send the chunk
		_, err := q.session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
			Content: msg[:end],
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{}, // Prevents pinging
			},
		})
		if err != nil {
			return errors.Wrapf(err, "failed to send message to channel %q", channelId)
		}

		// Move to the next chunk
		msg = msg[end:]
	}

	return nil
}

func (q *QBot) mustPostWithoutTags(channelId, msg string) {
	jmust.Must[any](q.postWithoutTags, channelId, msg)
}
