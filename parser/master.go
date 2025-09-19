package parser

import (
	"fmt"
	"regexp"

	"github.com/charmbracelet/log"
)

func GetMasterPlaylistName(file string) ([]string, error) {
	regex := regexp.MustCompile(`(?:)?([^/\r\n]+?\.m3u8)`)
	matches := regex.FindStringSubmatch(file)

	if len(matches) <= 0 {
		err := fmt.Errorf("no .m3u8 file found")
		log.Error(err)
		return nil, err
	}

	return matches, nil
}
