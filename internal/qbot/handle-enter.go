package qbot

import (
	"fmt"
	"time"
)

func (q *QBot) handleEnter(cmd Cmd) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	// Check if the user is the active player
	if q.currentUser == nil || q.currentUser.UserID != cmd.Message.Author.ID {
		// First, check if the user is even in the queue
		inQueue := false
		position := 0
		for i, item := range q.queue {
			if item.UserID == cmd.Message.Author.ID {
				inQueue = true
				position = i + 1
				break
			}
		}

		if inQueue {
			// User is in queue but not active yet
			q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are currently in the queue at position %d. Please wait for your turn before using `!enter`.", cmd.Message.Author.ID, position))
		} else {
			// User is not in the queue at all
			q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are not currently in the queue. Use `!queue` to join.", cmd.Message.Author.ID))
		}
		return nil
	}

	// If the user is already marked as entered, prevent redundant entry
	if q.currentUser.Entered {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you have already entered your bracket. Please use `!full` when your bracket is complete.", cmd.Message.Author.ID))
		return nil
	}

	// Mark user as having entered
	q.currentUser.Entered = true
	q.currentUser.AddedAt = time.Now()
	q.currentUser.Warned = false

	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, noted that you've entered your bracket. Please type `!full` once your bracket is full.", cmd.Message.Author.ID))

	return nil
}
