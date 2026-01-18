package util

import (
	"encoding/base64"
	"fmt"
)

func Base64Decode(str string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", fmt.Errorf("failed decoding base64: %w", err)
	}
	return string(decoded), nil
}
