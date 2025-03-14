package qbot

func (q *QBot) handleCommands(cmd Cmd) error {
	commandsMessage := "" +
		"\n**Commands:**\n" +
		"The following are all supported commands. To use these commands, just type them out and send them with no other message content, no need to mention me or add anything else to the message.\n" +
		"`!queue` - Join the waiting queue.\n" +
		"`!enter` - Signal that you've entered your bracket (when you are the active player).\n" +
		"`!full` - Signal that your bracket is full (when you are the active player).\n" +
		"`!view` - View the current state of the queue.\n" +
		"`!leave` - Leave the queue or active slot.\n" +
		"`!position` - Check your position in the queue.\n" +
		"`!moretime` - Extend your current timer if you need more time before marking your bracket as full.\n" +
		"`!help` - Show help message.\n" +
		"`!commands` - Show bot commands.\n" +
		"\n**Moderator Commands:**\n" +
		"`!skip` - Force-skip the current active user.\n" +
		"`!reset` - Clear the entire queue.\n" +
		"`!remove @user` - Remove a specific user from the queue.\n" +
		"`!version` - Display bot version\n"

	q.mustPost(cmd.Message.ChannelID, commandsMessage)

	return nil
}
