package client

import (
	"fmt"
	"strings"
)

type Post struct {
	title         string
	author        string
	subreddit     string
	friendlyDate  string
	postUrl       string
	commentsUrl   string
	totalComments string
	totalLikes    string
}

type Posts []Post

func (p Post) Title() string {
	return fmt.Sprintf(" %s %s", p.totalLikes, p.title)
}

func (p Post) Description() string {
	var sb strings.Builder
	if strings.TrimSpace(p.subreddit) != "" {
		sb.WriteString(p.subreddit)
		sb.WriteString("  ")
	}

	if strings.TrimSpace(p.totalComments) == "" {
		fmt.Fprintf(&sb, "%d comments  ", 0)
	} else {
		fmt.Fprintf(&sb, "%s comments  ", p.totalComments)
	}

	fmt.Fprintf(&sb, "submitted %s by %s", p.friendlyDate, p.author)
	return sb.String()
}

func (p Post) FilterValue() string {
	return p.title
}

type Comment struct {
	author       string
	text         string
	points       string
	friendlyDate string
	children     []*Comment
	hidden       bool
	depth        int
}

func (c Comment) Title() string {
	return c.FormatDepth(fmt.Sprintf("%s  %s points  %s", c.author, c.points, c.friendlyDate))
}

func (c Comment) Description() string {
	return c.FormatDepth(c.text)
}

func (c Comment) FilterValue() string {
	return c.author
}

func (c Comment) FormatDepth(s string) string {
	var sb strings.Builder
	for range c.depth {
		sb.WriteString("  ")
	}
	sb.WriteString(s)

	return sb.String()
}
