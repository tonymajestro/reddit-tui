package components

import (
	"reddittui/client"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultListTitle = "reddit.com"

var listStyle = lipgloss.NewStyle().Margin(1, 2)

type CommentsPage struct {
	redditClient client.RedditClient
	focus        bool
	w, h         int
}

func NewCommentsPage() CommentsPage {
	return CommentsPage{}
}

func (c *CommentsPage) SetSize(w, h int) {
	c.w = w
	c.h = h
}

func (c CommentsPage) IsFocused() bool {
	return c.focus
}

func (c CommentsPage) Init() tea.Cmd {
	return nil
}

func (c CommentsPage) Update(msg tea.Msg) (CommentsPage, tea.Cmd) {
	return c, nil
}

func (c CommentsPage) View() string {
	return "comments"
}
