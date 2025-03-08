package qbot

import (
	"github.com/pkg/errors"
	"strings"
)

// handleDevError is a developer function for testing errors. It will intentionally throw an error that echoes back what
// you throw into it.
func (q *QBot) handleDevError(cmd Cmd) error {
	if !q.isModerator(cmd.Message) {
		return nil
	}

	msg := strings.Join(cmd.Args, " ")
	err := errors.Errorf("intentional dev error triggered with payload %q", msg)

	return err
}
