package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
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
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230"))

	blurButtonStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("7"))

	focusButtonStyle = lipgloss.NewStyle().
				Bold(true).
				Padding(0, 1).
				Foreground(lipgloss.Color("#EE6FF8"))

	buttonsContainerStyle = lipgloss.NewStyle().MarginBottom(2)

	quitContainerStyle = lipgloss.NewStyle().Margin(1, 2)
)

type keymap struct {
	Left  key.Binding
	Right key.Binding
	Yes   key.Binding
	No    key.Binding
	Tab   key.Binding
	Enter key.Binding
	Quit  key.Binding
	Help  key.Binding
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Enter, k.Help}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right, k.Yes, k.No},
		{k.Tab, k.Help, k.Enter},
	}
}

var keys = keymap{
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

type QuitPage struct {
	keys        keymap
	help        help.Model
	yesStyle    lipgloss.Style
	noStyle     lipgloss.Style
	focus       bool
	yesSelected bool
}

func NewQuitPage() QuitPage {
	quitPage := QuitPage{
		keys:     keys,
		help:     help.New(),
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

	buttonsView := buttonsContainerStyle.Render(fmt.Sprintf("    %s  %s", noView, yesView))
	helpView := q.help.View(q.keys)
	joinedView := lipgloss.JoinVertical(lipgloss.Left, promptView, buttonsView, helpView)

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
			return q, GoBack

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
				return q, GoBack
			}

		case key.Matches(msg, q.keys.Help):
			q.help.ShowAll = !q.help.ShowAll
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
