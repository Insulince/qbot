package qbot

import (
	"fmt"
	"github.com/Insulince/qbot/internal/version"
	"github.com/bwmarrin/discordgo"
)

// handleVersion displays the bots current version.
func (q *QBot) handleVersion(m *discordgo.MessageCreate, _ []string) error {
	v := version.MustGet()

	versionMessage := fmt.Sprintf(""+
		"**Version**\n"+
		"%s\n",
		v)

	q.mustPost(m.ChannelID, versionMessage)

	return nil
}
