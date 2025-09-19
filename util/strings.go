package util

import (
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
