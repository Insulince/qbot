package qbot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseMention(t *testing.T) {
	for _, c := range []struct {
		description    string
		input          string
		expectedUserId string
		expectedOk     bool
	}{
		{"standard mention", "<@392896333670776832>", "392896333670776832", true},
		{"mention with nickname prefix", "<@!392896333670776832>", "392896333670776832", true},
		{"empty string", "", "", false},
		{"missing opening angle bracket", "@392896333670776832>", "", false},
		{"missing closing angle bracket", "<@392896333670776832", "", false},
		{"missing @ sign", "<392896333670776832>", "", false},
		{"empty user id", "<@>", "", false},
		{"empty user id with nickname prefix", "<@!>", "", false},
		{"plain user id no brackets", "392896333670776832", "", false},
	} {
		t.Run(c.description, func(t *testing.T) {
			userId, ok := parseMention(c.input)
			require.Equal(t, c.expectedOk, ok)
			require.Equal(t, c.expectedUserId, userId)
		})
	}
}
