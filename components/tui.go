package components

import (
	"fmt"
	"reddittui/client"
	"reddittui/components/comments"
	"reddittui/components/messages"
	"reddittui/components/modal"
	"reddittui/components/posts"
	"reddittui/utils"

	tea "github.com/charmbracelet/bubbletea"
)

const defaultLoadingMessage = "loading reddit.com..."

type (
	pageType int
)

const (
	HomePage pageType = iota
	SubredditPage
	CommentsPage
)

type RedditTui struct {
	homePage      posts.PostsPage
	subredditPage posts.PostsPage
	commentsPage  comments.CommentsPage
	modalManager  modal.ModalManager
	popup         bool
	initializing  bool
	page          pageType
	prevPage      pageType
	loadingPage   pageType
}

func NewRedditTui() RedditTui {
	redditClient := client.New()

	homePage := posts.NewPostsPage(redditClient, true)
	subredditPage := posts.NewPostsPage(redditClient, false)
	commentsPage := comments.NewCommentsPage(redditClient)

	modalManager := modal.NewModalManager()

	return RedditTui{
		homePage:      homePage,
		subredditPage: subredditPage,
		commentsPage:  commentsPage,
		modalManager:  modalManager,
		initializing:  true,
	}
}

func (r RedditTui) Init() tea.Cmd {
	return messages.LoadSubreddit("neovim")
}

func (r RedditTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case messages.OpenModalMsg:
		r.openModal()
		return r, nil

	case messages.LoadingCompleteMsg:
		r.popup = false
		r.setPage(r.loadingPage)
		r.focusActivePage()
		r.modalManager.Blur()
		return r, nil

	case messages.ExitModalMsg:
		r.popup = false
		r.focusActivePage()
		r.modalManager.Blur()
		return r, nil

	case messages.GoBackMsg:
		r.goBack()
		return r, nil

	case messages.LoadHomeMsg:
		if r.page == HomePage && !r.initializing {
			return r, nil
		}

		r.initializing = false
		r.loadingPage = HomePage
		cmds = append(cmds, messages.ShowSpinnerModal(defaultLoadingMessage))

	case messages.LoadSubredditMsg:
		subreddit := string(msg)
		if r.page == SubredditPage && r.subredditPage.Subreddit == subreddit {
			return r, nil
		}

		r.loadingPage = SubredditPage
		loadingMsg := fmt.Sprintf("loading %s...", utils.NormalizeSubreddit(subreddit))
		cmds = append(cmds, messages.ShowSpinnerModal(loadingMsg))

	case messages.LoadCommentsMsg:
		r.loadingPage = CommentsPage
		loadingMsg := "loading comments..."
		cmds = append(cmds, messages.ShowSpinnerModal(loadingMsg))

	case tea.WindowSizeMsg:
		r.homePage.SetSize(msg.Width, msg.Height)
		r.subredditPage.SetSize(msg.Width, msg.Height)
		r.commentsPage.SetSize(msg.Width, msg.Height)
		r.modalManager.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return r, tea.Quit
		}
	}

	r.modalManager, cmd = r.modalManager.Update(msg)
	cmds = append(cmds, cmd)

	r.homePage, cmd = r.homePage.Update(msg)
	cmds = append(cmds, cmd)

	r.subredditPage, cmd = r.subredditPage.Update(msg)
	cmds = append(cmds, cmd)

	r.commentsPage, cmd = r.commentsPage.Update(msg)
	cmds = append(cmds, cmd)

	return r, tea.Batch(cmds...)
}

func (r RedditTui) View() string {
	if r.popup {
		switch r.page {
		case HomePage:
			return r.modalManager.View(r.homePage)
		case SubredditPage:
			return r.modalManager.View(r.subredditPage)
		case CommentsPage:
			return r.modalManager.View(r.commentsPage)
		}
	}

	switch r.page {
	case HomePage:
		return r.homePage.View()
	case SubredditPage:
		return r.subredditPage.View()
	case CommentsPage:
		return r.commentsPage.View()
	}

	panic("Unexpected page")
}

func (r *RedditTui) goBack() {
	switch r.page {
	case CommentsPage:
		if r.prevPage == HomePage {
			r.setPage(HomePage)
		} else {
			r.setPage(SubredditPage)
		}
	default:
		r.setPage(HomePage)
	}

	r.focusActivePage()
}

func (r *RedditTui) setPage(page pageType) {
	r.page, r.prevPage = page, r.page
}

func (r *RedditTui) openModal() {
	r.popup = true
	r.homePage.Blur()
	r.subredditPage.Blur()
	r.commentsPage.Blur()
}

func (r *RedditTui) focusActivePage() {
	switch r.page {
	case HomePage:
		r.homePage.Focus()
		r.subredditPage.Blur()
		r.commentsPage.Blur()
	case SubredditPage:
		r.homePage.Blur()
		r.subredditPage.Focus()
		r.commentsPage.Blur()
	case CommentsPage:
		r.homePage.Blur()
		r.subredditPage.Blur()
		r.commentsPage.Focus()
	}
}
