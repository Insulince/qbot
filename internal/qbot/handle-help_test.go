package qbot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_resolveAlias(t *testing.T) {
	for _, c := range []struct {
		description string
		input       string
		expected    string
	}{
		{"primary command passes through unchanged", "!queue", "!queue"},
		{"!q resolves to !queue", "!q", "!queue"},
		{"!enqueue resolves to !queue", "!enqueue", "!queue"},
		{"!join resolves to !queue", "!join", "!queue"},
		{"!e resolves to !enter", "!e", "!enter"},
		{"!enterbracket resolves to !enter", "!enterbracket", "!enter"},
		{"!f resolves to !full", "!f", "!full"},
		{"!bracketfull resolves to !full", "!bracketfull", "!full"},
		{"!leavequeue resolves to !leave", "!leavequeue", "!leave"},
		{"!viewqueue resolves to !view", "!viewqueue", "!view"},
		{"!currentposition resolves to !position", "!currentposition", "!position"},
		{"!extend resolves to !moretime", "!extend", "!moretime"},
		{"!moar resolves to !moretime", "!moar", "!moretime"},
		{"!wave resolves to !submitwave", "!wave", "!submitwave"},
		{"!waves resolves to !submitwave", "!waves", "!submitwave"},
		{"!w resolves to !submitwave", "!w", "!submitwave"},
		{"!submitwaves resolves to !submitwave", "!submitwaves", "!submitwave"},
		{"!lb resolves to !leaderboard", "!lb", "!leaderboard"},
		{"!skipcurrent resolves to !skip", "!skipcurrent", "!skip"},
		{"!resetqueue resolves to !reset", "!resetqueue", "!reset"},
		{"!removeplayer resolves to !remove", "!removeplayer", "!remove"},
		{"unknown command passes through unchanged", "!unknown", "!unknown"},
	} {
		t.Run(c.description, func(t *testing.T) {
			require.Equal(t, c.expected, resolveAlias(c.input))
		})
	}
}
