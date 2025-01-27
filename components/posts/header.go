package posts

import (
	"reddittui/components/header"
	"reddittui/utils"

	"github.com/charmbracelet/lipgloss"
)

type PostsHeader struct {
	DescriptionStyle lipgloss.Style
	Title            string
	Description      string
	W                int
}

func NewPostsHeader() PostsHeader {
	return PostsHeader{DescriptionStyle: header.DefaultDescriptionStyle}
}

func (h *PostsHeader) SetSize(width, height int) {
	h.W = width - header.HeaderContainerStyle.GetHorizontalFrameSize()
	h.DescriptionStyle = h.DescriptionStyle.Width(h.W)
}

func (h PostsHeader) View() string {
	titleView := header.TitleStyle.Render(utils.TruncateString(h.Title, h.W))
	descriptionView := h.DescriptionStyle.Render(h.Description)

	joinedView := lipgloss.JoinVertical(lipgloss.Left, titleView, descriptionView)
	return header.HeaderContainerStyle.Render(joinedView)
}

func (h *PostsHeader) SetContent(subreddit, postTitle string) {
	h.Title = utils.NormalizeSubreddit(subreddit)
	h.Description = postTitle
}
