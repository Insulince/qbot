package qbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func (q *QBot) handleEnter(s *discordgo.Session, m *discordgo.MessageCreate) {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	// Check if the user is the active player
	if q.currentUser == nil || q.currentUser.UserID != m.Author.ID {
		// First, check if the user is even in the queue
		inQueue := false
		position := 0
		for i, item := range q.queue {
			if item.UserID == m.Author.ID {
				inQueue = true
				position = i + 1
				break
			}
		}

		if inQueue {
			// User is in queue but not active yet
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				"<@%s>, you are currently in the queue at position %d. Please wait for your turn before using `!enter`.",
				m.Author.ID, position))
		} else {
			// User is not in the queue at all
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
				"<@%s>, you are not currently in the queue. Use `!queue` to join.",
				m.Author.ID))
		}
		return
	}

	// If the user is already marked as entered, prevent redundant entry
	if q.currentUser.Entered {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
			"<@%s>, you have already entered your bracket. Please use `!full` when your bracket is complete.",
			m.Author.ID))
		return
	}

	// Mark user as having entered
	q.currentUser.Entered = true
	q.currentUser.AddedAt = time.Now()
	q.currentUser.Warned = false

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(
		"<@%s>, noted that you've entered your bracket. Please type `!full` once your bracket is full.",
		m.Author.ID))
}
