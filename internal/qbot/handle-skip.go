package qbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

// handleSkip allows a moderator to force-skip the current active user.
func (q *QBot) handleSkip(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isModerator(s, m) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command. Moderator role required.")
		return
	}

	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	if q.currentUser == nil {
		s.ChannelMessageSend(m.ChannelID, "There is no active user to skip.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Moderator <@%s> has skipped <@%s>.", m.Author.ID, q.currentUser.UserID))
	q.currentUser = nil
	if len(q.queue) > 0 {
		next := q.queue[0]
		q.queue = q.queue[1:]
		next.AddedAt = time.Now()
		next.Warned = false
		q.currentUser = &next
		s.ChannelMessageSend(next.ChannelID, fmt.Sprintf("<@%s>, it's now your turn! Please type `!enter` once you join your bracket.", next.UserID))
	}
}
