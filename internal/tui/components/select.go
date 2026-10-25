package components

import (
	"errors"
	"slices"

	"github.com/pterm/pterm"
)

func SingleSelect[T any](options []T, enableFilter bool, label func(T) string) (T, error) {
	var null T

	labels := make([]string, len(options))
	for i, o := range options {
		labels[i] = label(o)
	}

	selectInput := pterm.DefaultInteractiveSelect.
		WithOptions(labels).
		WithFilter(enableFilter).
		WithMaxHeight(15)

	selected, err := selectInput.Show()
	if err != nil {
		return null, err
	}

	index := slices.Index(labels, selected)
	return options[index], nil
}

func MultiSelect[T any](options []T, enableFilter bool, label func(T) string) ([]T, error) {
	if len(options) == 0 {
		return nil, errors.New("no options to select from")
	}

	labels := make([]string, len(options))
	for i, o := range options {
		labels[i] = label(o)
	}

	selectedLabels, err := pterm.DefaultInteractiveMultiselect.
		WithOptions(labels).
		WithFilter(enableFilter).
		WithMaxHeight(15).
		Show()
	if err != nil {
		return nil, err
	}

	selected := make([]T, 0, len(selectedLabels))
	for _, l := range selectedLabels {
		idx := slices.Index(labels, l)
		selected = append(selected, options[idx])
	}

	return selected, nil
}
