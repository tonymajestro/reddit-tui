package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))

type redditSpinner struct {
	model spinner.Model
	w, h  int
}

func NewSpinner() redditSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return redditSpinner{
		model: s,
	}
}

func (s redditSpinner) Init() tea.Cmd {
	return s.model.Tick
}

func (s redditSpinner) Update(msg tea.Msg) (redditSpinner, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		s.w, s.h = msg.Width-h, msg.Height-v
	case spinner.TickMsg:
		var cmd tea.Cmd
		s.model, cmd = s.model.Update(msg)
		return s, cmd
	}

	return s, nil
}

func (s redditSpinner) View() string {
	var sb strings.Builder
	for range s.h / 2 {
		sb.WriteString("\n")
	}

	loadingMessage := fmt.Sprintf("%s %s", s.model.View(), "loading reddit.com...")
	var line strings.Builder
	for range s.w/2 - 12 {
		line.WriteString(" ")
	}
	line.WriteString(loadingMessage)

	sb.WriteString(line.String())
	return sb.String()
}
