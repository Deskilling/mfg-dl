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
	Download string
	Temp     string
	Shader   string `toml:"shader" comment:"This will only be used if you enable shaders"`
}

type Extra struct {
	FilePattern    string
	MaxConcurrency int
	LogLevel       int `toml:"loglevel" comment:"Debug: -4, Info: 0, Warn: 4, Error: 8, Fatal: 12"`
}

type Shader struct {
	Enable       bool
	Autoupdate   bool
	Shader       string
	CRF          int
	Preset       string
	AudioCopy    bool
	Width        string
	Height       string
	ExtraOptions string
}

type Config struct {
	Tui      Tui
	Location Location
	Extra    Extra
	Shader   Shader
}

var defaultConfig = Config{
	Tui: Tui{
		Basic: false,
	},

	Location: Location{
		Download: "./downloads",
		Temp:     "./temp",
		Shader:   "./shader",
	},

	Extra: Extra{
		FilePattern:    "{location}/{name}/{name}-Season{season}-Episode{episode}-{language}.mp4",
		MaxConcurrency: 16,
		LogLevel:       0,
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
