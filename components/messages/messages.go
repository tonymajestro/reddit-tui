package messages

import (
	"reddittui/client"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	LoadCommentsMsg   client.Post
	UpdateCommentsMsg client.Comments
	LoadHomeMsg       struct{}
	UpdatePostsMsg    client.Posts
	GoBackMsg         struct{}
	GoHomeMsg         struct{}
	GoSubredditMsg    string
)

func GoBack() tea.Msg { return GoBackMsg{} }
func GoHome() tea.Msg { return GoHomeMsg{} }

func GoSubreddit(subreddit string) tea.Cmd {
	return func() tea.Msg {
		return GoSubredditMsg(subreddit)
	}
}
