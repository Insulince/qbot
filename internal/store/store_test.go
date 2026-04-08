package store_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/Insulince/qbot/internal/store"
)

const testSchema = `
CREATE TABLE IF NOT EXISTS tournaments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    short_name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS tournament_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tournament_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    username TEXT NOT NULL,
    waves INTEGER NOT NULL,
    display_name TEXT NOT NULL DEFAULT 'UNKNOWN',
    guild_id TEXT NOT NULL DEFAULT 'UNKNOWN',
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE,
    UNIQUE (tournament_id, user_id)
);
`

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	_, err = db.Exec(testSchema)
	require.NoError(t, err)
	s := store.MustNew(db)
	t.Cleanup(func() { _ = s.Close() })
	return s
}

// --- Tournament tests ---

func TestStore_InsertTournament(t *testing.T) {
	t.Run("inserts successfully", func(t *testing.T) {
		s := newTestStore(t)
		err := s.InsertTournament("Sunday Weekly", "2026-04-06")
		require.NoError(t, err)
	})

	t.Run("duplicate name fails", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))
		err := s.InsertTournament("Sunday Weekly", "2026-04-07")
		require.Error(t, err)
	})

	t.Run("duplicate short_name fails", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))
		err := s.InsertTournament("Different Name", "2026-04-06")
		require.Error(t, err)
	})
}

func TestStore_GetTournamentByShortName(t *testing.T) {
	t.Run("returns tournament when found", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))

		tournament, err := s.GetTournamentByShortName("2026-04-06")
		require.NoError(t, err)
		require.Equal(t, "Sunday Weekly", tournament.Name)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		s := newTestStore(t)
		_, err := s.GetTournamentByShortName("9999-01-01")
		require.Error(t, err)
	})
}

func TestStore_GetLatestTournament(t *testing.T) {
	t.Run("returns the tournament with the highest id", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("First", "2026-01-01"))
		require.NoError(t, s.InsertTournament("Second", "2026-01-08"))
		require.NoError(t, s.InsertTournament("Third", "2026-01-15"))

		tournament, err := s.GetLatestTournament()
		require.NoError(t, err)

		// Confirm the ID is greater than previous ones by checking it's the third inserted
		third, err := s.GetTournamentByShortName("2026-01-15")
		require.NoError(t, err)
		require.Equal(t, third.Id, tournament.Id)
	})
}

func TestStore_ListTournaments(t *testing.T) {
	t.Run("returns tournaments in descending order", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("First", "2026-01-01"))
		require.NoError(t, s.InsertTournament("Second", "2026-01-08"))
		require.NoError(t, s.InsertTournament("Third", "2026-01-15"))

		results, err := s.ListTournaments(10, 0)
		require.NoError(t, err)
		require.Len(t, results, 3)
		require.Equal(t, "Third", results[0].Name)
		require.Equal(t, "Second", results[1].Name)
		require.Equal(t, "First", results[2].Name)
	})

	t.Run("respects limit", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("First", "2026-01-01"))
		require.NoError(t, s.InsertTournament("Second", "2026-01-08"))
		require.NoError(t, s.InsertTournament("Third", "2026-01-15"))

		results, err := s.ListTournaments(2, 0)
		require.NoError(t, err)
		require.Len(t, results, 2)
		require.Equal(t, "Third", results[0].Name)
		require.Equal(t, "Second", results[1].Name)
	})

	t.Run("respects offset", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("First", "2026-01-01"))
		require.NoError(t, s.InsertTournament("Second", "2026-01-08"))
		require.NoError(t, s.InsertTournament("Third", "2026-01-15"))

		results, err := s.ListTournaments(10, 2)
		require.NoError(t, err)
		require.Len(t, results, 1)
		require.Equal(t, "First", results[0].Name)
	})

	t.Run("returns empty slice when no tournaments", func(t *testing.T) {
		s := newTestStore(t)
		results, err := s.ListTournaments(10, 0)
		require.NoError(t, err)
		require.Empty(t, results)
	})
}

