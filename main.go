package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://old.reddit.com"

type titleMsg []string

type titleErr struct {
	error
}

type post struct {
	title        string
	author       string
	subreddit    string
	friendlyDate string
	postUrl      string
	commentsUrl  string
}

type model struct {
	choices []string
	cursor  int
	err     error
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return getPosts()
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case titleErr:
		m.err = msg.error
		return m, tea.Quit
	case titleMsg:
		m.choices = msg
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "reddit.com"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s - %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
