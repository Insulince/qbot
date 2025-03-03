package qbot

import (
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
	"Saturday 19:00": "@everyone 1-hour warning! Tournament runs close soon",
	"Saturday 20:00": "@everyone Tournament runs are now closed",
	"Sunday 00:00":   "@everyone Tournament has ended, rewards are being distributed",

	// Wednesday tourney
	"Wednesday 00:00": "@everyone Tournament signups are now open! Use `!queue` to join",
	"Wednesday 19:00": "@everyone 1-hour warning! Tournament runs close soon",
	"Wednesday 20:00": "@everyone Tournament runs are now closed",
	"Thursday 00:00":  "@everyone Tournament has ended, rewards are being distributed",
}

// Function to check the time and send messages
func startScheduler(session *discordgo.Session) {
	for {
		now := time.Now().UTC()
		key := fmt.Sprintf("%s %02d:%02d", now.Weekday(), now.Hour(), now.Minute())

		if msg, exists := schedule[key]; exists {
			session.ChannelMessageSend(announcementChannelID, msg)
			log.Printf("[%s] Scheduled message sent: %q\n", key, msg)
		}

		time.Sleep(60 * time.Second) // Check every minute
	}
}
