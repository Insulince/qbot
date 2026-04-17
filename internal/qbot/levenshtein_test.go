package qbot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_levenshtein(t *testing.T) {
	for _, c := range []struct {
		description string
		a, b        string
		expected    int
	}{
		{"identical strings", "!queue", "!queue", 0},
		{"empty strings", "", "", 0},
		{"one empty", "!queue", "", 6},
		{"other empty", "", "!queue", 6},
		{"single insertion", "!queue", "!queu", 1},
		{"single deletion", "!queu", "!queue", 1},
		{"single substitution", "!queve", "!queue", 1},
		{"transposition counts as two edits", "!quuee", "!queue", 2},
		{"completely different", "!queue", "!help", 4},
		{"unicode runes counted correctly", "!queueé", "!queue", 1},
	} {
		t.Run(c.description, func(t *testing.T) {
			require.Equal(t, c.expected, levenshtein(c.a, c.b))
		})
	}
}

func Test_suggestCommand(t *testing.T) {
	for _, c := range []struct {
		description string
		input       string
		expected    string
	}{
		{"exact match", "!queue", "!queue"},
		{"one-character typo", "!queeu", "!queue"},
		{"transposed letters", "!queeu", "!queue"},
		// "queue" is only 1 edit from "!queue" (insert "!"), so the half-length guard doesn't suppress it
		{"missing exclamation prefix — still close", "queue", "!queue"},
		{"completely unrelated short input — half-length guard rejects", "!ab", ""},
		{"close to submitwave", "!submitewave", "!submitwave"},
		{"close to leaderboard", "!leaderborad", "!leaderboard"},
		{"close to moretime", "!moretiime", "!moretime"},
		{"exact match for alias-excluded command", "!clearwave", "!clearwave"},
	} {
		t.Run(c.description, func(t *testing.T) {
			require.Equal(t, c.expected, suggestCommand(c.input))
		})
	}
}
