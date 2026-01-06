package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	white    = lipgloss.ANSIColor(7)
	grey     = lipgloss.ANSIColor(8)
	blue     = lipgloss.ANSIColor(4)
	red      = lipgloss.ANSIColor(1)
	errStyle = lipgloss.NewStyle().Foreground(red).Bold(true)
)

func fillSpaceCentered(content string, w, h int) string {
	return fillSpace(content, w, h, lipgloss.Center, lipgloss.Center)
}

func fillSpace(content string, w, h int, alignH, alignV lipgloss.Position) string {
	return lipgloss.NewStyle().Width(w).Height(h).Align(alignH, alignV).Render(content)
}
