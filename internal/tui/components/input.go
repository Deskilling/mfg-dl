package components

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"charm.land/log/v2"
)

var Reader *bufio.Reader = bufio.NewReader(os.Stdin)

func ReadString(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Error("Failed reading string")
		return ""
	}
	return strings.TrimSpace(input)
}

func ReadInt(reader *bufio.Reader, prompt string) int {
	for {
		input := ReadString(reader, prompt)
		val, err := strconv.Atoi(input)
		if err != nil {
			log.Error("input an integer")
			continue
		}

		return val
	}
}
