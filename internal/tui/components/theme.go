package components

import (
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

var Theme huh.ThemeFunc

func init() {
	Theme = func(isDark bool) *huh.Styles {
		t := huh.ThemeBase(isDark)
		lightDark := lipgloss.LightDark(isDark)

		var (
			normalFg = lightDark(lipgloss.Color("252"), lipgloss.Color("235"))
			indigo   = lightDark(lipgloss.Color("#5A56E0"), lipgloss.Color("#7571F9"))
			cream    = lightDark(lipgloss.Color("#FFFDF5"), lipgloss.Color("#FFFDF5"))
			fuchsia  = lipgloss.Color("#F780E2")
			green    = lightDark(lipgloss.Color("#02BA84"), lipgloss.Color("#02BF87"))
			red      = lightDark(lipgloss.Color("#FF4672"), lipgloss.Color("#ED567A"))
		)

		t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("238"))
		t.Focused.Card = t.Focused.Base
		t.Focused.Title = t.Focused.Title.Foreground(indigo).Bold(true)
		t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(indigo).Bold(true).MarginBottom(1)
		t.Focused.Directory = t.Focused.Directory.Foreground(indigo)
		t.Focused.Description = t.Focused.Description.Foreground(lightDark(lipgloss.Color(""), lipgloss.Color("243")))
		t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
		t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
		t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(fuchsia)
		t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(fuchsia)
		t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(fuchsia)
		t.Focused.Option = t.Focused.Option.Foreground(normalFg)
		t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(fuchsia)
		t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
		t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("#02CF92"), lipgloss.Color("#02A877"))).SetString("✓ ")
		t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color(""), lipgloss.Color("243"))).SetString("• ")
		t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(normalFg)
		t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(cream).Background(fuchsia)
		t.Focused.Next = t.Focused.FocusedButton
		t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(normalFg).Background(lightDark(lipgloss.Color("237"), lipgloss.Color("252")))

		t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(green)
		t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lightDark(lipgloss.Color("248"), lipgloss.Color("238")))
		t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(fuchsia)

		t.Blurred = t.Focused
		t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
		t.Blurred.Card = t.Blurred.Base
		t.Blurred.NextIndicator = lipgloss.NewStyle()
		t.Blurred.PrevIndicator = lipgloss.NewStyle()

		t.Group.Title = t.Focused.Title
		t.Group.Description = t.Focused.Description
		return t
	}
}
