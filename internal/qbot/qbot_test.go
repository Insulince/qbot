package qbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func Test_interpretMessage(t *testing.T) {
	t.Run("should properly interpret command + args", func(t *testing.T) {
		content := "!submitwave 1234"
		m := &discordgo.MessageCreate{Message: &discordgo.Message{Content: content}}

		cmd, err := interpretMessage(m)
		require.NoError(t, err)

		segments := strings.Split(content, " ")
		expectedCommand := segments[0]
		expectedArgs := segments[1:]
		t.Logf("content: %s\nsegments: %#v\nlen(segments): %d\ncommand: %s\nargs: %s\n", content, segments, len(segments), cmd.Command, cmd.Args)

		require.Equal(t, expectedCommand, cmd.Command)
		require.Equal(t, expectedArgs, cmd.Args)
	})
}
