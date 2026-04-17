package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func validConfig() Config {
	return Config{
		DiscordBotToken:      "token",
		EnterTimeoutMinutes:  15,
		FullTimeoutMinutes:   30,
		WarnThresholdMinutes: 2,
		ErrorChannelId:       "chan-errors",
		NotificationRoleId:   "role-notify",
		DbDriver:             "sqlite",
		DbFile:               "/data/qbot.db",
		Guilds: map[string]GuildConfig{
			"guild1": {
				Name:                  "Test Server",
				AnnouncementChannelId: "chan-announce",
				AudienceIdentifier:    "@everyone",
			},
		},
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Run("valid config returns no error", func(t *testing.T) {
		require.NoError(t, validConfig().Validate())
	})

	t.Run("missing discordBotToken", func(t *testing.T) {
		cfg := validConfig()
		cfg.DiscordBotToken = ""
		require.Error(t, cfg.Validate())
	})

	t.Run("zero enterTimeoutMinutes", func(t *testing.T) {
		cfg := validConfig()
		cfg.EnterTimeoutMinutes = 0
		require.Error(t, cfg.Validate())
	})

	t.Run("negative enterTimeoutMinutes", func(t *testing.T) {
		cfg := validConfig()
		cfg.EnterTimeoutMinutes = -5
		require.Error(t, cfg.Validate())
	})

	t.Run("zero fullTimeoutMinutes", func(t *testing.T) {
		cfg := validConfig()
		cfg.FullTimeoutMinutes = 0
		require.Error(t, cfg.Validate())
	})

	t.Run("zero warnThresholdMinutes", func(t *testing.T) {
		cfg := validConfig()
		cfg.WarnThresholdMinutes = 0
		require.Error(t, cfg.Validate())
	})

	t.Run("warnThreshold equal to enterTimeout is invalid", func(t *testing.T) {
		cfg := validConfig()
		cfg.WarnThresholdMinutes = cfg.EnterTimeoutMinutes
		require.Error(t, cfg.Validate())
	})

	t.Run("warnThreshold greater than enterTimeout is invalid", func(t *testing.T) {
		cfg := validConfig()
		cfg.WarnThresholdMinutes = cfg.EnterTimeoutMinutes + 1
		require.Error(t, cfg.Validate())
	})

	t.Run("missing errorChannelId", func(t *testing.T) {
		cfg := validConfig()
		cfg.ErrorChannelId = ""
		require.Error(t, cfg.Validate())
	})

	t.Run("missing notificationRoleId", func(t *testing.T) {
		cfg := validConfig()
		cfg.NotificationRoleId = ""
		require.Error(t, cfg.Validate())
	})

	t.Run("missing dbDriver", func(t *testing.T) {
		cfg := validConfig()
		cfg.DbDriver = ""
		require.Error(t, cfg.Validate())
	})

	t.Run("missing dbFile", func(t *testing.T) {
		cfg := validConfig()
		cfg.DbFile = ""
		require.Error(t, cfg.Validate())
	})

	t.Run("empty guilds map", func(t *testing.T) {
		cfg := validConfig()
		cfg.Guilds = map[string]GuildConfig{}
		require.Error(t, cfg.Validate())
	})

	t.Run("nil guilds map", func(t *testing.T) {
		cfg := validConfig()
		cfg.Guilds = nil
		require.Error(t, cfg.Validate())
	})

	t.Run("guild missing name", func(t *testing.T) {
		cfg := validConfig()
		cfg.Guilds["guild1"] = GuildConfig{
			AnnouncementChannelId: "chan-announce",
			AudienceIdentifier:    "@everyone",
		}
		require.Error(t, cfg.Validate())
	})

	t.Run("guild missing announcementChannelId", func(t *testing.T) {
		cfg := validConfig()
		cfg.Guilds["guild1"] = GuildConfig{
			Name:               "Test Server",
			AudienceIdentifier: "@everyone",
		}
		require.Error(t, cfg.Validate())
	})

	t.Run("guild missing audienceIdentifier", func(t *testing.T) {
		cfg := validConfig()
		cfg.Guilds["guild1"] = GuildConfig{
			Name:                  "Test Server",
			AnnouncementChannelId: "chan-announce",
		}
		require.Error(t, cfg.Validate())
	})

	t.Run("moderatorRoleName is optional — empty is valid", func(t *testing.T) {
		cfg := validConfig()
		cfg.Guilds["guild1"] = GuildConfig{
			Name:                  "Test Server",
			AnnouncementChannelId: "chan-announce",
			AudienceIdentifier:    "@everyone",
			ModeratorRoleName:     "", // explicitly empty — falls back to "Moderator"
		}
		require.NoError(t, cfg.Validate())
	})
}
