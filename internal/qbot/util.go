package qbot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

// Send error message to Discord
func (q *QBot) reportError(err error) {
	if err == nil {
		return
	}

	const (
		errorChannelId     = "1342728552705163336" // #q-testing
		notificationRoleId = "1348155892377321554" // @Q-dev
	)

	// Format error message
	errorMessage := fmt.Sprintf("🚨 **Error in Q** 🚨\n<@&%s>: %s", notificationRoleId, err.Error())

	// Send error message to the private error channel
	q.mustPost(errorChannelId, errorMessage)

	// Also log error to stdout for redundancy
	log.Println(err)
}

func (q *QBot) Go(fn func() error) {
	go func() {
		if err := fn(); err != nil {
			q.reportError(err)
		}
	}()
}

// isModerator checks whether the invoking member has a role named "Moderator".
func (q *QBot) isModerator(m *discordgo.MessageCreate) bool {
	// Must be in a guild and have member info.
	if m.GuildID == "" || m.Member == nil {
		return false
	}

	if strings.ToLower(m.Author.GlobalName) == "insulince" { // dev privilege LOL
		return true
	}

	// First try using cached guild data.
	guild, err := q.session.State.Guild(m.GuildID)
	var roles []*discordgo.Role
	if err != nil {
		// Fallback: fetch guild roles.
		roles, err = q.session.GuildRoles(m.GuildID)
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
func (q *QBot) timeoutChecker() error {
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
				q.mustPost(q.currentUser.ChannelID, fmt.Sprintf("<@%s>, you have two minutes left (%s). Please update your status or use `!moretime` to extend the deadline.", q.currentUser.UserID, phase))
				q.currentUser.Warned = true
			}

			// Timeout and promote the next user if the allowed timeout is exceeded.
			if elapsed > allowedTimeout {
				q.mustPost(q.currentUser.ChannelID, fmt.Sprintf("<@%s> timed out (%s). Moving to the next person in the queue.", q.currentUser.UserID, phase))
				q.currentUser = nil
				if len(q.queue) > 0 {
					next := q.queue[0]
					q.queue = q.queue[1:]
					next.AddedAt = time.Now()
					next.Warned = false
					q.currentUser = &next
					q.mustPost(next.ChannelID, fmt.Sprintf("<@%s>, it's now your turn! Please type `!enter` once you join your bracket.", next.UserID))
				}
			}
		}
		q.queueMutex.Unlock()
	}

	return nil
}
