package components

import (
	"log"
	"reddittui/client"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadCommentsMsg client.Post

type showCommentsMsg struct {
	title string
	items []list.Item
}

func ShowComments(post client.Post) tea.Cmd {
	return func() tea.Msg {
		return loadCommentsMsg(post)
	}
}

var listStyle = lipgloss.NewStyle().Margin(1, 2)

type CommentsPage struct {
	comments     []client.Comment
	itemsList    list.Model
	spinner      RedditSpinner
	redditClient client.RedditClient
	focus        bool
	w, h         int
}

func NewCommentsPage() CommentsPage {
	comments := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	comments.Title = ""
	comments.SetShowStatusBar(true)

	spinner := NewRedditSpinner()
	spinner.Focus()

	redditClient := client.New()

	return CommentsPage{
		itemsList:    comments,
		spinner:      spinner,
		redditClient: redditClient,
	}
}

func (c *CommentsPage) SetSize(w, h int) {
	c.w = w
	c.h = h
	c.itemsList.SetSize(w, h)
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
	return c.spinner.Init()
}

func (c CommentsPage) Update(msg tea.Msg) (CommentsPage, tea.Cmd) {
	var (
		spinnerCmd  tea.Cmd
		commentsCmd tea.Cmd
	)

	if !c.focus {
		return c, nil
	}

	switch msg := msg.(type) {

	case showCommentsMsg:
		c.HideLoading()
		c.itemsList.Title = msg.title
		c.itemsList.ResetSelected()
		return c, c.itemsList.SetItems(msg.items)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "backspace":
			return c, ReturnToPosts
		}
	}

	c.spinner, spinnerCmd = c.spinner.Update(msg)
	c.itemsList, commentsCmd = c.itemsList.Update(msg)

	return c, tea.Batch(spinnerCmd, commentsCmd)
}

func (c CommentsPage) View() string {
	if c.spinner.IsFocused() {
		return appStyle.Render(c.spinner.View())
	} else {
		return appStyle.Render(c.itemsList.View())
	}
}

func (c *CommentsPage) ShowLoading(message string) {
	c.spinner.SetLoadingMessage(message)
	c.spinner.Focus()
}

func (c *CommentsPage) HideLoading() {
	c.spinner.Blur()
}

func (c *CommentsPage) ShowComments() {
	c.Focus()
	c.spinner.Blur()
}

func (c CommentsPage) LoadComments(title, url string) tea.Cmd {
	return func() tea.Msg {
		comments, err := c.redditClient.GetComments(url)
		if err != nil {
			log.Printf("Error: %v", err)
			return err
		}

		items := getCommentListItems(comments)
		return showCommentsMsg{
			items: items,
			title: title,
		}
	}
}

func getCommentListItems(comments client.Comments) []list.Item {
	var items []list.Item
	for _, c := range comments {
		items = append(items, c)
	}
	return items
}
