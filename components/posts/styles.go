package posts

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var postsListStyle = lipgloss.NewStyle().MarginRight(4)

func NewPostsDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()

	listStyle := delegate.Styles
	listStyle.NormalTitle = listStyle.NormalTitle.Bold(false)
	listStyle.SelectedTitle = listStyle.SelectedTitle.Bold(true)
	delegate.Styles = listStyle

	return delegate
}
