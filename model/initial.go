package model

import (
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/javiorfo/bitsmuggler/config"
)

var configuration = config.GetConfiguration()
var once sync.Once

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color(configuration.TableBorderColor)).Render

var downloadStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(configuration.DownloadColor)).Render
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Bold(true).Render

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Search movie..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 70

	columnsMovies := []table.Column{
		{Title: "YEAR", Width: 5},
		{Title: "NAME", Width: 50},
		{Title: "SIZE", Width: 10},
		{Title: "GENRE", Width: 35},
		{Title: "RATING", Width: 7},
		{Title: "DURATION", Width: 12},
		{Title: "RESOLUTION", Width: 10},
		{Title: "LANGUAGE", Width: 15},
	}
	tMovies := table.New(table.WithColumns(columnsMovies), table.WithFocused(true), table.WithHeight(21), tableKeymaps())

	columnsSubs := []table.Column{
		{Title: "NAME", Width: 100},
		{Title: "DATE", Width: 10},
		{Title: "DOWNLOADS", Width: 10},
	}
	tSubs := table.New(table.WithColumns(columnsSubs), table.WithHeight(10), tableKeymaps())

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(configuration.TableBorderColor)).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(configuration.TableSelectionFgColor)).
		Background(lipgloss.Color(configuration.TableSelectionBgColor)).
		Bold(false)

	tMovies.SetStyles(s)
	tSubs.SetStyles(s)

	sp := spinner.New()
	sp.Spinner = spinner.Points
	sp.Spinner.FPS = time.Second / 4
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(configuration.SpinnerColor))

	return model{
		tableMovies:    tMovies,
		initSearch:     configuration.YTSInitSearch,
		tableSubs:      tSubs,
		textInput:      ti,
		spinner:        sp,
		loading:        true,
		page:           1,
		cancelDownload: make(chan struct{}),
	}
}

func tableKeymaps() table.Option {
	return table.WithKeyMap(table.KeyMap{
		LineUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑ ", "up"),
		),
		LineDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓ ", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("pgdn", "page down"),
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
			key.WithHelp("home", "go to start"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("end"),
			key.WithHelp("end", "go to end"),
		),
	})
}
