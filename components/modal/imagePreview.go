package modal

import (
	"reddittui/components/messages"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ImagePreviewModal struct {
	Content string
	w       int
	h       int
}

func ImagePreviewContentSize(viewportW, viewportH int) (int, int) {
	modalWidth := int(float64(viewportW) * 0.97)
	modalHeight := int(float64(viewportH) * 0.97)

	// PlaceModal caps rendered modal width to maxModalWidthPercentage of max width.
	cappedModalWidth := min(modalWidth, int(float64(modalWidth)*maxModalWidthPercentage))

	// imageStyle has border(2) + padding(0,1) = 4 cols, border = 2 rows
	// View() has title(1) + help(1) + empty(1) = 3 rows
	contentWidth := max(1, cappedModalWidth-4)
	contentHeight := max(1, modalHeight-2-3)
	return contentWidth, contentHeight
}

func NewImagePreviewModal() ImagePreviewModal {
	return ImagePreviewModal{}
}

func (m ImagePreviewModal) View() string {
	title := lipgloss.NewStyle().Bold(true).Render("Image preview")

	contentStyle := lipgloss.NewStyle().Width(m.w).Height(m.h).MaxHeight(m.h)
	content := contentStyle.Render(m.Content)

	container := lipgloss.JoinVertical(lipgloss.Left, title, "", content)
	return lipgloss.NewStyle().Width(m.w).Render(container)
}

func (m ImagePreviewModal) Update(msg tea.Msg) (ImagePreviewModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "enter", "h", "backspace":
			m.Content = ""
			return m, messages.ExitModal
		}
	}

	return m, nil
}

func (m *ImagePreviewModal) SetSize(w, h int) {
	m.w = max(1, w)
	m.h = max(1, h)
}
