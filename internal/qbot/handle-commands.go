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
