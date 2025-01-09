package reddit

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	selectSubredditMsg struct{}
	hideSubredditMsg   struct{}
)

func ShowSubredditInput() tea.Cmd {
	return func() tea.Msg {
		return selectSubredditMsg{}
	}
}

func HideSubredditInput() tea.Cmd {
	return func() tea.Msg {
		return hideSubredditMsg{}
	}
}

type SubredditInput struct {
	model textinput.Model
	w, h  int
	focus bool
}

func NewSubredditInput() SubredditInput {
	model := textinput.New()
	model.Placeholder = "find subreddit"

	return SubredditInput{model: model}
}

func (s *SubredditInput) Focus() tea.Cmd {
	s.focus = true
	s.model.Reset()
	return s.model.Focus()
}

func (s *SubredditInput) Blur() {
	s.focus = false
	s.model.Blur()
}

func (s SubredditInput) Init() tea.Cmd {
	return textinput.Blink
}

func (s SubredditInput) Update(msg tea.Msg) (SubredditInput, tea.Cmd) {
	if !s.focus {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			s.Blur()
			return s, focusListPage
		case tea.KeyCtrlC:
			return s, tea.Quit
		case tea.KeyEnter:
			s.Blur()
			return s, fetchSubredditPosts(s.model.Value())
		}
	}

	var cmd tea.Cmd
	s.model, cmd = s.model.Update(msg)
	return s, cmd
}

func (s SubredditInput) View() string {
	return fmt.Sprintf("\n\n%s\n\n", s.model.View())
}
