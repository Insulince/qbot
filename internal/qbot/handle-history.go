package qbot

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func (q *QBot) handleHistory(cmd Cmd) error {
	if len(cmd.Args) == 0 {
		// If no args provided, just grab the default most recent list of tourneys
		if err := q.getTourneyList(cmd); err != nil {
			return errors.Wrap(err, "getting tourney list")
		}
		return nil
	}

	// If an arg was provided, we need to infer which approach is being requested, list of tourneys, or specific tourney leaderboard.
	if _, err := strconv.Atoi(cmd.Args[0]); err != nil {
		// If it can't be parsed, this is requesting a date, so fetch a specific tourney's leaderboard
		err := q.getTourneysHistory(cmd)
		if err != nil {
			return errors.Wrap(err, "getting tourney history")
		}
		return nil
	}
	// If it can be parsed, then the arg is a limit, fetch a list of tourneys
	if err := q.getTourneyList(cmd); err != nil {
		return errors.Wrap(err, "getting tourney list")
	}

	return nil
}

func (q *QBot) getTourneyList(cmd Cmd) error {
	const (
		defaultLimit = 10
		maxLimit     = 30

		defaultOffset = 0
		maxOffset     = 50000 // 480 years worth of tourneys, should get the job done.
	)
	limit := defaultLimit
	offset := defaultOffset

	// Parse limit from user args, if given
	if len(cmd.Args) >= 1 {
		userLimitStr := cmd.Args[0]
		userLimit, err := strconv.Atoi(userLimitStr)
		if err != nil || userLimit <= 0 || userLimit > maxLimit {
			q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("Invalid limit `%s`. Please use a positive integer between 1 and %d.", userLimitStr, maxLimit))
			return errors.Wrapf(err, "parsing limit %q", userLimitStr)
		}
		limit = userLimit
	}
	if len(cmd.Args) == 2 {
		userOffsetStr := cmd.Args[1]
		userOffset, err := strconv.Atoi(userOffsetStr)
		if err != nil || userOffset < 0 {
			q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("Invalid offset `%s`. Please use an integer between 0 and %d.", userOffsetStr, maxOffset))
			return errors.Wrapf(err, "parsing offset %q", userOffsetStr)
		}
		offset = userOffset
	}
	if len(cmd.Args) > 2 {
		q.mustPost(cmd.Message.ChannelID, "Invalid number of arguments, 0-2 arguments are supported (limit and offset).")
		return nil
	}

	// Get the total count of tournaments for "more" count
	totalCount, err := q.store.CountTournaments()
	if err != nil {
		return errors.Wrap(err, "counting tournaments")
	}

	tournaments, err := q.store.ListTournaments(limit, offset)
	if err != nil {
		return errors.Wrap(err, "listing tournaments")
	}

	var lines []string
	for _, tournament := range tournaments {
		entrants, maxWaves, averageWaves, err := q.store.GetTournamentStats(cmd.GuildId, tournament.Id)
		if err != nil {
			return errors.Wrap(err, "getting tournament stats")
		}

		// Query for winner display name
		winnerName := "_No winner_"
		winnerDisplayName := "_No winner_"
		winnerWaveCount := 0
		if maxWaves != nil && entrants > 0 {
			tournamentWinner, err := q.store.GetTournamentWinner(cmd.GuildId, tournament.Id, *maxWaves)
			if err != nil {
				return errors.Wrap(err, "get tournament winner")
			}

			winnerName = fmt.Sprintf("<@%s>", tournamentWinner.UserId)
			winnerDisplayName = tournamentWinner.DisplayName
			winnerWaveCount = int(*maxWaves)
		}

		avgWaveText := "_N/A_"
		if averageWaves != nil {
			avgWaveText = fmt.Sprintf("%d", int(math.Round(*averageWaves)))
		}

		line := fmt.Sprintf("• **%s** (`%s`) — Entrants: `%2d` (Avg `%4s`) | Winner: **%s** (%s) (`%4d`)", tournament.Name, tournament.ShortName, entrants, avgWaveText, winnerName, winnerDisplayName, winnerWaveCount)
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		q.mustPost(cmd.Message.ChannelID, "_No tournaments found._")
		return nil
	}

	more := totalCount - limit - offset
	if more < 0 {
		more = 0
	}

	moreMsg := ""
	if more > 0 {
		moreMsg = fmt.Sprintf("\n_... %d more tournaments not shown..._", more)
	}

	var latestMsg string
	if offset == 0 {
		latestMsg = fmt.Sprintf("Latest %d Tournaments", len(lines))
	} else {
		latestMsg = fmt.Sprintf("Most recent %d through %d Tournaments", offset+1, offset+len(lines))
	}
	msg := fmt.Sprintf("📋 **%s**\n%s%s\n\n_Use `!history <YYYY-MM-DD>` to view results for a specific tournament._", latestMsg, strings.Join(lines, "\n"), moreMsg)

	q.mustPostWithoutTags(cmd.Message.ChannelID, msg)
	return nil
}

func (q *QBot) getTourneysHistory(cmd Cmd) error {
	if len(cmd.Args) != 1 {
		q.mustPost(cmd.Message.ChannelID, "Invalid number of arguments. Please provide only a single argument, the date of the tournament you want to view.")
		return nil
	}
	givenShortName := cmd.Args[0]

	tournamentShortName, err := parseTournamentShortName(givenShortName)
	if err != nil {
		q.mustPost(cmd.Message.ChannelID, errors.Wrap(err, "❌ parse tournament identifier").Error())
		return nil
	}

	tournament, err := q.store.GetTournamentByShortName(tournamentShortName)
	if err != nil {
		return errors.Wrap(err, "getting tournament by short name")
	}

	tournamentEntries, err := q.store.GetTournamentEntries(cmd.GuildId, tournament.Id)
	if err != nil {
		return errors.Wrap(err, "getting tournament entries")
	}

	leaderboardMsg := fmt.Sprintf("🏆 **Tournament %s Leaderboard** 🏆\n", tournament.Name)
	var entries []string
	for i, tournamentEntry := range tournamentEntries {

		entry := fmt.Sprintf("%d. **<@%s>** (%s) - Wave %d", i+1, tournamentEntry.UserId, tournamentEntry.DisplayName, tournamentEntry.Waves)
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
	default:
		return "", fmt.Errorf("invalid format %q, `YYYY-MM-DD` or `MM-DD` expected for past tournaments", input)
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
		if month > currentMonth {
			if yearPart == "" {
				year -= 1
			} else {
				return "", errors.Errorf("date is in the future")
			}
		}
	}

	return fmt.Sprintf("%04d-%02d-%02d", year, month, day), nil
}
