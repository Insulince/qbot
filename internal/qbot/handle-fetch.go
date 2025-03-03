package qbot

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// handleFetch allows a moderator to fetch data from the database.
func (q *QBot) handleFetch(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isModerator(s, m) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command. Moderator role required.")
		return
	}

	db, err := sql.Open("sqlite", "/var/lib/litefs/qbot.db")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Errorf("error opening database: %w", err).Error())
		return
	}

	rows, err := db.Query("SELECT data FROM tbl")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Errorf("error fetching users: %w", err).Error())
		return
	}
	defer rows.Close()

	var data []string
	for rows.Next() {
		var datum string
		if err := rows.Scan(&datum); err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Errorf("error scanning row: %w", err).Error())
			return
		}
		data = append(data, fmt.Sprintf("%q", datum))
	}

	dataStr := "[" + strings.Join(data, ", ") + "]"
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Fetched data: %s\n", dataStr))
	log.Printf("Fetched data: %s\n", dataStr)
}
