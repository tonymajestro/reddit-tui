package header

import (
	"fmt"
	"reddittui/components/colors"

	"github.com/charmbracelet/lipgloss"
)

var headerContainerStyle = lipgloss.NewStyle().MarginBottom(2)

var (
	defaultDescriptionStyle = lipgloss.NewStyle().Foreground(colors.Text).Bold(true)

	titleStyle = lipgloss.NewStyle().
			MarginBottom(1).
			Padding(0, 2).
			Height(1).
			Background(colors.Indigo).
			Foreground(colors.Sand)
)

func trim(s string, w int) string {
	if len(s) <= w {
		return s
	}

	return fmt.Sprintf("%s...", s[:w-3])
}
