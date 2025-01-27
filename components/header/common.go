package header

import (
	"reddittui/components/colors"

	"github.com/charmbracelet/lipgloss"
)

var HeaderContainerStyle = lipgloss.NewStyle().MarginBottom(2)

var (
	TitleStyle = lipgloss.NewStyle().
			MarginBottom(1).
			Padding(0, 2).
			Height(1).
			Background(colors.AdaptiveColors(colors.Blue, colors.Indigo)).
			Foreground(colors.AdaptiveColors(colors.White, colors.Sand))

	DefaultDescriptionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colors.AdaptiveColor(colors.Text))
)
