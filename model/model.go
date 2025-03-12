package model

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/javiorfo/bitsmuggler/yts"
	"github.com/javiorfo/nilo"
	"github.com/javiorfo/steams"
)

type model struct {
	table        table.Model
	textInput    textinput.Model
	spinner      spinner.Model
	loading      bool
	filterText   string
	isLoading    bool
	total        int
	page         int
	totalPages   int
	filteredRows []table.Row
	movies       []yts.Movie
}

func (m *model) totalAndPages() string {
	return fmt.Sprintf(" %d movie/s - Page %d/%d ", m.total, m.page, m.totalPages)
}

func (m *model) updateTable(total int, movies []yts.Movie) {
	m.loading = false
	m.movies = movies
	rows := moviesToRows(movies)
	m.filteredRows = rows
	m.table.SetRows(m.filteredRows)
	m.total = total

	if total > 20 {
		m.totalPages = int(math.Ceil(float64(total) / float64(20)))
	} else {
		m.totalPages = 1
	}
}

type OnResponseMsg struct {
	total  int
	movies []yts.Movie
}

func (m *model) request(page int) func() tea.Msg {
	return func() tea.Msg {
		m.filterText = m.textInput.Value()

		genre := getFilter(m.filterText, "genre:").OrElse("all")
		rating := getFilter(m.filterText, "rating:").OrElse("0")
		year := getFilter(m.filterText, "year:").OrElse("0")
		order := getFilter(m.filterText, "order:").OrElse("rating")
		filter := trimFilter(m.filterText, "genre:", "rating:", "year:", "order:")

		total, movies := yts.GetMovies(Config.YTSHost, filter, genre, rating, year, order, page)

		return OnResponseMsg{total: total, movies: movies}
	}
}

func (m model) Init() tea.Cmd { return tea.Batch(m.spinner.Tick, m.request(1)) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "ctrl+n":
			if m.page < m.totalPages {
				m.loading = true
				m.page += 1
			}
			return m, tea.Batch(m.spinner.Tick, m.request(m.page))
		case "ctrl+b":
			if m.page > 1 {
				m.loading = true
				m.page -= 1
			}
			return m, tea.Batch(m.spinner.Tick, m.request(m.page))
		case "enter":
			m.loading = true
			m.page = 1
			return m, tea.Batch(m.spinner.Tick, m.request(m.page))
		case "ctrl+a":
			index := steams.OfSlice(m.table.Rows()).Position(func(tr table.Row) bool {
				return slices.Compare(tr, m.table.SelectedRow()) == 0
			}).OrElse(-1)
			return m, tea.Batch(
				tea.Printf("Movie %s", m.movies[index].GetTorrent("1080").File),
			)
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case OnResponseMsg:
		m.updateTable(msg.total, msg.movies)
		return m, cmd
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var sp string
	if m.loading {
		sp = m.spinner.View() + " searching movies\n"
	}
	return baseStyle.Render(m.textInput.View()) + "\n" + sp + baseStyle.Render(m.totalAndPages()) + "\n" + baseStyle.Render(m.table.View()) + "\n"
}

func moviesToRows(movies []yts.Movie) []table.Row {
	rows := steams.Mapping(steams.OfSlice(movies), func(m yts.Movie) table.Row {
		torrent := m.GetTorrent("1080")
		return table.Row{m.Year, m.Name, torrent.Size, m.Genre, m.Rate, torrent.Duration, torrent.Resolution, torrent.Language}
	}).Collect()
	return rows
}

func getFilter(input, filter string) nilo.Optional[string] {
	if strings.Contains(input, filter) {
		text := input[strings.LastIndex(input, filter)+len(filter):]
		if text != "" {
			return nilo.Of(text)
		}
	}
	return nilo.Empty[string]()
}

func trimFilter(input string, filters ...string) string {
	text := input
	for _, v := range filters {
		if strings.Contains(text, v) {
			text = text[:strings.LastIndex(input, v)]
		}
	}
	return text
}
