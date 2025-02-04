package comments

import "github.com/charmbracelet/bubbles/key"

type viewportKeyMap struct {
	CursorUp         key.Binding
	CursorDown       key.Binding
	GoToStart        key.Binding
	GoToEnd          key.Binding
	OpenPost         key.Binding
	GoHome           key.Binding
	CollapseComments key.Binding
	ShowFullHelp     key.Binding
	CloseFullHelp    key.Binding
	Quit             key.Binding
	ForceQuit        key.Binding
}

var commentsKeys = viewportKeyMap{
	CursorUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	CursorDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	GoToStart: key.NewBinding(
		key.WithKeys("home", "g"),
		key.WithHelp("g/home", "go to start"),
	),
	GoToEnd: key.NewBinding(
		key.WithKeys("end", "G"),
		key.WithHelp("G/end", "go to end"),
	),
	OpenPost: key.NewBinding(
		key.WithKeys("o", "O"),
		key.WithHelp("o", "open post"),
	),
	GoHome: key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "go home"),
	),
	CollapseComments: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "collapse comments"),
	),
	ShowFullHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "more"),
	),
	CloseFullHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "close help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q", "quit"),
	),
	ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
}

func (k viewportKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.CursorUp, k.CursorDown, k.OpenPost, k.GoHome, k.ShowFullHelp}
}

func (k viewportKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CursorUp, k.CursorDown, k.GoToStart, k.GoToEnd, k.OpenPost},
		{k.GoHome, k.CollapseComments, k.Quit, k.CloseFullHelp},
	}
}
