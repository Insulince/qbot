package qbot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// Schedule of announcements (UTC time)
var schedule = map[string]string{
	// Saturday tourney
	"Saturday 00:00": "@everyone Tournament signups are now open! Use `!queue` to join",
	"Saturday 23:00": "@everyone 1-hour warning! Tournament runs close soon",
	"Sunday 00:00":   "@everyone Tournament runs are now closed",
	"Sunday 04:00":   "@everyone Tournament has ended, rewards are being distributed",

	// Wednesday tourney
	"Wednesday 00:00": "@everyone Tournament signups are now open! Use `!queue` to join",
	"Wednesday 23:00": "@everyone 1-hour warning! Tournament runs close soon",
	"Thursday 00:00":  "@everyone Tournament runs are now closed",
	"Thursday 04:00":  "@everyone Tournament has ended, rewards are being distributed",
}

// Function to check the time and send messages
func (q *QBot) startScheduler() error {
	for {
		now := time.Now().UTC()
		key := fmt.Sprintf("%s %02d:%02d", now.Weekday(), now.Hour(), now.Minute())

		if msg, exists := schedule[key]; exists {
			if (now.Weekday() == time.Saturday || now.Weekday() == time.Wednesday) && now.Hour() == 0 && now.Minute() == 0 {
				if err := q.createNewTournament(); err != nil {
					return errors.Wrap(err, "creating new tournament")
				}
			}

			err := q.announceMessage(key, msg)
			if err != nil {
				return errors.Wrap(err, "announce message")
			}

			if (now.Weekday() == time.Sunday || now.Weekday() == time.Thursday) && now.Hour() == 4 && now.Minute() == 0 {
				for guildId := range q.guilds {
					if err := q.handleLeaderboard(Cmd{GuildId: guildId}, true); err != nil {
						return errors.Wrap(err, "handling leaderboard")
					}
				}
			}
		}

		time.Sleep(60 * time.Second) // Check every minute
	}
}

func (q *QBot) announceMessage(key, msg string) error {
	const everyonePath = "/app/assets/everyone.png"

	for _, g := range q.guilds {
		file, err := os.Open(everyonePath)
		if err != nil {
			q.mustPost(g.AnnouncementChannelId, "❌ Error: Could not open everyone image.")
			return errors.Wrapf(err, "open everyone image %q", everyonePath)
		}
		defer jmust.MustClose(file)

		// Create a message with the image
		message := &discordgo.MessageSend{
			Content: msg,
			Files: []*discordgo.File{
				{
					Name:   "everyone.png",
					Reader: file,
				},
			},
		}

		// Send the message with the image
		if _, err = q.session.ChannelMessageSendComplex(g.AnnouncementChannelId, message); err != nil {
			q.mustPost(g.AnnouncementChannelId, "❌ Error: Failed to send announcement message and image.")
			return errors.Wrapf(err, "send announcement message and image")
		}

		log.Printf("[%s] Scheduled message sent to guild %q: %q\n", key, g.Name, msg)
	}

	return nil
}

// Function to create a new tournament
func (q *QBot) createNewTournament() error {
	// Generate short_name in yyyy-mm-dd format
	now := time.Now().UTC()
	shortName := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())

	// Generate the full tournament name
	fullName := fmt.Sprintf("%d %s %d", now.Day(), strings.ToUpper(now.Month().String()), now.Year())

	if err := q.store.InsertTournament(fullName, shortName); err != nil {
		return errors.Wrap(err, "inserting new tournament")
	}

	log.Println("✅ New tournament created:", fullName)
	return nil
}
