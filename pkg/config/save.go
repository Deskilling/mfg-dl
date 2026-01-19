package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func Save(path string, cfg any) (err error) {
	if path == "" || cfg == nil {
		return fmt.Errorf("invalid path")
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
