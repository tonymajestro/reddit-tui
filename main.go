package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type post struct {
	title         string
	author        string
	subreddit     string
	friendlyDate  string
	postUrl       string
	commentsUrl   string
	totalComments string
	totalLikes    string
}

func (p post) Title() string {
	return p.title
}

func (p post) Description() string {
	return fmt.Sprintf(" %s  %s  %s comments  %s", p.totalLikes, p.subreddit, p.totalComments, p.friendlyDate)
}

func (p post) FilterValue() string {
	return p.title
}

type posts []post

type model struct {
	list    list.Model
	spinner spinner.Model
	loading bool
	w, h    int
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		posts, err := getPosts()
		if err != nil {
			fmt.Printf("Could not load reddit posts: %v", err)
			os.Exit(1)
		}

		var items []list.Item
		for _, p := range posts {
			items = append(items, p)
		}
		return items
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []list.Item:
		m.list.SetItems(msg)
		m.loading = false
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "c", "C":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.w, m.h = msg.Width-h, msg.Height-v
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if !m.loading {
		return docStyle.Render(m.list.View())
	}

	return ViewSpinner(m.spinner, m.w, m.h)
}

func main() {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "reddit.com"

	spin := NewSpinner()

	m := model{
		list:    l,
		spinner: spin,
		loading: true,
	}

	m.list.Title = "reddit.com"

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
