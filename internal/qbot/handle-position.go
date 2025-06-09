package qbot

import (
	"fmt"
)

// handlePosition tells a user their position in the queue.
func (q *QBot) handlePosition(cmd Cmd) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	position := -1
	var found bool
	for i, item := range q.queue {
		if item.UserID == cmd.Message.Author.ID {
			position = i
			found = true
			break
		}
	}
	if !found {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, you are not in the queue.", cmd.Message.Author.ID))
		return nil
	}

	var turnIndicator string
	if position == 0 {
		turnIndicator = " (it is currently your turn)"
	}
	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, your position in the queue is %d%s.", cmd.Message.Author.ID, position+1, turnIndicator))

	return nil
}
