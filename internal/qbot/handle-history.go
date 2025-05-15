package qbot

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/pkg/errors"
)

/*
SELECT
	id,
	name
FROM tournaments
WHERE short_name = "2025-05-14";

23|14 MAY 2025

SELECT
	user_id,
	waves
FROM tournament_entries
WHERE tournament_id = 23
ORDER BY waves DESC;
*/

func (q *QBot) handleHistory(cmd Cmd) error {
	if !q.isModerator(cmd.Message) {
		return nil
	}

	if len(cmd.Args) != 1 {
		q.mustPost(cmd.Message.ChannelID, "Usage: `!history <tournament-identifier>`")
		return nil
	}
	givenShortName := cmd.Args[0]

	tournamentShortName, err := parseTournamentShortName(givenShortName)
	if err != nil {
		q.mustPost(cmd.Message.ChannelID, errors.Wrap(err, "❌ parse tournament identifier").Error())
		return nil
	}

	const fetchTournamentSql = `
SELECT
	id,
	name
FROM tournaments
WHERE short_name = ?;
`
	var tournamentId int64
	var tournamentName string
	if err := q.db.QueryRow(fetchTournamentSql, tournamentShortName).Scan(&tournamentId, &tournamentName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("_No tournament found for %q_", tournamentShortName))
			return nil
		}
		q.mustPost(cmd.Message.ChannelID, "Error retrieving tournament history.")
		return errors.Wrap(err, "query row")
	}

	const fetchTournamentEntriesSql = `
SELECT
	user_id,
	waves
FROM tournament_entries
WHERE tournament_id = ?
ORDER BY waves DESC;
`
	tournamentEntriesRows, err := q.db.Query(fetchTournamentEntriesSql, tournamentId)
	if err != nil {
		q.mustPost(cmd.Message.ChannelID, "Error retrieving tournament entries.")
		return errors.Wrap(err, "query")
	}
	defer jmust.MustClose(tournamentEntriesRows)

	leaderboardMsg := fmt.Sprintf("🏆 **Tournament %s Leaderboard** 🏆\n", tournamentName)
	var entries []string
	for i := 1; tournamentEntriesRows.Next(); i++ {
		var userId string
		var waves int
		if err := tournamentEntriesRows.Scan(&userId, &waves); err != nil {
			return errors.Wrap(err, "scanning tournament entries rows")
		}

		entry := fmt.Sprintf("%d. **<@%s>** - Wave %d", i, userId, waves)
		entries = append(entries, entry)
	}

	entriesMsg := "_No entries._"
	if len(entries) > 0 {
		entriesMsg = strings.Join(entries, "\n")
	}
	leaderboardMsg += entriesMsg

	q.mustPostWithoutTags(cmd.Message.ChannelID, leaderboardMsg)

	return nil
}

// Parse user input for !history command
func parseTournamentShortName(input string) (string, error) {
	input = strings.TrimSpace(input)

	// Normalize separators
	const separators = "-/\\._"
	for _, sep := range separators {
		input = strings.ReplaceAll(input, string(sep), "-")
	}

	// Handle raw 6-digit (YYMMDD) or 8-digit (YYYYMMDD) formats
	if !strings.Contains(input, "-") {
		switch len(input) {
		case 6: // Convert YYMMDD → YYYY-MM-DD
			yy := input[:2]
			mm := input[2:4]
			dd := input[4:6]
			input = fmt.Sprintf("20%02s-%02s-%02s", yy, mm, dd)
		case 8: // Convert YYYYMMDD → YYYY-MM-DD
			yyyy := input[:4]
			mm := input[4:6]
			dd := input[6:8]
			input = fmt.Sprintf("%04s-%02s-%02s", yyyy, mm, dd)
		case 2, 1:
			// These cases are when just a single day is provided as either a single digit number, double digit number, or 0 padded single digit number.
		default:
			return "", errors.Errorf("Invalid tournament short name: %s", input)
		}
	}

	now := time.Now().UTC()
	currentYear := now.Year()
	currentMonth := int(now.Month())
	currentDay := now.Day()

	year := currentYear
	month := currentMonth
	day := currentDay

	// Split into parts
	parts := strings.Split(input, "-")
	var yearPart, monthPart, dayPart string
	switch len(parts) {
	case 3:
		yearPart = parts[0]
		monthPart = parts[1]
		dayPart = parts[2]
		if yearPart == "" || monthPart == "" || dayPart == "" {
			return "", errors.Errorf("Invalid tournament short name: %s", input)
		}
	case 2:
		monthPart = parts[0]
		dayPart = parts[1]
		if monthPart == "" || dayPart == "" {
			return "", errors.Errorf("Invalid tournament short name: %s", input)
		}
	case 1:
		dayPart = parts[0]
		if dayPart == "" {
			return "", errors.Errorf("Invalid tournament short name: %s", input)
		}
	default:
		return "", fmt.Errorf("invalid format %q, `YYYY-MM-DD`, `MM-DD`, or `DD` expected for past tournaments", input)
	}

	var err error
	if yearPart != "" {
		if len(yearPart) == 2 {
			yearPart = fmt.Sprintf("20%s", yearPart)
		}
		if len(yearPart) != 4 {
			return "", errors.Errorf("Invalid year: %q", yearPart)
		}
		year, err = strconv.Atoi(yearPart)
		if err != nil {
			return "", errors.Wrapf(err, "cannot parse year %q", yearPart)
		}
		if year < 2021 || year > currentYear {
			return "", errors.Errorf("year %d is out of range", year)
		}
	}
	if monthPart != "" {
		if len(monthPart) > 2 {
			return "", errors.Errorf("Invalid month: %q", monthPart)
		}
		month, err = strconv.Atoi(monthPart)
		if err != nil {
			return "", errors.Wrapf(err, "cannot parse month %q", monthPart)
		}
		if month < 1 || month > 12 {
			return "", errors.Errorf("month %d is out of range", month)
		}
		if month > currentMonth {
			if yearPart == "" {
				year -= 1
			} else {
				return "", errors.Errorf("date is in the future")
			}
		}
	}
	if len(dayPart) > 2 {
		return "", errors.Errorf("Invalid day: %q", dayPart)
	}
	day, err = strconv.Atoi(dayPart)
	if err != nil {
		return "", errors.Wrapf(err, "cannot parse day %q", dayPart)
	}
	if day < 1 || day > 31 {
		return "", errors.Errorf("day %d is out of range", day)
	}
	if day > currentDay {
		if monthPart == "" {
			// 30 -> 2025-02-30
			month -= 1
			if month == 0 {
				month = 12
				year -= 1
			}
		} else {
			if month > currentMonth {
				if yearPart == "" {
					year -= 1
				} else {
					return "", errors.Errorf("date is in the future")
				}
			}
		}
	}

	return fmt.Sprintf("%04d-%02d-%02d", year, month, day), nil
}
