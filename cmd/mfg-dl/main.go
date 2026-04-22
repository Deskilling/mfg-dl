package main

import (
	"fmt"

	"mfg-dl/internal/core"
	"mfg-dl/internal/ffmpeg"
	"mfg-dl/internal/server"
	"mfg-dl/internal/sites"
	"mfg-dl/internal/tui/service"
	"mfg-dl/internal/tui/tmdb"
	"mfg-dl/internal/util"
	"mfg-dl/pkg/filesystem"

	"charm.land/log/v2"
)

func init() {
	core.InitLogger(0)

	var err error
	err = filesystem.ChangeExecPath()
	if err != nil {
		log.Fatal("Failed to execute Programm in the same path", "err", err)
	}

	err = core.InitConfig()
	if err != nil {
		log.Fatal("Failed to init config", "err", err)
	} else {
		log.Info("Loaded config")
	}

	log.SetLevel(log.Level(core.GetConfig().Debug.LogLevel))

	err = filesystem.RemoveDirectory(core.GetConfig().Location.Temp + "/segments")
	if err != nil {
		log.Error("Failed deleting directory", "dir", core.GetConfig().Location.Temp+"/segments", "err", err)
	} else {
		log.Info("Cleaned temp segments")
	}

	err = ffmpeg.CheckInstalled()
	if err != nil {
		fmt.Printf("\n")
		log.Fatal("FFmpeg is required to use this application", "err", err)
	}

	core.CreateCleanupTask()
	sites.Init()
}

func main() {
	if core.GetConfig().Debug.LogLevel != -4 {
		util.ClearTerminal()
	}

	var err error = nil
	switch core.GetConfig().Tui.Mode {
	case 0:
		err = tmdb.Tui()
	case 1:
		err = service.Tui()
	case 2:
		err = server.Start()
	default:
		log.Fatal("Invaild TUI Mode")
	}

	if err != nil {
		log.Error(err)
	}

	fmt.Scanln()
	filesystem.RemoveDirectory(core.GetConfig().Location.Temp + "/segments")
}
