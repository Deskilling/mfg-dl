package util

import (
	"mfg-dl/pkg/filesystem"

	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
)

var cfg Config

type general struct {
	MaxRoutines int `comment:"maximum of gorutines at once"`
}

type location struct {
	Output  string
	Tempdir string
}

type Config struct {
	General  general
	Location location
}

var dCfg = Config{
	General: general{
		MaxRoutines: 64,
	},
	Location: location{
		Output:  "./output/",
		Tempdir: "/temp/",
	},
}

func DefaultConfig() {
	w, err := toml.Marshal(&dCfg)
	if err != nil {
		return
	}

	log.Debug("", "w", string(w))

	if filesystem.WriteFile("config.toml", w) != nil {
		log.Error("failed writing default config")
		return
	}

	log.Info("Created Default config")
}

func checkConfig(cfg *Config) {
	if cfg.General.MaxRoutines <= 0 {
		log.Warn("invalid setting for General.MaxRoutines", "old", cfg.General.MaxRoutines, "new", dCfg.General.MaxRoutines)
		cfg.General.MaxRoutines = dCfg.General.MaxRoutines
	}
}

func ReadConfig() *Config {
	if filesystem.ExistPath("config.toml") {
		c, err := filesystem.ReadFile("config.toml")
		if err != nil {
			log.Error("failed reading file")
			return nil
		}

		if len(c) > 0 {
			err = toml.Unmarshal([]byte(c), &cfg)
			if err != nil {
				return nil
			}

			checkConfig(&cfg)
			return &cfg
		}
	}

	DefaultConfig()
	cfg = dCfg
	return &cfg
}

func GetSettings() *Config {
	return &cfg
}
