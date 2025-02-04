package modal

import (
	"fmt"
	"reddittui/components/colors"
	"reddittui/components/messages"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	quitMsg  = "Are you sure you want to quit?"
	yesNoMsg = "(y/n)"
)

var (
	quitTitleStyle = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Text)).Italic(true)
	quitYesNoStyle = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Text)).Bold(true)
)

type QuitModal struct{}

func NewQuitModal() QuitModal {
	return QuitModal{}
}

func (q QuitModal) View() string {
	titleView := quitTitleStyle.Render(quitMsg)
	yesNoView := quitYesNoStyle.Render(yesNoMsg)
	return fmt.Sprintf("%s  %s", titleView, yesNoView)
}

func (q QuitModal) Update(msg tea.Msg) (QuitModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "q", "Q", "esc":
			return q, tea.Quit

		default:
			return q, messages.ExitModal
		}
	}

	return q, nil
}
