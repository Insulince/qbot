package qbot

import (
	"github.com/bwmarrin/discordgo"
)

// handleReset clears the entire queue (moderator only).
func (q *QBot) handleReset(m *discordgo.MessageCreate, _ []string) error {
	if !q.isModerator(m) {
		q.mustPost(m.ChannelID, "You do not have permission to use this command. Moderator role required.")
		return nil
	}

	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()
	q.currentUser = nil
	q.queue = []QueueItem{}
	q.mustPost(m.ChannelID, "Queue has been reset.")

	return nil
}
