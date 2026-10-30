package main

import (
	"fmt"
	"log/slog"
	"os"

	"mfg-dl/internal/core"
	"mfg-dl/internal/ffmpeg"
	"mfg-dl/internal/sites"
	"mfg-dl/internal/tui"
	"mfg-dl/internal/util"

	"github.com/Deskilling/gopkg/pkg/filesystem"
)

func init() {
	core.InitLogger(0)

	var err error
	err = filesystem.ChangeExecPath()
	if err != nil {
		slog.Error("Failed to execute Programm in the same path", "err", err)
		os.Exit(1)
	}

	err = core.InitConfig()
	if err != nil {
		slog.Error("Failed to init config", "err", err)
		os.Exit(1)
	} else {
		slog.Info("Loaded config")
	}

	core.InitLogger(core.GetConfig().Debug.LogLevel)

	err = filesystem.RemoveDirectory(core.GetConfig().Location.Temp + "/segments")
	if err != nil {
		slog.Error("Failed deleting directory", "dir", core.GetConfig().Location.Temp+"/segments", "err", err)
	} else {
		slog.Info("Cleaned temp segments")
	}

	err = ffmpeg.CheckInstalled()
	if err != nil {
		fmt.Printf("\n")
		slog.Error("FFmpeg is required to use this application", "err", err)
		os.Exit(1)
	}

	core.CreateCleanupTask()
	sites.Init()
}

func main() {
	if core.GetConfig().Debug.LogLevel != -4 {
		util.ClearTerminal()
	}

	slog.Error("finished", "err", tui.Tui())
	return

	/*
		var err error = nil
		switch core.GetConfig().Tui.Mode {
		case 0:
			fmt.Println("gg")
		case 1:
			fmt.Println("gg")
		case 2:
			err = server.Start()
		default:
			slog.Error("Invaild TUI Mode")
			os.Exit(1)
		}

		if err != nil {
			slog.Error("TUI error", "err", err)
		}

		fmt.Scanln()
		filesystem.RemoveDirectory(core.GetConfig().Location.Temp + "/segments")
	*/
}
