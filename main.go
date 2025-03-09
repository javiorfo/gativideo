package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/javiorfo/bitsmuggler/yts"
	"github.com/javiorfo/steams"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table        table.Model
	textInput    textinput.Model
	filterText   string
	pages        string
	total        string
	filteredRows []table.Row
}

func getRows(keyword string) (int, []table.Row) {
	total, movies := yts.GetMovies("https://en.yts-official.mx", keyword, "all", "all", "0", "0", "rating", 1)
	rows := steams.Mapping(steams.OfSlice(movies), func(m yts.Movie) table.Row {
		torrent := m.Torrents[0]
		return table.Row{m.Year, m.Name, torrent.Size, m.Genre, m.Rate, torrent.Duration, torrent.Resolution, torrent.Language}
	}).Collect()
	return total, rows
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Filter table..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 49

	columns := []table.Column{
		{Title: "YEAR", Width: 5},
		{Title: "NAME", Width: 50},
		{Title: "SIZE", Width: 10},
		{Title: "GENRE", Width: 20},
		{Title: "RATE", Width: 4},
		{Title: "DURATION", Width: 12},
		{Title: "RESOLUTION", Width: 10},
		{Title: "LANGUAGE", Width: 12},
	}

	total, rows := getRows("")

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
		table.WithKeyMap(table.KeyMap{
			LineUp: key.NewBinding(
				key.WithKeys("up"),
				key.WithHelp("↑/k", "up"),
			),
			LineDown: key.NewBinding(
				key.WithKeys("down"),
				key.WithHelp("↓/j", "down"),
			),
			PageUp: key.NewBinding(
				key.WithKeys("pgup"),
				key.WithHelp("b/pgup", "page up"),
			),
			PageDown: key.NewBinding(
				key.WithKeys("pgdown"),
				key.WithHelp("f/pgdn", "page down"),
			),
			HalfPageUp: key.NewBinding(
				key.WithKeys("ctrl+u"),
				key.WithHelp("u", "½ page up"),
			),
			HalfPageDown: key.NewBinding(
				key.WithKeys("ctrl+d"),
				key.WithHelp("d", "½ page down"),
			),
			GotoTop: key.NewBinding(
				key.WithKeys("home"),
				key.WithHelp("g/home", "go to start"),
			),
			GotoBottom: key.NewBinding(
				key.WithKeys("end"),
				key.WithHelp("G/end", "go to end"),
			),
		}),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("240")).
		Bold(false)
	t.SetStyles(s)

	return model{
		table:        t,
		textInput:    ti,
		filterText:   "",
		pages:        fmt.Sprintf("Page 1/%d ", total/20),
		total:        fmt.Sprintf(" %d movies - ", total),
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
			/* 			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			) */
			m.filterText = m.textInput.Value()
			_, rows := getRows(m.filterText)
			m.filteredRows = rows
			m.table.SetRows(m.filteredRows)

			return m, cmd
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	/* 	if m.filterText == "" {
		return baseStyle.Render(m.textInput.View()) + "\n"
	} */
	return baseStyle.Render(m.textInput.View()) + "\n" + baseStyle.Render(m.total+m.pages) + "\n" + baseStyle.Render(m.table.View()) + "\n"
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
	// 	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
