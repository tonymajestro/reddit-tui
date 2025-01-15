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

type Comment struct {
	Author    string
	Text      string
	Points    string
	Timestamp string
	Children  []*Comment
	Hidden    bool
	Depth     int
}

type Comments []Comment

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

func (c Comment) Title() string {
	return formatDepth(c.Text, c.Depth)
}

func (c Comment) Description() string {
	desc := fmt.Sprintf("%s  by %s  %s", c.Points, c.Author, c.Timestamp)
	return formatDepth(desc, c.Depth)
}

func (c Comment) FilterValue() string {
	return c.Author
}

func formatDepth(s string, depth int) string {
	var sb strings.Builder
	for range depth {
		sb.WriteString("  ")
	}
	sb.WriteString(s)

	return sb.String()
}
