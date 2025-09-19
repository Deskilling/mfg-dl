package main

import (
	"fmt"
	"mfg-dl/filesystem"
	"mfg-dl/util"
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
	fmt.Println("Currently not Implemented !!!")
}
