package components

import (
	"errors"
	"strconv"

	"charm.land/huh/v2"
)

func ReadString(prompt string) (string, error) {
	var input string
	err := huh.NewInput().
		Title(prompt).
		Value(&input).
		WithTheme(Theme).
		Run()
	if err != nil {
		return "", err
	}
	return input, nil
}

func ReadInt(prompt string) (int, error) {
	var input string
	err := huh.NewInput().
		Title(prompt).
		Validate(func(s string) error {
			_, err := strconv.Atoi(s)
			if err != nil {
				return errors.New("input an integer")
			}
			return nil
		}).
		Value(&input).
		WithTheme(Theme).
		Run()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(input)
}
