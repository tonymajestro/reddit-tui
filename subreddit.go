package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	selectSubredditMsg struct{}
	hideSubredditMsg   struct{}
)

func showSubredditInput() tea.Cmd {
	return func() tea.Msg {
		return selectSubredditMsg{}
	}
}

func hideSubredditInput() tea.Cmd {
	return func() tea.Msg {
		return hideSubredditMsg{}
	}
}

type subredditInput struct {
	model  textinput.Model
	w, h   int
	active bool
}

func NewSubredditInput() subredditInput {
	model := textinput.New()
	model.Placeholder = "find subreddit"

	return subredditInput{model: model}
}

func (s *subredditInput) enable() tea.Cmd {
	s.active = true
	s.model.Reset()
	return s.model.Focus()
}

func (s *subredditInput) disable() {
	s.active = false
	s.model.Blur()
}

func (s subredditInput) Init() tea.Cmd {
	return textinput.Blink
}

func (s subredditInput) Update(msg tea.Msg) (subredditInput, tea.Cmd) {
	switch msg := msg.(type) {
	case selectSubredditMsg:
		return s, s.enable()
	case hideSubredditMsg:
		s.disable()
		return s, showPostsList
	case tea.KeyMsg:
		if !s.active {
			break
		}

		switch msg.Type {
		case tea.KeyEsc:
			s.disable()
			return s, showPostsList
		case tea.KeyCtrlC:
			return s, tea.Quit
		case tea.KeyEnter:
			// Have to save subreddit because disable() resets the input
			s.disable()
			return s, getSubredditPage(s.model.Value())
		}
	}

	var cmd tea.Cmd
	s.model, cmd = s.model.Update(msg)
	return s, cmd
}

func (s subredditInput) View() string {
	return fmt.Sprintf("\n\n%s\n\n", s.model.View())
}
