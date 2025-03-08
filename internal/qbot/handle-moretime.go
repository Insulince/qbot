package qbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

func (q *QBot) handleMoreTime(m *discordgo.MessageCreate, _ []string) error {
	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()

	if q.currentUser == nil || q.currentUser.UserID != m.Author.ID {
		q.mustPost(m.ChannelID, fmt.Sprintf("<@%s>, it's not your turn.", m.Author.ID))
		return nil
	}

	// Determine which phase the user is in
	if !q.currentUser.Entered {
		// User has not yet entered, extend the "enter" timer
		q.currentUser.AddedAt = time.Now()
		q.currentUser.Warned = false
		q.mustPost(m.ChannelID, fmt.Sprintf("<@%s>, your time to enter has been extended. Please type `!enter` once you've joined your bracket.", m.Author.ID))
	} else {
		// User has entered but not yet confirmed full, extend the "full" timer
		q.currentUser.AddedAt = time.Now()
		q.currentUser.Warned = false
		q.mustPost(m.ChannelID, fmt.Sprintf("<@%s>, your time to confirm your full bracket has been extended.", m.Author.ID))
	}

	return nil
}
