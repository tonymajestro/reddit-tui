package posts

import (
	"fmt"
	"log/slog"
	"reddittui/client"
	"reddittui/client/common"
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
	postsErrorText           = "Could not load posts. Please try again in a few moments."
	subredditNotFoundText    = "Subreddit not found"
)

type PostsPage struct {
	Subreddit      string
	posts          model.Posts
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
	if home {
		header.SetContent(defaultHeaderTitle, defaultHeaderDescription)
	}

	containerStyle := styles.GlobalStyle

	return PostsPage{
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

	case messages.LoadMorePostsMsg:
		isHome := bool(msg)
		if p.Home == isHome {
			return p, p.loadMorePosts()
		}

	case messages.UpdatePostsMsg:
		posts := model.Posts(msg)
		if posts.IsHome == p.Home {
			p.updatePosts(posts)
			return p, messages.LoadingComplete
		}

	case messages.AddMorePostsMsg:
		posts := model.Posts(msg)
		if posts.IsHome == p.Home {
			p.addPosts(posts)
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
				post := p.posts.Posts[p.list.Index()]
				return messages.LoadCommentsMsg(post.CommentsUrl)
			}

			return p, loadCommentsCmd

		case "q", "Q":
			// Ignore q keystrokes to list.Modal. since it will default to sending a Quit message
			// instead of showing the quit modal. Tui component will correctly handle quit mesages
			return p, nil

		case "L":
			return p, messages.LoadMorePosts(p.Home)

		case "H":
			return p, messages.LoadHome

		case "esc", "backspace", "left", "h":
			return p, messages.GoBack
		}
	}

	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return p, cmd
}

func (p PostsPage) View() string {
	if len(p.posts.Posts) == 0 {
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
		posts, err := p.redditClient.GetHomePosts("")
		if err != nil {
			slog.Error(postsErrorText, "error", err)
			return messages.ShowErrorModalMsg{ErrorMsg: postsErrorText}
		}

		return messages.UpdatePostsMsg(posts)
	}
}

func (p *PostsPage) loadMorePosts() tea.Cmd {
	return func() tea.Msg {
		var (
			posts model.Posts
			err   error
		)

		if len(p.posts.After) == 0 {
			slog.Error(postsErrorText, "error", err)
			return messages.ShowErrorModalMsg{ErrorMsg: postsErrorText}
		}

		if p.posts.IsHome {
			posts, err = p.redditClient.GetHomePosts(p.posts.After)
		} else {
			posts, err = p.redditClient.GetSubredditPosts(p.Subreddit, p.posts.After)
		}

		if err != nil {
			slog.Error(postsErrorText, "error", err)
			return messages.ShowErrorModalMsg{ErrorMsg: postsErrorText}
		}

		return messages.AddMorePostsMsg(posts)
	}
}

func (p PostsPage) loadSubreddit(subreddit string) tea.Cmd {
	return func() tea.Msg {
		posts, err := p.redditClient.GetSubredditPosts(subreddit, "")
		if err == common.ErrNotFound {
			slog.Error(subredditNotFoundText, "error", err, "subreddit", subreddit)
			return messages.ShowErrorModalMsg{ErrorMsg: fmt.Sprintf("%s: %s", subredditNotFoundText, subreddit)}
		} else if err != nil {
			slog.Error(postsErrorText, "error", err)
			return messages.ShowErrorModalMsg{ErrorMsg: postsErrorText}
		}

		return messages.UpdatePostsMsg(posts)
	}
}

func (p *PostsPage) updatePosts(posts model.Posts) {
	p.posts = posts

	if posts.IsHome {
		p.header.SetContent(defaultHeaderTitle, defaultHeaderDescription)
	} else {
		p.header.SetContent(posts.Subreddit, posts.Description)
		p.Subreddit = posts.Subreddit
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

func (p *PostsPage) addPosts(posts model.Posts) {
	uniqueTitles := make(map[string]bool)

	p.posts.Posts = append(p.posts.Posts, posts.Posts...)
	p.posts.After = posts.After

	// Merge existing posts with new posts, avoiding duplicates
	var listItems []list.Item
	for _, p := range p.posts.Posts {
		if _, ok := uniqueTitles[p.PostTitle]; !ok {
			listItems = append(listItems, p)
			uniqueTitles[p.PostTitle] = true
		}
	}
	for _, p := range posts.Posts {
		if _, ok := uniqueTitles[p.PostTitle]; !ok {
			listItems = append(listItems, p)
			uniqueTitles[p.PostTitle] = true
		}
	}

	p.list.SetItems(listItems)

	// Need to set size again when content loads so padding and margins are correct
	p.resizeComponents()
}
