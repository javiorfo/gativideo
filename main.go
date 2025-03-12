package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javiorfo/bitsmuggler/model"
)

func main() {
	m := model.InitialModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal("Error running bitsmuggler:", err)
	}
}
