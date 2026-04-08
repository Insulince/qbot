package qbot

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseTournamentShortName(t *testing.T) {
	testCache := map[string]struct{}{}

	const separators = "-/\\._"

	// Cases that include an explicit year — always resolve to 2025-03-07.
	const expectedWithYear = "2025-03-07"
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
		{"yyyy0m0d", "20250307"},
		{"yy0m0d", "250307"},
	} {
		for _, sep := range strings.Split(separators, "") {
			description := strings.Replace(c.description, "-", sep, -1)
			input := strings.Replace(c.input, "-", sep, -1)

			_, ok := testCache[input]
			if ok {
				continue
			}
			testCache[input] = struct{}{}

			t.Run(fmt.Sprintf("when date format is %q, should result in %q", description, expectedWithYear), func(t *testing.T) {
				shortName, err := parseTournamentShortName(input)
				require.NoError(t, err)
				require.Equal(t, expectedWithYear, shortName)
			})
		}
	}

	// Cases with no year — resolved using the current year (2026), so March 7
	// (which is in the past relative to April 2026) resolves to 2026-03-07.
	const expectedWithoutYear = "2026-03-07"
	for _, c := range []struct {
		description string
		input       string
	}{
		{"0m-0d", "03-07"},
		{"0m-d", "03-7"},
		{"m-0d", "3-07"},
		{"m-d", "3-7"},
	} {
		for _, sep := range strings.Split(separators, "") {
			description := strings.Replace(c.description, "-", sep, -1)
			input := strings.Replace(c.input, "-", sep, -1)

			_, ok := testCache[input]
			if ok {
				continue
			}
			testCache[input] = struct{}{}

			t.Run(fmt.Sprintf("when date format is %q, should result in %q", description, expectedWithoutYear), func(t *testing.T) {
				shortName, err := parseTournamentShortName(input)
				require.NoError(t, err)
				require.Equal(t, expectedWithoutYear, shortName)
			})
		}
	}

	t.Run("yyyy-mm-dd", func(t *testing.T) {
		shortName, err := parseTournamentShortName("2025-03-07")
		require.NoError(t, err)
		require.Equal(t, "2025-03-07", shortName)
	})
	t.Run("yy-mm-dd", func(t *testing.T) {
		shortName, err := parseTournamentShortName("25-03-07")
		require.NoError(t, err)
		require.Equal(t, "2025-03-07", shortName)
	})
	t.Run("mm-dd", func(t *testing.T) {
		shortName, err := parseTournamentShortName("03-07")
		require.NoError(t, err)
		require.Equal(t, "2026-03-07", shortName)
	})
	t.Run("dd should not be supported", func(t *testing.T) {
		shortName, err := parseTournamentShortName("07")
		require.ErrorContains(t, err, "invalid format \"07\"")
		require.Zero(t, shortName)
	})
	t.Run("yyyymmdd", func(t *testing.T) {
		shortName, err := parseTournamentShortName("20250307")
		require.NoError(t, err)
		require.Equal(t, "2025-03-07", shortName)
	})
	t.Run("yymmdd", func(t *testing.T) {
		shortName, err := parseTournamentShortName("250307")
		require.NoError(t, err)
		require.Equal(t, "2025-03-07", shortName)
	})
}
