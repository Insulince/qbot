package qbot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Insulince/qbot/internal/models"
	"github.com/pkg/errors"
)

const progressMaxLimit = 30

// Handle !progress <@user> [args...]
func (q *QBot) handleProgress(cmd Cmd) error {
	if len(cmd.Args) == 0 {
		q.mustPost(cmd.Message.ChannelID, "Usage: `!progress <@user> [limit] [offset]` or `!progress <@user> <start-date> [end-date]`")
		return nil
	}

	// First arg is always the user mention.
	userId, ok := parseMention(cmd.Args[0])
	if !ok {
		q.mustPost(cmd.Message.ChannelID, "❌ Could not parse user mention. Please use the format: `!progress <@user>`")
		return nil
	}

	args := cmd.Args[1:]

	switch len(args) {
	case 0:
		return q.getPlayerProgress(cmd, userId, 10, 0)
	case 1:
		if limit, err := strconv.Atoi(args[0]); err == nil {
			// Single numeric arg: treat as limit.
			if limit <= 0 || limit > progressMaxLimit {
				q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("Invalid limit `%s`. Please use a positive integer between 1 and %d.", args[0], progressMaxLimit))
				return nil
			}
			return q.getPlayerProgress(cmd, userId, limit, 0)
		}
		// Single date arg: from that date to today.
		startShortName, err := parseTournamentShortName(args[0])
		if err != nil {
			q.mustPost(cmd.Message.ChannelID, errors.Wrap(err, "❌ parse start date").Error())
			return nil
		}
		endShortName := time.Now().UTC().Format("2006-01-02")
		return q.getPlayerProgressInDateRange(cmd, userId, startShortName, endShortName)
	case 2:
		limitInt, limitErr := strconv.Atoi(args[0])
		offsetInt, offsetErr := strconv.Atoi(args[1])
		if limitErr == nil && offsetErr == nil {
			// Two numeric args: limit + offset.
			if limitInt <= 0 || limitInt > progressMaxLimit {
				q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("Invalid limit `%s`. Please use a positive integer between 1 and %d.", args[0], progressMaxLimit))
				return nil
			}
			if offsetInt < 0 {
				q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("Invalid offset `%s`. Please use a non-negative integer.", args[1]))
				return nil
			}
			return q.getPlayerProgress(cmd, userId, limitInt, offsetInt)
		}
		// Two date args: explicit date range.
		startShortName, err := parseTournamentShortName(args[0])
		if err != nil {
			q.mustPost(cmd.Message.ChannelID, errors.Wrap(err, "❌ parse start date").Error())
			return nil
		}
		endShortName, err := parseTournamentShortName(args[1])
		if err != nil {
			q.mustPost(cmd.Message.ChannelID, errors.Wrap(err, "❌ parse end date").Error())
			return nil
		}
		if endShortName < startShortName {
			q.mustPost(cmd.Message.ChannelID, "❌ End date must be on or after start date.")
			return nil
		}
		return q.getPlayerProgressInDateRange(cmd, userId, startShortName, endShortName)
	default:
		q.mustPost(cmd.Message.ChannelID, "Invalid number of arguments. See `!commands` for usage.")
		return nil
	}
}

func (q *QBot) getPlayerProgress(cmd Cmd, userId string, limit, offset int) error {
	totalCount, err := q.store.CountUserTournamentEntries(userId)
	if err != nil {
		return errors.Wrap(err, "counting user tournament entries")
	}
	if totalCount == 0 {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("_No tournament results found for <@%s>._", userId))
		return nil
	}

	results, err := q.store.GetUserTournamentEntries(userId, limit, offset)
	if err != nil {
		return errors.Wrap(err, "getting user tournament entries")
	}

	var header string
	if offset == 0 {
		header = fmt.Sprintf("Latest %d tournament result(s)", len(results))
	} else {
		header = fmt.Sprintf("Tournament results %d through %d", offset+1, offset+len(results))
	}

	more := totalCount - limit - offset
	if more < 0 {
		more = 0
	}
	moreMsg := ""
	if more > 0 {
		moreMsg = fmt.Sprintf("\n_... %d more result(s) not shown..._", more)
	}

	msg := fmt.Sprintf("📊 **Player Progress: <@%s>** — %s\n%s%s", userId, header, formatProgressResults(results), moreMsg)
	q.mustPostWithoutTags(cmd.Message.ChannelID, msg)
	return nil
}

func (q *QBot) getPlayerProgressInDateRange(cmd Cmd, userId, startShortName, endShortName string) error {
	results, err := q.store.GetUserTournamentEntriesInDateRange(userId, startShortName, endShortName)
	if err != nil {
		return errors.Wrap(err, "getting user tournament entries in date range")
	}
	if len(results) == 0 {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("_No tournament results found for <@%s> in that date range._", userId))
		return nil
	}

	header := fmt.Sprintf("Results from `%s` to `%s`", startShortName, endShortName)
	msg := fmt.Sprintf("📊 **Player Progress: <@%s>** — %s\n%s", userId, header, formatProgressResults(results))
	q.mustPostWithoutTags(cmd.Message.ChannelID, msg)
	return nil
}

func formatProgressResults(results []*models.PlayerTournamentResult) string {
	var lines []string
	for _, r := range results {
		line := fmt.Sprintf("• **%s** (`%s`) — Wave `%d` | Rank `%d`/`%d`", r.TournamentName, r.TournamentShortName, r.Waves, r.Rank, r.TotalEntrants)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
