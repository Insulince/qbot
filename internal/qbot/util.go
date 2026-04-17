package qbot

import (
	"context"
	"fmt"
	"log/slog"
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
	slog.Error("error reported", "err", err)
}

func (q *QBot) Go(fn func() error) {
	go func() {
		defer func() {
			if v := recover(); v != nil {
				msg := fmt.Sprintf("🚨 **Panic in Q** 🚨\n<@&%s>: panic: %v", q.notificationRoleId, v)
				slog.Error("panic in goroutine", "panic", v)
				if err := q.post(q.errorChannelId, msg); err != nil {
					slog.Error("failed to report panic to Discord", "err", err)
				}
				time.Sleep(2 * time.Second) // allow Discord to receive the message before crashing
				panic(v)
			}
		}()
		if err := fn(); err != nil {
			q.reportError(err)
		}
	}()
}

// GoWithRestart launches a long-running goroutine. Unlike Go, if fn panics or returns
// an error it is reported and fn is restarted after a brief pause rather than crashing.
// The goroutine exits cleanly when ctx is cancelled.
func (q *QBot) GoWithRestart(ctx context.Context, name string, fn func() error) {
	go func() {
		for {
			func() {
				defer func() {
					if v := recover(); v != nil {
						q.reportError(fmt.Errorf("panic in %s (restarting in 5s): %v", name, v))
					}
				}()
				if err := fn(); err != nil {
					q.reportError(fmt.Errorf("%s exited with error (restarting in 5s): %w", name, err))
				}
			}()

			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
		}
	}()
}

// cachedGuildRoles holds a guild's role list along with when it was fetched.
type cachedGuildRoles struct {
	roles     []*discordgo.Role
	fetchedAt time.Time
}

const rolesCacheTTL = 10 * time.Minute

// getGuildRoles returns the role list for a guild, using a local cache to avoid
// hitting the Discord API on every moderator check. Entries older than rolesCacheTTL
// are refreshed transparently.
func (q *QBot) getGuildRoles(guildID string) ([]*discordgo.Role, error) {
	q.guildRolesMu.RLock()
	cached, ok := q.guildRolesCache[guildID]
	q.guildRolesMu.RUnlock()

	if ok && time.Since(cached.fetchedAt) < rolesCacheTTL {
		return cached.roles, nil
	}

	roles, err := q.session.GuildRoles(guildID)
	if err != nil {
		return nil, errors.Wrap(err, "fetch guild roles")
	}

	q.guildRolesMu.Lock()
	q.guildRolesCache[guildID] = cachedGuildRoles{roles: roles, fetchedAt: time.Now()}
	q.guildRolesMu.Unlock()

	return roles, nil
}

// isModerator checks whether the invoking member has a role matching the guild's configured moderator role name.
func (q *QBot) isModerator(m *discordgo.MessageCreate) bool {
	// TODO(Insulince): Rework this.

	// Must be in a guild and have member info.
	if m.GuildID == "" || m.Member == nil {
		return false
	}

	if strings.ToLower(m.Author.GlobalName) == "insulince" { // dev privilege LOL
		return true
	}

	g, err := q.guilds.Get(m.GuildID)
	if err != nil {
		return false
	}
	moderatorRoleName := g.ModeratorRoleName
	if moderatorRoleName == "" {
		moderatorRoleName = "Moderator" // fallback default
	}

	roles, err := q.getGuildRoles(m.GuildID)
	if err != nil {
		return false
	}

	for _, roleID := range m.Member.Roles {
		for _, role := range roles {
			if role.ID == roleID && role.Name == moderatorRoleName {
				return true
			}
		}
	}
	return false
}

// timeoutChecker periodically checks if the active user has timed out or is nearing timeout.
func (q *QBot) timeoutChecker(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
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
	}
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

// parseMention extracts a Discord user ID from a mention string like <@123456> or <@!123456>.
func parseMention(s string) (string, bool) {
	if len(s) < 4 || s[0] != '<' || s[1] != '@' || s[len(s)-1] != '>' {
		return "", false
	}
	inner := s[2 : len(s)-1]
	if len(inner) > 0 && inner[0] == '!' {
		inner = inner[1:]
	}
	if len(inner) == 0 {
		return "", false
	}
	return inner, true
}

func (q *QBot) mustAnnounceStart() {
	q.mustPost(q.errorChannelId, fmt.Sprintf("Q %s started at %v!", version.MustGet(), time.Now().UTC()))
}
