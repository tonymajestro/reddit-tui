package common

import (
	"fmt"
	"reddittui/components/colors"
	"reddittui/components/messages"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	prompt        = "Are you sure you want to quit? (y/n)"
	noButtonText  = "[ no ]"
	yesButtonText = "[ yes ]"
)

var (
	promptStyle = lipgloss.NewStyle().
			MarginBottom(1).
			Padding(0, 2).
			Height(1).
			Background(colors.AdaptiveColors(colors.Blue, colors.Indigo)).
			Foreground(colors.AdaptiveColors(colors.White, colors.Sand))

	blurButtonStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Faint(true).
			Foreground(colors.AdaptiveColor(colors.Text))

	focusButtonStyle = lipgloss.NewStyle().
				Bold(true).
				Padding(0, 1).
				Foreground(colors.AdaptiveColor(colors.Blue))

	buttonsContainerStyle = lipgloss.NewStyle().MarginBottom(2)
	quitContainerStyle    = lipgloss.NewStyle().Margin(1, 2)
)

type quitKeyMap struct {
	Left  key.Binding
	Right key.Binding
	Yes   key.Binding
	No    key.Binding
	Tab   key.Binding
	Enter key.Binding
	Quit  key.Binding
	Help  key.Binding
}

func (k quitKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Enter, k.Help}
}

func (k quitKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right, k.Yes, k.No},
		{k.Tab, k.Help, k.Enter},
	}
}

type QuitPage struct {
	keys        quitKeyMap
	yesStyle    lipgloss.Style
	noStyle     lipgloss.Style
	focus       bool
	yesSelected bool
}

func NewQuitPage() QuitPage {
	keys := quitKeyMap{
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move right"),
		),
		Yes: key.NewBinding(
			key.WithKeys("y", "Y"),
			key.WithHelp("y", "quit"),
		),
		No: key.NewBinding(
			key.WithKeys("n", "N"),
			key.WithHelp("n", "back"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next selection"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm selection"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}

	quitPage := QuitPage{
		keys:     keys,
		yesStyle: focusButtonStyle,
		noStyle:  blurButtonStyle,
		focus:    false,
	}
	quitPage.SelectNo()

	return quitPage
}

func (q QuitPage) View() string {
	promptView := promptStyle.Render(prompt)
	noView := q.noStyle.Render(noButtonText)
	yesView := q.yesStyle.Render(yesButtonText)

	buttonsView := buttonsContainerStyle.Render(fmt.Sprintf("        %s  %s", noView, yesView))
	joinedView := lipgloss.JoinVertical(lipgloss.Left, promptView, buttonsView)

	return quitContainerStyle.Render(joinedView)
}

func (q QuitPage) Update(msg tea.Msg) (QuitPage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, q.keys.Yes):
			return q, tea.Quit

		case key.Matches(msg, q.keys.No):
			q.Blur()
			return q, messages.GoBack

		case key.Matches(msg, q.keys.Left):
			q.SelectNo()

		case key.Matches(msg, q.keys.Right):
			q.SelectYes()

		case key.Matches(msg, q.keys.Tab):
			if q.yesSelected {
				q.SelectNo()
			} else {
				q.SelectYes()
			}

		case key.Matches(msg, q.keys.Enter):
			if q.yesSelected {
				return q, tea.Quit
			} else {
				q.Blur()
				return q, messages.GoBack
			}
		}
	}

	return q, nil
}

func (q *QuitPage) Focus() {
	q.SelectNo()
	q.focus = true
}

func (q *QuitPage) Blur() {
	q.focus = false
}

func (q QuitPage) IsFocused() bool {
	return q.focus
}

func (q *QuitPage) SelectNo() {
	q.noStyle = focusButtonStyle
	q.yesStyle = blurButtonStyle
	q.yesSelected = false
}

func (q *QuitPage) SelectYes() {
	q.noStyle = blurButtonStyle
	q.yesStyle = focusButtonStyle
	q.yesSelected = true
}
