package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/pkg/sdk"
)

type Model struct {
	SDKClient *sdk.Client
	data      *data
	keys      keyMap

	dimensions
	initFlags
}

type data struct {
	boards []*models.Board
	recips []*models.WebexRecipient
	rules  []*models.NotifierRuleFull
	fwds   []*models.NotifierForwardFull
}

type initFlags struct {
	gotDimensions bool
	gotData       bool
}

func NewModel(sl *sdk.Client) *Model {
	return &Model{
		SDKClient: sl,
		keys:      defaultKeyMap,
		data: &data{
			boards: []*models.Board{},
			recips: []*models.WebexRecipient{},
			rules:  []*models.NotifierRuleFull{},
			fwds:   []*models.NotifierForwardFull{},
		},
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, m.calculateDimensions(msg.Width, msg.Height)
	case dimensionsCalculatedMsg:
		m.dimensions = msg.dimensions
		m.gotDimensions = true
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *Model) View() string {
	if !m.gotDimensions {
		return "Initializing..."
	}

	return m.mainFlexbox().Render(m.windowW, m.windowH)
}
