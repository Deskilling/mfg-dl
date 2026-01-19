package core

import (
	"mfg-dl/pkg/config"
)

type Config struct {
}

var defaultConfig Config = Config{}

func InitConfig() (cfg any, err error) {
	cfg, err = config.Load("./config.toml")
	if err != nil {
		err = config.Default("./config.toml", defaultConfig)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	return cfg, nil
}
