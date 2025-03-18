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
	"github.com/javiorfo/bitsmuggler/opensubs"
	"github.com/javiorfo/bitsmuggler/torr"
	"github.com/javiorfo/bitsmuggler/yts"
	"github.com/javiorfo/nilo"
	"github.com/javiorfo/steams"
)

type model struct {
	tableMovies       table.Model
	filteredRows      []table.Row
	movies            []yts.Movie
	total             int
	page              int
	totalPages        int
	tableSubs         table.Model
	subtitles         []opensubs.Subtitle
	showSubs          bool
	textInput         textinput.Model
	filterText        string
	spinner           spinner.Model
	loading           bool
	downloader        *torr.Downloader
	movieDownloadInfo string
	cancelDownload    chan struct{}
}

func (m *model) totalAndPages() string {
	return fmt.Sprintf(" %d movie/s - Page %d/%d ", m.total, m.page, m.totalPages)
}

func (m *model) updateTable(total int, movies []yts.Movie) {
	m.loading = false
	m.movies = movies

	rows := steams.Mapping(steams.OfSlice(movies), func(m yts.Movie) table.Row {
		ts := m.GetTechSpec()
		return table.Row{m.Year, m.Name, ts.Size, m.Genre, m.Rate, ts.Duration, ts.Resolution, ts.Language}
	}).Collect()

	m.filteredRows = rows
	m.tableMovies.SetRows(m.filteredRows)
	m.total = total

	if total > 20 {
		m.totalPages = int(math.Ceil(float64(total) / float64(20)))
	} else {
		m.totalPages = 1
	}
}

func (m *model) updateTableSubs(year, movie string) error {
	m.toggleTables()

	subs, err := opensubs.GetSubs(year, movie)
	if err != nil {
		return err
	}

	m.subtitles = subs
	rows := steams.Mapping(steams.OfSlice(m.subtitles), func(s opensubs.Subtitle) table.Row {
		return table.Row{s.Name, s.Date, s.Downloads}
	}).Collect()

	m.tableSubs.SetRows(rows)
	return nil
}

func (m *model) toggleTables() {
	if m.showSubs {
		m.showSubs = false
		m.tableSubs.Blur()
		m.tableMovies.Focus()
		return
	}

	if configuration.OpenSubsEnable {
		m.showSubs = true
		m.tableMovies.Blur()
		m.tableSubs.Focus()
	}
}

type onYTSResponseMsg struct {
	total  int
	movies []yts.Movie
}

func (m *model) request(page int) func() tea.Msg {
	return func() tea.Msg {
		m.filterText = m.textInput.Value()

		genre := filter(m.filterText, "genre:").OrElse("all")
		rating := filter(m.filterText, "rating:").OrElse("0")
		year := filter(m.filterText, "year:").OrElse("0")
		order := filter(m.filterText, "order:").OrElse(configuration.YTSOrderBy)
		filter := trimFilter(m.filterText, "genre:", "rating:", "year:", "order:")

		total, movies := yts.GetMovies(configuration.YTSHost, filter, genre, rating, year, order, page)

		return onYTSResponseMsg{total: total, movies: movies}
	}
}

func (m *model) getTorrentFileLink() nilo.Optional[string] {
	return steams.OfSlice(m.tableMovies.Rows()).Position(func(tr table.Row) bool {
		return slices.Compare(tr, m.tableMovies.SelectedRow()) == 0
	}).MapToString(func(i int) string {
		return m.movies[i].GetTechSpec().TorrentFileLink
	})
}

func (m *model) getSubtitleCode() nilo.Optional[string] {
	return steams.OfSlice(m.tableSubs.Rows()).Position(func(tr table.Row) bool {
		return slices.Compare(tr, m.tableSubs.SelectedRow()) == 0
	}).MapToString(func(i int) string {
		return m.subtitles[i].GetDownloadSubtitleCode()
	})
}

func (m *model) download() tea.Msg {
	m.downloader = &torr.Downloader{}
	m.downloader.Scan(m.cancelDownload)
	return nil
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.request(1), m.download)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+r":
			if m.downloader != nil && m.downloader.Downloading {
				once.Do(func() { close(m.cancelDownload) })
			}

			return m, cmd
		case "ctrl+n":
			if m.page < m.totalPages {
				m.loading = true
				m.page += 1
			}

			if m.showSubs {
				m.toggleTables()
			}

			return m, tea.Batch(m.spinner.Tick, m.request(m.page))
		case "ctrl+p":
			if m.page > 1 {
				m.loading = true
				m.page -= 1
			}

			if m.showSubs {
				m.toggleTables()
			}

			return m, tea.Batch(m.spinner.Tick, m.request(m.page))
		case "enter":
			m.loading = true

			if m.showSubs {
				m.toggleTables()
			}

			m.page = 1
			return m, tea.Batch(m.spinner.Tick, m.request(m.page))
		case "ctrl+s":
			year := m.tableMovies.SelectedRow()[0]
			name := m.tableMovies.SelectedRow()[1]
			err := m.updateTableSubs(year, name)

			if err != nil {
				return m, tea.Batch(tea.Println(errorStyle(fmt.Sprintf("  %s", err.Error()))))
			}

			return m, cmd
		case "ctrl+d":
			tFile := m.getTorrentFileLink()
			if m.showSubs && tFile.IsPresent() {
				subCode := m.getSubtitleCode()
				movieTorrentName, err := torr.MovieTorrentName(tFile.Get())

				if err != nil {
					return m, tea.Batch(tea.Println(errorStyle(fmt.Sprintf("  %s", err.Error()))))
				}

				if subCode.IsPresent() {
					err := opensubs.DownloadSubtitle(subCode.Get(), movieTorrentName)
					print := downloadStyle(fmt.Sprintf("  %s.srt Downloaded!", movieTorrentName))
					if err != nil {
						print = errorStyle(fmt.Sprintf("  %s", err.Error()))
					}
					return m, tea.Batch(tea.Println(print))
				}
				return m, tea.Batch(tea.Println(errorStyle("  Subtitle could not be downloaded")))
			}

			if tFile.IsPresent() && (m.downloader == nil || !m.downloader.Downloading) {
				m.downloader = &torr.Downloader{}
				err := m.downloader.DownloadTorrentFile(tFile.Get())
				if err != nil {
					return m, tea.Batch(tea.Println(errorStyle(fmt.Sprintf("  %s", err.Error()))))
				}
				m.downloader.Scan(m.cancelDownload)
			}

			return m, cmd
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case onYTSResponseMsg:
		m.updateTable(msg.total, msg.movies)
		return m, cmd
	case torr.OnDownload:
		m.movieDownloadInfo = string(msg)
		return m, cmd
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.tableMovies, cmd = m.tableMovies.Update(msg)
	m.tableSubs, cmd = m.tableSubs.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if torr.Status != "" {
		m.movieDownloadInfo = downloadStyle(string(torr.Status)) + "\n"
	}

	var sp string
	if m.loading {
		sp = m.spinner.View() + " searching movies\n"
	}

	var tableToShow string
	if m.showSubs {
		tableToShow = baseStyle(m.tableSubs.View()) + "\n"
	} else {
		tableToShow = baseStyle(m.tableMovies.View()) + "\n"
	}

	return m.movieDownloadInfo + baseStyle(m.textInput.View()) + "\n" + sp + baseStyle(m.totalAndPages()) + "\n" + tableToShow
}

func filter(input, filter string) nilo.Optional[string] {
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
