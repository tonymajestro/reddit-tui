package client

import (
	"fmt"
	"strings"
)

const baseUrl = "https://old.reddit.com"

func GetHomeUrl() string {
	return baseUrl
}

func GetSubredditUrl(subreddit string) string {
	return fmt.Sprintf("%s/r/%s", baseUrl, subreddit)
}

func GetPostUrl(post string) string {
	if strings.Contains(post, "http") || strings.Contains(post, "reddit.com") {
		// Use old.reddit.com url for loading post and comment data
		index := strings.Index(post, "reddit.com")
		if index == -1 {
			return post
		}

		rest := post[index+len("reddit.com"):]
		return baseUrl + rest
	} else {
		// User passed in post ID, build URL from ID
		return fmt.Sprintf("%s/%s", baseUrl, post)
	}
}
