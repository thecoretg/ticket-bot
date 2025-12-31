package tui

import tea "github.com/charmbracelet/bubbletea"

type (
	dimensions struct {
		windowW int
		windowH int
	}

	dimensionsCalculatedMsg struct{ dimensions }
)

func (m *Model) calculateDimensions(w, h int) tea.Cmd {
	return func() tea.Msg {
		d := &dimensions{}
		d.windowW = w
		d.windowH = h

		return dimensionsCalculatedMsg{dimensions: *d}
	}
}
