package qbot

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"strings"
	"time"
)

func handleHistory(session *discordgo.Session, message *discordgo.MessageCreate, args []string, db *sql.DB) {
	if !isModerator(session, message) {
		return
	}

	if len(args) < 2 {
		session.ChannelMessageSend(message.ChannelID, "Usage: `!history <tournament-identifier>`")
		return
	}

	tournamentShortName, err := parseTournamentShortName(args[1])
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("❌ %s", err.Error()))
		return
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
	if err := db.QueryRow(fetchTournamentSql, tournamentShortName).Scan(&tournamentId, &tournamentName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("_No tournament found for %q_", tournamentShortName))
			return
		}
		log.Println("Error retrieving tournament history:", err)
		session.ChannelMessageSend(message.ChannelID, "Error retrieving tournament history.")
		return
	}

	const fetchTournamentEntriesSql = `
SELECT
	user_id,
	waves
FROM tournament_entries
WHERE tournament_id = ?
ORDER BY waves DESC;
`
	tournamentEntriesRows, err := db.Query(fetchTournamentEntriesSql, tournamentId)
	if err != nil {
		log.Println("Error retrieving tournament history:", err)
		session.ChannelMessageSend(message.ChannelID, "Error retrieving tournament history.")
		return
	}
	defer tournamentEntriesRows.Close()

	leaderboardMsg := fmt.Sprintf("🏆 **Tournament %s Leaderboard** 🏆\n", tournamentName)
	var entries []string
	for i := 1; tournamentEntriesRows.Next(); i++ {
		var userId string
		var waves int
		tournamentEntriesRows.Scan(&userId, &waves)
		entry := fmt.Sprintf("%d. **<@%s>** - Wave %d\n", i, userId, waves)
		entries = append(entries, entry)
	}

	entriesMsg := "_No entries._"
	if len(entries) > 0 {
		entriesMsg = strings.Join(entries, "\n")
	}
	leaderboardMsg += entriesMsg

	session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Content: leaderboardMsg,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{}, // Prevents pinging
		},
	})
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
		return "", fmt.Errorf("Invalid format %q. Use `YYYY-MM-DD`, `MM-DD`, `DD`, or `-N` for past tournaments.", input)
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

// Fetches the N-th most recent tournament from the database
func fetchNthMostRecentTournament(db *sql.DB, offset int) (string, error) {
	query := "SELECT short_id FROM tournaments ORDER BY start DESC LIMIT 1 OFFSET ?"
	var tournamentID string
	err := db.QueryRow(query, -offset-1).Scan(&tournamentID)
	if err != nil {
		return "", fmt.Errorf("No tournament found for offset %d", offset)
	}
	return tournamentID, nil
}
