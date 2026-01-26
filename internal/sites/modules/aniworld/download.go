package aniworld

import (
	"fmt"
	"mfg-dl/internal/sites/model"
)

func Download(stream model.Stream) (err error) {
	url := BaseURL + stream.Href
	fmt.Println(url)

	return nil
}
