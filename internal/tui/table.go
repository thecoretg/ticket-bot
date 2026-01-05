package tui

import "github.com/charmbracelet/bubbles/table"

func newTable() table.Model {
	t := table.New(
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Selected = s.Selected.
		Foreground(grey).
		Background(blue)

	t.SetStyles(s)
	return t
}
