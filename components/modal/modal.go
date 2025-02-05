package modal

import (
	"log/slog"
	"reddittui/components/colors"
	"reddittui/components/messages"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SessionState int

const (
	defaultState SessionState = iota
	loading
	searching
	quitting
	showingError
)

var modalStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder(), true).
	BorderForeground(colors.AdaptiveColor(colors.Blue)).
	Padding(1, 2).
	Margin(1, 1)

type ModalManager struct {
	quit       QuitModal
	search     SubredditSearchModal
	spinner    SpinnerModal
	errorModal ErrorModal
	state      SessionState
	style      lipgloss.Style
}

func NewModalManager() ModalManager {
	return ModalManager{
		quit:       NewQuitModal(),
		search:     NewSubredditSearchModal(),
		spinner:    NewSpinnerModal(),
		errorModal: NewErrorModal(),
		style:      modalStyle,
	}
}

func (m ModalManager) Init() tea.Cmd {
	return nil
}

func (m ModalManager) Update(msg tea.Msg) (ModalManager, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.state != defaultState {
		m, cmd = m.handleFocusedMessages(msg)
		cmds = append(cmds, cmd)
	}

	m, cmd = m.handleGlobalMessages(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ModalManager) handleGlobalMessages(msg tea.Msg) (ModalManager, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case messages.ShowSpinnerModalMsg:
		loadingMsg := string(msg)
		return m, m.setLoading(loadingMsg)

	case messages.ShowErrorModalMsg:
		return m, m.setError(string(msg))

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			if m.state == defaultState {
				return m, m.setQuitting()
			}
		case "s", "S":
			return m, m.setSearching()
		}
	}

	return m, nil
}

func (m ModalManager) handleFocusedMessages(msg tea.Msg) (ModalManager, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case loading:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case quitting:
		m.quit, cmd = m.quit.Update(msg)
		return m, cmd
	case searching:
		m.search, cmd = m.search.Update(msg)
		return m, cmd
	case showingError:
		m.errorModal, cmd = m.errorModal.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m ModalManager) View(background Viewer) string {
	switch m.state {
	case loading:
		return PlaceModal(m.spinner, background, lipgloss.Center, lipgloss.Center, m.style)
	case quitting:
		return PlaceModal(m.quit, background, lipgloss.Center, lipgloss.Center, m.style)
	case searching:
		return PlaceModal(m.search, background, lipgloss.Center, lipgloss.Center, m.style)
	case showingError:
		return PlaceModal(m.errorModal, background, lipgloss.Center, lipgloss.Center, m.style)
	default:
		// This sometimes happens when loading completes before the loading modal finishes rendering
		return ""
	}
}

func (m *ModalManager) SetSize(w, h int) {
	m.search.SetSize(w, h)

	modalSize := int((float64(w) * (2)) / 3.0)
	m.style = m.style.MaxWidth(modalSize)
}

func (m *ModalManager) Blur() {
	m.state = defaultState
	m.search.Blur()
}

func (m *ModalManager) setLoading(message string) tea.Cmd {
	m.state = loading
	m.spinner.SetLoading(message)
	return m.spinner.Tick
}

func (m *ModalManager) setSearching() tea.Cmd {
	m.state = searching
	m.search.Focus()
	return messages.OpenModal
}

func (m *ModalManager) setQuitting() tea.Cmd {
	m.state = quitting
	return messages.OpenModal
}

func (m *ModalManager) setError(errorMsg string) tea.Cmd {
	m.state = showingError
	slog.Info("Setting error modal")
	m.errorModal.ErrorMsg = errorMsg
	return messages.OpenModal
}
