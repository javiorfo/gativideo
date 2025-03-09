package model

import (
	"fmt"
	"math"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/javiorfo/bitsmuggler/yts"
	"github.com/javiorfo/steams"
)

type model struct {
	table        table.Model
	textInput    textinput.Model
	filterText   string
	total        int
	page         int
	totalPages   int
	filteredRows []table.Row
}

func (m *model) totalAndPages() string {
	return fmt.Sprintf(" %d movies - Page %d/%d ", m.total, m.page, m.totalPages)
}

func (m *model) updateTable(page int) {
	m.filterText = m.textInput.Value()
	total, rows := getRows(m.filterText, page)
	m.filteredRows = rows
	m.table.SetRows(m.filteredRows)
	m.total = total

	if total > 20 {
		m.totalPages = int(math.Ceil(float64(total) / float64(20)))
	} else {
		m.totalPages = 1
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
		case "ctrl+n":
			if m.page < m.totalPages {
				m.page += 1
				m.updateTable(m.page)
			}
			return m, cmd
		case "ctrl+b":
			if m.page > 1 {
				m.page -= 1
				m.updateTable(m.page)
			}
			return m, cmd
		case "enter":
			m.updateTable(1)
			return m, cmd
		case "ctrl+a":
			return m, tea.Batch(
				tea.Printf("Movie %s", m.table.SelectedRow()[1]),
			)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.textInput.View()) + "\n" + baseStyle.Render(m.totalAndPages()) + "\n" + baseStyle.Render(m.table.View()) + "\n"
}

func getRows(keyword string, page int) (int, []table.Row) {
	total, movies := yts.GetMovies("https://en.yts-official.mx", keyword, "all", "all", "0", "0", "rating", page)
	rows := steams.Mapping(steams.OfSlice(movies), func(m yts.Movie) table.Row {
		torrent := m.GetTorrent("1080")
		return table.Row{m.Year, m.Name, torrent.Size, m.Genre, m.Rate, torrent.Duration, torrent.Resolution, torrent.Language}
	}).Collect()
	return total, rows
}
