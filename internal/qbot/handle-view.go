package qbot

import (
	"fmt"
	"time"
)

// handleView displays the current queue.
func (q *QBot) handleView(cmd Cmd) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	msg := "**Queue**:\n"
	if len(q.queue) == 0 {
		msg += "The queue is empty."
	} else {
		for i, item := range q.queue {
			if i == 0 {
				status := "Not Entered"
				timeoutDur := q.enterTimeout
				if item.Entered {
					status = "Entered"
					timeoutDur = q.fullTimeout
				}
				msg += fmt.Sprintf("%d. <@%s> [**ACTIVE**] (Status: %s, since %s, timeout period: %s)\n", i+1, item.UserID, status, item.AddedAt.Format(time.RFC822), timeoutDur)
			} else {
				msg += fmt.Sprintf("%d. <@%s> (queued at %s)\n", i+1, item.UserID, item.AddedAt.Format(time.RFC822))
			}
		}
	}

	q.mustPostWithoutTags(cmd.Message.ChannelID, msg)

	return nil
}
