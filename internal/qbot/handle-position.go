package qbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// handlePosition tells a user their position in the queue or if they're active.
func (q *QBot) handlePosition(s *discordgo.Session, m *discordgo.MessageCreate) {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	if q.currentUser != nil && q.currentUser.UserID == m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, you are currently active.", m.Author.ID))
		return
	}

	position := -1
	for i, item := range q.queue {
		if item.UserID == m.Author.ID {
			position = i + 1
			break
		}
	}

	if position > 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, your position in the queue is %d.", m.Author.ID, position))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, you are not in the queue.", m.Author.ID))
	}
}
