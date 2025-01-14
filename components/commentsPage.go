package components

import (
	"fmt"
	"log"
	"reddittui/client"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type loadCommentsMsg struct {
	post      client.Post
	subreddit string
}

type showCommentsMsg struct {
	comments  []client.Comment
	title     string
	subreddit string
}

type CommentsPage struct {
	comments       []client.Comment
	listModel      list.Model
	spinnerModel   spinner.Model
	redditClient   client.RedditClient
	loading        bool
	searching      bool
	loadingMessage string
	w, h           int
	focus          bool
}

func NewCommentsPage() CommentsPage {
	items := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	items.SetStatusBarItemName("comment", "comments")

	redditClient := client.New()

	return CommentsPage{
		comments:     []client.Comment{},
		listModel:    items,
		redditClient: redditClient,
	}
}

func (c *CommentsPage) SetSize(w, h int) {
	c.w = w
	c.h = h
	c.listModel.SetSize(w, h)
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
	return c.spinnerModel.Tick
}

func (c CommentsPage) Update(msg tea.Msg) (CommentsPage, tea.Cmd) {
	switch msg := msg.(type) {

	case showCommentsMsg:
		cmd := c.DisplayComments(msg.comments, msg.title, msg.subreddit)
		return c, cmd

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "backspace":
			return c, ReturnToPosts
		}
	}

	var cmd tea.Cmd

	if c.loading {
		c.spinnerModel, cmd = c.spinnerModel.Update(msg)
		return c, cmd
	} else {
		c.listModel, cmd = c.listModel.Update(msg)
		return c, cmd
	}
}

func (c CommentsPage) View() string {
	if c.loading {
		return appStyle.Render(c.GetSpinnerView())
	} else {
		return appStyle.Render(c.listModel.View())
	}
}

func (c *CommentsPage) ShowLoading(message string) {
	spinnerModel := spinner.New()
	spinnerModel.Spinner = spinner.Dot
	spinnerModel.Style = spinnerStyle

	c.spinnerModel = spinnerModel
	c.loadingMessage = message
	c.loading = true
}

func (c *CommentsPage) HideLoading() {
	c.loading = false
}

func (c CommentsPage) GetSpinnerView() string {
	return fmt.Sprintf("%s %s", c.spinnerModel.View(), c.loadingMessage)
}

func (c *CommentsPage) LoadComments(url, title, subreddit string) tea.Cmd {
	c.ShowLoading(fmt.Sprintf("loading comments for \"%s\"...", title))

	loadCommentsCmd := func() tea.Msg {
		comments, err := c.redditClient.GetComments(url)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		return showCommentsMsg{
			comments:  comments,
			title:     title,
			subreddit: subreddit,
		}
	}

	return tea.Batch(loadCommentsCmd, c.spinnerModel.Tick)
}

func (c *CommentsPage) DisplayComments(comments []client.Comment, title, subreddit string) tea.Cmd {
	c.HideLoading()

	c.listModel.Title = fmt.Sprintf("%s  |  %s", subreddit, title)

	c.comments = comments
	c.listModel.ResetSelected()
	return c.listModel.SetItems(getCommentListItems(comments))
}

func getCommentListItems(comments client.Comments) []list.Item {
	var items []list.Item
	for _, c := range comments {
		items = append(items, c)
	}
	return items
}
