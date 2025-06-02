package qbot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Cmd struct {
	Message *discordgo.MessageCreate
	Command string
	Args    []string
	GuildId string
}

var _ fmt.Stringer = Cmd{}

var emptyCmd = Cmd{}

func (cmd Cmd) String() string {
	return fmt.Sprintf("%s %s", cmd.Command, strings.Join(cmd.Args, " "))
}

// IsEmpty reports whether cmd represents the empty command.
func (cmd Cmd) IsEmpty() bool {
	return cmd.Command == ""
}
