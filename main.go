package main

import (
	"mfg-dl/filesystem"
	"mfg-dl/tui"
	"mfg-dl/util"
	"os"
)

func init() {
	/*
		DebugLevel = -4
		InfoLevel = 0
		WarnLevel = 4
		ErrorLevel = 8
		FatalLevel = 12
	*/

	// TODO Implement cli arguments parsing
	util.InitLogger(-4)
	filesystem.InitExecDir()
}

func main() {
	defer os.RemoveAll("./temp")
	tui.Start()
}
