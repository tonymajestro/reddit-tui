package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	defaultTitleStyle = lipgloss.NewStyle().
				MarginBottom(1).
				Padding(0, 2).
				Height(1).
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230"))

	defaultDescriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	headerContainerStyle    = lipgloss.NewStyle().Margin(0, 2, 1, 2)
)

type Header struct {
	TitleStyle       lipgloss.Style
	DescriptionStyle lipgloss.Style
	Title            string
	Description      string
	W                int
}

func NewHeader() Header {
	return Header{
		TitleStyle:       defaultTitleStyle,
		DescriptionStyle: defaultDescriptionStyle,
	}
}

func (h *Header) SetTitle(title string) {
	h.Title = title
}

func (h *Header) SetDescription(description string) {
	h.Description = description
}

func (h *Header) SetSize(width, height int) {
	h.W = width - headerContainerStyle.GetHorizontalFrameSize()
	h.DescriptionStyle = h.DescriptionStyle.Width(h.W)
}

func (h Header) View() string {
	titleView := h.GetTitleView()
	descriptionView := h.GetDescriptionView()
	joinedView := lipgloss.JoinVertical(lipgloss.Left, titleView, descriptionView)
	return headerContainerStyle.Render(joinedView)
}

func (h Header) GetTitleView() string {
	return h.TitleStyle.Render(trim(h.Title, h.W))
}

func (h Header) GetDescriptionView() string {
	return h.DescriptionStyle.Render(h.Description)
}

func trim(s string, w int) string {
	if len(s) <= w {
		return s
	}

	return fmt.Sprintf("%s...", s[:w-3])
}
