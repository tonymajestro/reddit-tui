package components

import (
	tea "github.com/charmbracelet/bubbletea"
)

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
	} else {
		return r.commentsPage.View()
	}
}
