package tui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"mfg-dl/globals"

	"github.com/charmbracelet/log"
)

func GetUserInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func Start() {
	sites := globals.Sites()

	for i := range sites {
		fmt.Printf("[%v] %s\n", i+1, sites[i])
	}

	input := GetUserInput("Enter: ")
	userSite, err := strconv.Atoi(input)
	if err != nil {
		log.Error(err)
		return
	}

	if userSite < 1 || userSite > len(sites) {
		log.Error("invalid site selection")
		return
	}

	if sites[userSite-1] == "aniworld" {
		Aniworld()
	}
}
