package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func Load(path string, cfg any) error {
	if path == "" || cfg == nil {
		return fmt.Errorf("invalid arguments")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if len(bytes.TrimSpace(data)) == 0 {
		return fmt.Errorf("data is empty, using default")
	}

	return toml.Unmarshal(data, cfg)
}
