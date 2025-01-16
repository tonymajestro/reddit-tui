package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle().Margin(1, 2)

type returnToPostsMsg struct{}

func ReturnToPosts() tea.Msg {
	return returnToPostsMsg{}
}

type RedditTui struct {
	postsPage    PostsPage
	commentsPage CommentsPage
}

func NewRedditTui() RedditTui {
	postsPage := NewPostsPage()
	commentsPage := NewCommentsPage()

	postsPage.Focus()
	commentsPage.Blur()

	return RedditTui{postsPage, commentsPage}
}

func (r RedditTui) Init() tea.Cmd {
	return r.postsPage.Init()
}

func (r RedditTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case returnToPostsMsg:
		r.commentsPage.Blur()
		r.postsPage.Focus()

	case loadCommentsMsg:
		r.commentsPage.Focus()
		r.postsPage.Blur()

		cmd := r.commentsPage.LoadComments(msg.CommentsUrl, msg.PostTitle)
		return r, cmd

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc", "q":
			if !r.commentsPage.IsFocused() && !r.postsPage.search.Searching {
				return r, tea.Quit
			}
		case "ctrl+c":
			return r, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		newW, newH := msg.Width-h-2, msg.Height-v
		r.postsPage.SetSize(newW, newH)
		r.commentsPage.SetSize(newW, newH)
	}

	var cmd tea.Cmd
	if r.postsPage.IsFocused() {
		r.postsPage, cmd = r.postsPage.Update(msg)
		return r, cmd
	} else {
		r.commentsPage, cmd = r.commentsPage.Update(msg)
		return r, cmd
	}
}

func (r RedditTui) View() string {
	if r.postsPage.IsFocused() {
		return appStyle.Render(r.postsPage.View())
	} else {
		return appStyle.Render(r.commentsPage.View())
	}
}
