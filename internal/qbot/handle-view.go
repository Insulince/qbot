package qbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

// handleView displays the current queue state.
func (q *QBot) handleView(m *discordgo.MessageCreate, _ []string) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	msg := "**Queue State**:\n"
	if q.currentUser != nil {
		status := "Not Entered"
		timeoutDur := q.enterTimeout
		if q.currentUser.Entered {
			status = "Entered"
			timeoutDur = q.fullTimeout
		}
		msg += fmt.Sprintf("**Active:** <@%s> (Status: %s, since %s, timeout period: %s)\n", q.currentUser.UserID, status, q.currentUser.AddedAt.Format(time.RFC822), timeoutDur)
	} else {
		msg += "No active user at the moment.\n"
	}

	if len(q.queue) == 0 {
		msg += "The waiting queue is empty."
	} else {
		msg += "**Waiting Queue:**\n"
		for i, item := range q.queue {
			msg += fmt.Sprintf("%d. <@%s> (queued at %s)\n", i+1, item.UserID, item.AddedAt.Format(time.RFC822))
		}
	}

	q.mustPostWithoutTags(m.ChannelID, msg)

	return nil
}
