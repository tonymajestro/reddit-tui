package header

import (
	"fmt"
	"reddittui/components/colors"

	"github.com/charmbracelet/lipgloss"
)

var headerContainerStyle = lipgloss.NewStyle().MarginBottom(2)

var (
	titleStyle = lipgloss.NewStyle().
			MarginBottom(1).
			Padding(0, 2).
			Height(1).
			Background(colors.AdaptiveColors(colors.Blue, colors.Indigo)).
			Foreground(colors.AdaptiveColor(colors.Sand))

	defaultDescriptionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colors.AdaptiveColor(colors.Text))
)

func trim(s string, w int) string {
	if len(s) <= w {
		return s
	}

	return fmt.Sprintf("%s...", s[:w-3])
}
