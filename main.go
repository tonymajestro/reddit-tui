package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

const url = "https://old.reddit.com"

type post struct {
	title        string
	author       string
	subreddit    string
	friendlyDate string
	postUrl      string
	commentsUrl  string
}

func (p post) Title() string {
	return p.title
}

func (p post) Description() string {
	return fmt.Sprintf("%s %s", p.subreddit, p.friendlyDate)
}

func (p post) FilterValue() string {
	return p.title
}

type posts []post

type model struct {
	list   list.Model
	cursor int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	posts, err := getPosts()
	if err != nil {
		fmt.Printf("Could not load reddit posts: %v", err)
		os.Exit(1)
	}

	var items []list.Item
	for _, p := range posts {
		items = append(items, p)
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "reddit.com"

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
