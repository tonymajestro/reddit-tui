package comments

import (
	"log/slog"
	"reddittui/client"
	"reddittui/components/messages"
	"reddittui/components/styles"
	"reddittui/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var commentsErrorText = "Could not load comments. Please try again in a few moments."

type CommentsPage struct {
	redditClient   client.RedditClient
	header         CommentsHeader
	pager          CommentsViewport
	containerStyle lipgloss.Style
	postUrl        string
	focus          bool
}

func NewCommentsPage(redditClient client.RedditClient) CommentsPage {
	header := NewCommentsHeader()
	vp := NewCommentsViewport()

	return CommentsPage{
		redditClient:   redditClient,
		header:         header,
		pager:          vp,
		containerStyle: styles.GlobalStyle,
	}
}

func (c CommentsPage) Init() tea.Cmd {
	return nil
}

func (c CommentsPage) Update(msg tea.Msg) (CommentsPage, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if c.focus {
		c, cmd = c.handleFocusedMessages(msg)
		cmds = append(cmds, cmd)
	}

	c, cmd = c.handleGlobalMessages(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c CommentsPage) handleGlobalMessages(msg tea.Msg) (CommentsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.LoadCommentsMsg:
		url := string(msg)
		return c, c.loadComments(url)
	case messages.UpdateCommentsMsg:
		c.updateComments(model.Comments(msg))
		return c, messages.LoadingComplete
	}

	return c, nil
}

func (c CommentsPage) handleFocusedMessages(msg tea.Msg) (CommentsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "H":
			return c, messages.LoadHome

		case "B", "escape", "backspace", "left", "h":
			return c, messages.GoBack

		case "o", "O":
			return c, messages.OpenUrl(c.postUrl)
		}
	}

	var cmd tea.Cmd
	c.pager, cmd = c.pager.Update(msg)
	return c, cmd
}

func (c CommentsPage) View() string {
	headerView := c.header.View()
	pagerView := c.pager.View()
	joined := lipgloss.JoinVertical(lipgloss.Center, headerView, pagerView)
	return c.containerStyle.Render(joined)
}

func (c *CommentsPage) SetSize(w, h int) {
	c.containerStyle = c.containerStyle.Width(w).Height(h)
	c.resizeComponents()
}

func (c *CommentsPage) Focus() {
	c.focus = true
}

func (c *CommentsPage) Blur() {
	c.focus = false
}

func (c *CommentsPage) resizeComponents() {
	var (
		w            = c.containerStyle.GetWidth() - c.containerStyle.GetHorizontalFrameSize()
		h            = c.containerStyle.GetHeight() - c.containerStyle.GetVerticalFrameSize()
		headerHeight = lipgloss.Height(c.header.View())
		pagerHeight  = h - headerHeight
	)

	c.header.SetSize(w, h)
	c.pager.SetSize(w, pagerHeight)
}

func (c *CommentsPage) loadComments(url string) tea.Cmd {
	return func() tea.Msg {
		comments, err := c.redditClient.GetComments(url)
		if err != nil {
			slog.Error(commentsErrorText, "error", err)
			return messages.ShowErrorModalMsg{ErrorMsg: commentsErrorText}
		}

		return messages.UpdateCommentsMsg(comments)
	}
}

func (c *CommentsPage) updateComments(comments model.Comments) {
	c.header.SetContent(comments)
	c.pager.SetContent(comments)
	c.postUrl = comments.PostUrl

	// Need to resize components when content loads so padding and margins are correct
	c.resizeComponents()
}
