package main

import (
	"fmt"
	"os"
	"reddittui/components"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	reddit := components.NewRedditTui()

	p := tea.NewProgram(reddit, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
