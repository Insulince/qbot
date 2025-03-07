package qbot

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func Test_parseTournamentShortName(t *testing.T) {
	testCache := map[string]struct{}{}

	const expected = "2025-03-07"
	const separators = "-/\\._"
	for _, c := range []struct {
		description string
		input       string
	}{
		{"yyyy-0m-0d", "2025-03-07"},
		{"yyyy-0m-d", "2025-03-7"},
		{"yyyy-m-0d", "2025-3-07"},
		{"yyyy-m-d", "2025-3-7"},
		{"yy-0m-0d", "25-03-07"},
		{"yy-0m-d", "25-03-7"},
		{"yy-m-0d", "25-3-07"},
		{"yy-m-d", "25-3-7"},
		{"0m-0d", "03-07"},
		{"0m-d", "03-7"},
		{"m-0d", "3-07"},
		{"m-d", "3-7"},
		{"0d", "07"},
		{"d", "07"},
		{"yyyy0m0d", "20250307"},
		{"yy0m0d", "250307"},
	} {
		for _, sep := range strings.Split(separators, "") {
			description := strings.Replace(c.description, "-", sep, -1)
			input := strings.Replace(c.input, "-", sep, -1)

			// To prevent duplicate runs of the same test.
			_, ok := testCache[input]
			if ok {
				continue
			}
			testCache[input] = struct{}{}

			t.Run(fmt.Sprintf("when date format is %q, should result in %q", description, expected), func(t *testing.T) {
				shortName, err := parseTournamentShortName(input)
				require.NoError(t, err)
				require.Equal(t, expected, shortName)
			})
		}
	}

	t.Run("yyyy-mm-dd", func(t *testing.T) {
		input := "2025-03-07"
		shortName, err := parseTournamentShortName(input)
		require.NoError(t, err)
		require.Equal(t, "2025-03-07", shortName)
	})
	t.Run("yy-mm-dd", func(t *testing.T) {
		input := "25-03-07"
		shortName, err := parseTournamentShortName(input)
		require.NoError(t, err)
		require.Equal(t, "2025-03-07", shortName)
	})
	t.Run("mm-dd", func(t *testing.T) {
		input := "03-07"
		shortName, err := parseTournamentShortName(input)
		require.NoError(t, err)
		require.Equal(t, "2025-03-07", shortName)
	})
	t.Run("dd", func(t *testing.T) {
		input := "07"
		shortName, err := parseTournamentShortName(input)
		require.NoError(t, err)
		require.Equal(t, "2025-03-07", shortName)
	})
	t.Run("yyyymmdd", func(t *testing.T) {
		input := "20250307"
		shortName, err := parseTournamentShortName(input)
		require.NoError(t, err)
		require.Equal(t, "2025-03-07", shortName)
	})
	t.Run("yymmdd", func(t *testing.T) {
		input := "250307"
		shortName, err := parseTournamentShortName(input)
		require.NoError(t, err)
		require.Equal(t, "2025-03-07", shortName)
	})
}
