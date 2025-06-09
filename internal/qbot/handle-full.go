package qbot

import (
	"fmt"
	"time"
)

// handleFull signals that the active user's bracket is full.
func (q *QBot) handleFull(cmd Cmd) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	// Check if the user is the active player (first in queue)
	if len(q.queue) == 0 || q.queue[0].UserID != cmd.Message.Author.ID {
		// Check if the user is in the queue at all
		inQueue := false
		position := 0
		for i, item := range q.queue {
			if item.UserID == cmd.Message.Author.ID {
				inQueue = true
				position = i
				break
			}
		}

		if inQueue {
			// User is in the queue but not active yet
			q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are currently in the queue at position %d. Please wait for your turn before using `!full`.", cmd.Message.Author.ID, position))
		} else {
			// User is not in the queue at all
			q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are not currently in the queue. Use `!queue` to join.", cmd.Message.Author.ID))
		}
		return nil
	}

	if !q.queue[0].Entered {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, please signal that you've entered your bracket first using `!enter`.", cmd.Message.Author.ID))
		return nil
	}

	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, noted that your bracket is now full. Removing you from the queue.", cmd.Message.Author.ID))

	// Remove the active user from the queue
	q.queue = q.queue[1:]

	// If there's a new active user, notify them
	if len(q.queue) > 0 {
		q.mustPostWithoutTags(cmd.Message.ChannelID, fmt.Sprintf("Continuing to next user in queue, <@%s> (may be in different server)", q.queue[0].UserID))
		// Reset the timer for the new active user
		q.queue[0].AddedAt = time.Now()
		q.queue[0].Warned = false
		q.sendPass(q.queue[0].ChannelID, q.queue[0].UserID, fmt.Sprintf("<@%s>, it's now your turn! Please type `!enter` once you join your bracket.", q.queue[0].UserID))
	}

	return nil
}
