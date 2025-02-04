package messages

import (
	"reddittui/client"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	InitMsg            struct{}
	GoBackMsg          struct{}
	LoadCommentsMsg    client.Post
	LoadHomeMsg        struct{}
	LoadSubredditMsg   string
	UpdateCommentsMsg  client.Comments
	UpdatePostsMsg     client.Posts
	LoadingCompleteMsg struct{}

	OpenModalMsg        struct{}
	ExitModalMsg        struct{}
	ShowSpinnerModalMsg string

	ErrorMsg          string
	ShowErrorModalMsg string

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

func LoadSubreddit(subreddit string) tea.Cmd {
	return func() tea.Msg {
		return LoadSubredditMsg(subreddit)
	}
}

func LoadComments(post client.Post) tea.Cmd {
	return func() tea.Msg {
		return LoadCommentsMsg(post)
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

func Error(errorMsg string) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg(errorMsg)
	}
}

func ShowErrorModal(errorMsg string) tea.Cmd {
	return func() tea.Msg {
		return ShowErrorModalMsg(errorMsg)
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
