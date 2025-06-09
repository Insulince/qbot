package qbot

import (
	"fmt"
	"time"
)

// handleFull signals that the active user's bracket is full.
func (q *QBot) handleFull(cmd Cmd) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	if len(q.queue) == 0 || q.queue[0].UserID != cmd.Message.Author.ID {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, it's not your turn.", cmd.Message.Author.ID))
		return nil
	}

	if !q.queue[0].Entered {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, please signal that you've entered your bracket first using `!enter`.", cmd.Message.Author.ID))
		return nil
	}

	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, your bracket is now full. Removing you from the queue.", cmd.Message.Author.ID))

	// Remove the active user from the queue
	q.queue = q.queue[1:]

	// If there's a new active user, notify them
	if len(q.queue) > 0 {
		// Reset the timer for the new active user
		q.queue[0].AddedAt = time.Now()
		q.queue[0].Warned = false
		q.mustPost(q.queue[0].ChannelID, fmt.Sprintf("<@%s>, it's now your turn! Please type `!enter` once you join your bracket.", q.queue[0].UserID))
	}

	return nil
}
