package qbot

func (q *QBot) handleCommands(cmd Cmd) error {
	commandsMessage := "" +
		"**Commands**\n" +
		"\n" +
		"🎮 **Queue Management**\n" +
		"• `!queue` - Join the queue\n" +
		"• `!enter` - (used after `!queue` on your turn) Indicates you have entered a tournament bracket and are monitoring its fill-progress\n" +
		"• `!full` - (used after `!full` on your turn) Indicates your bracket is full and removes you from the queue\n" +
		"• `!leave` - Leave the current queue\n" +
		"• `!view` - View the current queue\n" +
		"• `!position` - Check your position in queue\n" +
		"• `!moretime` - Add more time to current queue step\n" +
		"\n" +
		"🏆 **Statistics & Records**\n" +
		"• `!leaderboard` - View the current tournament leaderboard\n" +
		"• `!submitwave <number>` - Submit your wave count for the current tournament\n" +
		"• `!history [YYYY-MM-DD]` - View tournament history\n" +
		"  └ Usage: `!history` (recent tournaments)\n" +
		"  └ Usage: `!history 2024-03-15` (specific tournament)\n" +
		"• `!progress <@user> [args]` - View a player's results across tournaments\n" +
		"  └ Usage: `!progress <@user>` (latest 10)\n" +
		"  └ Usage: `!progress <@user> 5` (latest 5)\n" +
		"  └ Usage: `!progress <@user> 10 5` (limit/offset)\n" +
		"  └ Usage: `!progress <@user> 2024-01-01` (since date)\n" +
		"  └ Usage: `!progress <@user> 2024-01-01 2024-06-01` (date range)\n" +
		"\n" +
		"🛡️ **Moderator Commands**\n" +
		"• `!forcewave <@user> <number>` - Force-overwrite a player's wave count for the current tournament\n" +
		"• `!remove <@user>` - Remove a user from the queue\n" +
		"• `!reset` - Reset the queue\n" +
		"• `!skip` - Skip the current user in the queue\n" +
		"\n" +
		"ℹ️ **Help & Information**\n" +
		"• `!help` - Display help message + commands\n" +
		"• `!commands` - Display summary of available commands\n" +
		"• `!version` - Show bot version and uptime\n" +
		"\n" +
		"🎭 **Memes**\n" +
		"• `!owned` - Get owned\n" +
		"\n" +
		"**Notes:**\n" +
		"- Commands are not case-sensitive, _but I am not a fan of shouting_."

	q.mustPost(cmd.Message.ChannelID, commandsMessage)

	return nil
}