func TestStore_CountTournaments(t *testing.T) {
	t.Run("returns 0 when empty", func(t *testing.T) {
		s := newTestStore(t)
		count, err := s.CountTournaments()
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("returns correct count", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("First", "2026-01-01"))
		require.NoError(t, s.InsertTournament("Second", "2026-01-08"))

		count, err := s.CountTournaments()
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}

// --- TournamentEntry tests ---

func TestStore_InsertTournamentEntry(t *testing.T) {
	t.Run("inserts successfully", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))
		tournament, err := s.GetTournamentByShortName("2026-04-06")
		require.NoError(t, err)

		err = s.InsertTournamentEntry("guild1", tournament.Id, "user1", "username1", "Display1", 1500)
		require.NoError(t, err)
	})

	t.Run("upserts waves on conflict", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))
		tournament, err := s.GetTournamentByShortName("2026-04-06")
		require.NoError(t, err)

		require.NoError(t, s.InsertTournamentEntry("guild1", tournament.Id, "user1", "username1", "Display1", 1500))
		require.NoError(t, s.InsertTournamentEntry("guild1", tournament.Id, "user1", "username1", "Display1", 1800))

		entries, err := s.GetTournamentEntries(tournament.Id)
		require.NoError(t, err)
		require.Len(t, entries, 1)
		require.Equal(t, 1800, entries[0].Waves)
	})
}

func TestStore_GetTournamentEntries(t *testing.T) {
	t.Run("returns entries sorted by waves descending", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))
		tournament, err := s.GetTournamentByShortName("2026-04-06")
		require.NoError(t, err)

		require.NoError(t, s.InsertTournamentEntry("guild1", tournament.Id, "user1", "u1", "Alice", 1200))
		require.NoError(t, s.InsertTournamentEntry("guild1", tournament.Id, "user2", "u2", "Bob", 1800))
		require.NoError(t, s.InsertTournamentEntry("guild1", tournament.Id, "user3", "u3", "Carol", 1500))

		entries, err := s.GetTournamentEntries(tournament.Id)
		require.NoError(t, err)
		require.Len(t, entries, 3)
		require.Equal(t, 1800, entries[0].Waves)
		require.Equal(t, 1500, entries[1].Waves)
		require.Equal(t, 1200, entries[2].Waves)
	})

	t.Run("returns empty slice when no entries", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))
		tournament, err := s.GetTournamentByShortName("2026-04-06")
		require.NoError(t, err)

		entries, err := s.GetTournamentEntries(tournament.Id)
		require.NoError(t, err)
		require.Empty(t, entries)
	})
}

func TestStore_GetLatestTournamentEntries(t *testing.T) {
	t.Run("returns only entries for the latest tournament", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("First", "2026-01-01"))
		first, err := s.GetTournamentByShortName("2026-01-01")
		require.NoError(t, err)
		require.NoError(t, s.InsertTournamentEntry("g", first.Id, "user1", "u1", "Alice", 1000))

		require.NoError(t, s.InsertTournament("Second", "2026-01-08"))
		second, err := s.GetTournamentByShortName("2026-01-08")
		require.NoError(t, err)
		require.NoError(t, s.InsertTournamentEntry("g", second.Id, "user2", "u2", "Bob", 2000))

		entries, err := s.GetLatestTournamentEntries()
		require.NoError(t, err)
		require.Len(t, entries, 1)
		require.Equal(t, "user2", entries[0].UserId)
	})
}

func TestStore_GetTournamentWinner(t *testing.T) {
	t.Run("returns the player with the highest waves", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))
		tournament, err := s.GetTournamentByShortName("2026-04-06")
		require.NoError(t, err)

		require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user1", "u1", "Alice", 1200))
		require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user2", "u2", "Bob", 1800))

		winner, err := s.GetTournamentWinner(tournament.Id, 1800)
		require.NoError(t, err)
		require.Equal(t, "user2", winner.UserId)
		require.Equal(t, "Bob", winner.DisplayName)
	})
}

