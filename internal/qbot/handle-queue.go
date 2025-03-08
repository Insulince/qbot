package qbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

// handleQueue adds a user to the waiting queue.
func (q *QBot) handleQueue(m *discordgo.MessageCreate, _ []string) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	// Check if user is already active.
	if q.currentUser != nil && q.currentUser.UserID == m.Author.ID {
		q.mustPost(m.ChannelID, fmt.Sprintf("<@%s>, you are already active. Please use `!enter` or `!full` as appropriate.", m.Author.ID))
		return nil
	}
	// Check if user is already in the waiting queue.
	for _, item := range q.queue {
		if item.UserID == m.Author.ID {
			q.mustPost(m.ChannelID, fmt.Sprintf("<@%s>, you are already in the queue.", m.Author.ID))
			return nil
		}
	}

	// If no one is active, promote immediately.
	if q.currentUser == nil {
		q.currentUser = &QueueItem{
			UserID:    m.Author.ID,
			AddedAt:   time.Now(), // start the enter timeout period
			ChannelID: m.ChannelID,
			Entered:   false,
			Warned:    false,
		}
		q.mustPost(m.ChannelID, fmt.Sprintf("<@%s>, it's your turn! Please type `!enter` once you join your bracket.", m.Author.ID))
	} else {
		// Otherwise, add the user to the waiting queue.
		q.queue = append(q.queue, QueueItem{
			UserID:    m.Author.ID,
			AddedAt:   time.Now(), // not used until promoted
			ChannelID: m.ChannelID,
			Entered:   false,
			Warned:    false,
		})
		q.mustPost(m.ChannelID, fmt.Sprintf("<@%s>, you've been added to the queue. Your position is %d.", m.Author.ID, len(q.queue)))
	}

	return nil
}
