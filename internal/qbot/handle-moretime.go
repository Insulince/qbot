package qbot

import (
	"fmt"
	"time"
)

func (q *QBot) handleMoreTime(cmd Cmd) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	if q.currentUser == nil || q.currentUser.UserID != cmd.Message.Author.ID {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, it's not your turn.", cmd.Message.Author.ID))
		return nil
	}

	// Determine which phase the user is in
	if !q.currentUser.Entered {
		// User has not yet entered, extend the "enter" timer
		q.currentUser.AddedAt = time.Now()
		q.currentUser.Warned = false
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, your time to enter has been extended. Please type `!enter` once you've joined your bracket.", cmd.Message.Author.ID))
	} else {
		// User has entered but not yet confirmed full, extend the "full" timer
		q.currentUser.AddedAt = time.Now()
		q.currentUser.Warned = false
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("<@%s>, your time to confirm your full bracket has been extended.", cmd.Message.Author.ID))
	}

	return nil
}
