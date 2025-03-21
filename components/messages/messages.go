package messages

import (
	"reddittui/model"

	tea "github.com/charmbracelet/bubbletea"
)

type ErrorModalMsg struct {
	ErrorMsg string
	OnClose  tea.Cmd
}

type (
	InitMsg            struct{}
	GoBackMsg          struct{}
	LoadCommentsMsg    string
	LoadHomeMsg        struct{}
	LoadMorePostsMsg   bool
	LoadSubredditMsg   string
	UpdateCommentsMsg  model.Comments
	UpdatePostsMsg     model.Posts
	AddMorePostsMsg    model.Posts
	LoadingCompleteMsg struct{}

	OpenModalMsg        struct{}
	ExitModalMsg        struct{}
	ShowSpinnerModalMsg string

	ShowErrorModalMsg ErrorModalMsg

	OpenUrlMsg string
)

func Init() tea.Msg {
	return InitMsg{}
}

func GoBack() tea.Msg {
	return GoBackMsg{}
}

func LoadHome() tea.Msg {
	return LoadHomeMsg{}
}

func LoadMorePosts(home bool) tea.Cmd {
	return func() tea.Msg {
		return LoadMorePostsMsg(home)
	}
}

func LoadSubreddit(subreddit string) tea.Cmd {
	return func() tea.Msg {
		return LoadSubredditMsg(subreddit)
	}
}

func LoadComments(url string) tea.Cmd {
	return func() tea.Msg {
		return LoadCommentsMsg(url)
	}
}

func LoadingComplete() tea.Msg {
	return LoadingCompleteMsg{}
}

func OpenModal() tea.Msg {
	return OpenModalMsg{}
}

func ExitModal() tea.Msg {
	return ExitModalMsg{}
}

func ShowSpinnerModal(loadingMsg string) tea.Cmd {
	return func() tea.Msg {
		return ShowSpinnerModalMsg(loadingMsg)
	}
}

func ShowErrorModal(errorMsg string) tea.Cmd {
	return func() tea.Msg {
		return ShowErrorModalMsg{ErrorMsg: errorMsg}
	}
}

func ShowErrorModalWithCallback(errorMsg string, callback tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return ShowErrorModalMsg{ErrorMsg: errorMsg, OnClose: callback}
	}
}

func HideSpinnerModal() tea.Msg {
	return ExitModalMsg{}
}

func OpenUrl(url string) tea.Cmd {
	return func() tea.Msg {
		return OpenUrlMsg(url)
	}
}
