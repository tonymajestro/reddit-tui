package components

import (
	"fmt"
	"reddittui/client"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	viewportStyle = lipgloss.NewStyle().Margin(0, 4, 1, 4)
	authorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	dateStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	pointsStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
)

type viewportKeyMap struct {
	CursorUp      key.Binding
	CursorDown    key.Binding
	GoToStart     key.Binding
	GoToEnd       key.Binding
	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding
	Quit          key.Binding
	ForceQuit     key.Binding
}

func (k viewportKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.CursorUp, k.CursorDown, k.ShowFullHelp, k.Quit}
}

func (k viewportKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CursorUp, k.CursorDown, k.GoToStart, k.GoToEnd},
		{k.Quit, k.CloseFullHelp},
	}
}

type CommentsViewport struct {
	viewport viewport.Model
	postText string
	comments []client.Comment
	keyMap   viewportKeyMap
	help     help.Model
	w, h     int
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
		// Toggle help.
		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),

		// Quitting.
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	}

	return CommentsViewport{
		viewport: viewport.New(0, 0),
		keyMap:   keys,
		help:     help.New(),
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

func (c *CommentsViewport) SetContent(postText string, comments []client.Comment) {
	c.postText = postText
	c.comments = comments
	c.viewport.SetYOffset(0)
	c.ResizeComponents()
}

func (c *CommentsViewport) ResizeComponents() {
	var (
		content    strings.Builder
		helpHeight = lipgloss.Height(c.help.View(c.keyMap))
	)

	c.viewport.Width = c.w
	c.viewport.Height = c.h - helpHeight - 1

	content.WriteString(c.postText)
	content.WriteString("\n")

	for _, comment := range c.comments {
		content.WriteString(c.formatComment(comment))
	}

	c.viewport.SetContent(content.String())
}

// Format comment, adding padding to the entry according to the comment's depth
func (c *CommentsViewport) formatComment(comment client.Comment) string {
	var (
		tabWidth     = comment.Depth * 2
		commentWidth = c.w - tabWidth
	)

	authorView := authorStyle.Render(comment.Author)
	dateView := dateStyle.Render(comment.Timestamp)
	authorAndPointsView := formatLine(fmt.Sprintf("%s • %s", authorView, dateView), commentWidth, comment.Depth)

	commentTextView := formatLine(comment.Text, commentWidth, comment.Depth)
	pointsView := formatLine(pointsStyle.Render(comment.Points), commentWidth, comment.Depth)

	return fmt.Sprintf("%s\n%s\n%s\n\n", authorAndPointsView, commentTextView, pointsView)
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
				//
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
