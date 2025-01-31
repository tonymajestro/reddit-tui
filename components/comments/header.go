package comments

import (
	"fmt"
	"reddittui/client"
	"reddittui/components/colors"
	"reddittui/utils"
	"strconv"

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

type CommentsHeader struct {
	DescriptionStyle lipgloss.Style
	Title            string
	Description      string
	Author           string
	Timestamp        string
	Points           int
	TotalComments    int
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
	titleView := titleStyle.Render(utils.TruncateString(h.Title, h.W))
	descriptionView := h.DescriptionStyle.Render(h.Description)

	authorView := postAuthorStyle.Render(h.Author)
	timestampView := postTimestampStyle.Render(fmt.Sprintf("submitted %s by", h.Timestamp))
	authorTimestampView := fmt.Sprintf("%s %s", timestampView, authorView)

	postPointsView := postPointsStyle.Render(utils.GetSingularPlural(h.Points, "point", "points"))
	totalCommentsView := totalCommentsStyle.Render(utils.GetSingularPlural(h.TotalComments, "comment", "comments"))
	pointsAndCommentsView := fmt.Sprintf("%s • %s", postPointsView, totalCommentsView)

	joinedView := lipgloss.JoinVertical(lipgloss.Left, titleView, descriptionView, authorTimestampView, pointsAndCommentsView)

	return headerContainerStyle.Render(joinedView)
}

func (h *CommentsHeader) SetContent(comments client.Comments) {
	h.Title = utils.NormalizeSubreddit(comments.Subreddit)
	h.Description = comments.PostTitle
	h.Author = comments.PostAuthor
	h.TotalComments = len(comments.Comments)
	h.Timestamp = comments.PostTimestamp

	if points, err := strconv.Atoi(comments.PostPoints); err != nil {
		h.Points = 0
	} else {
		h.Points = points
	}
}
