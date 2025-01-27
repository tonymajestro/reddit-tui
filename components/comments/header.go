package comments

import (
	"fmt"
	"reddittui/client"
	"reddittui/components/header"
	"reddittui/utils"

	"github.com/charmbracelet/lipgloss"
)

type CommentsHeader struct {
	DescriptionStyle lipgloss.Style
	Title            string
	Description      string
	Author           string
	Timestamp        string
	Points           string
	W                int
}

func NewCommentsHeader() CommentsHeader {
	return CommentsHeader{DescriptionStyle: header.DefaultDescriptionStyle}
}

func (h *CommentsHeader) SetSize(width, height int) {
	h.W = width - header.HeaderContainerStyle.GetHorizontalFrameSize()
	h.DescriptionStyle = h.DescriptionStyle.Width(h.W)
}

func (h CommentsHeader) View() string {
	titleView := header.TitleStyle.Render(utils.TruncateString(h.Title, h.W))
	descriptionView := h.DescriptionStyle.Render(h.Description)

	authorView := postAuthorStyle.Render(h.Author)
	timestampView := postTimestampStyle.Render(fmt.Sprintf("submitted %s by", h.Timestamp))
	authorTimestampView := fmt.Sprintf("%s %s", timestampView, authorView)

	postPointsView := postPointsStyle.Render(fmt.Sprintf("%s points", h.Points))
	joinedView := lipgloss.JoinVertical(lipgloss.Left, titleView, descriptionView, authorTimestampView, postPointsView)

	return header.HeaderContainerStyle.Render(joinedView)
}

func (h *CommentsHeader) SetContent(comments client.Comments) {
	h.Title = utils.NormalizeSubreddit(comments.Subreddit)
	h.Description = comments.PostTitle
	h.Author = comments.PostAuthor
	h.Points = comments.PostPoints
	h.Timestamp = comments.PostTimestamp
}
