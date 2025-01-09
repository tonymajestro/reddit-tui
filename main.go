package main

import (
	"fmt"
	"os"
	"reddittui/reddit"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	postsList := reddit.NewRedditTui()

	p := tea.NewProgram(postsList, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
