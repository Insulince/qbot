package qbot

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// handleHelp displays general help, or detailed help for a specific command if one is given.
// Usage: !help [command]
func (q *QBot) handleHelp(cmd Cmd) error {
	if len(cmd.Args) > 0 {
		return q.handleHelpForCommand(cmd, cmd.Args[0])
	}

	helpMessage := "" +
		"**About**\n" +
		"Q helps mediate a tourney queue to avoid joining the same bracket as other players.\n" +
		"You should use Q when you are planning to join a bracket and want to avoid joining with other players in this server and/or want to inform other players that you are about to join.\n" +
		"The main flow is:\n" +
		"1. Enter the queue (`!queue`).\n" +
		"2. Wait for your turn (Q will ping you when its your turn).\n" +
		"3. When its your turn, enter your bracket and tell Q you have done so (`!enter`).\n" +
		"4. Watch for your bracket to be full, then tell Q so it can move on to the next player in the queue (`!full`).\n" +
		"**Disclaimer**\n" +
		"- Q is purely voluntary—there is **no obligation** to use it.\n" +
		"- Players are not at fault for ignoring or bypassing the queue.\n" +
		"- The system is meant as a **courtesy tool** to help organize tournament entries for those who choose to participate.\n" +
		"\n**Notes/Tips**\n" +
		"- If you take too long to `!enter` when you are at the front of the queue, you will be assumed AFK and kicked out of the queue and the next player will be given a turn. There is a 2 minute warning before this happens.\n" +
		"- If you take too long to `!full` to signal your bracket is filled after doing `!enter`, you will be assumed AFK and kicked out of the queue and the next player will be given a turn. There is a 2 minute warning before this happens.\n" +
		"- If you need more time because your bracket isn't full yet but you are close to being timed out, use `!moretime` to give yourself more time for this purpose.\n" +
		"\n_Tip: Use `!help <command>` for detailed help on any specific command._\n"

	q.mustPost(cmd.Message.ChannelID, helpMessage)

	if err := q.handleCommands(cmd); err != nil {
		return errors.Wrap(err, "handle commands")
	}

	if err := q.handleVersion(cmd); err != nil {
		return errors.Wrap(err, "handle version")
	}

	return nil
}

