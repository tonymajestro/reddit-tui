package components

import (
	"log"
	"reddittui/client"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadCommentsMsg client.Post

type updateCommentsMsg client.Comments

type CommentsPage struct {
	redditClient client.RedditClient
	comments     []client.Comment
	header       Header
	postText     PostText
	list         list.Model
	spinner      Spinner
	focus        bool
	w, h         int
}

func NewCommentsPage() CommentsPage {
	items := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	items.SetShowTitle(false)
	items.SetShowStatusBar(false)

	redditClient := client.New()
	header := NewHeader()
	postText := NewPostText()

	return CommentsPage{
		comments:     []client.Comment{},
		list:         items,
		redditClient: redditClient,
		header:       header,
		postText:     postText,
	}
}

func (c *CommentsPage) SetSize(w, h int) {
	c.w = w
	c.h = h

	c.ResizeComponents()
}

func (c *CommentsPage) ResizeComponents() {
	headerHeight := lipgloss.Height(c.header.View())
	postTextViewHeight := lipgloss.Height(c.postText.View())
	listHeight := c.h - headerHeight - postTextViewHeight

	c.header.SetSize(c.w, c.h)
	c.postText.SetSize(c.w, c.h)
	c.list.SetSize(c.w, listHeight)
}

func (c CommentsPage) IsFocused() bool {
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
		cmd := c.UpdateComments(client.Comments(msg))
		return c, cmd

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "backspace":
			return c, ReturnToPosts
		}
	}

	var cmd tea.Cmd

	if c.spinner.Loading {
		c.spinner, cmd = c.spinner.Update(msg)
		return c, cmd
	} else {
		c.list, cmd = c.list.Update(msg)
		return c, cmd
	}
}

func (c CommentsPage) View() string {
	if c.spinner.Loading {
		return c.spinner.View()
	}

	headerView := c.header.View()
	postTextView := c.postText.View()
	listView := c.list.View()

	return lipgloss.JoinVertical(lipgloss.Left, headerView, postTextView, listView)
}

func (c *CommentsPage) LoadComments(url, title string) tea.Cmd {
	c.spinner.SetLoading(true)
	c.spinner.SetLoadingMessage("loading comments...")

	loadCommentsCmd := func() tea.Msg {
		comments, err := c.redditClient.GetComments(url)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		return updateCommentsMsg(comments)
	}

	return tea.Batch(loadCommentsCmd, c.spinner.Tick)
}

func (c *CommentsPage) UpdateComments(comments client.Comments) tea.Cmd {
	c.spinner.SetLoading(false)

	c.header.SetTitle(normalizeSubreddit(comments.Subreddit))
	c.header.SetDescription(comments.PostTitle)

	c.postText.Contents = comments.Text
	c.comments = comments.Comments
	c.list.ResetSelected()

	// Need to resize components when content loads so padding and margins are correct
	c.ResizeComponents()

	var listItems []list.Item
	for _, c := range comments.Comments {
		listItems = append(listItems, c)
	}

	return c.list.SetItems(listItems)
}
