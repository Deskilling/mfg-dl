package core

import (
	"mfg-dl/pkg/config"

	"github.com/charmbracelet/log"
)

const configLocation string = "./config.toml"

type Tui struct {
	Basic bool
}

type Location struct {
	FilePattern string `toml:"filepattern" comment:"customizes the output filename for video files\n\nAvailable placeholders:\n{location} download directory\n{name} show name\n{season} season number\n{title} episode title\n{title2} alternative title (if available)\n{episode} episode number\n{language} language\n{hoster} stream hoster"`
	Download    string `toml:"download" comment:"base download directory"`
	Temp        string `toml:"temp" comment:"directory for temporary files"`
	Shader      string `toml:"shader" comment:"directory for shaders, only used if shaders are enabled"`
}

type Extra struct {
	MaxVideoConcurrency int  `toml:"maxvideoconcurrency" comment:"maximum number of concurrent video downloads"`
	FfmpegDownload      bool `toml:"ffmpegdownload" comment:"use ffmpeg for HLS streams, usually slower but more stable, only enable if you run into issues"`
	LogLevel            int  `toml:"loglevel" comment:"log level: Debug (-4), Info (0), Warn (4), Error (8), Fatal (12)"`
}

type Downloads struct {
	MaxSegmentConcurrency int `toml:"maxsegmentconcurrency" comment:"maximum number of concurrent segment downloads"`
	MaxRetires            int `toml:"maxsegmentretires" comment:"max retries for each segment"`
	RetryDelay            int `toml:"retrydelay" comments:"retry delay of the segments in sec"`
}

type Shader struct {
	Enable       bool   `toml:"enable" comment:"enable shader processing after download"`
	Autoupdate   bool   `toml:"autoupdate" comment:"automatically update the shader before use"`
	Shader       string `toml:"shader" comment:"name or path of the shader to use"`
	CRF          int    `toml:"crf" comment:"ffmpeg CRF value (lower = better quality, bigger file size)"`
	Preset       string `toml:"preset" comment:"ffmpeg preset (ultrafast, fast, medium, slow, etc.)"`
	AudioCopy    bool   `toml:"audiocopy" comment:"copy audio stream without re-encoding"`
	Width        string `toml:"width" comment:"output width (keep empty to preserve original)"`
	Height       string `toml:"height" comment:"output height (keep empty to preserve original)"`
	ExtraOptions string `toml:"extraoptions" comment:"additional raw ffmpeg options (advanced use only)"`
}

type Config struct {
	Tui       Tui
	Location  Location
	Downloads Downloads
	Extra     Extra
	Shader    Shader
}

var defaultConfig = Config{
	Tui: Tui{
		Basic: false,
	},

	Location: Location{
		FilePattern: "{location}/{name}/{name}-Season{season}-Episode{episode}-{language}.mp4",
		Download:    "./downloads",
		Temp:        "./temp",
		Shader:      "./shader",
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

	Shader: Shader{
		Enable:     false,
		Autoupdate: false,
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
