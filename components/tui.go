package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle().Margin(1, 2)

type (
	goBackMsg      struct{}
	goHomeMsg      struct{}
	goSubredditMsg string
)

func GoBack() tea.Msg { return goBackMsg{} }
func GoHome() tea.Msg { return goHomeMsg{} }

func GoSubreddit(subreddit string) tea.Cmd {
	return func() tea.Msg {
		return goSubredditMsg(subreddit)
	}
}

type RedditTui struct {
	postsPage    PostsPage
	commentsPage CommentsPage
	quitPage     QuitPage
	focusStack   FocusStack
}

func NewRedditTui() RedditTui {
	postsPage := NewPostsPage()
	commentsPage := NewCommentsPage()
	quitPage := NewQuitPage()

	postsPage.Focus()
	commentsPage.Blur()

	focusStack := FocusStack{Home}

	return RedditTui{
		postsPage:    postsPage,
		commentsPage: commentsPage,
		quitPage:     quitPage,
		focusStack:   focusStack,
	}
}

func (r RedditTui) Init() tea.Cmd {
	return r.postsPage.Init()
}

func (r RedditTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case goBackMsg:
		if len(r.focusStack) == 1 {
			return r, nil
		}

		prevPage := r.focusStack.Pop()

		switch r.focusStack.Peek() {
		case Home:
			r.Focus(Home)

			if prevPage == Quit {
				return r, nil
			} else {
				return r, GoHome
			}

		case Subreddit:
			r.Focus(Subreddit)
			return r, nil

		case Comments:
			r.Focus(Comments)
			return r, nil

		default:
			return r, tea.Quit
		}

	case loadCommentsMsg:
		r.Focus(Comments)
		r.focusStack.Push(Comments)

		return r, r.commentsPage.LoadComments(msg.CommentsUrl, msg.PostTitle)

	case goHomeMsg:
		r.Focus(Home)
		r.focusStack.Clear()
		r.focusStack.Push(Home)
		return r, r.postsPage.LoadHome()

	case goSubredditMsg:
		r.Focus(Subreddit)
		r.focusStack.Push(Subreddit)
		return r, r.postsPage.LoadSubreddit(string(msg))

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			if !r.quitPage.IsFocused() && !r.postsPage.IsSearching() {
				r.PromptQuit()
				return r, nil
			}
		case "ctrl+c":
			return r, tea.Quit

		case "h":
			if !r.quitPage.IsFocused() && !r.postsPage.IsSearching() {
				return r, GoHome
			}

		case "backspace":
			if r.CanGoBack() {
				return r, GoBack
			}
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
	} else if r.commentsPage.IsFocused() {
		r.commentsPage, cmd = r.commentsPage.Update(msg)
		return r, cmd
	} else {
		r.quitPage, cmd = r.quitPage.Update(msg)
		return r, cmd
	}
}

func (r RedditTui) View() string {
	if r.postsPage.IsFocused() {
		return appStyle.Render(r.postsPage.View())
	} else if r.commentsPage.IsFocused() {
		return appStyle.Render(r.commentsPage.View())
	} else {
		return appStyle.Render(r.quitPage.View())
	}
}

func (r *RedditTui) Focus(page PageType) {
	r.quitPage.Blur()
	r.postsPage.Blur()
	r.commentsPage.Blur()

	switch page {
	case Home, Subreddit:
		r.postsPage.Focus()
	case Comments:
		r.commentsPage.Focus()
	case Quit:
		r.quitPage.Focus()
	}
}

func (r *RedditTui) PromptQuit() {
	r.Focus(Quit)
	r.focusStack.Push(Quit)
}

func (r RedditTui) CanGoBack() bool {
	return !r.quitPage.IsFocused() && (!r.postsPage.IsFocused() || !r.postsPage.IsSearching())
}
