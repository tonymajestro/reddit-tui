package components

import (
	"fmt"
	"reddittui/client"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle().Margin(1, 2)

type RedditTui struct {
	postsPage    PostsPage
	commentsPage CommentsPage
}

func NewRedditTui() RedditTui {
	postsPage := NewPostsPage()
	postsPage.Focus()

	commentsPage := NewCommentsPage()

	return RedditTui{
		postsPage:    postsPage,
		commentsPage: commentsPage,
	}
}

func (r RedditTui) Init() tea.Cmd {
	return tea.Batch(r.postsPage.Init(), r.postsPage.LoadHome())
}

func (r RedditTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		postsCmd    tea.Cmd
		commentsCmd tea.Cmd
	)

	switch msg := msg.(type) {

	case returnToPostsMsg:
		r.postsPage.Focus()
		r.postsPage.maximizePostsList()
		r.postsPage.HideSearch()

	case loadCommentsMsg:
		post := client.Post(msg)
		r.postsPage.Blur()
		r.commentsPage.Focus()

		r.commentsPage.ShowLoading(fmt.Sprintf("loading comments for post %s...", post.PostTitle))
		return r, r.commentsPage.LoadComments(post.PostTitle, post.CommentsUrl)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			if !r.postsPage.subredditSearch.IsFocused() {
				return r, tea.Quit
			}
		case "q", "ctrl+c":
			return r, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		newW, newH := msg.Width-h, msg.Height-v
		r.postsPage.SetSize(newW, newH)
		r.commentsPage.SetSize(newW, newH)
	}

	r.postsPage, postsCmd = r.postsPage.Update(msg)

	r.commentsPage, commentsCmd = r.commentsPage.Update(msg)
	return r, tea.Batch(postsCmd, commentsCmd)
}

func (r RedditTui) View() string {
	if r.postsPage.IsFocused() {
		return r.postsPage.View()
	} else if r.commentsPage.IsFocused() {
		return r.commentsPage.View()
	}

	return ""
}
