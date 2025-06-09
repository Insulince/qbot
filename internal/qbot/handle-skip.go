package qbot

import (
	"fmt"
	"time"
)

// handleSkip allows a moderator to force-skip the current active user.
func (q *QBot) handleSkip(cmd Cmd) error {
	if !q.isModerator(cmd.Message) {
		q.mustPost(cmd.Message.ChannelID, "You do not have permission to use this command. Moderator role required.")
		return nil
	}

	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	if len(q.queue) == 0 {
		q.mustPost(cmd.Message.ChannelID, "Cannot skip, the queue is empty.")
		return nil
	}

	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("Moderator <@%s> has skipped <@%s>.", cmd.Message.Author.ID, q.queue[0].UserID))

	// Remove the active user
	q.queue = q.queue[1:]

	// If there's a new active user, notify them
	if len(q.queue) > 0 {
		// Reset the timer for the new active user
		q.queue[0].AddedAt = time.Now()
		q.queue[0].Warned = false
		q.sendPass(q.queue[0].ChannelID, q.queue[0].UserID, fmt.Sprintf("<@%s>, it's now your turn! Please type `!enter` once you join your bracket.", q.queue[0].UserID))
	}

	return nil
}
