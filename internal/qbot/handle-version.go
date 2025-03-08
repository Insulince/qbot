package qbot

import (
	"fmt"
	"github.com/Insulince/qbot/internal/version"
)

// handleVersion displays the bots current version.
func (q *QBot) handleVersion(cmd Cmd) error {
	v := version.MustGet()

	versionMessage := fmt.Sprintf(""+
		"**Version**\n"+
		"%s\n",
		v)

	q.mustPost(cmd.Message.ChannelID, versionMessage)

	return nil
}
