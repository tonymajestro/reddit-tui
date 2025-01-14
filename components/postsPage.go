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

const defaultListTitle = "reddit.com"

type loadHomeMsg struct{}

type displayPostsMsg struct {
	posts     []client.Post
	title     string
	subreddit string
}

var (
	spinnerStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	searchStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	searchContainerStyle = lipgloss.NewStyle().Margin(1, 2)
)

type PostsPage struct {
	posts          []client.Post
	listModel      list.Model
	spinnerModel   spinner.Model
	searchModel    textinput.Model
	redditClient   client.RedditClient
	loading        bool
	searching      bool
	loadingMessage string
	subreddit      string
	w, h           int
	focus          bool
}

func NewPostsPage() PostsPage {
	items := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	items.Title = defaultListTitle
	items.SetStatusBarItemName("post", "posts")

	searchModel := textinput.New()
	searchModel.ShowSuggestions = true
	searchModel.SetSuggestions(subredditSuggestions)
	searchModel.CharLimit = 30

	redditClient := client.New()

	return PostsPage{
		posts:        []client.Post{},
		listModel:    items,
		searchModel:  searchModel,
		redditClient: redditClient,
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
		p.DisplayPosts(msg.posts, msg.title, msg.subreddit)
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
			}
		case "ctrl+c":
			return p, tea.Quit
		case "s", "S":
			if !p.searching && !p.loading {
				return p, p.ShowSearch()
			}
		case "c", "C":
			if !p.searching && !p.loading {
				loadCommentsCmd := func() tea.Msg {
					post := p.posts[p.listModel.Index()]
					return loadCommentsMsg{post, p.subreddit}
				}

				return p, loadCommentsCmd
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
	if p.loading {
		return appStyle.Render(p.GetSpinnerView())
	} else if p.searching {
		searchView := p.GetSearchView()
		joinedView := lipgloss.JoinVertical(lipgloss.Left, searchView, p.listModel.View())
		return appStyle.Render(joinedView)
	} else {
		return appStyle.Render(p.listModel.View())
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
	p.listModel.SetSize(w, h)
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
	p.shrinkPostsList()
	p.searching = true
	p.searchModel.Reset()
	return p.searchModel.Focus()
}

func (p *PostsPage) HideSearch() {
	p.maximizePostsList()
	p.searching = false
	p.searchModel.Blur()
}

func (p PostsPage) GetSpinnerView() string {
	return fmt.Sprintf("%s %s", p.spinnerModel.View(), p.loadingMessage)
}

func (p PostsPage) GetSearchView() string {
	selectionView := searchStyle.Render(fmt.Sprintf("Choose a subreddit:\n%s", p.searchModel.View()))
	return searchContainerStyle.Render(selectionView)
}

func (p *PostsPage) LoadHome() tea.Cmd {
	p.ShowLoading("loading reddit.com...")

	getPostsCmd := func() tea.Msg {
		posts, _ := p.redditClient.GetHomePosts()
		return displayPostsMsg{
			posts:     posts,
			title:     defaultListTitle,
			subreddit: defaultListTitle,
		}
	}

	return tea.Batch(getPostsCmd, p.spinnerModel.Tick)
}

func (p *PostsPage) DisplayPosts(posts []client.Post, title, subreddit string) {
	p.HideLoading()
	p.maximizePostsList()

	p.listModel.Title = title
	p.subreddit = subreddit

	p.posts = posts
	p.listModel.ResetSelected()
	p.listModel.SetItems(getPostListItems(posts))
}

func (p *PostsPage) LoadSubreddit(subreddit string) tea.Cmd {
	p.ShowLoading(fmt.Sprintf("loading r/%s...", subreddit))

	getPostsCmd := func() tea.Msg {
		posts, _ := p.redditClient.GetSubredditPosts(subreddit)
		return displayPostsMsg{
			posts:     posts,
			title:     subreddit,
			subreddit: subreddit,
		}
	}

	return tea.Batch(getPostsCmd, p.spinnerModel.Tick)
}

func (p *PostsPage) shrinkPostsList() {
	_, h := lipgloss.Size(p.GetSearchView())
	p.listModel.SetHeight(p.listModel.Height() - h)
}

func (p *PostsPage) maximizePostsList() {
	p.listModel.SetSize(p.w, p.h)
}

func getPostListItems(posts []client.Post) []list.Item {
	var items []list.Item
	for _, p := range posts {
		items = append(items, p)
	}
	return items
}
