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

const version = "v0.3.2"

type CliArgs struct {
	subreddit   string
	postId      string
	showVersion bool
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
	flag.BoolVar(&args.showVersion, "version", false, "Version")
	flag.Parse()

	if args.showVersion {
		fmt.Printf("reddittui version %s\n", version)
		os.Exit(0)
	}

	reddit := components.NewRedditTui(configuration, args.subreddit, args.postId)
	p := tea.NewProgram(reddit, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		slog.Error("Error running reddittui, see logfile for details", "error", err)
		os.Exit(1)
	}
}
