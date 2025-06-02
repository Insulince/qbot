package qbot

import (
	"fmt"
	"time"

	"github.com/Insulince/qbot/internal/version"
)

// handleVersion displays the bot's current version and uptime.
func (q *QBot) handleVersion(cmd Cmd) error {
	v := version.MustGet()

	uptime := time.Since(q.started).Round(time.Second)

	versionMessage := fmt.Sprintf(
		"**Version**\n%s\n**Uptime**\n%s",
		v, uptime.String(),
	)

	q.mustPost(cmd.Message.ChannelID, versionMessage)

	return nil
}
