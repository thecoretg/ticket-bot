package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/pkg/sdk"
)

type (
	rulesModel struct {
		SDKClient   *sdk.Client
		rulesLoaded bool
		table       table.Model
		form        *huh.Form
		spinner     spinner.Model
		formResult  *rulesFormResult
		status      subModelStatus

		availWidth  int
		availHeight int
		rules       []models.NotifierRuleFull
	}

	ruleFormDataMsg struct {
		boards []models.Board
		recips []models.WebexRecipient
	}

	rulesFormResult struct {
		board models.Board
		recip models.WebexRecipient
	}

	refreshRulesMsg struct{}
	gotRulesMsg     struct{ rules []models.NotifierRuleFull }
)

func newRulesModel(cl *sdk.Client) *rulesModel {
	s := spinner.New()
	s.Spinner = spinner.Line
	return &rulesModel{
		SDKClient:  cl,
		rules:      []models.NotifierRuleFull{},
		table:      newTable(),
		formResult: &rulesFormResult{},
		spinner:    s,
	}
}

func (rm *rulesModel) Init() tea.Cmd {
	return tea.Batch(rm.spinner.Tick, rm.getRules())
}

func (rm *rulesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, allKeys.newItem) && rm.status == statusMain:
			rm.status = statusLoadingFormData
			return rm, tea.Batch(rm.prepareForm())
		case key.Matches(msg, allKeys.deleteItem) && rm.status == statusMain:
			if len(rm.rules) > 0 {
				rule := rm.rules[rm.table.Cursor()]
				rm.status = statusRefresh
				return rm, tea.Batch(rm.deleteRule(rule.ID))
			}
		}

	case resizeModelsMsg:
		rm.availWidth = msg.w
		rm.availHeight = msg.h
		rm.setModuleDimensions()

	case refreshRulesMsg:
		return rm, rm.getRules()

	case gotRulesMsg:
		rm.rules = msg.rules
		rm.rulesLoaded = true
		rm.status = statusMain
		return rm, tea.Batch(rm.setRows())

	case ruleFormDataMsg:
		rm.formResult = &rulesFormResult{}
		rm.form = ruleEntryForm(msg.boards, msg.recips, rm.formResult, rm.availHeight)
		rm.status = statusEntry
		return rm, rm.form.Init()
	case errMsg:
		rm.status = statusMain
	}

	var cmds []tea.Cmd
	switch rm.status {

	case statusEntry:
		form, cmd := rm.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			rm.form = f
		}

		cmds = append(cmds, cmd)
		switch rm.form.State {
		case huh.StateAborted:
			rm.status = statusMain

		case huh.StateCompleted:
			res := rm.formResult
			rule := &models.NotifierRule{
				CwBoardID:        res.board.ID,
				WebexRecipientID: res.recip.ID,
				NotifyEnabled:    true,
			}
			rm.status = statusRefresh
			cmds = append(cmds, rm.submitRule(rule))
		}

	default:
		var cmd tea.Cmd
		rm.table, cmd = rm.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	rm.spinner, cmd = rm.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return rm, tea.Batch(cmds...)
}

func (rm *rulesModel) View() string {
	switch rm.status {
	case statusInit:
		return fillSpaceCentered(useSpinner(rm.spinner, "Loading rules..."), rm.availWidth, rm.availHeight)
	case statusRefresh:
		return fillSpaceCentered(useSpinner(rm.spinner, "Refreshing..."), rm.availWidth, rm.availHeight)
	case statusMain:
		return rm.table.View()
	case statusLoadingFormData:
		return fillSpaceCentered(useSpinner(rm.spinner, "Loading form data..."), rm.availWidth, rm.availHeight)
	case statusEntry:
		return rm.form.View()
	}

	return rm.table.View()
}

func (rm *rulesModel) Status() subModelStatus {
	return rm.status
}

func (rm *rulesModel) Form() *huh.Form {
	return rm.form
}

func (rm *rulesModel) setModuleDimensions() {
	rm.setTableDimensions(rm.availWidth, rm.availHeight)
}

func (rm *rulesModel) setTableDimensions(w, h int) {
	t := &rm.table
	enableW := 8
	boardW := 20
	remainingW := w - enableW - boardW
	recipW := remainingW
	t.SetColumns([]table.Column{
		{Title: "ENABLED", Width: enableW},
		{Title: "BOARD", Width: boardW},
		{Title: "RECIPIENT", Width: recipW},
	})

	t.SetRows(rulesToRows(rm.rules))
	t.SetHeight(h)
}

func (rm *rulesModel) prepareForm() tea.Cmd {
	return func() tea.Msg {
		boards, err := rm.SDKClient.ListBoards()
		if err != nil {
			return errMsg{fmt.Errorf("listing boards: %w", err)}
		}

		if len(boards) == 0 {
			return noBoardsAvailMsg
		}
		sortBoards(boards)

		recips, err := rm.SDKClient.ListRecipients()
		if err != nil {
			return errMsg{fmt.Errorf("listing webex recipients: %w", err)}
		}

		if len(recips) == 0 {
			return noRecipsAvailMsg
		}
		sortRecips(recips)

		return ruleFormDataMsg{
			boards: boards,
			recips: recips,
		}
	}
}

func (rm *rulesModel) submitRule(rule *models.NotifierRule) tea.Cmd {
	return func() tea.Msg {
		_, err := rm.SDKClient.CreateNotifierRule(rule)
		if err != nil {
			return errMsg{fmt.Errorf("creating notifier rule: %w", err)}
		}

		return refreshRulesMsg{}
	}
}

func (rm *rulesModel) deleteRule(id int) tea.Cmd {
	return func() tea.Msg {
		if err := rm.SDKClient.DeleteNotifierRule(id); err != nil {
			return errMsg{fmt.Errorf("deleting notifier rule: %w", err)}
		}

		return refreshRulesMsg{}
	}
}

func (rm *rulesModel) getRules() tea.Cmd {
	return func() tea.Msg {
		rules, err := rm.SDKClient.ListNotifierRules()
		if err != nil {
			return errMsg{fmt.Errorf("getting rules: %w", err)}
		}

		return gotRulesMsg{rules: rules}
	}
}

func (rm *rulesModel) setRows() tea.Cmd {
	sortRules(rm.rules)
	rm.table.SetRows(rulesToRows(rm.rules))
	rm.table.SetCursor(0)
	return nil
}

func rulesToRows(rules []models.NotifierRuleFull) []table.Row {
	if len(rules) == 0 {
		return []table.Row{
			{
				"NO", "RULEZ", "FOUND",
			},
		}
	}
	var rows []table.Row
	for _, r := range rules {
		recip := fmt.Sprintf("%s (%s)", r.RecipientName, r.RecipientType)
		rows = append(rows, []string{boolToIcon(r.Enabled), r.BoardName, recip})
	}

	return rows
}

func ruleEntryForm(boards []models.Board, recips []models.WebexRecipient, result *rulesFormResult, height int) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[models.Board]().
				Title("Connectwise Board").
				Options(boardsToFormOpts(boards)...).
				Value(&result.board),
		),
		huh.NewGroup(
			huh.NewSelect[models.WebexRecipient]().
				Title("Webex Recipient").
				Options(recipsToFormOpts(recips, nil)...).
				Value(&result.recip),
		),
	).WithTheme(huh.ThemeBase16()).WithHeight(height + 1).WithShowHelp(false) // add +1 to height to account for not showing help
}
