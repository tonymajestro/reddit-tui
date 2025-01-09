package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultListTitle = "reddit.com"

var listStyle = lipgloss.NewStyle().Margin(1, 2)

type (
	showPostsListMsg struct {
		items []list.Item
		title string
	}

	hidePostsListMsg struct{}
)

func showPostsList() tea.Msg {
	return showPostsListMsg{}
}

func hidePostsList() tea.Msg {
	return hidePostsListMsg{}
}

type post struct {
	title         string
	author        string
	subreddit     string
	friendlyDate  string
	postUrl       string
	commentsUrl   string
	totalComments string
	totalLikes    string
}

func (p post) Title() string {
	return p.title
}

func (p post) Description() string {
	return fmt.Sprintf(" %s  %s  %s comments  %s", p.totalLikes, p.subreddit, p.totalComments, p.friendlyDate)
}

func (p post) FilterValue() string {
	return p.title
}

type postsList struct {
	model   list.Model
	spinner redditSpinner
	active  bool
}

func newPostsList() postsList {
	model := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	model.Title = defaultListTitle

	spinner := NewRedditSpinner()
	spinner.enable()

	return postsList{model, spinner, false}
}

func getSubredditPage(subreddit string) tea.Cmd {
	return tea.Batch(
		showSpinner(subreddit),
		func() tea.Msg {
			posts, err := getSubredditPosts(subreddit)
			if err != nil {
				fmt.Printf("Could not load reddit posts: %v", err)
				os.Exit(1)
			}

			items := getListItems(posts)
			return showPostsListMsg{items: items, title: subreddit}
		})
}

func getHomePage() tea.Cmd {
	return tea.Batch(
		showSpinner(defaultSpinnerTitle),
		func() tea.Msg {
			posts, err := getHomePosts()
			if err != nil {
				fmt.Printf("Could not load reddit posts: %v", err)
				os.Exit(1)
			}

			items := getListItems(posts)
			return showPostsListMsg{items, defaultListTitle}
		})
}

func getListItems(posts []post) []list.Item {
	var items []list.Item
	for _, p := range posts {
		items = append(items, p)
	}
	return items
}

func (p *postsList) enable() {
	p.active = true
}

func (p *postsList) disable() {
	p.active = false
}

func (p postsList) Init() tea.Cmd {
	return tea.Batch(
		p.spinner.Init(),
		getHomePage())
}

func (p postsList) Update(msg tea.Msg) (postsList, tea.Cmd) {
	switch msg := msg.(type) {
	case showPostsListMsg:
		p.enable()
		p.spinner.disable()
		if len(msg.items) > 0 {
			p.model.Title = msg.title
			p.model.ResetSelected()
			return p, p.model.SetItems(msg.items)
		}
	case tea.KeyMsg:
		if !p.active {
			break
		}

		switch keypress := msg.String(); keypress {
		case "q", "esc", "ctrl+c":
			return p, tea.Quit
		case "c", "C":
			return p, tea.Quit
		case "s", "S":
			p.disable()
			return p, showSubredditInput()
		}
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		p.model.SetSize(msg.Width-h, msg.Height-v)
	}

	var listCmd, spinnerCmd tea.Cmd

	p.model, listCmd = p.model.Update(msg)
	p.spinner, spinnerCmd = p.spinner.Update(msg)
	return p, tea.Batch(listCmd, spinnerCmd)
}

func (p postsList) View() string {
	if p.spinner.active {
		return p.spinner.View()
	} else {
		return listStyle.Render(p.model.View())
	}
}
