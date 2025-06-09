package qbot

import (
	"fmt"
	"time"
)

func (q *QBot) handleEnter(cmd Cmd) error {
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
			q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are currently in the queue at position %d. Please wait for your turn before using `!enter`.", cmd.Message.Author.ID, position))
		} else {
			// User is not in the queue at all
			q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are not currently in the queue. Use `!queue` to join.", cmd.Message.Author.ID))
		}
		return nil
	}

	// If the user is already marked as entered, prevent redundant entry
	if q.queue[0].Entered {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you have already entered your bracket. Please use `!full` when your bracket is complete.", cmd.Message.Author.ID))
		return nil
	}

	// Mark the user as having entered
	q.queue[0].Entered = true
	q.queue[0].AddedAt = time.Now()
	q.queue[0].Warned = false

	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, noted that you've entered your bracket. Please type `!full` once your bracket is full.", cmd.Message.Author.ID))

	return nil
}
