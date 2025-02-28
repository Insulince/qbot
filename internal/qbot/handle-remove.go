package qbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

// handleRemove allows a moderator to remove a specific user from the queue.
func (q *QBot) handleRemove(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isModerator(s, m) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command. Moderator role required.")
		return
	}

	parts := strings.Fields(m.Content)
	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !remove @user")
		return
	}
	mention := parts[1]
	userID := ""
	if len(mention) >= 3 && mention[0] == '<' && mention[1] == '@' {
		mention = mention[2:]
		if mention[0] == '!' {
			mention = mention[1:]
		}
		if mention[len(mention)-1] == '>' {
			userID = mention[:len(mention)-1]
		}
	}
	if userID == "" {
		s.ChannelMessageSend(m.ChannelID, "Could not parse user mention.")
		return
	}

	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()
	removed := false
	if q.currentUser != nil && q.currentUser.UserID == userID {
		q.currentUser = nil
		removed = true
	}
	newQueue := []QueueItem{}
	for _, item := range q.queue {
		if item.UserID == userID {
			removed = true
		} else {
			newQueue = append(newQueue, item)
		}
	}
	q.queue = newQueue
	if removed {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User <@%s> has been removed from the queue.", userID))
		// If the active user was removed, promote the next one.
		if q.currentUser == nil && len(q.queue) > 0 {
			next := q.queue[0]
			q.queue = q.queue[1:]
			next.AddedAt = time.Now()
			next.Warned = false
			q.currentUser = &next
			s.ChannelMessageSend(next.ChannelID, fmt.Sprintf("<@%s>, it's now your turn! Please type `!enter` once you join your bracket.", next.UserID))
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "User not found in the queue.")
	}
}
