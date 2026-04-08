package qbot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmd_IsEmpty(t *testing.T) {
	t.Run("empty cmd is empty", func(t *testing.T) {
		require.True(t, emptyCmd.IsEmpty())
	})

	t.Run("cmd with command is not empty", func(t *testing.T) {
		cmd := Cmd{Command: "!queue"}
		require.False(t, cmd.IsEmpty())
	})

	t.Run("zero value is empty", func(t *testing.T) {
		var cmd Cmd
		require.True(t, cmd.IsEmpty())
	})
}

func TestCmd_String(t *testing.T) {
	t.Run("command only", func(t *testing.T) {
		cmd := Cmd{Command: "!queue"}
		require.Equal(t, "!queue ", cmd.String())
	})

	t.Run("command with single arg", func(t *testing.T) {
		cmd := Cmd{Command: "!submitwave", Args: []string{"1500"}}
		require.Equal(t, "!submitwave 1500", cmd.String())
	})

	t.Run("command with multiple args", func(t *testing.T) {
		cmd := Cmd{Command: "!history", Args: []string{"10", "5"}}
		require.Equal(t, "!history 10 5", cmd.String())
	})
}
