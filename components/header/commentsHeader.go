package header

import (
	"fmt"
	"reddittui/client"
	"reddittui/components/colors"
	"reddittui/utils"

	"github.com/charmbracelet/lipgloss"
)

var (
	postAuthorStyle    = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Blue))
	postPointsStyle    = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Purple))
	postTextStyle      = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Sand))
	postTimestampStyle = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Text)).Faint(true)
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
	return CommentsHeader{DescriptionStyle: defaultDescriptionStyle}
}

func (h *CommentsHeader) SetSize(width, height int) {
	h.W = width - headerContainerStyle.GetHorizontalFrameSize()
	h.DescriptionStyle = h.DescriptionStyle.Width(h.W)
}

func (h CommentsHeader) View() string {
	titleView := titleStyle.Render(trim(h.Title, h.W))
	descriptionView := h.DescriptionStyle.Render(h.Description)

	authorView := postAuthorStyle.Render(h.Author)
	timestampView := postTimestampStyle.Render(fmt.Sprintf("submitted %s by", h.Timestamp))
	authorTimestampView := fmt.Sprintf("%s %s", timestampView, authorView)

	postPointsView := postPointsStyle.Render(fmt.Sprintf("%s points", h.Points))
	joinedView := lipgloss.JoinVertical(lipgloss.Left, titleView, descriptionView, authorTimestampView, postPointsView)

	return headerContainerStyle.Render(joinedView)
}

func (h *CommentsHeader) SetContent(comments client.Comments) {
	h.Title = utils.NormalizeSubreddit(comments.Subreddit)
	h.Description = comments.PostTitle
	h.Author = comments.PostAuthor
	h.Points = comments.PostPoints
	h.Timestamp = comments.PostTimestamp
}
