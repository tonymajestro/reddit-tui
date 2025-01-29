package comments

import (
	"reddittui/components/colors"

	"github.com/charmbracelet/lipgloss"
)

var viewportStyle = lipgloss.NewStyle().Margin(0, 2, 1, 2)

var (
	commentAuthorStyle  = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Blue)).Bold(true)
	commentDateStyle    = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Lavender)).Italic(true)
	commentTextStyle    = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Text))
	popularPointsStyle  = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Purple))
	defaultPointsStyle  = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Purple))
	negativePointsStyle = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Red))
	collapsedStyle      = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Yellow))
)

var (
	postAuthorStyle    = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Blue))
	postPointsStyle    = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Purple))
	postTextStyle      = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Sand))
	postTimestampStyle = lipgloss.NewStyle().Foreground(colors.AdaptiveColor(colors.Text)).Faint(true)
)
