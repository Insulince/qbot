package qbot

import (
	"fmt"
	"time"
)

// handleLeave allows a user to remove themselves from the active slot or waiting queue.
func (q *QBot) handleLeave(cmd Cmd) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	// If the user is currently active.
	if q.currentUser != nil && q.currentUser.UserID == cmd.Message.Author.ID {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you have left the active slot.", cmd.Message.Author.ID))
		q.currentUser = nil
		// Promote the next user if available.
		if len(q.queue) > 0 {
			next := q.queue[0]
			q.queue = q.queue[1:]
			next.AddedAt = time.Now()
			next.Warned = false
			q.currentUser = &next
			q.mustPost(next.ChannelID, fmt.Sprintf("<@%s>, it's now your turn! Please type `!enter` once you join your bracket.", next.UserID))
		}
		return nil
	}

	// Check if the user is in the waiting queue.
	found := false
	var newQueue []QueueItem
	for _, item := range q.queue {
		if item.UserID == cmd.Message.Author.ID {
			found = true
		} else {
			newQueue = append(newQueue, item)
		}
	}
	q.queue = newQueue
	if found {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you have been removed from the queue.", cmd.Message.Author.ID))
	} else {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are not in the queue.", cmd.Message.Author.ID))
	}

	return nil
}
