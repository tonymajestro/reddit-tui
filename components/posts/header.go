package posts

import (
	"reddittui/components/colors"
	"reddittui/utils"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerContainerStyle = lipgloss.NewStyle().MarginBottom(2)
	titleStyle           = lipgloss.NewStyle().
				MarginBottom(1).
				Padding(0, 2).
				Height(1).
				Background(colors.AdaptiveColors(colors.Blue, colors.Indigo)).
				Foreground(colors.AdaptiveColors(colors.White, colors.Sand))

	defaultDescriptionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colors.AdaptiveColor(colors.Text))
)

type PostsHeader struct {
	DescriptionStyle lipgloss.Style
	Title            string
	Description      string
	W                int
}

func NewPostsHeader() PostsHeader {
	return PostsHeader{DescriptionStyle: defaultDescriptionStyle}
}

func (h *PostsHeader) SetSize(width, height int) {
	h.W = width - headerContainerStyle.GetHorizontalFrameSize()
	h.DescriptionStyle = h.DescriptionStyle.Width(h.W)
}

func (h PostsHeader) View() string {
	titleView := titleStyle.Render(utils.TruncateString(h.Title, h.W))
	descriptionView := h.DescriptionStyle.Render(h.Description)

	joinedView := lipgloss.JoinVertical(lipgloss.Left, titleView, descriptionView)
	return headerContainerStyle.Render(joinedView)
}

func (h *PostsHeader) SetContent(title, desc string) {
	h.Title = utils.NormalizeSubreddit(title)
	h.Description = desc
}
