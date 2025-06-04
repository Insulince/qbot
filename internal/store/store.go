package store

import (
	"database/sql"
	"io"

	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/pkg/errors"

	"github.com/Insulince/qbot/internal/models"
	"github.com/Insulince/qbot/internal/qbot"
)

type Store struct {
	db *sql.DB
}

var _ qbot.Store = Store{}

var _ io.Closer = Store{}

func New(db *sql.DB) (*Store, error) {
	s := new(Store)

	s.db = db

	return s, nil
}

func MustNew(db *sql.DB) *Store {
	return jmust.Must[*Store](New, db)[0]
}

func (s Store) Close() error {
	if err := s.db.Close(); err != nil {
		return errors.Wrap(err, "db close")
	}

	return nil
}

func (s Store) MustClose() {
	jmust.MustClose(s)
}

func (s Store) InsertTournament(name, shortName string) error {
	const newTournamentSql = `
INSERT INTO tournaments
    (name, short_name)
VALUES
    (?, ?)
;`

	if _, err := s.db.Exec(newTournamentSql, name, shortName); err != nil {
		return errors.Wrap(err, "❌ Error inserting new tournament")
	}

	return nil
}

func (s Store) GetLatestTournament() (*models.Tournament, error) {
	const getLatestTournamentSql = `
SELECT
    MAX(id)
FROM tournaments
;`

	var tournamentId int
	if err := s.db.QueryRow(getLatestTournamentSql).Scan(&tournamentId); err != nil {
		return nil, errors.Wrap(err, "get latest tournament")
	}

	tournament := new(models.Tournament)
	tournament.Id = int64(tournamentId)

	return tournament, nil
}

func (s Store) GetTournamentByShortName(shortName string) (*models.Tournament, error) {
	const fetchTournamentSql = `
SELECT
	id,
	name
FROM tournaments
WHERE TRUE
	AND short_name = ?
;`

	var tournamentId int64
	var tournamentName string
	if err := s.db.QueryRow(fetchTournamentSql, shortName).Scan(&tournamentId, &tournamentName); err != nil {
		return nil, errors.Wrap(err, "query row")
	}

	tournament := new(models.Tournament)
	tournament.Id = tournamentId
	tournament.Name = tournamentName

	return tournament, nil
}

func (s Store) ListTournaments(limit, offset int) ([]*models.Tournament, error) {
	const listTournamentsSql = `
SELECT
	id,
	name,
	short_name
FROM tournaments
ORDER BY id DESC
LIMIT ?
OFFSET ?
;`

	rows, err := s.db.Query(listTournamentsSql, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "query latest tournaments")
	}
	defer jmust.MustClose(rows)

	var tournaments []*models.Tournament
	for rows.Next() {
		var id int64
		var name, shortName string
		if err := rows.Scan(&id, &name, &shortName); err != nil {
			return nil, errors.Wrap(err, "scanning tournament list")
		}

		tournament := new(models.Tournament)
		tournament.Id = id
		tournament.Name = name
		tournament.ShortName = shortName
		tournaments = append(tournaments, tournament)
	}

	return tournaments, nil
}

func (s Store) CountTournaments() (int, error) {
	const countTournamentsSql = `
SELECT
	COUNT(*)
FROM tournaments
;`

	var totalCount int
	if err := s.db.QueryRow(countTournamentsSql).Scan(&totalCount); err != nil {
		return 0, errors.Wrap(err, "count tournaments")
	}

	return totalCount, nil
}

func (s Store) InsertTournamentEntry(guildId string, tournamentId int64, userId, username, displayName string, waves int) error {
	insertWaveSql := `
INSERT INTO tournament_entries
    (guild_id, tournament_id, user_id, username, waves, display_name)
VALUES
    (?, ?, ?, ?, ?, ?)
ON CONFLICT (tournament_id, user_id) DO UPDATE SET
    waves = excluded.waves,
    guild_id = excluded.guild_id
;`

	if _, err := s.db.Exec(insertWaveSql, guildId, tournamentId, userId, username, waves, displayName); err != nil {
		return errors.Wrap(err, "exec query")
	}

	return nil
}

