package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	BOLD_START = "\033[1m"
	BOLD_END   = "\033[0m"
)

type post struct {
	title        string
	author       string
	subreddit    string
	friendlyDate string
	postUrl      string
	commentsUrl  string
}

type model struct {
	err    error
	posts  []post
	cursor int
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		posts := getPosts()
		return posts
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case getPostsErr:
		m.err = msg
		return m, tea.Quit
	case postsMsg:
		m.posts = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.posts)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "reddit.com\n\n"

	for i, post := range m.posts {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += "      +------------------------------------------------------------------------------------------+"
		s += fmt.Sprintf("\n      |  %s%s%s\n", BOLD_START, post.title, BOLD_END)
		s += fmt.Sprintf(" %s    |%91s", cursor, "|")
		s += fmt.Sprintf("\n      |  %-20s %65s |\n", post.subreddit, post.friendlyDate)
	}

	s += "      +------------------------------------------------------------------------------------------+\n\n"
	s += "\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
