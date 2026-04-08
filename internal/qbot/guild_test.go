package qbot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGuilds_Get(t *testing.T) {
	guilds := Guilds{
		"guild1": {Name: "Test Server", AnnouncementChannelId: "chan1", AudienceIdentifier: "@everyone"},
	}

	t.Run("returns guild when found", func(t *testing.T) {
		g, err := guilds.Get("guild1")
		require.NoError(t, err)
		require.Equal(t, "Test Server", g.Name)
		require.Equal(t, "chan1", g.AnnouncementChannelId)
		require.Equal(t, "@everyone", g.AudienceIdentifier)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		_, err := guilds.Get("nonexistent")
		require.Error(t, err)
	})

	t.Run("empty guilds map returns error", func(t *testing.T) {
		empty := Guilds{}
		_, err := empty.Get("guild1")
		require.Error(t, err)
	})
}
