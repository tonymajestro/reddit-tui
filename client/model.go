package client

import (
	"fmt"
	"strings"
)

type Post struct {
	PostTitle     string
	Author        string
	Subreddit     string
	FriendlyDate  string
	PostUrl       string
	CommentsUrl   string
	TotalComments string
	TotalLikes    string
}

type Posts []Post

func (p Post) Title() string {
	return fmt.Sprintf(" %s %s", p.TotalLikes, p.PostTitle)
}

func (p Post) Description() string {
	var sb strings.Builder
	if strings.TrimSpace(p.Subreddit) != "" {
		sb.WriteString(p.Subreddit)
		sb.WriteString("  ")
	}

	if strings.TrimSpace(p.TotalComments) == "" {
		fmt.Fprintf(&sb, "%d comments  ", 0)
	} else {
		fmt.Fprintf(&sb, "%s comments  ", p.TotalComments)
	}

	fmt.Fprintf(&sb, "submitted %s by %s", p.FriendlyDate, p.Author)
	return sb.String()
}

func (p Post) FilterValue() string {
	return p.PostTitle
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

type Comments []Comment

func (c Comment) Title() string {
	return c.text
}

func (c Comment) Description() string {
	return fmt.Sprintf("%s points | by %s %s", c.points, c.author, c.friendlyDate)
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
