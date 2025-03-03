package qbot

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

// handleDeleteAll allows a moderator to delete all data from the database.
func (q *QBot) handleDeleteAll(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isModerator(s, m) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command. Moderator role required.")
		return
	}

	db, err := sql.Open("sqlite", "/var/lib/litefs/qbot.db")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Errorf("error opening database: %w", err).Error())
		return
	}

	if _, err := db.Query("DELETE FROM tbl"); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Errorf("error deleting all data: %w", err).Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Deleted all data.")
	log.Printf("Deleted all data.")
}
