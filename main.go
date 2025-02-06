package main

import (
	"fmt"
	"log/slog"
	"os"
	"reddittui/components"
	"reddittui/config"
	"reddittui/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logFile, err := utils.InitLogger()
	if err != nil {
		fmt.Printf("Could not open logfile: %v\n", err)
		defer logFile.Close()
	}

	configuration, _ := config.LoadConfig()

	reddit := components.NewRedditTui(configuration)
	p := tea.NewProgram(reddit, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		slog.Error("Error running reddit-tui", "error", err)
		os.Exit(1)
	}
}
