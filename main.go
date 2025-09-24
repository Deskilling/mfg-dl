package main

import (
	"fmt"
	"mfg-dl/filesystem"
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
	fmt.Println("WILL BE IMPLEMENTED SOON !!! :3")
}
