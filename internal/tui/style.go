package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	white     = lipgloss.ANSIColor(7)
	grey      = lipgloss.ANSIColor(8)
	blue      = lipgloss.ANSIColor(4)
	helpStyle = lipgloss.NewStyle().Foreground(white)
)
