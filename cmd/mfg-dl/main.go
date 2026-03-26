package main

import (
	"fmt"
	"os"

	"mfg-dl/internal/core"
	"mfg-dl/internal/ffmpeg"
	"mfg-dl/internal/tui/newsimple"
	"mfg-dl/internal/tui/simple"
	"mfg-dl/pkg/filesystem"

	"charm.land/log/v2"
)

func init() {
	var err error
	err = filesystem.ChangeExecDir()
	if err != nil {
		fmt.Printf("Failed to execute Programm in the same path")
		os.Exit(-1)
	}

	err = core.InitConfig()
	if err != nil {
		fmt.Printf("Failed to init config")
		os.Exit(-1)
	}

	core.InitLogger(log.Level(core.GetConfig().Extra.LogLevel))

	filesystem.RemoveDirectory(core.GetConfig().Location.Temp)

	err = ffmpeg.CheckInstalled()
	if err != nil {
		log.Error("ffmpeg is required to use this programm\nA tutorial for installation can befound here: https://gist.github.com/Deskilling/4994b58618a9ef896bf02481215e3291")
		os.Exit(0)
	}

	core.CreateCleanupTask()
}

func main() {
	if core.GetConfig().Tui.Basic {
		simple.SimpleTui()
	} else {
		newsimple.Tui()
	}

	log.Info("Execution finished press ENTER to quit")
	fmt.Scanln()

	filesystem.RemoveDirectory(core.GetConfig().Location.Temp)
}
