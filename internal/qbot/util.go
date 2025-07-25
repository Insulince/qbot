package qbot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/Insulince/qbot/internal/version"
)

// Send an error message to Discord
func (q *QBot) reportError(err error) {
	if err == nil {
		return
	}

	// Format error message
	errorMessage := fmt.Sprintf("🚨 **Error in Q** 🚨\n<@&%s>: %s", q.notificationRoleId, err.Error())

	// Send an error message to the private error channel
	q.mustPost(q.errorChannelId, errorMessage)

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
	// TODO(Insulince): Rework this.

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
		if len(q.queue) > 0 {
			// The first user in the queue is the active user
			activeUser := &q.queue[0]
			var allowedTimeout time.Duration
			var phase string

			if !activeUser.Entered {
				allowedTimeout = q.enterTimeout
				phase = "waiting to enter"
			} else {
				allowedTimeout = q.fullTimeout
				phase = "waiting to complete your bracket"
			}

			elapsed := time.Since(activeUser.AddedAt)

			// Send a warning if within the warning threshold and not yet warned.
			if !activeUser.Warned && elapsed >= allowedTimeout-q.warnThreshold {
				q.mustPost(activeUser.ChannelID, fmt.Sprintf("<@%s>, you have two minutes left (%s). Please update your status or use `!moretime` to extend the deadline.", activeUser.UserID, phase))
				activeUser.Warned = true
			}

			// Timeout and promote the next user if the allowed timeout is exceeded.
			if elapsed > allowedTimeout {
				q.mustPost(activeUser.ChannelID, fmt.Sprintf("<@%s> timed out (%s).", activeUser.UserID, phase))

				// Remove the active user
				q.queue = q.queue[1:]

				// If there's a new active user, notify them
				if len(q.queue) > 0 {
					q.mustPostWithoutTags(activeUser.ChannelID, fmt.Sprintf("Continuing to next user in queue, <@%s> (may be in different server)", q.queue[0].UserID))
					// Reset the timer for the new active user
					q.queue[0].AddedAt = time.Now()
					q.queue[0].Warned = false
					q.sendPass(q.queue[0].ChannelID, q.queue[0].UserID, fmt.Sprintf("<@%s>, it's now your turn! Please type `!enter` once you join your bracket.", q.queue[0].UserID))
				} else {
					q.mustPost(activeUser.ChannelID, fmt.Sprintf("The queue is empty. Use `!queue` to join!"))
				}
			}
		}
		q.queueMutex.Unlock()
	}

	return nil
}

func (q *QBot) GetDisplayName(msg *discordgo.MessageCreate) (string, error) {
	if msg.GuildID == "" {
		// Not in a guild, only global/username available
		if msg.Author.GlobalName != "" {
			return msg.Author.GlobalName, nil
		}
		return msg.Author.Username, nil
	}

	member, err := q.session.GuildMember(msg.GuildID, msg.Author.ID)
	if err != nil {
		return "", errors.Wrap(err, "get guild member")
	}

	if member.Nick != "" {
		return member.Nick, nil
	}

	if member.User.GlobalName != "" {
		return member.User.GlobalName, nil
	}

	return member.User.Username, nil
}

func (q *QBot) mustAnnounceStart() {
	q.mustPost(q.errorChannelId, fmt.Sprintf("Q %s started at %v!", version.MustGet(), time.Now().UTC()))
}
