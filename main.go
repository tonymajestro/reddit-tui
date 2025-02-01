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
	p := tea.NewProgram(reddit, tea.WithAltScreen(), tea.WithMouseCellMotion())

	slog.Info("Starting reddit-tui", "test", "123")

	if _, err := p.Run(); err != nil {
		slog.Error("Error running program", "error", err)
		os.Exit(1)
	}

	slog.Info("Exiting reddit-tui")
}
