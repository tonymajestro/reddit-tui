package components

import "github.com/charmbracelet/lipgloss"

var (
	postTextStyle          = lipgloss.NewStyle()
	postTextContainerStyle = lipgloss.NewStyle().Margin(0, 2)
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
	if len(p.Contents) == 0 {
		return ""
	}

	view := p.Style.Render(p.Contents)
	return postTextContainerStyle.Render(view)
}

func (p *PostText) SetSize(w, h int) {
	p.W = w - postTextContainerStyle.GetHorizontalFrameSize()
	p.Style = p.Style.Width(p.W)
}
