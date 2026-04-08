package qbot

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Insulince/qbot/internal/models"
)

func Test_formatProgressResults(t *testing.T) {
	t.Run("empty slice returns empty string", func(t *testing.T) {
		require.Equal(t, "", formatProgressResults(nil))
	})

	t.Run("single entry formats correctly", func(t *testing.T) {
		results := []*models.PlayerTournamentResult{
			{TournamentName: "Sunday Weekly", TournamentShortName: "2026-04-06", Waves: 1850, Rank: 2, TotalEntrants: 12},
		}
		output := formatProgressResults(results)
		require.Equal(t, "• **Sunday Weekly** (`2026-04-06`) — Wave `1850` | Rank `2`/`12`", output)
	})

	t.Run("multiple entries are newline separated", func(t *testing.T) {
		results := []*models.PlayerTournamentResult{
			{TournamentName: "Sunday Weekly", TournamentShortName: "2026-04-06", Waves: 1850, Rank: 2, TotalEntrants: 12},
			{TournamentName: "Wednesday Weekly", TournamentShortName: "2026-04-02", Waves: 1400, Rank: 4, TotalEntrants: 8},
		}
		output := formatProgressResults(results)
		require.Contains(t, output, "\n")
		require.Contains(t, output, "Sunday Weekly")
		require.Contains(t, output, "Wednesday Weekly")
	})

	t.Run("rank 1 with 1 entrant formats correctly", func(t *testing.T) {
		results := []*models.PlayerTournamentResult{
			{TournamentName: "Solo Tourney", TournamentShortName: "2026-03-01", Waves: 500, Rank: 1, TotalEntrants: 1},
		}
		output := formatProgressResults(results)
		require.Equal(t, "• **Solo Tourney** (`2026-03-01`) — Wave `500` | Rank `1`/`1`", output)
	})
}
