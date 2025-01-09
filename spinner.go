package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultSpinnerTitle = "reddit.com"

var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))

type (
	showSpinnerMsg struct {
		desc string
	}
	hideSpinnerMsg struct{}
)

func showSpinner(title string) tea.Cmd {
	return func() tea.Msg {
		return showSpinnerMsg{title}
	}
}

func hideSpinner() tea.Cmd {
	return func() tea.Msg {
		return hideSpinnerMsg{}
	}
}

type redditSpinner struct {
	model  spinner.Model
	title  string
	w, h   int
	active bool
}

func NewRedditSpinner() redditSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return redditSpinner{
		model: s,
	}
}

func (s *redditSpinner) enable() {
	s.active = true
}

func (s *redditSpinner) disable() {
	s.active = false
}

func (s redditSpinner) Init() tea.Cmd {
	return s.model.Tick
}

func (s redditSpinner) Update(msg tea.Msg) (redditSpinner, tea.Cmd) {
	switch msg := msg.(type) {

	case showSpinnerMsg:
		s.enable()
		s.title = msg.desc

	case hideSpinnerMsg:
		s.disable()

	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		s.w, s.h = msg.Width-h, msg.Height-v
	}

	var cmd tea.Cmd
	s.model, cmd = s.model.Update(msg)
	return s, cmd
}

func (s redditSpinner) View() string {
	var sb strings.Builder
	for range s.h / 2 {
		sb.WriteString("\n")
	}

	loadingMessage := fmt.Sprintf("%s loading %s...", s.model.View(), s.title)
	var line strings.Builder
	for range s.w/2 - 12 {
		line.WriteString(" ")
	}
	line.WriteString(loadingMessage)

	sb.WriteString(line.String())
	return sb.String()
}
