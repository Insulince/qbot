package qbot

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/require"
)

func Test_messageIsShouting(t *testing.T) {
	for _, c := range []struct {
		description string
		input       string
		expected    bool
	}{
		{"all caps letters", "HELLO WORLD", true},
		{"all caps with numbers", "HELLO123", true},
		{"all caps single word", "HELLO", true},
		{"mixed case", "Hello World", false},
		{"all lowercase", "hello world", false},
		{"no letters — numbers only", "12345", false},
		{"no letters — empty string", "", false},
		{"no letters — punctuation only", "!!! ???", false},
		{"exclamation command all caps", "!HELLO", true},
		{"exclamation command lowercase", "!hello", false},
	} {
		t.Run(c.description, func(t *testing.T) {
			require.Equal(t, c.expected, messageIsShouting(c.input))
		})
	}
}

func Test_isNumber(t *testing.T) {
	for _, c := range []struct {
		description string
		input       string
		expected    bool
	}{
		{"positive integer", "!1234", true},
		{"zero", "!0", true},
		{"negative integer", "!-5", true}, // strconv.Atoi parses negatives
		{"letters", "!abc", false},
		{"mixed alphanumeric", "!12abc", false},
		{"empty after prefix", "!", false},
		{"float", "!1.5", false},
	} {
		t.Run(c.description, func(t *testing.T) {
			require.Equal(t, c.expected, isNumber(c.input))
		})
	}
}

func Test_QBot_interpretMessage_extended(t *testing.T) {
	newQ := func() *QBot {
		q := new(QBot)
		q.guilds = Guilds{"": {}}
		return q
	}

	t.Run("lowercases the command", func(t *testing.T) {
		q := newQ()
		m := &discordgo.MessageCreate{Message: &discordgo.Message{Content: "!HISTORY"}}
		cmd, err := q.interpretMessage(m)
		require.NoError(t, err)
		require.Equal(t, "!history", cmd.Command)
	})

	t.Run("trims leading spaces", func(t *testing.T) {
		q := newQ()
		m := &discordgo.MessageCreate{Message: &discordgo.Message{Content: "   !queue"}}
		cmd, err := q.interpretMessage(m)
		require.NoError(t, err)
		require.Equal(t, "!queue", cmd.Command)
	})

	t.Run("parses multiple args", func(t *testing.T) {
		q := newQ()
		m := &discordgo.MessageCreate{Message: &discordgo.Message{Content: "!history 10 5"}}
		cmd, err := q.interpretMessage(m)
		require.NoError(t, err)
		require.Equal(t, "!history", cmd.Command)
		require.Equal(t, []string{"10", "5"}, cmd.Args)
	})

	t.Run("returns empty cmd for non-command message", func(t *testing.T) {
		q := newQ()
		m := &discordgo.MessageCreate{Message: &discordgo.Message{Content: "just a normal message"}}
		cmd, err := q.interpretMessage(m)
		require.NoError(t, err)
		require.True(t, cmd.IsEmpty())
	})

	t.Run("returns empty cmd for only exclamation marks", func(t *testing.T) {
		q := newQ()
		m := &discordgo.MessageCreate{Message: &discordgo.Message{Content: "!!!"}}
		cmd, err := q.interpretMessage(m)
		require.NoError(t, err)
		require.True(t, cmd.IsEmpty())
	})

	t.Run("handles ❗ emoji prefix", func(t *testing.T) {
		q := newQ()
		m := &discordgo.MessageCreate{Message: &discordgo.Message{Content: "❗queue"}}
		cmd, err := q.interpretMessage(m)
		require.NoError(t, err)
		require.Equal(t, "!queue", cmd.Command)
	})

	t.Run("handles ❕ emoji prefix", func(t *testing.T) {
		q := newQ()
		m := &discordgo.MessageCreate{Message: &discordgo.Message{Content: "❕queue"}}
		cmd, err := q.interpretMessage(m)
		require.NoError(t, err)
		require.Equal(t, "!queue", cmd.Command)
	})

	t.Run("no args produces empty args slice", func(t *testing.T) {
		q := newQ()
		m := &discordgo.MessageCreate{Message: &discordgo.Message{Content: "!queue"}}
		cmd, err := q.interpretMessage(m)
		require.NoError(t, err)
		require.Equal(t, "!queue", cmd.Command)
		require.Empty(t, cmd.Args)
	})
}
