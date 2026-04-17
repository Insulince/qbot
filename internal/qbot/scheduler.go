package qbot

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// scheduledEvent describes a single scheduled occurrence.
type scheduledEvent struct {
	message string       // announcement message to broadcast (empty if no announcement)
	action  func() error // optional action to run before the announcement (nil if none)
	enabled bool         // set to false to skip without removing the entry
}

func (q *QBot) startScheduler(ctx context.Context) error {
	// schedule is built inside startScheduler so action fields can reference q methods.
	schedule := map[string]scheduledEvent{
		// Saturday tourney
		"Saturday 00:00": {
			message: "Tournament signups are now open! Use `!queue` to join",
			action:  q.createNewTournament,
			enabled: true,
		},
		"Saturday 23:00": {
			message: "1-hour warning! Tournament runs close soon",
			enabled: true,
		},
		"Sunday 04:00": {
			action:  q.announceLeaderboard,
			enabled: false,
		},

		// Wednesday tourney
		"Wednesday 00:00": {
			message: "Tournament signups are now open! Use `!queue` to join",
			action:  q.createNewTournament,
			enabled: true,
		},
		"Wednesday 23:00": {
			message: "1-hour warning! Tournament runs close soon",
			enabled: true,
		},
		"Thursday 04:00": {
			action:  q.announceLeaderboard,
			enabled: false,
		},
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			now := time.Now().UTC()
			key := fmt.Sprintf("%s %02d:%02d", now.Weekday(), now.Hour(), now.Minute())

			event, exists := schedule[key]
			if !exists || !event.enabled {
				continue
			}

			if event.action != nil {
				if err := event.action(); err != nil {
					return errors.Wrap(err, "scheduled action")
				}
			}

			if event.message != "" {
				if err := q.announceMessage(key, event.message); err != nil {
					return errors.Wrap(err, "announce message")
				}
			}
		}
	}
}

func (q *QBot) announceLeaderboard() error {
	for guildId := range q.guilds {
		if err := q.handleLeaderboard(Cmd{GuildId: guildId}, true); err != nil {
			return errors.Wrapf(err, "handling leaderboard for guild %s", guildId)
		}
	}
	return nil
}

func (q *QBot) announceMessage(key, msg string) error {
	data, err := os.ReadFile(assetEveryone)
	if err != nil {
		return errors.Wrapf(err, "read everyone image %q", assetEveryone)
	}

	for _, g := range q.guilds {
		var mentionMsg string
		if g.AudienceIdentifier == "everyone" {
			mentionMsg = fmt.Sprintf("@%s %s", g.AudienceIdentifier, msg)
		} else {
			mentionMsg = fmt.Sprintf("<@&%s> %s", g.AudienceIdentifier, msg)
		}

		reader := bytes.NewReader(data)

		message := &discordgo.MessageSend{
			Content: mentionMsg,
			Files: []*discordgo.File{
				{
					Name:   "everyone.png",
					Reader: reader,
				},
			},
		}

		if _, err = q.session.ChannelMessageSendComplex(g.AnnouncementChannelId, message); err != nil {
			q.mustPost(g.AnnouncementChannelId, "❌ Error: Failed to send announcement message and image.")
			return errors.Wrapf(err, "send announcement message and image")
		}

		slog.Info("scheduled message sent", "event", key, "guild", g.Name, "message", mentionMsg)
	}

	return nil
}

func (q *QBot) createNewTournament() error {
	now := time.Now().UTC()
	shortName := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	fullName := fmt.Sprintf("%d %s %d", now.Day(), strings.ToUpper(now.Month().String()), now.Year())

	if err := q.store.InsertTournament(fullName, shortName); err != nil {
		return errors.Wrap(err, "inserting new tournament")
	}

	slog.Info("new tournament created", "name", fullName)
	return nil
}
