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
	configuration, _ := config.LoadConfig()

	logFile, err := utils.InitLogger(configuration.Core.LogLevel)
	if err != nil {
		fmt.Printf("Could not open logfile: %v\n", err)
	}

	defer logFile.Close()

	reddit := components.NewRedditTui(configuration)
	p := tea.NewProgram(reddit, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		slog.Error("Error running reddittui", "error", err)
		os.Exit(1)
	}
}
