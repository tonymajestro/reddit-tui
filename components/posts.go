package components

import (
	"fmt"
	"reddittui/client"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultHeaderTitle       = "reddit.com"
	defaultHeaderDescription = "The front page of the internet"
	defaultLoadingMessage    = "loading reddit.com..."
	searchHelpText           = "Select a subreddit:"
	searchPlaceholder        = "subreddit"
)

type (
	loadHomeMsg    struct{}
	updatePostsMsg client.Posts
)

type postsKeyMap struct {
	Home   key.Binding
	Search key.Binding
	Back   key.Binding
}

func (k postsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Home, k.Search}
}

func (k postsKeyMap) FullHelp() []key.Binding {
	return []key.Binding{k.Home, k.Search, k.Back}
}

type PostsPage struct {
	posts        []client.Post
	redditClient client.RedditClient
	header       Header
	list         list.Model
	spinner      Spinner
	search       SubredditSearch
	w, h         int
	focus        bool
	home         bool
}

func NewPostsPage() PostsPage {
	keys := postsKeyMap{
		Home: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "home")),
		Search: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "subreddit search")),
		Back: key.NewBinding(
			key.WithKeys("bs"),
			key.WithHelp("bs", "back")),
	}

	items := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	items.SetShowTitle(false)
	items.SetShowStatusBar(false)
	items.SetFilteringEnabled(false)
	items.AdditionalShortHelpKeys = keys.ShortHelp
	items.AdditionalFullHelpKeys = keys.FullHelp

	header := NewHeader()
	search := NewSubredditSearch()
	spinner := NewSpinner()

	redditClient := client.New()

	return PostsPage{
		posts:        []client.Post{},
		list:         items,
		search:       search,
		redditClient: redditClient,
		spinner:      spinner,
		header:       header,
		home:         true,
	}
}

func (p PostsPage) Init() tea.Cmd {
	return tea.Batch(p.spinner.Init(), func() tea.Msg {
		return loadHomeMsg{}
	})
}

func (p PostsPage) Update(msg tea.Msg) (PostsPage, tea.Cmd) {
	switch msg := msg.(type) {

	case loadHomeMsg:
		return p, p.LoadHome()

	case updatePostsMsg:
		p.UpdatePosts(client.Posts(msg))
		return p, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			if p.search.Searching {
				p.HideSearch()
				return p, nil
			}
		case "enter":
			if p.search.Searching {
				p.HideSearch()
				return p, GoSubreddit(p.search.Value())
			} else if !p.spinner.Loading {
				loadCommentsCmd := func() tea.Msg {
					post := p.posts[p.list.Index()]
					return loadCommentsMsg(post)
				}

				return p, loadCommentsCmd
			}
		case "s", "S":
			if !p.search.Searching && !p.spinner.Loading {
				return p, p.ShowSearch()
			}
		}
	}

	var cmd tea.Cmd

	if p.spinner.Loading {
		p.spinner, cmd = p.spinner.Update(msg)
		return p, cmd
	} else if p.search.Searching {
		p.search.Model, cmd = p.search.Update(msg)
		return p, cmd
	} else {
		p.list, cmd = p.list.Update(msg)
		return p, cmd
	}
}

func (p PostsPage) View() string {
	if p.spinner.Loading {
		return p.spinner.View()
	}

	headerView := p.header.View()
	listView := p.list.View()

	if p.search.Searching {
		searchView := p.search.View()
		return lipgloss.JoinVertical(lipgloss.Left, searchView, headerView, listView)
	}

	return lipgloss.JoinVertical(lipgloss.Left, headerView, listView)
}

func (p *PostsPage) IsFocused() bool {
	return p.focus
}

func (p *PostsPage) Focus() {
	p.focus = true
}

func (p *PostsPage) Blur() {
	p.focus = false
}

func (p *PostsPage) IsSearching() bool {
	return p.search.Searching
}

func (p *PostsPage) IsLoading() bool {
	return p.spinner.Loading
}

func (p *PostsPage) IsHome() bool {
	return p.home
}

func (p *PostsPage) SetSize(w, h int) {
	p.w = w
	p.h = h

	p.ResizeComponents()
}

func (p *PostsPage) ResizeComponents() {
	p.header.SetSize(p.w, p.h)
	p.resizeList()
}

func (p *PostsPage) ShowSearch() tea.Cmd {
	p.search.SetSearching(true)
	p.resizeList()

	return p.search.Focus()
}

func (p *PostsPage) HideSearch() {
	p.search.SetSearching(false)
	p.resizeList()
}

func (p *PostsPage) LoadHome() tea.Cmd {
	p.spinner.SetLoading(true)
	p.spinner.SetLoadingMessage(defaultLoadingMessage)

	getPostsCmd := func() tea.Msg {
		posts, _ := p.redditClient.GetHomePosts()
		return updatePostsMsg(posts)
	}

	return tea.Batch(getPostsCmd, p.spinner.Tick)
}

func (p *PostsPage) UpdatePosts(posts client.Posts) {
	p.spinner.SetLoading(false)

	p.posts = posts.Posts

	if posts.IsHome {
		p.header.SetTitle(defaultHeaderTitle)
		p.home = true
	} else {
		p.header.SetTitle(normalizeSubreddit(posts.Subreddit))
		p.home = false
	}

	p.header.SetDescription(posts.Description)

	p.list.ResetSelected()

	var listItems []list.Item
	for _, p := range posts.Posts {
		listItems = append(listItems, p)
	}
	p.list.SetItems(listItems)

	// Need to set size again when content loads so padding and margins are correct
	p.ResizeComponents()
}

func (p *PostsPage) LoadSubreddit(subreddit string) tea.Cmd {
	p.spinner.SetLoading(true)
	p.spinner.SetLoadingMessage(fmt.Sprintf("loading %s...", normalizeSubreddit(subreddit)))

	getPostsCmd := func() tea.Msg {
		posts, _ := p.redditClient.GetSubredditPosts(subreddit)
		return updatePostsMsg(posts)
	}

	return tea.Batch(getPostsCmd, p.spinner.Tick)
}

func (p *PostsPage) resizeList() {
	var listHeight int

	headerHeight := lipgloss.Height(p.header.View())

	if p.search.Searching {
		searchHeight := lipgloss.Height(p.search.View())
		listHeight = p.h - headerHeight - searchHeight
	} else {
		listHeight = p.h - headerHeight
	}

	p.list.SetSize(p.w, listHeight)
}
