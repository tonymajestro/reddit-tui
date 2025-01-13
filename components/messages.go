package components

import (
	"reddittui/client"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type loadCommentsMsg struct {
	post client.Post
}

type showCommentsMsg struct {
	title string
	items []list.Item
}

type returnToPostsMsg struct{}

type acceptSearchMsg struct {
	subreddit string
}

type showPostsMsg struct {
	posts []client.Post
	title string
	items []list.Item
}

func ShowComments(post client.Post) tea.Cmd {
	return func() tea.Msg {
		return loadCommentsMsg{post}
	}
}

func ReturnToPosts() tea.Msg {
	return returnToPostsMsg{}
}

func AcceptSearch(subreddit string) tea.Cmd {
	return func() tea.Msg {
		return acceptSearchMsg{subreddit}
	}
}
