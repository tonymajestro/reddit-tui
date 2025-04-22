package components

import (
	"fmt"
	"log/slog"
	"reddittui/client"
	"reddittui/components/comments"
	"reddittui/components/messages"
	"reddittui/components/modal"
	"reddittui/components/posts"
	"reddittui/config"
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
	initCmd       tea.Cmd
}

func NewRedditTui(configuration config.Config, subreddit, post string) RedditTui {
	redditClient := client.NewRedditClient(configuration)

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
		initCmd:       getInitCmd(redditClient.BaseUrl, subreddit, post),
	}
}

func getInitCmd(baseUrl, subreddit, post string) tea.Cmd {
	if len(subreddit) != 0 {
		return messages.LoadSubreddit(subreddit)
	} else if len(post) != 0 {
		url, err := client.GetPostUrl(baseUrl, post)
		if err != nil {
			panic(fmt.Sprintf("Could not load post %s: %v", post, err))
		}

		return messages.LoadComments(url)
	} else {
		return messages.LoadHome
	}
}

func (r RedditTui) Init() tea.Cmd {
	return messages.LoadHome
}

func (r RedditTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case messages.ShowErrorModalMsg:
		if r.initializing && msg.OnClose == nil {
			slog.Error("Error during initialization")
			if r.loadingPage == HomePage {
				errorMsg := "Could not initialize reddittui. Check the logfile for details."
				return r, messages.ShowErrorModalWithCallback(errorMsg, tea.Quit)
			}

			var errorMsg string
			if r.loadingPage == SubredditPage {
				errorMsg = "Error loading subreddit. Returning to home page..."
			} else {
				errorMsg = "Error loading post. Returning to home page..."
			}

			return r, messages.ShowErrorModalWithCallback(errorMsg, messages.LoadHome)
		}

	case messages.OpenModalMsg:
		r.focusModal()
		return r, nil

	case messages.LoadingCompleteMsg:
		cmd = r.completeLoading()
		return r, cmd

	case messages.ExitModalMsg:
		r.popup = false
		r.focusActivePage()
		cmd = r.modalManager.Blur()
		return r, cmd

	case messages.GoBackMsg:
		r.goBack()
		return r, nil

	case messages.LoadHomeMsg:
		if r.page == HomePage && !r.initializing {
			return r, r.modalManager.Blur()
		}

		r.focusModal()
		r.loadingPage = HomePage

		cmd = r.modalManager.SetLoading(defaultLoadingMessage)
		cmds = append(cmds, cmd)

	case messages.LoadSubredditMsg:
		subreddit := string(msg)
		if r.page == SubredditPage && r.subredditPage.Subreddit == subreddit {
			return r, nil
		}

		r.focusModal()
		r.loadingPage = SubredditPage

		loadingMsg := fmt.Sprintf("loading %s...", utils.NormalizeSubreddit(subreddit))
		cmd = r.modalManager.SetLoading(loadingMsg)
		cmds = append(cmds, cmd)

	case messages.LoadMorePostsMsg:
		r.focusModal()

		cmd = r.modalManager.SetLoading("loading posts...")
		cmds = append(cmds, cmd)

	case messages.LoadCommentsMsg:
		r.focusModal()
		r.loadingPage = CommentsPage

		cmd = r.modalManager.SetLoading("loading comments...")
		cmds = append(cmds, cmd)

	case messages.OpenUrlMsg:
		url := string(msg)
		if err := utils.OpenUrl(url); err != nil {
			slog.Error("Error opening url in browser", "url", url, "error", err.Error())
			cmd = r.modalManager.SetError(fmt.Sprintf("Could not open url %s in browser", url))
			cmds = append(cmds, cmd)
		}

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

	return ""
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

func (r *RedditTui) completeLoading() tea.Cmd {
	initializing := r.initializing

	r.initializing = false
	r.popup = false
	r.setPage(r.loadingPage)
	r.focusActivePage()

	if initializing {
		r.initializing = false
		return r.initCmd
	}

	return r.modalManager.Blur()
}

func (r *RedditTui) focusModal() {
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
