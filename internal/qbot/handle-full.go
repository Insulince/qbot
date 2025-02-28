package qbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

// handleFull signals that the active user's bracket is full.
func (q *QBot) handleFull(s *discordgo.Session, m *discordgo.MessageCreate) {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	if q.currentUser == nil || q.currentUser.UserID != m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, it's not your turn.", m.Author.ID))
		return
	}

	if !q.currentUser.Entered {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, please signal that you've entered your bracket first using `!enter`.", m.Author.ID))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, your bracket is now full. Removing you from the active position.", m.Author.ID))
	q.currentUser = nil

	// Promote the next user in the waiting queue, if any.
	if len(q.queue) > 0 {
		next := q.queue[0]
		q.queue = q.queue[1:]
		next.AddedAt = time.Now()
		next.Warned = false
		q.currentUser = &next
		s.ChannelMessageSend(next.ChannelID, fmt.Sprintf("<@%s>, it's now your turn! Please type `!enter` once you join your bracket.", next.UserID))
	}
}
