package qbot

import (
	"fmt"
	"strings"
	"time"
)

// handleRemove allows a moderator to remove a specific user from the queue.
func (q *QBot) handleRemove(cmd Cmd) error {
	if !q.isModerator(cmd.Message) {
		q.mustPost(cmd.Message.ChannelID, "You do not have permission to use this command. Moderator role required.")
		return nil
	}

	parts := strings.Fields(cmd.Message.Content)
	if len(parts) < 2 {
		q.mustPost(cmd.Message.ChannelID, "Usage: !remove @user")
		return nil
	}
	mention := parts[1]
	userID := ""
	if len(mention) >= 3 && mention[0] == '<' && mention[1] == '@' {
		mention = mention[2:]
		if mention[0] == '!' {
			mention = mention[1:]
		}
		if mention[len(mention)-1] == '>' {
			userID = mention[:len(mention)-1]
		}
	}
	if userID == "" {
		q.mustPost(cmd.Message.ChannelID, "Could not parse user mention.")
		return nil
	}

	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	// Check if the user is in the queue
	removed := false
	wasActive := false
	var newQueue []QueueItem

	for i, item := range q.queue {
		if item.UserID == userID {
			removed = true
			if i == 0 {
				wasActive = true
			}
		} else {
			newQueue = append(newQueue, item)
		}
	}

	q.queue = newQueue

	if removed {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s> has been removed from the queue.", userID))

		// If the active user was removed, notify the new active user
		if wasActive && len(q.queue) > 0 {
			// Reset the timer for the new active user
			q.queue[0].AddedAt = time.Now()
			q.queue[0].Warned = false
			q.mustPost(q.queue[0].ChannelID, fmt.Sprintf("<@%s>, it's now your turn! Please type `!enter` once you join your bracket.", q.queue[0].UserID))
		}
	} else {
		q.mustPost(cmd.Message.ChannelID, "User not found in the queue.")
	}

	return nil
}
