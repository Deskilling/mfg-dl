package util

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
)

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func RemoveAfterSymbol(input string, symbol string) string {
	index := strings.Index(input, symbol)
	if index != -1 {
		return input[:index]
	}
	return input
}

func ShiftChars(str string, shift int) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, char := range str {
		b.WriteRune(char - rune(shift))
	}
	return b.String()
}

func ParseAnimePath(path string) (name string, season int, episode int, err error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 5 {
		return "", 0, 0, fmt.Errorf("invalid path format")
	}

	name = parts[2]

	if _, err = fmt.Sscanf(parts[3], "staffel-%d", &season); err != nil {
		season = 0
		if _, err = fmt.Sscanf(parts[3], "filme"); err != nil {
			return "", 0, 0, fmt.Errorf("invalid season format")
		}
	}

	if _, err = fmt.Sscanf(parts[4], "episode-%d", &episode); err != nil {
		if _, err = fmt.Sscanf(parts[4], "film-%d", &episode); err != nil {
			return "", 0, 0, fmt.Errorf("invalid episode format")
		}
	}

	return name, season, episode, nil
}

func ParseMultipleInts(input string) []int {
	parts := strings.Split(input, ",")
	var result []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		v, err := strconv.Atoi(p)
		if err == nil && v > 0 {
			result = append(result, v)
		}
	}
	return result
}

func Hash256String(string string) string {
	h := sha256.Sum256([]byte(string))
	return fmt.Sprintf("%x", h[:])
}
