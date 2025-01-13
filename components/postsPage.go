package components

import (
	"fmt"
	"log"
	"reddittui/client"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultListTitle = "reddit.com"

type PostsPage struct {
	posts           []client.Post
	itemsList       list.Model
	spinner         RedditSpinner
	subredditSearch SubredditSearch
	redditClient    client.RedditClient
	w, h            int
	focus           bool
}

func NewPostsPage() PostsPage {
	items := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	items.Title = defaultListTitle
	items.SetShowStatusBar(false)

	spinner := NewRedditSpinner()
	spinner.Focus()

	subredditSearch := NewSubredditSearch()

	redditClient := client.New()

	return PostsPage{
		posts:           []client.Post{},
		itemsList:       items,
		spinner:         spinner,
		subredditSearch: subredditSearch,
		redditClient:    redditClient,
		focus:           false,
	}
}

func (p PostsPage) Init() tea.Cmd {
	return p.spinner.Init()
}

func (p PostsPage) Update(msg tea.Msg) (PostsPage, tea.Cmd) {
	var (
		spinnerCmd tea.Cmd
		searchCmd  tea.Cmd
		postsCmd   tea.Cmd
	)

	if !p.focus {
		return p, nil
	}

	switch msg := msg.(type) {

	case acceptSearchMsg:
		p.ShowLoading(fmt.Sprintf("loading r/%s...", msg.subreddit))
		return p, tea.Batch(p.LoadSubreddit(msg.subreddit), p.spinner.Focus())

	case showPostsMsg:
		p.HideLoading()
		p.maximizePostsList()
		p.posts = msg.posts
		p.itemsList.Title = msg.title
		p.itemsList.ResetSelected()
		return p, p.itemsList.SetItems(msg.items)

	case tea.KeyMsg:
		if p.spinner.IsFocused() || p.subredditSearch.IsFocused() {
			break
		}

		switch keypress := msg.String(); keypress {
		case "s", "S":
			p.ShowSearch()
			return p, nil
		case "c", "C":
			return p, ShowComments(p.posts[p.itemsList.Index()])
		}
	}

	p.spinner, spinnerCmd = p.spinner.Update(msg)
	p.subredditSearch, searchCmd = p.subredditSearch.Update(msg)
	p.itemsList, postsCmd = p.itemsList.Update(msg)

	return p, tea.Batch(spinnerCmd, searchCmd, postsCmd)
}

func (p PostsPage) View() string {
	if p.spinner.IsFocused() {
		return appStyle.Render(p.spinner.View())
	} else if p.subredditSearch.IsFocused() {
		searchView := p.subredditSearch.View()
		joinedView := lipgloss.JoinVertical(lipgloss.Left, searchView, p.itemsList.View())
		return appStyle.Render(joinedView)
	} else {
		return appStyle.Render(p.itemsList.View())
	}
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
	p.itemsList.SetSize(w, h)
}

func (p *PostsPage) ShowLoading(message string) {
	p.spinner.SetLoadingMessage(message)
	p.spinner.Focus()
}

func (p *PostsPage) HideLoading() {
	p.spinner.Blur()
}

func (p *PostsPage) ShowSearch() {
	p.shrinkPostsList()
	p.subredditSearch.Focus()
}

func (p *PostsPage) HideSearch() {
	p.subredditSearch.Blur()
}

func (p *PostsPage) ShowPosts() {
	p.Focus()
	p.spinner.Blur()
	p.subredditSearch.Blur()
}

func (p PostsPage) LoadHome() tea.Cmd {
	return func() tea.Msg {
		posts, err := p.redditClient.GetHomePosts()
		if err != nil {
			log.Printf("Error: %v", err)
			return err
		}

		items := getPostListItems(posts)
		return showPostsMsg{
			posts: posts,
			items: items,
			title: defaultListTitle,
		}
	}
}

func (p PostsPage) LoadSubreddit(subreddit string) tea.Cmd {
	return func() tea.Msg {
		posts, _ := p.redditClient.GetSubredditPosts(subreddit)
		items := getPostListItems(posts)
		return showPostsMsg{
			posts: posts,
			items: items,
			title: subreddit,
		}
	}
}

func (p *PostsPage) shrinkPostsList() {
	_, h := lipgloss.Size(p.subredditSearch.View())
	p.itemsList.SetHeight(p.itemsList.Height() - h)
}

func (p *PostsPage) maximizePostsList() {
	p.itemsList.SetSize(p.w, p.h)
}

func getPostListItems(posts []client.Post) []list.Item {
	var items []list.Item
	for _, p := range posts {
		items = append(items, p)
	}
	return items
}
