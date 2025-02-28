package qbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

// handleLeave allows a user to remove themselves from the active slot or waiting queue.
func (q *QBot) handleLeave(s *discordgo.Session, m *discordgo.MessageCreate) {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	// If the user is currently active.
	if q.currentUser != nil && q.currentUser.UserID == m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, you have left the active slot.", m.Author.ID))
		q.currentUser = nil
		// Promote the next user if available.
		if len(q.queue) > 0 {
			next := q.queue[0]
			q.queue = q.queue[1:]
			next.AddedAt = time.Now()
			next.Warned = false
			q.currentUser = &next
			s.ChannelMessageSend(next.ChannelID, fmt.Sprintf("<@%s>, it's now your turn! Please type `!enter` once you join your bracket.", next.UserID))
		}
		return
	}

	// Check if the user is in the waiting queue.
	found := false
	newQueue := []QueueItem{}
	for _, item := range q.queue {
		if item.UserID == m.Author.ID {
			found = true
		} else {
			newQueue = append(newQueue, item)
		}
	}
	q.queue = newQueue
	if found {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, you have been removed from the queue.", m.Author.ID))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, you are not in the queue.", m.Author.ID))
	}
}
