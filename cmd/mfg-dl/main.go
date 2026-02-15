package main

import (
	"mfg-dl/internal/core"
	"mfg-dl/internal/tui"

	"github.com/charmbracelet/log"
)

func load() (err error) {
	core.InitLogger(log.DebugLevel)
	core.InitConfig()

	return nil
}

func main() {
	load()

	if core.GetConfig().Tui.Basic {
		tui.SimpleTui()
	} else {
		log.Fatal("Not Implemented")
	}

}
