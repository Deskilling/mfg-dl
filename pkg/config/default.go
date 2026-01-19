package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

func Default(path string, cfg any) (err error) {
	if cfg == nil || path == "" {
		return fmt.Errorf("invalid path or cfg")
	}

	err = os.MkdirAll(filepath.Dir(path), 0666)
	if err != nil {
		return nil
	}

	data, err := toml.Marshal(&cfg)
	if err != nil {
		return err
	}
	os.WriteFile(path, data, 0666)

	return nil
}
