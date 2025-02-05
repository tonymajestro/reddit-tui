package model

import (
	"fmt"
	"strings"
	"time"
)

type Comment struct {
	Author    string `json:"author"`
	Text      string `json:"text"`
	Points    string `json:"points"`
	Timestamp string `json:"timestamp"`
	Depth     int    `json:"depth"`
}

type Comments struct {
	PostTitle     string    `json:"title"`
	PostAuthor    string    `json:"author"`
	Subreddit     string    `json:"subreddit"`
	PostPoints    string    `json:"points"`
	PostText      string    `json:"text"`
	PostUrl       string    `json:"url"`
	PostTimestamp string    `json:"timestamp"`
	Expiry        time.Time `json:"expiry"`
	Comments      []Comment `json:"comments"`
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
	var results strings.Builder
	for range depth {
		results.WriteString("  ")
	}
	results.WriteString(s)

	return results.String()
}
