package reddit

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultLoadingMessage = "loading reddit.com..."

var (
	spinnerStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	spinnerContainerStyle = lipgloss.NewStyle().Margin(2, 4)
)

type RedditSpinner struct {
	loadingMessage string
	model          spinner.Model
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
	var cmd tea.Cmd
	s.model, cmd = s.model.Update(msg)
	return s, cmd
}

func (s RedditSpinner) View() string {
	spinnerView := fmt.Sprintf("%s %s", s.model.View(), s.loadingMessage)
	return spinnerContainerStyle.Render(spinnerView)
}