// commandHelpTexts maps each primary command alias to its detailed description.
var commandHelpTexts = map[string]string{
	"!queue": "**!queue** (aliases: `!q`, `!enqueue`, `!join`)\n" +
		"Join the tournament queue. You'll get a 🟢 pass image if you're first and can go right away, or a 🔴 block image showing who's ahead of you.\n" +
		"You can only be in the queue once at a time. If you're already in it, Q will remind you of your position.",

	"!enter": "**!enter** (aliases: `!e`, `!enterbracket`)\n" +
		"Signal that you've joined a bracket. Only valid when you're first in queue.\n" +
		"This resets your timeout timer — you then have the `fullTimeout` window to call `!full` before being skipped.",

	"!full": "**!full** (aliases: `!f`, `!bracketfull`)\n" +
		"Signal that your bracket is full and you're done. This removes you from the queue and notifies the next player.\n" +
		"You must have called `!enter` first.",

	"!leave": "**!leave** (alias: `!leavequeue`)\n" +
		"Remove yourself from the queue at any point. Has no effect if you're not in the queue.",

	"!view": "**!view** (alias: `!viewqueue`)\n" +
		"Display the current queue — who's active and who's waiting.",

	"!position": "**!position** (alias: `!currentposition`)\n" +
		"Show your current position in the queue and who is ahead of you.",

	"!moretime": "**!moretime** (aliases: `!extend`, `!moar`)\n" +
		"Reset your timeout timer. Use this if your bracket is taking longer than expected to fill.\n" +
		"Can be used repeatedly — there is no limit.",

	"!submitwave": "**!submitwave <waves>** (aliases: `!wave`, `!waves`, `!w`, `!<number>`)\n" +
		"Submit your wave count for the current tournament. Valid range: 1–100000.\n" +
		"Re-submitting overwrites your previous entry. Shorthand: `!1500` is the same as `!submitwave 1500`.",

	"!clearwave": "**!clearwave**\n" +
		"Remove your own wave entry from the current tournament leaderboard.\n" +
		"Only affects your own entry — you cannot clear someone else's.",

	"!leaderboard": "**!leaderboard** (alias: `!lb`)\n" +
		"Show the current tournament leaderboard sorted by wave count.",

	"!history": "**!history [limit] [offset]** or **!history <date>**\n" +
		"List past tournaments or show the leaderboard for a specific one.\n" +
		"Date formats: `YYYY-MM-DD`, `MM-DD`, `YYMMDD`, `YYYYMMDD`. Dates cannot be in the future.",

	"!progress": "**!progress <@user> [limit] [offset]** or **!progress <@user> <start-date> <end-date>**\n" +
		"Show a player's tournament history — wave counts, rank, and total entrants per tournament.",

	"!skip": "**!skip** _(moderator only)_\n" +
		"Force-skip the current active user, removing them from the queue and promoting the next player.",

	"!reset": "**!reset** _(moderator only)_\n" +
		"Clear the entire queue. Use with caution — this removes everyone.",

	"!remove": "**!remove <@user>** _(moderator only)_\n" +
		"Remove a specific player from the queue without affecting anyone else.",

	"!forcewave": "**!forcewave <@user> <waves>** _(moderator only)_\n" +
		"Set or override another player's wave count for the current tournament.",

	"!shame": "**!shame [@user]**\n" +
		"Post the shame GIF. Optionally tag a specific user to direct it at them. 🔔",

	"!help": "**!help [command]**\n" +
		"Show general help, or `!help <command>` for detailed help on a specific command.",

	"!commands": "**!commands**\n" +
		"List all available commands with brief descriptions.",

	"!version": "**!version**\n" +
		"Show the current bot version and uptime.",

	"!owned": "**!owned**\n" +
		"You already know what this does.",

	"!dev": "**!dev <subcommand>** _(moderator only)_\n" +
		"Developer commands for testing. Subcommands: `error`, `panic`, `displayname`, `parsetourneyname`, `guild`.",
}

func (q *QBot) handleHelpForCommand(cmd Cmd, arg string) error {
	command := strings.ToLower(arg)
	if !strings.HasPrefix(command, "!") {
		command = "!" + command
	}

	// Resolve aliases to their primary command for lookup.
	command = resolveAlias(command)

	text, ok := commandHelpTexts[command]
	if !ok {
		q.mustPost(cmd.Message.ChannelID, fmt.Sprintf("No detailed help available for `%s`. Use `!help` for general help or `!commands` for a full list.", arg))
		return nil
	}

	q.mustPost(cmd.Message.ChannelID, text)
	return nil
}

// aliasMap maps command aliases back to their primary name for help lookup.
var aliasMap = map[string]string{
	"!q": "!queue", "!enqueue": "!queue", "!join": "!queue",
	"!e": "!enter", "!enterbracket": "!enter",
	"!f": "!full", "!bracketfull": "!full",
	"!leavequeue": "!leave",
	"!viewqueue":  "!view",
	"!currentposition": "!position",
	"!extend": "!moretime", "!moar": "!moretime",
	"!wave": "!submitwave", "!waves": "!submitwave", "!w": "!submitwave",
	"!lb":          "!leaderboard",
	"!skipcurrent": "!skip",
	"!resetqueue":  "!reset",
	"!removeplayer": "!remove",
	"!submitwaves": "!submitwave",
}

// resolveAlias maps a command alias to its primary name for help lookup.
func resolveAlias(cmd string) string {
	if primary, ok := aliasMap[cmd]; ok {
		return primary
	}
	return cmd
}
