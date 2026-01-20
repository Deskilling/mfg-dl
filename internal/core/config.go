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
}

type Config struct {
	Tui      Tui
	Location Location
}

var defaultConfig = Config{
	Tui: Tui{
		Basic: false,
	},

	Location: Location{
		Download: "./downloads",
		Temp:     "./temp",
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
