package components

import (
	"fmt"
	"reddittui/client"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultListTitle       = "reddit.com"
	defaultListDescription = "The front page of the internet"
	defaultLoadingMessage  = "loading reddit.com..."
	searchHelpText         = "Select a subreddit:"
	searchPlaceholder      = "subreddit"
)

var (
	spinnerStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	spinnerContainerStyle = lipgloss.NewStyle().Margin(2, 2)
	searchStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	searchContainerStyle  = lipgloss.NewStyle().Margin(1, 2, 2, 2)
	descriptionStyle      = lipgloss.NewStyle().Margin(0, 2, 0, 2)
)

var titleStyle = lipgloss.NewStyle().
	Margin(0, 2, 1, 2).
	Padding(0, 2).
	Background(lipgloss.Color("62")).
	Foreground(lipgloss.Color("230"))

type (
	loadHomeMsg     struct{}
	displayPostsMsg client.Posts
)

type PostsPage struct {
	posts          []client.Post
	listModel      list.Model
	spinnerModel   spinner.Model
	searchModel    textinput.Model
	redditClient   client.RedditClient
	loadingMessage string
	subreddit      string
	description    string
	home           bool
	loading        bool
	searching      bool
	ready          bool
	w, h           int
	focus          bool
}

func NewPostsPage() PostsPage {
	items := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	items.SetShowTitle(false)
	items.SetShowStatusBar(false)

	searchModel := textinput.New()
	searchModel.Placeholder = searchPlaceholder
	searchModel.ShowSuggestions = true
	searchModel.SetSuggestions(subredditSuggestions)
	searchModel.CharLimit = 30

	redditClient := client.New()

	return PostsPage{
		posts:        []client.Post{},
		listModel:    items,
		searchModel:  searchModel,
		redditClient: redditClient,
		home:         true,
	}
}

func (p PostsPage) Init() tea.Cmd {
	return func() tea.Msg {
		return loadHomeMsg{}
	}
}

func (p PostsPage) Update(msg tea.Msg) (PostsPage, tea.Cmd) {
	switch msg := msg.(type) {

	case loadHomeMsg:
		return p, p.LoadHome()

	case displayPostsMsg:
		p.DisplayPosts(client.Posts(msg))
		return p, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			if p.searching {
				p.HideSearch()
				return p, nil
			} else {
				return p, tea.Quit
			}
		case "enter":
			if p.searching {
				p.HideSearch()
				return p, p.LoadSubreddit(p.searchModel.Value())
			} else if !p.loading {
				loadCommentsCmd := func() tea.Msg {
					post := p.posts[p.listModel.Index()]
					return loadCommentsMsg{post, p.subreddit}
				}

				return p, loadCommentsCmd
			}
		case "ctrl+c":
			return p, tea.Quit
		case "s", "S":
			if !p.searching && !p.loading {
				return p, p.ShowSearch()
			}
		}
	}

	var cmd tea.Cmd

	if p.loading {
		p.spinnerModel, cmd = p.spinnerModel.Update(msg)
		return p, cmd
	} else if p.searching {
		p.searchModel, cmd = p.searchModel.Update(msg)
		return p, cmd
	} else {
		p.listModel, cmd = p.listModel.Update(msg)
		return p, cmd
	}
}

func (p PostsPage) View() string {
	var joinedView string

	if p.loading {
		return appStyle.Render(p.GetSpinnerView())
	}

	titleView := p.GetTitleView()
	descriptionView := p.GetDescriptionView()
	listView := p.listModel.View()

	if p.searching {
		searchView := p.GetSearchView()
		joinedView = lipgloss.JoinVertical(lipgloss.Left, searchView, titleView, descriptionView, listView)
	} else {
		joinedView = lipgloss.JoinVertical(lipgloss.Left, titleView, descriptionView, listView)
	}

	return appStyle.Render(joinedView)
}

func (p PostsPage) IsFocused() bool {
	return p.focus
}

func (p *PostsPage) Focus() {
	p.focus = true
}

func (p *PostsPage) Blur() {
	p.focus = false
}

func (p *PostsPage) SetSize(w, h int) {
	p.w = w
	p.h = h

	p.updateListHeight()
}

func (p *PostsPage) ShowLoading(message string) {
	spinnerModel := spinner.New()
	spinnerModel.Spinner = spinner.Dot
	spinnerModel.Style = spinnerStyle

	p.spinnerModel = spinnerModel
	p.loadingMessage = message
	p.loading = true
}

func (p *PostsPage) HideLoading() {
	p.loading = false
}

func (p *PostsPage) ShowSearch() tea.Cmd {
	p.searching = true
	p.searchModel.Reset()
	p.updateListHeight()

	return p.searchModel.Focus()
}

func (p *PostsPage) HideSearch() {
	p.searching = false
	p.searchModel.Blur()
	p.updateListHeight()
}

func (p PostsPage) GetTitleView() string {
	if p.home {
		return titleStyle.Render(defaultListTitle)
	}

	return titleStyle.Render(fmt.Sprintf("r/%s", p.subreddit))
}

func (p PostsPage) GetDescriptionView() string {
	if p.home {
		return descriptionStyle.Render(defaultListDescription)
	}

	return descriptionStyle.Render(p.description)
}

func (p PostsPage) GetSpinnerView() string {
	spinnerView := fmt.Sprintf("%s %s", p.spinnerModel.View(), p.loadingMessage)
	return spinnerContainerStyle.Render(spinnerView)
}

func (p PostsPage) GetSearchView() string {
	selectionView := searchStyle.Render(fmt.Sprintf("%s\n%s", searchHelpText, p.searchModel.View()))
	return searchContainerStyle.Render(selectionView)
}

func (p *PostsPage) LoadHome() tea.Cmd {
	p.ShowLoading(defaultLoadingMessage)

	getPostsCmd := func() tea.Msg {
		posts, _ := p.redditClient.GetHomePosts()
		return displayPostsMsg(posts)
	}

	return tea.Batch(getPostsCmd, p.spinnerModel.Tick)
}

func (p *PostsPage) DisplayPosts(posts client.Posts) {
	p.HideLoading()

	p.subreddit = posts.Subreddit
	p.posts = posts.Posts
	p.description = posts.Description
	p.home = posts.IsHome

	p.listModel.ResetSelected()
	p.listModel.SetItems(getPostListItems(posts.Posts))
}

func (p *PostsPage) LoadSubreddit(subreddit string) tea.Cmd {
	p.ShowLoading(fmt.Sprintf("loading r/%s...", subreddit))

	getPostsCmd := func() tea.Msg {
		posts, _ := p.redditClient.GetSubredditPosts(subreddit)
		return displayPostsMsg(posts)
	}

	return tea.Batch(getPostsCmd, p.spinnerModel.Tick)
}

func (p *PostsPage) updateListHeight() {
	var listHeight int

	titleHeight := lipgloss.Height(p.GetTitleView())
	descriptionHeight := lipgloss.Height(p.GetDescriptionView())

	if p.searching {
		searchHeight := lipgloss.Height(p.GetSearchView())
		listHeight = p.h - titleHeight - descriptionHeight - searchHeight
	} else {
		listHeight = p.h - titleHeight - descriptionHeight
	}

	p.listModel.SetSize(p.w, listHeight)
}

func getPostListItems(posts []client.Post) []list.Item {
	var items []list.Item
	for _, p := range posts {
		items = append(items, p)
	}
	return items
}
