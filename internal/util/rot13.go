package util

import "strings"

func Rot13(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, char := range str {
		if 'a' <= char && char <= 'z' {
			b.WriteRune('a' + ((char - 'a' + 13) % 26))
		} else if 'A' <= char && char <= 'Z' {
			b.WriteRune('A' + ((char - 'A' + 13) % 26))
		} else {
			b.WriteRune(char)
		}
	}
	return b.String()
}
