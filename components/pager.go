package components

import (
	"fmt"
	"reddittui/client"
	"reddittui/components/colors"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var viewportStyle = lipgloss.NewStyle().Margin(0, 2, 1, 2)

var (
	commentAuthorStyle  = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Blue)).Bold(true)
	commentDateStyle    = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Lavender)).Italic(true)
	commentTextStyle    = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Text))
	popularPointsStyle  = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Green))
	defaultPointsStyle  = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Orange))
	negativePointsStyle = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Red))
	collapsedStyle      = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Yellow))
)

type viewportKeyMap struct {
	CursorUp         key.Binding
	CursorDown       key.Binding
	GoToStart        key.Binding
	GoToEnd          key.Binding
	CollapseComments key.Binding
	ShowFullHelp     key.Binding
	CloseFullHelp    key.Binding
	Quit             key.Binding
	ForceQuit        key.Binding
}

func (k viewportKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.CursorUp, k.CursorDown, k.CollapseComments, k.ShowFullHelp, k.Quit}
}

func (k viewportKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CursorUp, k.CursorDown, k.GoToStart, k.GoToEnd},
		{k.CollapseComments, k.Quit, k.CloseFullHelp},
	}
}

type CommentsViewport struct {
	viewport      viewport.Model
	postText      string
	comments      []client.Comment
	keyMap        viewportKeyMap
	help          help.Model
	collapsed     bool
	viewportLines []string
	w, h          int
}

func NewCommentsViewport() CommentsViewport {
	keys := viewportKeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		GoToStart: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GoToEnd: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
		CollapseComments: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "collapse comments"),
		),
		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	}

	return CommentsViewport{
		viewport:  viewport.New(0, 0),
		keyMap:    keys,
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
}

func (c *CommentsViewport) SetContent(comments client.Comments) {
	c.postText = comments.PostText
	c.comments = comments.Comments

	c.collapsed = false
	c.ResizeComponents()
	c.viewport.SetYOffset(0)

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
	}

	for i := range len(c.comments) - 1 {
		comment := c.comments[i]
		content.WriteString(c.formatComment(comment, i))
	}

	return content.String()
}

func (c *CommentsViewport) SetViewportContent() {
	content := c.GetViewportView()
	c.viewport.SetContent(content)
	c.viewportLines = strings.Split(content, "\n")
}

// Format comment, adding padding to the entry according to the comment's depth
func (c *CommentsViewport) formatComment(comment client.Comment, i int) string {
	var (
		authorAndDateView      string
		commentTextView        string
		pointsView             string
		pointsAndCollapsedView string
		tabWidth               = comment.Depth * 2
		commentWidth           = c.w - tabWidth
	)

	if c.collapsed && comment.Depth > 0 {
		return ""
	}

	authorView := commentAuthorStyle.Render(comment.Author)
	dateView := commentDateStyle.Render(comment.Timestamp)

	authorAndDateView = formatLine(fmt.Sprintf("%s • %s", authorView, dateView), commentWidth, comment.Depth)
	commentTextView = formatLine(commentTextStyle.Render(comment.Text), commentWidth, comment.Depth)

	padding := strings.Repeat("  ", comment.Depth)
	pointsView = renderPoints(comment.Points)
	pointsAndCollapsedView = fmt.Sprintf("%s%s", padding, pointsView)

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
			collapsedView := collapsedStyle.Render("(1 comment hidden)")
			pointsAndCollapsedView = fmt.Sprintf("%s%s  %s", padding, pointsView, collapsedView)
		} else if children > 1 {
			collapsedView := collapsedStyle.Render(fmt.Sprintf("(%d comments hidden)", children))
			pointsAndCollapsedView = fmt.Sprintf("%s%s  %s", padding, pointsView, collapsedView)
		}
	}

	return fmt.Sprintf("%s\n%s\n%s\n\n", authorAndDateView, commentTextView, pointsAndCollapsedView)
}

// Format string according to padding and width rules.
//
// If the string is longer than 'width', it will be split into multiple lines that are
// no more than 'width' wide. Each line will also have padded whitespace according to
// the 'depth' argument.
func formatLine(s string, width, depth int) string {
	var (
		lines   strings.Builder
		lineW   = depth
		padding = strings.Repeat("  ", depth)
	)

	lines.WriteString(padding)

	for _, word := range strings.Fields(s) {
		runes := []rune(word)

		if lineW+len(runes) > width {
			// Word doesn't fit on current line.
			// Add linebreak and padding and write word to next line

			if lineW > depth {
				lines.WriteRune('\n')
				lines.WriteString(padding)
				lines.WriteString(word)
				lineW = depth + len(runes)
			} else {
				// Edge case where first word on line doesn't fit
				// Hack: assume the word will fit on two lines
				// To-do: split the word into the correct number of lines

				left, right := runes[:width-depth], runes[width-depth:]

				lines.WriteString(string(left))
				lines.WriteString("-\n")
				lines.WriteString(padding)
				lines.WriteString(string(right))
				lines.WriteRune(' ')

				lineW = depth + len(right) + 1
			}
		} else {
			// Word fits on current line, write it to buffer
			lines.WriteString(word)
			lines.WriteRune(' ')
			lineW += len(runes) + 1
		}
	}

	return lines.String()
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
	pos, title, text := c.findAnchor()
	offset := pos - c.viewport.YOffset

	c.collapsed = !c.collapsed
	c.SetViewportContent()

	newPos := c.findComment(title, text)
	c.viewport.SetYOffset(newPos - offset)
}

// Find comment closest to the center of the screen to act as an anchor when toggling
// child comments.
func (c *CommentsViewport) findAnchor() (int, string, string) {
	// Don't use actual center of viewport since the header takes up some amount of space and
	// users probably look closer to the top of the screen rather than the bottom
	midPoint := c.viewport.YOffset + int(float64(c.viewport.Height)*0.4)

	findAnchorHelper := func(offset int) int {
		for i := midPoint; i >= 0 && i < len(c.viewportLines); i += offset {
			line := c.viewportLines[i]
			if len(line) > 0 && line[0] == ' ' {
				continue
			}

			split := strings.Split(line, "•")
			if len(split) == 2 {
				return i
			}
		}

		return -1
	}

	upPos := findAnchorHelper(-1)
	downPos := findAnchorHelper(1)

	if upPos >= 0 && downPos < 0 {
		return upPos, c.viewportLines[upPos], c.viewportLines[upPos+1]
	} else if upPos < 0 && downPos >= 0 {
		return downPos, c.viewportLines[downPos], c.viewportLines[downPos+1]
	}

	upDiff, downDiff := midPoint-upPos, downPos-midPoint
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
