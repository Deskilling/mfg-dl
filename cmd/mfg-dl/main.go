package main

import (
	"mfg-dl/internal/core"

	"github.com/charmbracelet/log"
)

func load() (err error) {
	core.InitLogger(log.InfoLevel)
	return nil
}

func main() {
	load()
}
