package components

import (
	"fmt"
	"strings"
)

// this wont work for like multiple routines, but that is an issue for later
func PrintProgress(current, total int) {
	width := 40
	filled := (current * width) / total

	var builder strings.Builder
	// 3 bytes
	builder.Grow(width * 3)
	for i := range width {
		if i < filled {
			builder.WriteRune('█')
		} else {
			builder.WriteRune('░')
		}
	}

	percent := (current * 100) / total

	fmt.Printf("\r[%s] %d%% (%d/%d)", builder.String(), percent, current, total)
}
