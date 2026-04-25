package core

import (
	"fmt"

	"github.com/Deskilling/gopkg/pkg/config"

	"github.com/Deskilling/gopkg/pkg/filesystem"

	"charm.land/log/v2"
)

const configLocation string = "./config.toml"
const configVersion string = "0.1"

type Tui struct {
	Mode int `toml:"tui" comment:"0 tmdb (recommendend)\n1 services directly"`
}

type Location struct {
	FilePattern string `toml:"filepattern" comment:"customizes the output filename for video files\n\nAvailable placeholders:\n{location} download directory\n{name} show name\n{season} season number\n{title} episode title\n{title2} alternative title (if available)\n{episode} episode number\n{language} language\n{hoster} stream hoster"`
	Download    string `toml:"download" comment:"base download directory"`
	Temp        string `toml:"temp" comment:"directory for temporary files"`
}

type Downloads struct {
	FfmpegDownload      bool `toml:"ffmpegdownload" comment:"use ffmpeg for HLS streams, usually slower but more stable, only enable if you run into issues"`
	MaxVideoConcurrency int  `toml:"maxvideoconcurrency" comment:"maximum number of concurrent video downloads"`

	MaxSegmentConcurrency int `toml:"maxsegmentconcurrency" comment:"maximum number of concurrent segment downloads"`
	MaxRetires            int `toml:"maxsegmentretires" comment:"max retries for each segment"`
	RetryDelay            int `toml:"retrydelay" comment:"retry delay of the segments in sec"`
}

type Cache struct {
	EnableCache  bool `toml:"cacheEnabled" comment:"Caches all requested webpages for x time"`
	Minutes      int  `toml:"cacheMinutes" comment:"How long to keep each cached page"`
	CleanMinutes int  `toml:"cacheClean" comment:"How often in Minutes to check the cache"`
}

type Debug struct {
	LogLevel    int  `toml:"loglevel" comment:"log level: Debug (-4), Info (0), Warn (4), Error (8), Fatal (12)"`
	Sha256Cache bool `toml:"sha256cache" comment:"Hash cache filenames with sha256"`
}

type Config struct {
	Tui       Tui
	Services  map[string]bool `toml:"services" comment:"used to disable/enable modules"`
	Location  Location
	Downloads Downloads
	Cache     Cache
	Debug     Debug

	Version string
}

var defaultConfig = Config{
	Tui: Tui{
		Mode: 0,
	},

	Services: map[string]bool{
		"Aniworld": true,
	},

	Location: Location{
		FilePattern: "{location}/{name}/Season{season}/Episode-{episode}-{language}-{title}.mp4",
		Download:    "./downloads",
		Temp:        "./temp",
	},

	Downloads: Downloads{
		FfmpegDownload:        false,
		MaxVideoConcurrency:   4,
		MaxSegmentConcurrency: 16,
		MaxRetires:            3,
		RetryDelay:            3,
	},

	Cache: Cache{
		EnableCache:  true,
		Minutes:      60,
		CleanMinutes: 30,
	},

	Debug: Debug{
		LogLevel:    0,
		Sha256Cache: true,
	},

	Version: configVersion,
}

var cfg Config

func InitConfig() error {
	err := config.Load(configLocation, &cfg)
	if err != nil {
		cfg = defaultConfig

		log.Info("Creating new config", "path", configLocation)
		err = config.Save(configLocation, cfg)
		if err != nil {
			return fmt.Errorf("Failed to save config: %w", err)
		}
	}

	if GetConfig().Version != configVersion {
		content, err := filesystem.ReadFile(configLocation)
		if err != nil {
			return err
		}

		err = filesystem.WriteFile(configLocation+".old", content)
		if err != nil {
			return err
		}

		err = filesystem.DeleteFile(configLocation)
		if err != nil {
			return err
		}

		log.Warnf("Renamed old config: %s -> %s", configLocation, configLocation+".old")
		return InitConfig()
	}

	return nil
}

func GetConfig() *Config {
	return &cfg
}