func TestStore_GetTournamentStats(t *testing.T) {
	t.Run("returns correct stats for populated tournament", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))
		tournament, err := s.GetTournamentByShortName("2026-04-06")
		require.NoError(t, err)

		require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user1", "u1", "Alice", 1000))
		require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user2", "u2", "Bob", 2000))

		entrants, maxWaves, avgWaves, err := s.GetTournamentStats(tournament.Id)
		require.NoError(t, err)
		require.Equal(t, 2, entrants)
		require.NotNil(t, maxWaves)
		require.Equal(t, int64(2000), *maxWaves)
		require.NotNil(t, avgWaves)
		require.InDelta(t, 1500.0, *avgWaves, 0.01)
	})

	t.Run("returns nil wave stats for empty tournament", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))
		tournament, err := s.GetTournamentByShortName("2026-04-06")
		require.NoError(t, err)

		entrants, maxWaves, avgWaves, err := s.GetTournamentStats(tournament.Id)
		require.NoError(t, err)
		require.Equal(t, 0, entrants)
		require.Nil(t, maxWaves)
		require.Nil(t, avgWaves)
	})
}

// --- User progress tests ---

func TestStore_GetUserTournamentEntries(t *testing.T) {
	t.Run("returns user entries with correct rank", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))
		tournament, err := s.GetTournamentByShortName("2026-04-06")
		require.NoError(t, err)

		require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user1", "u1", "Alice", 1800)) // rank 1
		require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user2", "u2", "Bob", 1200))   // rank 2
		require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user3", "u3", "Carol", 900))  // rank 3

		results, err := s.GetUserTournamentEntries("user2", 10, 0)
		require.NoError(t, err)
		require.Len(t, results, 1)
		require.Equal(t, 1200, results[0].Waves)
		require.Equal(t, 2, results[0].Rank)
		require.Equal(t, 3, results[0].TotalEntrants)
	})

	t.Run("returns entries across multiple tournaments in descending order", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("First", "2026-01-01"))
		first, err := s.GetTournamentByShortName("2026-01-01")
		require.NoError(t, err)
		require.NoError(t, s.InsertTournamentEntry("g", first.Id, "user1", "u1", "Alice", 1000))

		require.NoError(t, s.InsertTournament("Second", "2026-01-08"))
		second, err := s.GetTournamentByShortName("2026-01-08")
		require.NoError(t, err)
		require.NoError(t, s.InsertTournamentEntry("g", second.Id, "user1", "u1", "Alice", 1500))

		results, err := s.GetUserTournamentEntries("user1", 10, 0)
		require.NoError(t, err)
		require.Len(t, results, 2)
		require.Equal(t, "2026-01-08", results[0].TournamentShortName) // most recent first
		require.Equal(t, "2026-01-01", results[1].TournamentShortName)
	})

	t.Run("respects limit and offset", func(t *testing.T) {
		s := newTestStore(t)
		for i, date := range []string{"2026-01-01", "2026-01-08", "2026-01-15"} {
			require.NoError(t, s.InsertTournament("Tournament "+date, date))
			tournament, err := s.GetTournamentByShortName(date)
			require.NoError(t, err)
			require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user1", "u1", "Alice", 1000+i*100))
		}

		results, err := s.GetUserTournamentEntries("user1", 2, 1)
		require.NoError(t, err)
		require.Len(t, results, 2)
		require.Equal(t, "2026-01-08", results[0].TournamentShortName)
		require.Equal(t, "2026-01-01", results[1].TournamentShortName)
	})

	t.Run("returns empty slice when user has no entries", func(t *testing.T) {
		s := newTestStore(t)
		results, err := s.GetUserTournamentEntries("nonexistent", 10, 0)
		require.NoError(t, err)
		require.Empty(t, results)
	})

	t.Run("rank 1 when user has highest waves", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))
		tournament, err := s.GetTournamentByShortName("2026-04-06")
		require.NoError(t, err)

		require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user1", "u1", "Alice", 2000))
		require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user2", "u2", "Bob", 1000))

		results, err := s.GetUserTournamentEntries("user1", 10, 0)
		require.NoError(t, err)
		require.Len(t, results, 1)
		require.Equal(t, 1, results[0].Rank)
	})
}

