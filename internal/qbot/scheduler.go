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

var announcementChannelID = "1343112046404833351" // #tournament-queue

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

			err := q.announceMessage(msg)
			if err != nil {
				return errors.Wrap(err, "announce message")
			}
			log.Printf("[%s] Scheduled message sent: %q\n", key, msg)

			if (now.Weekday() == time.Sunday || now.Weekday() == time.Thursday) && now.Hour() == 4 && now.Minute() == 0 {
				if err := q.handleLeaderboard(emptyCmd, true); err != nil {
					return errors.Wrap(err, "handling leaderboard")
				}
			}
		}

		time.Sleep(60 * time.Second) // Check every minute
	}
}

func (q *QBot) announceMessage(msg string) error {
	const everyonePath = "/app/assets/everyone.png"

	file, err := os.Open(everyonePath)
	if err != nil {
		q.mustPost(announcementChannelID, "❌ Error: Could not open everyone image.")
		return errors.Wrapf(err, "open everyone image %q", everyonePath)
	}
	defer jmust.MustClose(file)

	// Create a message with the modified image
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
	if _, err = q.session.ChannelMessageSendComplex(announcementChannelID, message); err != nil {
		q.mustPost(announcementChannelID, "❌ Error: Failed to send announcement message and image.")
		return errors.Wrapf(err, "send announcement message and image")
	}

	return nil
}

// Function to create a new tournament
func (q *QBot) createNewTournament() error {
	// Generate short_name in yyyy-mm-dd format
	now := time.Now().UTC()
	shortName := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())

	// Generate full tournament name
	fullName := fmt.Sprintf("%d %s %d", now.Day(), strings.ToUpper(now.Month().String()), now.Year())

	// Insert new tournament into the database
	const newTournamentSql = `
INSERT INTO tournaments
    (name, short_name)
VALUES
    (?, ?);
`
	if _, err := q.db.Exec(newTournamentSql, fullName, shortName); err != nil {
		return errors.Wrap(err, "❌ Error inserting new tournament")
	}

	log.Println("✅ New tournament created:", fullName)
	return nil
}
