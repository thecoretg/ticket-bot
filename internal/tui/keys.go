package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	quit key.Binding
}

var defaultKeyMap = keyMap{
	quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.quit},
	}
}

type rulesModelKeys struct {
	switchFwds key.Binding
}

var defaultRulesKeys = rulesModelKeys{
	switchFwds: key.NewBinding(
		key.WithKeys("ctrl+f"),
		key.WithHelp("ctrl+f", "switch to forwards"),
	),
}

func (rk rulesModelKeys) ShortHelp() []key.Binding {
	return []key.Binding{rk.switchFwds}
}

func (rk rulesModelKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{rk.switchFwds},
	}
}

type fwdsModelKeys struct {
	switchRules key.Binding
}

var defaultFwdsKeys = fwdsModelKeys{
	switchRules: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "switch to rules"),
	),
}

func (fk fwdsModelKeys) ShortHelp() []key.Binding {
	return []key.Binding{fk.switchRules}
}

func (fk fwdsModelKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{fk.switchRules},
	}
}
