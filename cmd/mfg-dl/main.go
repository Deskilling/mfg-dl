package main

import (
	"mfg-dl/internal/core"
	"mfg-dl/internal/tui"
	"mfg-dl/pkg/filesystem"

	"github.com/charmbracelet/log"
)

func init() {
	core.InitConfig()
	core.InitLogger(log.Level(core.GetConfig().Extra.LogLevel))
	filesystem.RemoveDirectory(core.GetConfig().Location.Temp)
}

func main() {
	if core.GetConfig().Tui.Basic {
		tui.SimpleTui()
	} else {
		log.Fatal("Not Implemented")
	}
}
