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

	if q.currentUser == nil {
		q.mustPost(cmd.Message.ChannelID, "There is no active user to skip.")
		return nil
	}

	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("Moderator <@%s> has skipped <@%s>.", cmd.Message.Author.ID, q.currentUser.UserID))
	q.currentUser = nil
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
