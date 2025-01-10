package reddit

import (
	"fmt"
	"reddittui/client"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultListTitle = "reddit.com"

var listStyle = lipgloss.NewStyle().Margin(1, 2)

type RedditTui struct {
	redditClient   client.RedditClient
	postsList      list.Model
	spinner        RedditSpinner
	subredditInput SubredditInput
	focus          bool
	w, h           int
}

func NewRedditTui() RedditTui {
	postsList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	postsList.Title = defaultListTitle
	postsList.SetShowStatusBar(false)

	spinner := NewRedditSpinner()
	spinner.Focus()

	subredditInput := NewSubredditInput()
	redditClient := client.New()

	return RedditTui{
		redditClient:   redditClient,
		postsList:      postsList,
		spinner:        spinner,
		subredditInput: subredditInput,
	}
}

func (p *RedditTui) Focus() {
	p.focus = true
}

func (p *RedditTui) Blur() {
	p.focus = false
}

func (p RedditTui) Init() tea.Cmd {
	return tea.Batch(
		fetchHomePosts,
		p.spinner.Init(),
		p.subredditInput.Init())
}

func (p RedditTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		listCmd      tea.Cmd
		spinnerCmd   tea.Cmd
		subredditCmd tea.Cmd
	)

	switch msg := msg.(type) {

	case fetchPostsMsg:
		p.Blur()
		p.spinner.Focus()

		if msg.home {
			return p, p.loadHomePage
		} else {
			p.spinner.loadingMessage = fmt.Sprintf("loading r/%s", msg.subreddit)
			return p, p.loadSubredditPage(msg.subreddit)
		}

	case showPostsMsg:
		p.Focus()
		p.spinner.Blur()
		p.maximizePostsList()

		if !msg.noFetch {
			p.postsList.Title = msg.title
			p.postsList.ResetSelected()
			return p, p.postsList.SetItems(msg.items)
		}

	case tea.KeyMsg:
		if !p.focus {
			break
		}

		switch keypress := msg.String(); keypress {
		case "q", "esc", "ctrl+c":
			return p, tea.Quit
		case "c", "C":
			return p, tea.Quit
		case "s", "S":
			p.Blur()
			p.shrinkPostsList()
			cmd := p.subredditInput.Focus()
			return p, cmd
		}
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		newW, newH := msg.Width-h, msg.Height-v
		p.w, p.h = newW, newH
		p.postsList.SetSize(newW, newH)
	}

	if p.focus {
		p.postsList, listCmd = p.postsList.Update(msg)
	}

	p.spinner, spinnerCmd = p.spinner.Update(msg)
	p.subredditInput, subredditCmd = p.subredditInput.Update(msg)

	return p, tea.Batch(listCmd, spinnerCmd, subredditCmd)
}

func (p RedditTui) View() string {
	if p.spinner.focus {
		return p.spinner.View()
	} else if p.subredditInput.focus {
		return lipgloss.JoinVertical(lipgloss.Left, p.subredditInput.View(), listStyle.Render(p.postsList.View()))
	} else {
		return listStyle.Render(p.postsList.View())
	}
}

type fetchPostsMsg struct {
	subreddit string
	home      bool
}

type showPostsMsg struct {
	title   string
	items   []list.Item
	noFetch bool
}

func fetchSubredditPosts(subreddit string) tea.Cmd {
	return func() tea.Msg {
		return fetchPostsMsg{
			home:      false,
			subreddit: subreddit,
		}
	}
}

func fetchHomePosts() tea.Msg {
	return fetchPostsMsg{home: true}
}

func focusListPage() tea.Msg {
	return showPostsMsg{noFetch: true}
}

func getListItems(posts []client.Post) []list.Item {
	var items []list.Item
	for _, p := range posts {
		items = append(items, p)
	}
	return items
}

func (p RedditTui) loadSubredditPage(subreddit string) tea.Cmd {
	return func() tea.Msg {
		posts, _ := p.redditClient.GetSubredditPosts(subreddit)
		items := getListItems(posts)
		return showPostsMsg{
			items: items,
			title: subreddit,
		}
	}
}

func (p *RedditTui) shrinkPostsList() {
	_, h := lipgloss.Size(p.subredditInput.View())
	p.postsList.SetHeight(p.postsList.Height() - h)
}

func (p *RedditTui) maximizePostsList() {
	p.postsList.SetSize(p.w, p.h)
}

func (p RedditTui) loadHomePage() tea.Msg {
	posts, _ := p.redditClient.GetHomePosts()

	items := getListItems(posts)
	return showPostsMsg{
		items: items,
		title: defaultListTitle,
	}
}
