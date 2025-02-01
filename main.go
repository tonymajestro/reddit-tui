package main

import (
	"log/slog"
	"os"
	"reddittui/components"
	"reddittui/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	utils.InitLogger()

	reddit := components.NewRedditTui()
	p := tea.NewProgram(reddit, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		slog.Error("Error running reddit-tui", "error", err)
		os.Exit(1)
	}
}
