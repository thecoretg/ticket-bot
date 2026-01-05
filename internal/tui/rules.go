package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thecoretg/ticketbot/internal/models"
)

type rulesModel struct {
	keys          rulesModelKeys
	gotDimensions bool
	width         int
	height        int
	rulesLoaded   bool
	table         table.Model
	help          help.Model

	rules []models.NotifierRuleFull
}

func newRulesModel() *rulesModel {
	h := help.New()
	h.Styles.ShortDesc = helpStyle
	h.Styles.ShortKey = helpStyle
	return &rulesModel{
		keys:  defaultRulesKeys,
		rules: []models.NotifierRuleFull{},
		table: newTable(),
		help:  h,
	}
}

func (rm *rulesModel) Init() tea.Cmd {
	return nil
}

func (rm *rulesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		rm.width = msg.Width
		rm.height = msg.Height
		rm.gotDimensions = true
		hh := lipgloss.Height(rm.helpView())
		setRulesTableDimensions(&rm.table, rm.width, rm.height-hh)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, rm.keys.switchFwds):
			return rm, switchModel(modelTypeFwds)
		}
	case gotRulesMsg:
		rm.rules = msg.rules
		rm.rulesLoaded = true
		rm.table.SetRows(rulesToRows(rm.rules))
	}

	var cmd tea.Cmd
	rm.table, cmd = rm.table.Update(msg)

	return rm, cmd
}

func (rm *rulesModel) View() string {
	if !rm.gotDimensions {
		return "Initializing..."
	}

	if !rm.rulesLoaded {
		return "Loading rules..."
	}

	return lipgloss.JoinVertical(lipgloss.Top, rm.table.View(), rm.helpView())
}

func (rm *rulesModel) helpView() string {
	return rm.help.ShortHelpView(rm.keys.ShortHelp())
}

func setRulesTableDimensions(t *table.Model, w, h int) {
	enableW := 8
	boardW := 20
	remainingW := w - enableW - boardW
	recipW := remainingW
	t.SetColumns([]table.Column{
		{Title: "ENABLED", Width: enableW},
		{Title: "BOARD", Width: boardW},
		{Title: "RECIPIENT", Width: recipW},
	})
	t.SetHeight(h)
}

func rulesToRows(rules []models.NotifierRuleFull) []table.Row {
	var rows []table.Row
	for _, r := range rules {
		rows = append(rows, []string{boolToIcon(r.Enabled), r.BoardName, r.RecipientName})
	}

	return rows
}
