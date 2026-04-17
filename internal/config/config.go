package config

import (
	"encoding/json"
	"os"

	"github.com/Insulince/jlib/pkg/jmust"
	"github.com/pkg/errors"
)

const (
	EnvVarConfigFilePath = "CONFIG_FILE_PATH"
)

type (
	Config struct {
		DiscordBotToken string `json:"discordBotToken"`

		EnterTimeoutMinutes  int `json:"enterTimeoutMinutes"`
		FullTimeoutMinutes   int `json:"fullTimeoutMinutes"`
		WarnThresholdMinutes int `json:"warnThresholdMinutes"`

		ErrorChannelId     string `json:"errorChannelId"`
		NotificationRoleId string `json:"notificationRoleId"`

		DbDriver string `json:"dbDriver"`
		DbFile   string `json:"dbFile"`

		Guilds map[string]GuildConfig `json:"guilds"`
	}

	GuildConfig struct {
		Name                  string `json:"name"`
		AnnouncementChannelId string `json:"announcementChannelId"`
		AudienceIdentifier    string `json:"audienceIdentifier"`
		ModeratorRoleName     string `json:"moderatorRoleName"` // defaults to "Moderator" if empty
	}
)

// Validate returns an error if any required config field is missing or logically invalid.
func (c Config) Validate() error {
	if c.DiscordBotToken == "" {
		return errors.New("discordBotToken is required")
	}
	if c.EnterTimeoutMinutes <= 0 {
		return errors.Errorf("enterTimeoutMinutes must be positive, got %d", c.EnterTimeoutMinutes)
	}
	if c.FullTimeoutMinutes <= 0 {
		return errors.Errorf("fullTimeoutMinutes must be positive, got %d", c.FullTimeoutMinutes)
	}
	if c.WarnThresholdMinutes <= 0 {
		return errors.Errorf("warnThresholdMinutes must be positive, got %d", c.WarnThresholdMinutes)
	}
	if c.WarnThresholdMinutes >= c.EnterTimeoutMinutes {
		return errors.Errorf("warnThresholdMinutes (%d) must be less than enterTimeoutMinutes (%d)", c.WarnThresholdMinutes, c.EnterTimeoutMinutes)
	}
	if c.ErrorChannelId == "" {
		return errors.New("errorChannelId is required")
	}
	if c.NotificationRoleId == "" {
		return errors.New("notificationRoleId is required")
	}
	if c.DbDriver == "" {
		return errors.New("dbDriver is required")
	}
	if c.DbFile == "" {
		return errors.New("dbFile is required")
	}
	if len(c.Guilds) == 0 {
		return errors.New("at least one guild must be configured")
	}
	for id, g := range c.Guilds {
		if g.Name == "" {
			return errors.Errorf("guild %q: name is required", id)
		}
		if g.AnnouncementChannelId == "" {
			return errors.Errorf("guild %q: announcementChannelId is required", id)
		}
		if g.AudienceIdentifier == "" {
			return errors.Errorf("guild %q: audienceIdentifier is required", id)
		}
	}
	return nil
}

func GetConfig() (Config, error) {
	path, found := os.LookupEnv(EnvVarConfigFilePath)
	if !found {
		return Config{}, errors.Errorf("environment variable %q not set", EnvVarConfigFilePath)
	}

	f, err := os.Open(path)
	if err != nil {
		return Config{}, errors.Wrapf(err, "open config file %q", path)
	}
	defer jmust.MustClose(f)

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, errors.Wrapf(err, "decode config %q", path)
	}

	return cfg, nil
}

func MustGetConfig() Config {
	return jmust.Must[Config](GetConfig)[0]
}
