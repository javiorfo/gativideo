package config

import (
	"os"

	"github.com/pelletier/go-toml"
)

const (
	defaultYTSHost          = "https://en.yts-official.mx"
	defaultYTSQuality       = 1080
	defaultOpenSubsEnable   = false
	defaultOpenSubsLanguage = "es"
)

type config struct {
	YTSHost          string
	YTSQuality       int64
	OpenSubsEnable   bool
	OpenSubsLanguage string
}

func defaultConfig() config {
	return config{
		YTSHost:          defaultYTSHost,
		YTSQuality:       defaultYTSQuality,
		OpenSubsEnable:   defaultOpenSubsEnable,
		OpenSubsLanguage: defaultOpenSubsLanguage,
	}
}

func GetConfiguration() config {
	home, _ := os.UserHomeDir()
	tomlFile, err := toml.LoadFile(home + "/.config/bitsmuggler/config.toml")
	if err != nil {
		return defaultConfig()
	}

	return config{
		YTSHost:          tomlFile.GetDefault("yts.host", defaultYTSHost).(string),
		YTSQuality:       tomlFile.GetDefault("yts.quality", defaultYTSQuality).(int64),
		OpenSubsEnable:   tomlFile.GetDefault("opensubs.enable", defaultOpenSubsEnable).(bool),
		OpenSubsLanguage: tomlFile.GetDefault("opensubs.language", defaultOpenSubsLanguage).(string),
	}
}
