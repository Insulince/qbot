package qbot

import "github.com/bwmarrin/discordgo"

// handleReset clears the entire queue (moderator only).
func (q *QBot) handleReset(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isModerator(s, m) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command. Moderator role required.")
		return
	}

	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()
	q.currentUser = nil
	q.queue = []QueueItem{}
	s.ChannelMessageSend(m.ChannelID, "Queue has been reset.")
}
