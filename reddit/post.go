package reddit

import "fmt"

type post struct {
	title         string
	author        string
	subreddit     string
	friendlyDate  string
	postUrl       string
	commentsUrl   string
	totalComments string
	totalLikes    string
}

func (p post) Title() string {
	return p.title
}

func (p post) Description() string {
	return fmt.Sprintf(" %s  %s  %s comments  %s", p.totalLikes, p.subreddit, p.totalComments, p.friendlyDate)
}

func (p post) FilterValue() string {
	return p.title
}
