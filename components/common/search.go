package common

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

const (
	searchHelpText    = "Select a subreddit:"
	searchPlaceholder = "subreddit"
)

var (
	searchStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	searchContainerStyle = lipgloss.NewStyle().Margin(1, 2)
)

type SubredditSearch struct {
	textinput.Model
	Searching bool
}

func NewSubredditSearch() SubredditSearch {
	searchTextInput := textinput.New()
	searchTextInput.Placeholder = searchPlaceholder
	searchTextInput.ShowSuggestions = true
	searchTextInput.SetSuggestions(subredditSuggestions)
	searchTextInput.CharLimit = 30

	return SubredditSearch{Model: searchTextInput}
}

func (s SubredditSearch) View() string {
	view := searchStyle.Render(fmt.Sprintf("%s\n%s", searchHelpText, s.Model.View()))
	return searchContainerStyle.Render(view)
}

func (s *SubredditSearch) SetSearching(searching bool) {
	s.Searching = searching
	if searching {
		s.Focus()
		s.Reset()
	} else {
		s.Blur()
	}
}
