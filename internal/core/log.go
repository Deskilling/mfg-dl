package core

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func InitLogger(level log.Level) {
	log.SetLevel(level)
	log.SetPrefix("mfg-dl")

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
