package reddit

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultLoadingMessage = "loading reddit.com..."

var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))

type RedditSpinner struct {
	loadingMessage string
	model          spinner.Model
	w              int
	h              int
	focus          bool
}

func NewRedditSpinner() RedditSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return RedditSpinner{
		model:          s,
		loadingMessage: defaultLoadingMessage,
	}
}

func (s *RedditSpinner) SetLoadingMessage(msg string) {
	s.loadingMessage = msg
}

func (s *RedditSpinner) Focus() {
	s.focus = true
}

func (s *RedditSpinner) Blur() {
	s.focus = false
}

func (s RedditSpinner) Init() tea.Cmd {
	return s.model.Tick
}

func (s RedditSpinner) Update(msg tea.Msg) (RedditSpinner, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		s.w, s.h = msg.Width-h, msg.Height-v
	}

	var cmd tea.Cmd
	s.model, cmd = s.model.Update(msg)
	return s, cmd
}

func (s RedditSpinner) View() string {
	var (
		sb   strings.Builder
		line strings.Builder
	)

	for range s.h / 2 {
		sb.WriteString("\n")
	}

	loadingMessage := fmt.Sprintf("%s %s", s.model.View(), s.loadingMessage)
	for range s.w/2 - (len(loadingMessage) / 2) {
		line.WriteString(" ")
	}
	line.WriteString(loadingMessage)

	sb.WriteString(line.String())
	return sb.String()
}
