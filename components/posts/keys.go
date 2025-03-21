package posts

import "github.com/charmbracelet/bubbles/key"

type postsKeyMap struct {
	Home   key.Binding
	Search key.Binding
	Back   key.Binding
	Load   key.Binding
}

var postsKeys = postsKeyMap{
	Home: key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "home")),
	Search: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "subreddit search")),
	Back: key.NewBinding(
		key.WithKeys("bs"),
		key.WithHelp("bs", "back")),
	Load: key.NewBinding(
		key.WithKeys("L"),
		key.WithHelp("L", "load more posts")),
}

func (k postsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Home, k.Search, k.Load}
}

func (k postsKeyMap) FullHelp() []key.Binding {
	return []key.Binding{k.Home, k.Search, k.Back, k.Load}
}
