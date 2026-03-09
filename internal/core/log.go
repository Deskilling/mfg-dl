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

	style := log.DefaultStyles()

	style.Levels[log.DebugLevel] = lipgloss.NewStyle().SetString("DEBUG").Foreground(lipgloss.Color("63"))
	if level == log.DebugLevel {
		log.SetReportCaller(true)
	} else {
		log.SetReportCaller(false)
		log.SetReportTimestamp(false)
	}

	style.Levels[log.InfoLevel] = lipgloss.NewStyle().Padding(0, 1, 0, 0).SetString("INFO").Foreground(lipgloss.Color("86"))

	style.Levels[log.WarnLevel] = lipgloss.NewStyle().Padding(0, 1, 0, 0).SetString("WARN").Foreground(lipgloss.Color("192"))

	style.Levels[log.ErrorLevel] = lipgloss.NewStyle().SetString("ERROR").Foreground(lipgloss.Color("#FF0000"))
	style.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)

	style.Levels[log.FatalLevel] = lipgloss.NewStyle().SetString("FATAL").Foreground(lipgloss.Color("134"))

	log.SetStyles(style)
}
