package qbot

// handleReset clears the entire queue (moderator only).
func (q *QBot) handleReset(cmd Cmd) error {
	if !q.isModerator(cmd.Message) {
		q.mustPost(cmd.Message.ChannelID, "You do not have permission to use this command. Moderator role required.")
		return nil
	}

	q.queueMutex.Lock()
	defer q.queueMutex.Unlock()
	q.queue = []QueueItem{}
	q.mustPost(cmd.Message.ChannelID, "Queue has been reset.")

	return nil
}
