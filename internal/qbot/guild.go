package qbot

import (
	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/pkg/errors"
)

type Guilds map[string]Guild

type Guild struct {
	Name                  string
	AnnouncementChannelId string
	AudienceIdentifier    string
}

func (gs Guilds) Get(guildId string) (Guild, error) {
	g, ok := gs[guildId]
	if !ok {
		return Guild{}, errors.Errorf("guild not found wtih id %q", guildId)
	}
	return g, nil
}

func (gs Guilds) MustGet(guildId string) Guild {
	return jmust.Must[Guild](gs.Get, guildId)[0]
}