func (s Store) GetTournamentEntries(tournamentId int64) ([]*models.TournamentEntry, error) {
	const getTournamentEntriesSql = `
SELECT
    user_id,
    username,
    waves,
    display_name
FROM tournament_entries
WHERE TRUE
	AND tournament_id = ?
ORDER BY waves DESC
;`

	tournamentEntriesRows, err := s.db.Query(getTournamentEntriesSql, tournamentId)
	if err != nil {
		return nil, errors.Wrap(err, "get tournament entries")
	}
	defer jmust.MustClose(tournamentEntriesRows)

	var tournamentEntries []*models.TournamentEntry
	for tournamentEntriesRows.Next() {
		var userId, username, displayName string
		var waves int
		if err := tournamentEntriesRows.Scan(&userId, &username, &waves, &displayName); err != nil {
			return nil, errors.Wrap(err, "scanning tournament entries rows")
		}

		tournamentEntry := new(models.TournamentEntry)
		tournamentEntry.UserId = userId
		tournamentEntry.Username = username
		tournamentEntry.Waves = waves
		tournamentEntry.DisplayName = displayName
		tournamentEntries = append(tournamentEntries, tournamentEntry)
	}

	return tournamentEntries, nil
}

func (s Store) GetLatestTournamentEntries() ([]*models.TournamentEntry, error) {
	const getLatestTournamentEntriesSql = `
SELECT
    user_id,
    username,
    waves,
    display_name
FROM tournament_entries
WHERE TRUE
	AND tournament_id = (SELECT MAX(id) FROM tournaments)
ORDER BY waves DESC
;`

	rows, err := s.db.Query(getLatestTournamentEntriesSql)
	if err != nil {
		return nil, errors.Wrap(err, "querying latest tournament entries")
	}
	defer jmust.MustClose(rows)

	var tournamentEntries []*models.TournamentEntry
	for rows.Next() {
		var userId, username, displayName string
		var waves int
		if err := rows.Scan(&userId, &username, &waves, &displayName); err != nil {
			return nil, errors.Wrap(err, "scanning tournament entries rows")
		}

		tournamentEntry := new(models.TournamentEntry)
		tournamentEntry.UserId = userId
		tournamentEntry.Username = username
		tournamentEntry.Waves = waves
		tournamentEntry.DisplayName = displayName
		tournamentEntries = append(tournamentEntries, tournamentEntry)
	}

	return tournamentEntries, nil
}

func (s Store) GetTournamentWinner(tournamentId, maxWaves int64) (*models.TournamentEntry, error) {
	const getTournamentWinnerSql = `
SELECT
	user_id,
	display_name
FROM tournament_entries
WHERE TRUE
	AND tournament_id = ?
	AND waves = ?
LIMIT 1
;`

	var userId string
	var displayName string
	err := s.db.QueryRow(getTournamentWinnerSql, tournamentId, maxWaves).Scan(&userId, &displayName)
	if err != nil {
		return nil, errors.Wrap(err, "querying tournament winner")
	}

	tournamentEntry := new(models.TournamentEntry)
	tournamentEntry.UserId = userId
	tournamentEntry.DisplayName = displayName

	return tournamentEntry, nil
}

func (s Store) GetTournamentStats(tournamentId int64) (entrants int, maxWaves *int64, averageWaves *float64, _ error) {
	const tournamentStatsSql = `
SELECT
	COUNT(*) as entrants,
	MAX(waves) as max_waves,
	AVG(waves) as avg_waves
FROM tournament_entries
WHERE TRUE
	AND tournament_id = ?
;`

	var maxWavesNullable sql.NullInt64
	var averageWavesNullable sql.NullFloat64
	if err := s.db.QueryRow(tournamentStatsSql, tournamentId).Scan(&entrants, &maxWavesNullable, &averageWavesNullable); err != nil {
		return 0, nil, nil, errors.Wrap(err, "query tournament stats")
	}

	if maxWavesNullable.Valid {
		maxWaves = &maxWavesNullable.Int64
	}
	if averageWavesNullable.Valid {
		averageWaves = &averageWavesNullable.Float64
	}

	return entrants, maxWaves, averageWaves, nil
}
