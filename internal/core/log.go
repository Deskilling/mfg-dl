package core

import (
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"github.com/charmbracelet/colorprofile"
)

func InitLogger(level log.Level) {
	log.SetLevel(level)
	log.SetPrefix("mfg-dl")
	log.SetColorProfile(colorprofile.TrueColor)
	log.SetFormatter(log.TextFormatter)

	style := log.DefaultStyles()

	style.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DEBU").
		Padding(0, 1, 0, 1).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("62")).
		Italic(true).
		Bold(true)

	style.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("42")).
		Italic(true).
		Bold(true)

	style.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WARN").
		Padding(0, 1, 0, 1).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("192")).
		Italic(true).
		Bold(true)

	style.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERRO").
		Padding(0, 1, 0, 1).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("160")).
		Italic(true).
		Bold(true)

	style.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	style.Values["err"] = lipgloss.NewStyle().Bold(true)

	style.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("FATA").
		Padding(0, 1, 0, 1).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("92")).
		Italic(true).
		Bold(true)

	if level == log.DebugLevel {
		log.SetReportCaller(true)
	} else {
		log.SetReportCaller(false)
		log.SetReportTimestamp(false)
	}

	log.SetStyles(style)
}
