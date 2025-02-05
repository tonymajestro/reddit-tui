package posts

import (
	"log/slog"
	"os"
	"reddittui/client"
	"reddittui/components/messages"
	"reddittui/components/styles"
	"reddittui/model"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultHeaderTitle       = "reddit.com"
	defaultHeaderDescription = "The front page of the internet"
)

type PostsPage struct {
	Subreddit      string
	posts          []model.Post
	redditClient   client.RedditClient
	header         PostsHeader
	list           list.Model
	focus          bool
	Home           bool
	containerStyle lipgloss.Style
}

func NewPostsPage(redditClient client.RedditClient, home bool) PostsPage {
	items := list.New(nil, NewPostsDelegate(), 0, 0)
	items.SetShowTitle(false)
	items.SetShowStatusBar(false)
	items.KeyMap.NextPage.SetEnabled(false)
	items.KeyMap.PrevPage.SetEnabled(false)
	items.SetFilteringEnabled(false)
	items.AdditionalShortHelpKeys = postsKeys.ShortHelp
	items.AdditionalFullHelpKeys = postsKeys.FullHelp

	header := NewPostsHeader()
	header.SetContent(defaultHeaderTitle, defaultHeaderDescription)

	containerStyle := styles.GlobalStyle

	return PostsPage{
		posts:          []model.Post{},
		list:           items,
		redditClient:   redditClient,
		header:         header,
		Home:           home,
		containerStyle: containerStyle,
	}
}

func (p PostsPage) Init() tea.Cmd {
	return nil
}

func (p PostsPage) Update(msg tea.Msg) (PostsPage, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if p.focus {
		p, cmd = p.handleFocusedMessages(msg)
		cmds = append(cmds, cmd)
	}

	p, cmd = p.handleGlobalMessages(msg)
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p PostsPage) handleGlobalMessages(msg tea.Msg) (PostsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.LoadHomeMsg:
		if p.Home {
			return p, p.loadHome()
		}

	case messages.LoadSubredditMsg:
		if !p.Home {
			subreddit := string(msg)
			return p, p.loadSubreddit(subreddit)
		}

	case messages.UpdatePostsMsg:
		posts := model.Posts(msg)
		if posts.IsHome == p.Home {
			p.updatePosts(posts)
			return p, messages.LoadingComplete
		}
	}

	return p, nil
}

func (p PostsPage) handleFocusedMessages(msg tea.Msg) (PostsPage, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter", "right", "l":
			loadCommentsCmd := func() tea.Msg {
				post := p.posts[p.list.Index()]
				return messages.LoadCommentsMsg(post)
			}

			return p, loadCommentsCmd

		case "q", "Q":
			// Ignore q keystrokes to list.Modal. since it will default to sending a Quit message
			// instead of showing the quit modal. Tui component will correctly handle quit mesages
			return p, nil

		case "H":
			return p, messages.LoadHome

		case "b", "B", "escape", "backspace", "left", "h":
			return p, messages.GoBack
		}
	}

	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return p, cmd
}

func (p PostsPage) View() string {
	if len(p.posts) == 0 {
		return p.containerStyle.Render("")
	}

	headerView := p.header.View()
	listView := p.list.View()
	joined := lipgloss.JoinVertical(lipgloss.Left, headerView, listView)
	return p.containerStyle.Render(joined)
}

func (p *PostsPage) SetSize(w, h int) {
	p.containerStyle = p.containerStyle.Width(w).Height(h)
	p.resizeComponents()
}

func (p *PostsPage) Focus() {
	p.focus = true
}

func (p *PostsPage) Blur() {
	p.focus = false
}

func (p *PostsPage) resizeComponents() {
	var (
		w            = p.containerStyle.GetWidth() - p.containerStyle.GetHorizontalFrameSize()
		h            = p.containerStyle.GetHeight() - p.containerStyle.GetVerticalFrameSize()
		listWidth    = w - postsListStyle.GetHorizontalFrameSize()
		headerHeight = lipgloss.Height(p.header.View())
		listHeight   = h - headerHeight
	)

	p.header.SetSize(w, h)
	p.list.SetSize(listWidth, listHeight)
}

func (p *PostsPage) loadHome() tea.Cmd {
	return func() tea.Msg {
		slog.Info("Loading home page posts", "subreddit", "reddit.com")
		posts, err := p.redditClient.GetHomePosts()
		if err != nil {
			slog.Error("Error loading home page posts", "error", err)
			os.Exit(1)
		}
		return messages.UpdatePostsMsg(posts)
	}
}

func (p PostsPage) loadSubreddit(subreddit string) tea.Cmd {
	p.Subreddit = subreddit
	return func() tea.Msg {
		posts, err := p.redditClient.GetSubredditPosts(subreddit)
		if err != nil {
			slog.Error("Error loading home page posts", "error", err)
			os.Exit(1)
		}
		return messages.UpdatePostsMsg(posts)
	}
}

func (p *PostsPage) updatePosts(posts model.Posts) {
	p.posts = posts.Posts

	if posts.IsHome {
		p.header.SetContent(defaultHeaderTitle, defaultHeaderDescription)
	} else {
		p.header.SetContent(posts.Subreddit, posts.Description)
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
