package qbot

import (
	"fmt"
)

// handlePosition tells a user their position in the queue or if they're active.
func (q *QBot) handlePosition(cmd Cmd) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	if q.currentUser != nil && q.currentUser.UserID == cmd.Message.Author.ID {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are currently active.", cmd.Message.Author.ID))
		return nil
	}

	position := -1
	for i, item := range q.queue {
		if item.UserID == cmd.Message.Author.ID {
			position = i + 1
			break
		}
	}

	if position > 0 {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, your position in the queue is %d.", cmd.Message.Author.ID, position))
	} else {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are not in the queue.", cmd.Message.Author.ID))
	}

	return nil
}
