package main

import (
	"mfg-dl/internal/core"

	"github.com/charmbracelet/log"
)

func load() (err error) {
	core.InitLogger(log.DebugLevel)

	return nil
}

func main() {
	load()
}
