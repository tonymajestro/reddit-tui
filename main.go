package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reddittui/components"
	"reddittui/config"
	"reddittui/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type CliArgs struct {
	subreddit string
	postId    string
}

func main() {
	configuration, _ := config.LoadConfig()

	logFile, err := utils.InitLogger(configuration.Core.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open logfile: %v\n", err)
	}

	defer logFile.Close()

	var args CliArgs
	flag.StringVar(&args.postId, "post", "", "Post id")
	flag.StringVar(&args.subreddit, "subreddit", "", "Subreddit")
	flag.Parse()

	reddit := components.NewRedditTui(configuration, args.subreddit, args.postId)
	p := tea.NewProgram(reddit, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		slog.Error("Error running reddittui, see logfile for details", "error", err)
		os.Exit(1)
	}
}
