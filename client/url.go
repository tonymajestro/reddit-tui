package client

import "fmt"

const baseUrl = "https://old.reddit.com"

func GetHomeUrl() string {
	return baseUrl
}

func GetSubredditUrl(subreddit string) string {
	return fmt.Sprintf("%s/r/%s", baseUrl, subreddit)
}

func GetPostUrl(postId string) string {
	return fmt.Sprintf("%s/%s", baseUrl, postId)
}
