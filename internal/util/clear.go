package util

import "fmt"

func ClearTerminal() {
	fmt.Print("\033[H\033[2J")
}
