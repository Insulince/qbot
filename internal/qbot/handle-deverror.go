package qbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"strings"
)

// handleDevError is a developer function for testing errors. It will intentionally throw an error that echoes back what
// you throw into it.
func (q *QBot) handleDevError(m *discordgo.MessageCreate, args []string) error {
	if !q.isModerator(m) {
		return nil
	}

	msg := strings.Join(args, " ")
	err := errors.Errorf("intentional dev error triggered with payload %q", msg)

	return err
}
