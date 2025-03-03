package qbot

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// handleInsert allows a moderator to insert data into the database.
func (q *QBot) handleInsert(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isModerator(s, m) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command. Moderator role required.")
		return
	}

	parts := strings.Fields(m.Content)
	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !insert <data>")
		return
	}
	data := strings.Join(parts[1:], " ")
	if data == "" {
		s.ChannelMessageSend(m.ChannelID, "Data must not be empty.")
		return
	}

	db, err := sql.Open("sqlite", "/var/lib/litefs/qbot.db")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Errorf("error opening database: %w", err).Error())
		return
	}

	if _, err := db.Exec("INSERT INTO tbl (data) VALUES (?)", data); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Errorf("error inserting user: %w", err).Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Inserted data: %q\n", data))
	log.Printf("Inserted data: %q\n", data)
}
