package reddit

import (
	"reddittui/client"

	"github.com/charmbracelet/bubbles/list"
)

type CommentsPage struct {
	redditClient client.RedditClient
	postsList    list.Model
}
