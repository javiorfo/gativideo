package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table        table.Model
	textInput    textinput.Model
	filterText   string
	filteredRows []table.Row
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Filter table..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 49

	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "City", Width: 20},
		{Title: "Country", Width: 10},
		{Title: "Population", Width: 10},
	}

	rows := []table.Row{
		{"1", "Tokyo", "Japan", "37,274,000"},
		{"2", "Delhi", "India", "32,065,760"},
		{"3", "Shanghai", "China", "28,516,904"},
		{"4", "Dhaka", "Bangladesh", "22,478,116"},
		{"5", "São Paulo", "Brazil", "22,429,800"},
		{"6", "Mexico City", "Mexico", "22,085,140"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return model{
		table:        t,
		textInput:    ti,
		filterText:   "",
		filteredRows: rows,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}

	// Update the text input and filter rows
	m.textInput, cmd = m.textInput.Update(msg)
	m.filterText = m.textInput.Value()
	m.filteredRows = filterRows(m.table.Rows(), m.filterText)
	m.table.SetRows(m.filteredRows)

	// Update the table model
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.textInput.View()) + "\n" + baseStyle.Render(m.table.View()) + "\n"
}

func filterRows(rows []table.Row, filterText string) []table.Row {
	var filteredRows []table.Row
	for _, row := range rows {
		for _, cell := range row {
			if strings.ToLower(lipgloss.NewStyle().Render(cell)) == filterText || strings.Contains(strings.ToLower(lipgloss.NewStyle().Render(cell)), filterText) {
				filteredRows = append(filteredRows, row)
				break
			}
		}
	}
	return filteredRows
}

func main() {
	m := initialModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
