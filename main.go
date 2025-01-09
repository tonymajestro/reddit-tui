package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	postsList      postsList
	subredditInput subredditInput
	w, h           int
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.postsList.Init(),
		m.subredditInput.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var postsListCmd, subredditInputCmd tea.Cmd

	m.postsList, postsListCmd = m.postsList.Update(msg)
	m.subredditInput, subredditInputCmd = m.subredditInput.Update(msg)

	return m, tea.Batch(postsListCmd, subredditInputCmd)
}

func (m model) View() string {
	if m.subredditInput.active {
		return m.subredditInput.View()
	} else {
		return m.postsList.View()
	}
}

func main() {
	postsList := newPostsList()
	subredditInput := NewSubredditInput()

	m := model{
		postsList:      postsList,
		subredditInput: subredditInput,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
