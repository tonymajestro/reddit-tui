package modal

import (
	"fmt"
	"reddittui/components/colors"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerStyle     = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Purple))
	spinnerTextStyle = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Text)).Italic(true)
)

type SpinnerModal struct {
	spinner.Model
	LoadingMessage string
}

func NewSpinnerModal() SpinnerModal {
	model := spinner.New()
	model.Spinner = spinner.Dot
	model.Style = spinnerStyle

	return SpinnerModal{
		Model: model,
	}
}

func (s SpinnerModal) Init() tea.Cmd {
	return s.Tick
}

func (s SpinnerModal) Update(msg tea.Msg) (SpinnerModal, tea.Cmd) {
	var cmd tea.Cmd
	s.Model, cmd = s.Model.Update(msg)
	return s, cmd
}

func (s SpinnerModal) View() string {
	loadingTextView := spinnerTextStyle.Render(s.LoadingMessage)
	return fmt.Sprintf("%s %s", s.Model.View(), loadingTextView)
}

func (s *SpinnerModal) SetLoading(message string) {
	model := spinner.New()
	model.Spinner = spinner.Dot
	model.Style = spinnerStyle
	s.Model = model
	s.LoadingMessage = message
}
