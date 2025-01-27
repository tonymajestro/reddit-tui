package posts

import (
	"fmt"
	"reddittui/client"
	"reddittui/components/common"
	"reddittui/components/messages"
	"reddittui/utils"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultHeaderTitle       = "reddit.com"
	defaultHeaderDescription = "The front page of the internet"
	defaultLoadingMessage    = "loading reddit.com..."
)

type PostsPage struct {
	posts        []client.Post
	redditClient client.RedditClient
	header       PostsHeader
	list         list.Model
	spinner      common.Spinner
	search       common.SubredditSearch
	w, h         int
	focus        bool
	home         bool
}

func NewPostsPage() PostsPage {
	items := list.New(nil, NewPostsDelegate(), 0, 0)
	items.SetShowTitle(false)
	items.SetShowStatusBar(false)
	items.SetFilteringEnabled(false)
	items.AdditionalShortHelpKeys = postsKeys.ShortHelp
	items.AdditionalFullHelpKeys = postsKeys.FullHelp

	header := NewPostsHeader()
	search := common.NewSubredditSearch()
	spinner := common.NewSpinner()

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
		return messages.LoadHomeMsg{}
	})
}

func (p PostsPage) Update(msg tea.Msg) (PostsPage, tea.Cmd) {
	switch msg := msg.(type) {

	case messages.LoadHomeMsg:
		return p, p.LoadHome()

	case messages.UpdatePostsMsg:
		p.updatePosts(client.Posts(msg))
		return p, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			if p.search.Searching {
				p.hideSearch()
				return p, nil
			}
		case "enter":
			if p.search.Searching {
				p.hideSearch()
				return p, messages.GoSubreddit(p.search.Value())
			} else if !p.spinner.Loading {
				loadCommentsCmd := func() tea.Msg {
					post := p.posts[p.list.Index()]
					return messages.LoadCommentsMsg(post)
				}

				return p, loadCommentsCmd
			}
		case "s", "S":
			if !p.search.Searching && !p.spinner.Loading {
				return p, p.showSearch()
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

	p.resizeComponents()
}

func (p *PostsPage) resizeComponents() {
	p.header.SetSize(p.w, p.h)
	p.resizeList()
}

func (p *PostsPage) showSearch() tea.Cmd {
	p.search.SetSearching(true)
	p.resizeList()

	return p.search.Focus()
}

func (p *PostsPage) hideSearch() {
	p.search.SetSearching(false)
	p.resizeList()
}

func (p *PostsPage) LoadHome() tea.Cmd {
	p.spinner.SetLoading(true)
	p.spinner.LoadingMessage = defaultLoadingMessage

	getPostsCmd := func() tea.Msg {
		posts, _ := p.redditClient.GetHomePosts()
		return messages.UpdatePostsMsg(posts)
	}

	return tea.Batch(getPostsCmd, p.spinner.Tick)
}

func (p *PostsPage) updatePosts(posts client.Posts) {
	p.spinner.SetLoading(false)

	p.posts = posts.Posts

	if posts.IsHome {
		p.header.SetContent(defaultHeaderTitle, defaultHeaderDescription)
		p.home = true
	} else {
		p.header.SetContent(posts.Subreddit, posts.Description)
		p.home = false
	}

	p.list.ResetSelected()

	var listItems []list.Item
	for _, p := range posts.Posts {
		listItems = append(listItems, p)
	}
	p.list.SetItems(listItems)

	// Need to set size again when content loads so padding and margins are correct
	p.resizeComponents()
}

func (p *PostsPage) LoadSubreddit(subreddit string) tea.Cmd {
	p.spinner.SetLoading(true)
	p.spinner.LoadingMessage = fmt.Sprintf("loading %s...", utils.NormalizeSubreddit(subreddit))

	getPostsCmd := func() tea.Msg {
		posts, _ := p.redditClient.GetSubredditPosts(subreddit)
		return messages.UpdatePostsMsg(posts)
	}

	return tea.Batch(getPostsCmd, p.spinner.Tick)
}

func (p *PostsPage) resizeList() {
	var (
		listHeight   int
		listWidth    = p.w - postsListStyle.GetHorizontalFrameSize()
		headerHeight = lipgloss.Height(p.header.View())
	)

	if p.search.Searching {
		searchHeight := lipgloss.Height(p.search.View())
		listHeight = p.h - headerHeight - searchHeight
	} else {
		listHeight = p.h - headerHeight
	}

	p.list.SetSize(listWidth, listHeight)
}
