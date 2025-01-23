package components

import (
	"log"
	"reddittui/client"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadCommentsMsg client.Post

type updateCommentsMsg client.Comments

type CommentsPage struct {
	redditClient client.RedditClient
	header       Header
	pager        CommentsViewport
	spinner      Spinner
	focus        bool
	w, h         int
}

func NewCommentsPage() CommentsPage {
	redditClient := client.New()
	header := NewHeader()
	vp := NewCommentsViewport()

	return CommentsPage{
		redditClient: redditClient,
		header:       header,
		pager:        vp,
	}
}

func (c *CommentsPage) SetSize(w, h int) {
	c.w = w
	c.h = h

	c.ResizeComponents()
}

func (c *CommentsPage) ResizeComponents() {
	var (
		headerHeight = lipgloss.Height(c.header.View())
		pagerHeight  = c.h - headerHeight
	)

	c.header.SetSize(c.w, c.h)
	c.pager.SetSize(c.w, pagerHeight)
}

func (c *CommentsPage) IsFocused() bool {
	return c.focus
}

func (c *CommentsPage) Focus() {
	c.focus = true
}

func (c *CommentsPage) Blur() {
	c.focus = false
}

func (c CommentsPage) Init() tea.Cmd {
	return c.spinner.Tick
}

func (c CommentsPage) Update(msg tea.Msg) (CommentsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case updateCommentsMsg:
		c.UpdateComments(client.Comments(msg))
		return c, nil
	}

	var cmd tea.Cmd
	if c.spinner.Loading {
		c.spinner, cmd = c.spinner.Update(msg)
		return c, cmd
	} else {
		c.pager, cmd = c.pager.Update(msg)
		return c, cmd
	}
}

func (c CommentsPage) View() string {
	if c.spinner.Loading {
		return c.spinner.View()
	}

	headerView := c.header.View()
	pagerView := c.pager.View()

	return lipgloss.JoinVertical(lipgloss.Left, headerView, pagerView)
}

func (c *CommentsPage) LoadComments(url, title string) tea.Cmd {
	c.spinner.SetLoading(true)
	c.spinner.LoadingMessage = "loading comments..."

	loadCommentsCmd := func() tea.Msg {
		comments, err := c.redditClient.GetComments(url)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		return updateCommentsMsg(comments)
	}

	return tea.Batch(loadCommentsCmd, c.spinner.Tick)
}

func (c *CommentsPage) UpdateComments(comments client.Comments) {
	c.spinner.SetLoading(false)

	c.header.SetTitle(normalizeSubreddit(comments.Subreddit))
	c.header.SetDescription(comments.PostTitle)

	c.pager.SetContent(comments.Text, comments.Comments)

	// Need to resize components when content loads so padding and margins are correct
	c.ResizeComponents()
}
