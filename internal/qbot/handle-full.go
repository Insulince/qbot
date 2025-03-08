package qbot

import (
	"fmt"
	"time"
)

// handleFull signals that the active user's bracket is full.
func (q *QBot) handleFull(cmd Cmd) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	if q.currentUser == nil || q.currentUser.UserID != cmd.Message.Author.ID {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, it's not your turn.", cmd.Message.Author.ID))
		return nil
	}

	if !q.currentUser.Entered {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, please signal that you've entered your bracket first using `!enter`.", cmd.Message.Author.ID))
		return nil
	}

	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, your bracket is now full. Removing you from the active position.", cmd.Message.Author.ID))
	q.currentUser = nil

	// Promote the next user in the waiting queue, if any.
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
