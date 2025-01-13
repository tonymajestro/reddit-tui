package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	inputStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	inputContainerStyle = lipgloss.NewStyle().Margin(1, 2)
)

type SubredditSearch struct {
	model textinput.Model
	focus bool
}

func NewSubredditSearch() SubredditSearch {
	model := textinput.New()
	model.ShowSuggestions = true
	model.SetSuggestions(subredditSuggestions)
	model.CharLimit = 30

	return SubredditSearch{
		model: model,
	}
}

func (s SubredditSearch) IsFocused() bool {
	return s.focus
}

func (s *SubredditSearch) Focus() tea.Cmd {
	s.focus = true
	s.model.Reset()
	return s.model.Focus()
}

func (s *SubredditSearch) Blur() {
	s.focus = false
	s.model.Blur()
}

func (s SubredditSearch) Init() tea.Cmd {
	return textinput.Blink
}

func (s SubredditSearch) Update(msg tea.Msg) (SubredditSearch, tea.Cmd) {
	if !s.focus {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			s.Blur()
			return s, ReturnToPosts
		case "enter":
			s.Blur()
			return s, AcceptSearch(s.model.Value())
		case "ctrl+c":
			return s, tea.Quit
		}
	}

	var cmd tea.Cmd
	s.model, cmd = s.model.Update(msg)
	return s, cmd
}

func (s SubredditSearch) View() string {
	selectionView := inputStyle.Render(fmt.Sprintf("Choose a subreddit:\n%s", s.model.View()))
	return inputContainerStyle.Render(selectionView)
}
