package qbot

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/bwmarrin/discordgo"
)

// QueueItem represents an entry in the queue.
type QueueItem struct {
	UserID    string
	AddedAt   time.Time // timestamp for the current phase (enter or full)
	ChannelID string    // channel the user joined from
	Entered   bool      // whether the user has signaled they've joined a bracket
	Warned    bool      // whether a warning has been sent for the current phase
}

// handleQueue adds a user to the queue.
func (q *QBot) handleQueue(cmd Cmd) error {
	const passPath = "/app/assets/pass-smaller.png"
	const blockPath = "/app/assets/block-smaller.png"

	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	// Check if the user is already in the queue.
	for i, item := range q.queue {
		if item.UserID == cmd.Message.Author.ID {
			if i == 0 {
				q.sendImageMessage(cmd.Message.ChannelID, cmd.Message.Author.ID, fmt.Sprintf("<@%s>, you are already first in queue and **it is currently your turn**! Please use `!enter` or `!full` as appropriate.", cmd.Message.Author.ID), passPath, "pass.png")
			} else {
				q.sendImageMessage(cmd.Message.ChannelID, cmd.Message.Author.ID, fmt.Sprintf("<@%s> 🚨 **DO NOT JOIN YET!** 🚨 You are already in the queue in position %d. Please wait for your turn, you will be pinged here when the time comes.\n_Players ahead of you:_\n%s", cmd.Message.Author.ID, i+1, q.formatPlayersAhead(i)), blockPath, "block.png")
			}
			return nil
		}
	}

	// Add the user to the queue
	newItem := QueueItem{
		UserID:    cmd.Message.Author.ID,
		AddedAt:   time.Now(),
		ChannelID: cmd.Message.ChannelID,
		Entered:   false,
		Warned:    false,
	}

	q.queue = append(q.queue, newItem)

	if len(q.queue) == 1 {
		q.sendImageMessage(cmd.Message.ChannelID, cmd.Message.Author.ID, fmt.Sprintf("<@%s>, you've been added to the queue and you're first so **it is now your turn**! Type `!enter` once you join your bracket.", cmd.Message.Author.ID), passPath, "pass.png")
	} else {
		q.sendImageMessage(cmd.Message.ChannelID, cmd.Message.Author.ID, fmt.Sprintf("<@%s> 🚨 **DO NOT JOIN YET!** 🚨 You've been added to the queue in position %d.\nPlease wait for your turn, you will be pinged here when the time comes.\n_Players ahead of you:_\n%s", cmd.Message.Author.ID, len(q.queue), q.formatPlayersAhead(len(q.queue)-1)), blockPath, "block.png")
	}

	return nil
}

// formatPlayersAhead returns a formatted string of players ahead of the given position
func (q *QBot) formatPlayersAhead(position int) string {
	limit := 10
	if position > limit {
		position = limit
	}
	var players []string
	for i := 0; i < position; i++ {
		players = append(players, fmt.Sprintf("<@%s>", q.queue[i].UserID))
	}
	if position == limit && len(q.queue) > limit {
		players = append(players, fmt.Sprintf("_and %d more..._", len(q.queue)-limit))
	}
	return strings.Join(players, ", ")
}

func (q *QBot) sendImageMessage(channelID, userId, content, imagePath, imageName string) {
	// TODO(Insulince): This should be made more official and re-homed.

	file, err := os.Open(imagePath)
	if err != nil {
		q.mustPost(channelID, "❌ Error: Could not open image.")
		return
	}
	defer jmust.MustClose(file)

	_, err = q.session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: content,
		Files: []*discordgo.File{
			{
				Name:   imageName,
				Reader: file,
			},
		},
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Users: []string{userId}, // Only ping the author, not those in the queue in front of the author.
			Roles: []string{},
			Parse: []discordgo.AllowedMentionType{},
		},
	})
	if err != nil {
		q.mustPost(channelID, "❌ Error: Failed to send image.")
	}
}
