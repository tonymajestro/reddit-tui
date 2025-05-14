package comments

import (
	"fmt"
	"reddittui/model"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CommentsViewport struct {
	viewport      viewport.Model
	postText      string
	postUrl       string
	comments      []model.Comment
	keyMap        viewportKeyMap
	help          help.Model
	collapsed     bool
	viewportLines []string
	w, h          int
}

func NewCommentsViewport() CommentsViewport {
	return CommentsViewport{
		viewport:  viewport.New(0, 0),
		keyMap:    commentsKeys,
		help:      help.New(),
		collapsed: false,
	}
}

func (c CommentsViewport) Update(msg tea.Msg) (CommentsViewport, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.keyMap.GoToStart):
			c.viewport.GotoTop()
		case key.Matches(msg, c.keyMap.GoToEnd):
			c.viewport.GotoBottom()
		case key.Matches(msg, c.keyMap.CollapseComments):
			c.toggleCollapseComments()
		case key.Matches(msg, c.keyMap.ShowFullHelp),
			key.Matches(msg, c.keyMap.CloseFullHelp):
			c.help.ShowAll = !c.help.ShowAll
		}
	}

	var cmd tea.Cmd
	c.viewport, cmd = c.viewport.Update(msg)
	return c, cmd
}

func (c CommentsViewport) View() string {
	viewportView := viewportStyle.Render(c.viewport.View())
	helpView := c.help.View(c.keyMap)
	return lipgloss.JoinVertical(lipgloss.Left, viewportView, helpView)
}

func (c *CommentsViewport) SetSize(w, h int) {
	c.w = w - viewportStyle.GetHorizontalFrameSize()
	c.h = h

	c.ResizeComponents()
	c.SetViewportContent()
}

func (c *CommentsViewport) SetContent(comments model.Comments) {
	c.postText = comments.PostText
	c.postUrl = comments.PostUrl
	c.comments = comments.Comments

	c.collapsed = false
	c.viewport.SetYOffset(0)
	c.ResizeComponents()
	c.SetViewportContent()
}

func (c *CommentsViewport) ResizeComponents() {
	helpHeight := lipgloss.Height(c.help.View(c.keyMap))

	c.viewport.Width = c.w
	c.viewport.Height = c.h - helpHeight - 1
}

func (c *CommentsViewport) GetViewportView() string {
	var content strings.Builder

	if len(c.postText) > 0 {
		content.WriteString(c.postText)
		content.WriteString("\n")
	} else {
		content.WriteString(c.postUrl)
		content.WriteString("\n\n")
	}

	for i := range len(c.comments) {
		comment := c.comments[i]
		commentView := c.formatComment(comment, i)
		if len(commentView) > 0 {
			content.WriteString(commentView)
			content.WriteString("\n\n")
		}
	}

	return content.String()
}

func (c *CommentsViewport) SetViewportContent() {
	content := c.GetViewportView()
	c.viewport.SetContent(content)
	c.viewportLines = strings.Split(content, "\n")
}

// Format comment, adding padding to the entry according to the comment's depth
func (c *CommentsViewport) formatComment(comment model.Comment, i int) string {
	var (
		authorAndDateView          string
		pointsView                 string
		pointsAndCollapsedHintView string
		paddingW                   = comment.Depth * 2
		containerStyle             = lipgloss.NewStyle().PaddingLeft(paddingW).Width(c.w - paddingW)
	)

	if c.collapsed && comment.Depth > 0 {
		return ""
	}

	authorView := commentAuthorStyle.Render(comment.Author)
	dateView := commentDateStyle.Render(comment.Timestamp)
	authorAndDateView = fmt.Sprintf("%s • %s", authorView, dateView)
	pointsView = renderPoints(comment.Points)
	pointsAndCollapsedHintView = pointsView

	if c.collapsed {
		children := 0
		for j := i + 1; j < len(c.comments); j++ {
			nextComment := c.comments[j]
			if nextComment.Depth == 0 {
				break
			}
			children++
		}

		if children == 1 {
			collapsedHintView := collapsedStyle.Render("(1 comment hidden)")
			pointsAndCollapsedHintView = fmt.Sprintf("%s  %s", pointsView, collapsedHintView)
		} else if children > 1 {
			collapsedView := collapsedStyle.Render(fmt.Sprintf("(%d comments hidden)", children))
			pointsAndCollapsedHintView = fmt.Sprintf("%s  %s", pointsView, collapsedView)
		}
	}

	joined := lipgloss.JoinVertical(lipgloss.Left, authorAndDateView, comment.Text, pointsAndCollapsedHintView)
	return containerStyle.Render(joined)
}

func renderPoints(pointsString string) string {
	parts := strings.Fields(pointsString)
	if len(parts) != 2 {
		return defaultPointsStyle.Render(pointsString)
	}

	if strings.Contains(parts[0], "-") {
		return negativePointsStyle.Render(pointsString)
	} else if strings.Contains(parts[0], "k") {
		return popularPointsStyle.Render(pointsString)
	}

	points, err := strconv.Atoi(parts[0])
	if err != nil {
		return defaultPointsStyle.Render(pointsString)
	} else if points >= 1000 {
		return popularPointsStyle.Render(pointsString)
	}

	return defaultPointsStyle.Render(pointsString)
}

func (c *CommentsViewport) toggleCollapseComments() {
	pos, title, text := c.findAnchorComment()
	if pos < 0 {
		return
	}

	offset := pos - c.viewport.YOffset

	c.collapsed = !c.collapsed
	c.SetViewportContent()

	newPos := c.findComment(title, text)
	c.viewport.SetYOffset(newPos - offset)
}

// Find comment closest to the center of the screen to act as an anchor when toggling
// child comments.
func (c *CommentsViewport) findAnchorComment() (pos int, title string, text string) {
	findAnchorHelper := func(start, offset int) int {
		for i := start; i >= 0 && i < len(c.viewportLines); i += offset {
			line := c.viewportLines[i]
			if len(line) > 0 && line[0] == ' ' {
				continue
			}

			split := strings.Split(line, "•")
			if len(split) == 2 && strings.Contains(split[1], "ago") {
				return i
			}
		}

		return -1
	}

	// Don't use actual center of viewport since the header takes up some amount of space and
	// users probably look closer to the top of the screen rather than the bottom
	searchStart := c.viewport.YOffset + int(float64(c.viewport.Height)*0.4)

	if searchStart >= len(c.viewportLines) {
		searchStart = 0
	}

	// Look for the comment above and below the center of the screen. Calculate which comment is closer to
	// the center of the screen
	upPos := findAnchorHelper(searchStart, -1)
	downPos := findAnchorHelper(searchStart, 1)

	if upPos < 0 && downPos < 0 {
		return -1, "", ""
	} else if upPos >= 0 && downPos < 0 {
		return upPos, c.viewportLines[upPos], c.viewportLines[upPos+1]
	} else if upPos < 0 && downPos >= 0 {
		return downPos, c.viewportLines[downPos], c.viewportLines[downPos+1]
	}

	upDiff, downDiff := searchStart-upPos, downPos-searchStart
	if upDiff < downDiff {
		return upPos, c.viewportLines[upPos], c.viewportLines[upPos+1]
	}
	return downPos, c.viewportLines[downPos], c.viewportLines[downPos+1]
}

func (c *CommentsViewport) findComment(title, text string) int {
	for i := range len(c.viewportLines) - 1 {
		currTitle := c.viewportLines[i]
		currText := c.viewportLines[i+1]

		if currTitle == title && currText == text {
			return i
		}
	}

	return -1
}
