package core

import (
	"mfg-dl/pkg/config"

	"charm.land/log/v2"
)

const configLocation string = "./config.toml"

type Tui struct {
	Tmdb bool `toml:"tmdb" comment:"uses tmdb as a search, recommendend in most cases"`
}

type Location struct {
	FilePattern string `toml:"filepattern" comment:"customizes the output filename for video files\n\nAvailable placeholders:\n{location} download directory\n{name} show name\n{season} season number\n{title} episode title\n{title2} alternative title (if available)\n{episode} episode number\n{language} language\n{hoster} stream hoster"`
	Download    string `toml:"download" comment:"base download directory"`
	Temp        string `toml:"temp" comment:"directory for temporary files"`
	Cache       string `toml:"cache" comment:"cache for like search results"`
}

type Extra struct {
	MaxVideoConcurrency int  `toml:"maxvideoconcurrency" comment:"maximum number of concurrent video downloads"`
	FfmpegDownload      bool `toml:"ffmpegdownload" comment:"use ffmpeg for HLS streams, usually slower but more stable, only enable if you run into issues"`
	LogLevel            int  `toml:"loglevel" comment:"log level: Debug (-4), Info (0), Warn (4), Error (8), Fatal (12)"`
}

type Downloads struct {
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
	Sha256Cache bool `toml:"sha256cache" comment:"Hash cache filenames with sha256"`
}

type Config struct {
	Tui       Tui
	Location  Location
	Downloads Downloads
	Extra     Extra
	Cache     Cache
	Debug     Debug
}

var defaultConfig = Config{
	Tui: Tui{
		Tmdb: true,
	},

	Location: Location{
		FilePattern: "{location}/{name}/Season{season}/Episode-{episode}-{language}.mp4",
		Download:    "./downloads",
		Temp:        "./temp",
		Cache:       "./cache",
	},

	Downloads: Downloads{
		MaxSegmentConcurrency: 16,
		MaxRetires:            3,
		RetryDelay:            3,
	},

	Extra: Extra{
		MaxVideoConcurrency: 4,
		FfmpegDownload:      false,
		LogLevel:            0,
	},

	Cache: Cache{
		EnableCache:  true,
		Minutes:      60,
		CleanMinutes: 30,
	},

	Debug: Debug{
		Sha256Cache: true,
	},
}

var cfg Config

func InitConfig() error {
	err := config.Load(configLocation, &cfg)
	if err != nil {
		log.Error("Failed loading config", "err", err)
		cfg = defaultConfig

		err = config.Save(configLocation, cfg)
		if err != nil {
			log.Error("failed saving config", "err", err)
			return err
		}

		log.Info("created config at", "configLocation", configLocation)
	}
	log.Debug("loaded config")
	return nil
}

func GetConfig() *Config {
	return &cfg
}
