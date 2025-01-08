package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

func NewSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	return s
}

func ViewSpinner(s spinner.Model, w, h int) string {
	var sb strings.Builder
	for range h / 2 {
		sb.WriteString("\n")
	}

	loadingMessage := fmt.Sprintf("%s %s", s.View(), "loading reddit.com...")
	var line strings.Builder
	for range w/2 - 12 {
		line.WriteString(" ")
	}
	line.WriteString(loadingMessage)

	sb.WriteString(line.String())
	return sb.String()
}
