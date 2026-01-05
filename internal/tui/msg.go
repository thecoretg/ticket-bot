package tui

import tea "github.com/charmbracelet/bubbletea"

type (
	switchModelMsg struct{ modelType }
	modelType      int
)

const (
	modelTypeRules modelType = iota
	modelTypeFwds
)

func switchModel(m modelType) tea.Cmd {
	return func() tea.Msg {
		return switchModelMsg{m}
	}
}
