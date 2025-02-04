package modal

import (
	"fmt"
	"reddittui/components/colors"
	"reddittui/components/messages"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	defaultErrorStyle = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Red)).Bold(true)
	errorMsgStyle     = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Text))
)

type ErrorModal struct {
	ErrorMsg string
}

func NewErrorModal() ErrorModal {
	return ErrorModal{}
}

func (e ErrorModal) Init() tea.Cmd {
	return nil
}

func (e ErrorModal) Update(msg tea.Msg) (ErrorModal, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		// Press any key to exit modal
		return e, messages.ExitModal
	}

	return e, nil
}

func (e ErrorModal) View() string {
	defaultErrorView := defaultErrorStyle.Render("Error:")
	errorMsgView := errorMsgStyle.Render(e.ErrorMsg)
	return fmt.Sprintf("%s %s", defaultErrorView, errorMsgView)
}
