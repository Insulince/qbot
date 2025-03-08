package qbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// handleHelp displays a list of available commands.
func (q *QBot) handleHelp(m *discordgo.MessageCreate, args []string) error {
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
		"- If you need more time because your bracket isn't full yet but you are close to being timed out, use `!moretime` to give yourself more time for this purpose.\n"

	q.mustPost(m.ChannelID, helpMessage)

	if err := q.handleCommands(m, args); err != nil {
		return errors.Wrap(err, "handle commands")
	}

	if err := q.handleVersion(m, args); err != nil {
		return errors.Wrap(err, "handle version")
	}

	return nil
}
