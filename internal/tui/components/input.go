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

func ReadString(reader *bufio.Reader, prompt string) (input string, err error) {
	fmt.Print(prompt)
	input, err = reader.ReadString('\n')
	if err != nil {
		log.Error("Failed reading string")
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func ReadInt(reader *bufio.Reader, prompt string) (val int, err error) {
	for {
		input, err := ReadString(reader, prompt)
		if err != nil {
			return 0, err
		}
		val, err := strconv.Atoi(input)
		if err != nil {
			log.Error("input an integer")
			continue
		}

		return val, nil
	}
}
