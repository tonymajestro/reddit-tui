package modal

import (
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
	showingImagePreview
)

var modalStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder(), true).
	BorderForeground(colors.AdaptiveColor(colors.Blue)).
	Padding(1, 2).
	Margin(1, 1)

var imagePreviewModalStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder(), true).
	BorderForeground(colors.AdaptiveColor(colors.Blue)).
	Padding(0, 1).
	Margin(0, 0)

type ModalManager struct {
	quit         QuitModal
	search       SubredditSearchModal
	spinner      SpinnerModal
	errorModal   ErrorModal
	imagePreview ImagePreviewModal
	state        SessionState
	style        lipgloss.Style
	imageStyle   lipgloss.Style
	onClose      tea.Cmd
}

func NewModalManager() ModalManager {
	return ModalManager{
		quit:         NewQuitModal(),
		search:       NewSubredditSearchModal(),
		spinner:      NewSpinnerModal(),
		errorModal:   NewErrorModal(),
		imagePreview: NewImagePreviewModal(),
		style:        modalStyle,
		imageStyle:   imagePreviewModalStyle,
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
		return m, m.SetLoading(loadingMsg)

	case messages.ShowErrorModalMsg:
		return m, m.SetErrorWithCallback(msg.ErrorMsg, msg.OnClose)

	case messages.UpdateImagePreviewMsg:
		return m, m.SetImagePreview(string(msg))

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			if m.state == defaultState {
				return m, m.SetQuitting()
			}
		case "s", "S":
			return m, m.SetSearching()
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
	case showingImagePreview:
		m.imagePreview, cmd = m.imagePreview.Update(msg)
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
	case showingImagePreview:
		return PlaceModal(
			m.imagePreview,
			background,
			lipgloss.Center,
			lipgloss.Center,
			m.imageStyle,
		)
	default:
		// This sometimes happens when loading completes before the loading modal finishes rendering
		return ""
	}
}

func (m *ModalManager) SetSize(w, h int) {
	m.search.SetSize(w, h)
	m.imagePreview.SetSize(w, h)

	modalWidth := int(float64(w) * 0.97)
	modalHeight := int(float64(h) * 0.97)
	m.style = m.style.MaxWidth(modalWidth).MaxHeight(modalHeight)
	m.imageStyle = m.imageStyle.MaxWidth(modalWidth).MaxHeight(modalHeight)
}

func (m *ModalManager) Blur() tea.Cmd {
	m.state = defaultState
	m.search.Blur()

	onClose := m.onClose
	m.onClose = nil
	return onClose
}

func (m *ModalManager) SetLoading(message string) tea.Cmd {
	m.state = loading
	m.spinner.SetLoading(message)
	return m.spinner.Tick
}

func (m *ModalManager) SetSearching() tea.Cmd {
	m.state = searching
	m.search.Focus()
	return messages.OpenModal
}

func (m *ModalManager) SetQuitting() tea.Cmd {
	m.state = quitting
	return messages.OpenModal
}

func (m *ModalManager) SetError(errorMsg string) tea.Cmd {
	m.state = showingError
	m.errorModal.ErrorMsg = errorMsg
	return messages.OpenModal
}

func (m *ModalManager) SetErrorWithCallback(errorMsg string, onClose tea.Cmd) tea.Cmd {
	m.state = showingError
	m.onClose = onClose
	m.errorModal.ErrorMsg = errorMsg
	return messages.OpenModal
}

func (m *ModalManager) SetImagePreview(content string) tea.Cmd {
	m.state = showingImagePreview
	m.imagePreview.Content = content
	return messages.OpenModal
}

func (m ModalManager) IsImagePreviewOpen() bool {
	return m.state == showingImagePreview
}
