package client

import "fmt"

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

func (p Post) Title() string {
	return p.title
}

func (p Post) Description() string {
	return fmt.Sprintf(" %s  %s  %s comments  %s", p.totalLikes, p.subreddit, p.totalComments, p.friendlyDate)
}

func (p Post) FilterValue() string {
	return p.title
}