func TestStore_CountUserTournamentEntries(t *testing.T) {
	t.Run("returns 0 when user has no entries", func(t *testing.T) {
		s := newTestStore(t)
		count, err := s.CountUserTournamentEntries("nonexistent")
		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("returns correct count across tournaments", func(t *testing.T) {
		s := newTestStore(t)
		for _, date := range []string{"2026-01-01", "2026-01-08", "2026-01-15"} {
			require.NoError(t, s.InsertTournament("Tournament "+date, date))
			tournament, err := s.GetTournamentByShortName(date)
			require.NoError(t, err)
			require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user1", "u1", "Alice", 1000))
		}

		count, err := s.CountUserTournamentEntries("user1")
		require.NoError(t, err)
		require.Equal(t, 3, count)
	})

	t.Run("only counts entries for the specified user", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-04-06"))
		tournament, err := s.GetTournamentByShortName("2026-04-06")
		require.NoError(t, err)

		require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user1", "u1", "Alice", 1000))
		require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user2", "u2", "Bob", 2000))

		count, err := s.CountUserTournamentEntries("user1")
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
}

func TestStore_GetUserTournamentEntriesInDateRange(t *testing.T) {
	t.Run("returns entries within date range", func(t *testing.T) {
		s := newTestStore(t)
		for _, date := range []string{"2026-01-01", "2026-02-01", "2026-03-01", "2026-04-01"} {
			require.NoError(t, s.InsertTournament("Tournament "+date, date))
			tournament, err := s.GetTournamentByShortName(date)
			require.NoError(t, err)
			require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user1", "u1", "Alice", 1000))
		}

		results, err := s.GetUserTournamentEntriesInDateRange("user1", "2026-02-01", "2026-03-01")
		require.NoError(t, err)
		require.Len(t, results, 2)
		require.Equal(t, "2026-03-01", results[0].TournamentShortName)
		require.Equal(t, "2026-02-01", results[1].TournamentShortName)
	})

	t.Run("returns empty when no entries in range", func(t *testing.T) {
		s := newTestStore(t)
		require.NoError(t, s.InsertTournament("Sunday Weekly", "2026-01-01"))
		tournament, err := s.GetTournamentByShortName("2026-01-01")
		require.NoError(t, err)
		require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user1", "u1", "Alice", 1000))

		results, err := s.GetUserTournamentEntriesInDateRange("user1", "2026-06-01", "2026-12-01")
		require.NoError(t, err)
		require.Empty(t, results)
	})

	t.Run("inclusive bounds — includes start and end date", func(t *testing.T) {
		s := newTestStore(t)
		for _, date := range []string{"2026-01-01", "2026-02-01", "2026-03-01"} {
			require.NoError(t, s.InsertTournament("Tournament "+date, date))
			tournament, err := s.GetTournamentByShortName(date)
			require.NoError(t, err)
			require.NoError(t, s.InsertTournamentEntry("g", tournament.Id, "user1", "u1", "Alice", 1000))
		}

		results, err := s.GetUserTournamentEntriesInDateRange("user1", "2026-01-01", "2026-03-01")
		require.NoError(t, err)
		require.Len(t, results, 3)
	})

	t.Run("capped at 30 results", func(t *testing.T) {
		s := newTestStore(t)
		// Insert 35 tournaments with entries for user1
		dates := make([]string, 35)
		for i := range dates {
			dates[i] = fmt.Sprintf("2026-%02d-%02d", (i/28)+1, (i%28)+1)
		}
		// Use simpler sequential dates
		for i := 1; i <= 35; i++ {
			date := fmt.Sprintf("2020-%02d-01", i%12+1)
			name := fmt.Sprintf("Tournament %d", i)
			// avoid duplicate short_names by using unique dates
			date = fmt.Sprintf("2020-01-%02d", i)
			if i > 28 {
				date = fmt.Sprintf("2020-02-%02d", i-28)
			}
			_ = s.InsertTournament(name, date)
			tournament, err := s.GetTournamentByShortName(date)
			if err != nil {
				continue
			}
			_ = s.InsertTournamentEntry("g", tournament.Id, "user1", "u1", "Alice", 1000)
		}

		results, err := s.GetUserTournamentEntriesInDateRange("user1", "2000-01-01", "2099-12-31")
		require.NoError(t, err)
		require.LessOrEqual(t, len(results), 30)
	})
}
