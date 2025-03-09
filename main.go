package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javiorfo/bitsmuggler/model"
)

func main() {
	m := model.InitialModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		// 	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
