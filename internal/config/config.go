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
	}
)

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
