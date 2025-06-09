package qbot

import (
	"fmt"
	"time"
)

// handleLeave allows a user to remove themselves from the queue.
func (q *QBot) handleLeave(cmd Cmd) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	// Check if the user is in the queue
	found := false
	position := -1
	for i, item := range q.queue {
		if item.UserID == cmd.Message.Author.ID {
			found = true
			position = i
			break
		}
	}

	if !found {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are not in the queue.", cmd.Message.Author.ID))
		return nil
	}

	// If the user is the active user (first in queue)
	if position == 0 {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you have left the queue.", cmd.Message.Author.ID))

		// Remove the active user
		q.queue = q.queue[1:]

		// If there's a new active user, notify them
		if len(q.queue) > 0 {
			// Reset the timer for the new active user
			q.queue[0].AddedAt = time.Now()
			q.queue[0].Warned = false
			q.sendPass(q.queue[0].ChannelID, q.queue[0].UserID, fmt.Sprintf("<@%s>, it's now your turn! Please type `!enter` once you join your bracket.", q.queue[0].UserID))
		}
	} else {
		// Remove the user from elsewhere in the queue
		q.queue = append(q.queue[:position], q.queue[position+1:]...)
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you have been removed from the queue.", cmd.Message.Author.ID))
	}

	return nil
}
