package qbot

import (
	"fmt"
	"time"
)

// QueueItem represents an entry in the queue.
type QueueItem struct {
	UserID    string
	AddedAt   time.Time // timestamp for the current phase (enter or full)
	ChannelID string    // channel the user joined from
	Entered   bool      // whether the user has signaled they've joined a bracket
	Warned    bool      // whether a warning has been sent for the current phase
}

// handleQueue adds a user to the waiting queue.
func (q *QBot) handleQueue(cmd Cmd) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	// Check if the user is already active.
	if q.currentUser != nil && q.currentUser.UserID == cmd.Message.Author.ID {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are already active. Please use `!enter` or `!full` as appropriate.", cmd.Message.Author.ID))
		return nil
	}
	// Check if the user is already in the waiting queue.
	for _, item := range q.queue {
		if item.UserID == cmd.Message.Author.ID {
			q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are already in the queue.", cmd.Message.Author.ID))
			return nil
		}
	}

	// If no one is active, promote immediately.
	if q.currentUser == nil {
		q.currentUser = &QueueItem{
			UserID:    cmd.Message.Author.ID,
			AddedAt:   time.Now(), // start the enter timeout period
			ChannelID: cmd.Message.ChannelID,
			Entered:   false,
			Warned:    false,
		}
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, it's your turn! Please type `!enter` once you join your bracket.", cmd.Message.Author.ID))
	} else {
		// Otherwise, add the user to the waiting queue.
		q.queue = append(q.queue, QueueItem{
			UserID:    cmd.Message.Author.ID,
			AddedAt:   time.Now(), // not used until promoted
			ChannelID: cmd.Message.ChannelID,
			Entered:   false,
			Warned:    false,
		})
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you've been added to the queue. Your position is %d.", cmd.Message.Author.ID, len(q.queue)))
	}

	return nil
}
