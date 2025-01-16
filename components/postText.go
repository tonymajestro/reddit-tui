package components

import "github.com/charmbracelet/lipgloss"

var (
	postTextStyle          = lipgloss.NewStyle()
	postTextContainerStyle = lipgloss.NewStyle().Margin(0, 2, 1, 2)
)

type PostText struct {
	Style    lipgloss.Style
	Contents string
	W        int
}

func NewPostText() PostText {
	return PostText{Style: postTextStyle}
}

func (p PostText) View() string {
	view := p.Style.Render(p.Contents)
	return postTextContainerStyle.Render(view)
}

func (p *PostText) SetSize(w, h int) {
	p.W = w - postTextContainerStyle.GetHorizontalFrameSize()
	p.Style = p.Style.Width(p.W)
}
