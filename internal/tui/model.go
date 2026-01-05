package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/pkg/sdk"
)

type Model struct {
	SDKClient   *sdk.Client
	activeModel tea.Model
	allModels   allModels
	data        *data
	keys        keyMap
	width       int
	height      int
}

type allModels struct {
	rules *rulesModel
	fwds  *fwdsModel
}

type data struct {
	rules []models.NotifierRuleFull
}

func NewModel(sl *sdk.Client) *Model {
	rm := newRulesModel()
	fm := newFwdsModel()
	return &Model{
		SDKClient:   sl,
		keys:        defaultKeyMap,
		activeModel: rm,
		allModels: allModels{
			rules: rm,
			fwds:  fm,
		},
		data: &data{
			rules: []models.NotifierRuleFull{},
		},
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.getRules(), m.getFwds())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		}
	case switchModelMsg:
		switch msg.modelType {
		case modelTypeRules:
			if m.activeModel != m.allModels.rules {
				m.allModels.rules.table.SetCursor(0)
				m.activeModel = m.allModels.rules
			}
		case modelTypeFwds:
			if m.activeModel != m.allModels.fwds {
				m.allModels.fwds.table.SetCursor(0)
				m.activeModel = m.allModels.fwds
			}
		}
	}

	var cmds []tea.Cmd
	rules, cmd := m.allModels.rules.Update(msg)
	if r, ok := rules.(*rulesModel); ok {
		m.allModels.rules = r
	}
	cmds = append(cmds, cmd)

	fwds, cmd := m.allModels.fwds.Update(msg)
	if f, ok := fwds.(*fwdsModel); ok {
		m.allModels.fwds = f
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.activeModel.View()
}
