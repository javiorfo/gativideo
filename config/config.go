package config

import (
	"os"

	"github.com/pelletier/go-toml"
)

const defaultYTSHost = "https://en.yts-official.mx"
const defaultYTSQuality = 1080

type config struct {
	YTSHost    string
	YTSQuality int64
}

func defaultConfig() config {
	return config{
		YTSHost:    defaultYTSHost,
		YTSQuality: defaultYTSQuality,
	}
}

func GetConfig() config {
	home, _ := os.UserHomeDir()
	tomlFile, err := toml.LoadFile(home + "/.config/bitsmuggler/config.toml")
	if err != nil {
		return defaultConfig()
	}

	return config{
		YTSHost:    tomlFile.GetDefault("yts.host", defaultYTSHost).(string),
		YTSQuality: tomlFile.GetDefault("yts.quality", defaultYTSQuality).(int64),
	}
}
