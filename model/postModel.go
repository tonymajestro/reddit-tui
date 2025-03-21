package model

import (
	"fmt"
	"strings"
	"time"
)

type Post struct {
	PostTitle     string    `json:"title"`
	Author        string    `json:"author"`
	Subreddit     string    `json:"subreddit"`
	FriendlyDate  string    `json:"friendlyDate"`
	Expiry        time.Time `json:"expiry"`
	PostUrl       string    `json:"postUrl"`
	CommentsUrl   string    `json:"commentsUrl"`
	TotalComments string    `json:"totalComments"`
	TotalLikes    string    `json:"totalLikes"`
}

type Posts struct {
	Description string
	Subreddit   string
	IsHome      bool
	Posts       []Post
	After       string
	Expiry      time.Time
}

func (p Post) Title() string {
	return fmt.Sprintf(" %s  %s", p.TotalLikes, p.PostTitle)
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
