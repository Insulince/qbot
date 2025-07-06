package qbot

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// Handle !forcewave <@user> <wave>
func (q *QBot) handleForceWave(cmd Cmd) error {
	// Check if the user is a moderator
	if !q.isModerator(cmd.Message) {
		q.mustPost(cmd.Message.ChannelID, "You do not have permission to use this command. Moderator role required.")
		return nil
	}

	// Check for correct number of arguments
	if len(cmd.Args) != 2 {
		q.mustPost(cmd.Message.ChannelID, "Usage: `!forcewave <@user> <wave>`")
		return nil
	}

	// Parse the user mention
	mention := cmd.Args[0]
	targetUserId := ""
	if len(mention) >= 3 && mention[0] == '<' && mention[1] == '@' {
		mention = mention[2:]
		if mention[0] == '!' {
			mention = mention[1:]
		}
		if mention[len(mention)-1] == '>' {
			targetUserId = mention[:len(mention)-1]
		}
	}
	if targetUserId == "" {
		q.mustPost(cmd.Message.ChannelID, "Could not parse user mention. Please use the format: `!forcewave <@user> <wave>`")
		return nil
	}

	// Parse the wave count
	wavesStr := cmd.Args[1]
	waves, err := strconv.Atoi(wavesStr)
	if err != nil {
		return errors.Wrapf(err, "converting waves %q to int", wavesStr)
	}
	if waves < 1 || waves > 10000 {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("Please enter a valid wave number between 1 and 10000."))
		return nil
	}

	// Get the latest tournament
	tournament, err := q.store.GetLatestTournament()
	if err != nil {
		return errors.Wrap(err, "get latest tournament")
	}

	// Get the target user's information
	member, err := q.session.GuildMember(cmd.GuildId, targetUserId)
	if err != nil {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("Could not find user with ID %s in this server.", targetUserId))
		return nil
	}

	username := member.User.Username
	displayName := member.Nick
	if displayName == "" {
		if member.User.GlobalName != "" {
			displayName = member.User.GlobalName
		} else {
			displayName = username
		}
	}

	// Insert the tournament entry for the target user
	if err := q.store.InsertTournamentEntry(cmd.GuildId, tournament.Id, targetUserId, username, displayName, waves); err != nil {
		return errors.Wrap(err, "inserting tournament entry")
	}

	// Send confirmation message
	q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("✅ <@%s> set <@%s>'s waves to **%d**!", cmd.Message.Author.ID, targetUserId, waves))

	return nil
}
