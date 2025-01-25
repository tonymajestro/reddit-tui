package components

import (
	"fmt"
	"reddittui/components/colors"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerStyle          = lipgloss.NewStyle().Foreground(colors.Blue)
	spinnerContainerStyle = lipgloss.NewStyle().Margin(2, 2)
)

type Spinner struct {
	spinner.Model
	Style          lipgloss.Style
	LoadingMessage string
	Loading        bool
}

func NewSpinner() Spinner {
	model := spinner.New()
	model.Spinner = spinner.Dot
	model.Style = spinnerStyle

	return Spinner{
		Model:   model,
		Loading: true,
	}
}

func (s Spinner) Init() tea.Cmd {
	return s.Tick
}

func (s Spinner) Update(msg tea.Msg) (Spinner, tea.Cmd) {
	var cmd tea.Cmd
	s.Model, cmd = s.Model.Update(msg)
	return s, cmd
}

func (s Spinner) View() string {
	view := fmt.Sprintf("%s %s", s.Model.View(), s.LoadingMessage)
	return spinnerContainerStyle.Render(view)
}

func (s *Spinner) SetLoading(loading bool) {
	s.Loading = loading
	if loading {
		model := spinner.New()
		model.Spinner = spinner.Dot
		model.Style = spinnerStyle
		s.Model = model
	}
}
