package config

import (
	"os"

	"github.com/pelletier/go-toml"
)

const (
	defaultYTSHost               = "https://en.yts-official.mx"
	defaultYTSQuality            = 1080
	defaultYTSOrderBy            = "rating"
	defaultOpenSubsLanguage      = "es"
	defaultTableBorderColor      = "240"
	defaultTableSelectionFgColor = "15"
	defaultTableSelectionBgColor = "240"
	defaultSpinnerColor          = "15"
	defaultDownloadColor         = "250"
)

type config struct {
	YTSHost               string
	YTSQuality            int64
	YTSInitSearch         bool
	YTSOrderBy            string
	OpenSubsEnable        bool
	OpenSubsLanguage      string
	DownloadFolder        string
	TableBorderColor      string
	TableSelectionFgColor string
	TableSelectionBgColor string
	SpinnerColor          string
	DownloadColor         string
}

func defaultConfig() config {
	return config{
		YTSHost:               defaultYTSHost,
		YTSQuality:            defaultYTSQuality,
		YTSInitSearch:         true,
		YTSOrderBy:            defaultYTSOrderBy,
		OpenSubsLanguage:      defaultOpenSubsLanguage,
		TableBorderColor:      defaultTableBorderColor,
		TableSelectionFgColor: defaultTableSelectionFgColor,
		TableSelectionBgColor: defaultTableSelectionBgColor,
		SpinnerColor:          defaultSpinnerColor,
		DownloadColor:         defaultDownloadColor,
	}
}

func GetConfiguration() config {
	home, _ := os.UserHomeDir()
	tomlFile, err := toml.LoadFile(home + "/.config/bitsmuggler/config.toml")
	if err != nil {
		return defaultConfig()
	}

	return config{
		YTSHost:               tomlFile.GetDefault("yts.host", defaultYTSHost).(string),
		YTSQuality:            tomlFile.GetDefault("yts.quality", defaultYTSQuality).(int64),
		YTSInitSearch:         tomlFile.GetDefault("yts.init_search", true).(bool),
		YTSOrderBy:            tomlFile.GetDefault("yts.order_by", defaultYTSOrderBy).(string),
		OpenSubsEnable:        tomlFile.GetDefault("opensubs.enable", false).(bool),
		OpenSubsLanguage:      tomlFile.GetDefault("opensubs.language", defaultOpenSubsLanguage).(string),
		DownloadFolder:        tomlFile.GetDefault("general.download_folder", "").(string),
		TableBorderColor:      tomlFile.GetDefault("general.table_border_color", defaultTableBorderColor).(string),
		TableSelectionFgColor: tomlFile.GetDefault("general.table_selection_fg_color", defaultTableSelectionFgColor).(string),
		TableSelectionBgColor: tomlFile.GetDefault("general.table_selection_bg_color", defaultTableSelectionBgColor).(string),
		SpinnerColor:          tomlFile.GetDefault("general.spinner_color", defaultSpinnerColor).(string),
		DownloadColor:         tomlFile.GetDefault("general.download_color", defaultDownloadColor).(string),
	}
}
