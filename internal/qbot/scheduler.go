package qbot

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
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
func startScheduler(session *discordgo.Session) {
	for {
		now := time.Now().UTC()
		key := fmt.Sprintf("%s %02d:%02d", now.Weekday(), now.Hour(), now.Minute())

		if msg, exists := schedule[key]; exists {
			session.ChannelMessageSend(announcementChannelID, msg)
			log.Printf("[%s] Scheduled message sent: %q\n", key, msg)

			if (now.Weekday() == time.Sunday || now.Weekday() == time.Thursday) && now.Hour() == 0 && now.Minute() == 0 {
				db, err := sql.Open("sqlite", "/var/lib/litefs/qbot.db")
				if err != nil {
					log.Println(fmt.Errorf("error opening database: %w", err).Error())
					return
				}

				handleLeaderboard(session, nil, db)

				db.Close()
			}
		}

		time.Sleep(60 * time.Second) // Check every minute
	}
}
