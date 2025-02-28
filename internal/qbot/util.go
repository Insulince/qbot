package qbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

// isModerator checks whether the invoking member has a role named "Moderator".
func isModerator(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	// Must be in a guild and have member info.
	if m.GuildID == "" || m.Member == nil {
		return false
	}

	// First try using cached guild data.
	guild, err := s.State.Guild(m.GuildID)
	var roles []*discordgo.Role
	if err != nil {
		// Fallback: fetch guild roles.
		roles, err = s.GuildRoles(m.GuildID)
		if err != nil {
			return false
		}
	} else {
		roles = guild.Roles
	}

	// Check the member's role IDs against the guild roles.
	for _, roleID := range m.Member.Roles {
		for _, role := range roles {
			if role.ID == roleID && role.Name == "Moderator" {
				return true
			}
		}
	}
	return false
}

// timeoutChecker periodically checks if the active user has timed out or is nearing timeout.
func (q *QBot) timeoutChecker(s *discordgo.Session) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		q.queueMutex.Lock()
		if q.currentUser != nil {
			var allowedTimeout time.Duration
			var phase string

			if !q.currentUser.Entered {
				allowedTimeout = q.enterTimeout
				phase = "waiting to enter"
			} else {
				allowedTimeout = q.fullTimeout
				phase = "waiting to complete your bracket"
			}

			elapsed := time.Since(q.currentUser.AddedAt)

			// Send a warning if within the warning threshold and not yet warned.
			if !q.currentUser.Warned && elapsed >= allowedTimeout-q.warnThreshold {
				s.ChannelMessageSend(q.currentUser.ChannelID,
					fmt.Sprintf("<@%s>, you have two minutes left (%s). Please update your status.", q.currentUser.UserID, phase))
				q.currentUser.Warned = true
			}

			// Timeout and promote the next user if the allowed timeout is exceeded.
			if elapsed > allowedTimeout {
				s.ChannelMessageSend(q.currentUser.ChannelID,
					fmt.Sprintf("<@%s> timed out (%s). Moving to the next person in the queue.", q.currentUser.UserID, phase))
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
		}
		q.queueMutex.Unlock()
	}
}
