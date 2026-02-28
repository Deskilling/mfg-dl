package main

import (
	"fmt"

	"mfg-dl/internal/anime4k"
	"mfg-dl/internal/core"
	"mfg-dl/internal/tui/simple"
	"mfg-dl/pkg/filesystem"

	"github.com/charmbracelet/log"
)

func init() {
	core.InitConfig()
	core.InitLogger(log.Level(core.GetConfig().Extra.LogLevel))
	filesystem.RemoveDirectory(core.GetConfig().Location.Temp)
	anime4k.DownloadLatestRelease()
}

func main() {
	if core.GetConfig().Tui.Basic {
		simple.SimpleTui()
	} else {
		log.Fatal("Not Implemented")
	}

	filesystem.RemoveDirectory(core.GetConfig().Location.Temp)

	log.Info("Execution finished press ENTER to quit")
	fmt.Scanln()
}
