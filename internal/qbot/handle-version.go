package qbot

import (
	"fmt"
	"github.com/Insulince/qbot/internal/version"
	"github.com/bwmarrin/discordgo"
)

// handleVersion displays the bots current version.
func (q *QBot) handleVersion(s *discordgo.Session, m *discordgo.MessageCreate) {
	v := version.MustGet()
	versionMessage := fmt.Sprintf(""+
		"**Version**\n"+
		"%s\n",
		v)
	s.ChannelMessageSend(m.ChannelID, versionMessage)
}
