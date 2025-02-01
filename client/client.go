package client

import (
	"net/http"
	"time"
)

const (
	homeUrl        = "https://old.reddit.com"
	subredditUrl   = "https://old.reddit.com/r/"
	defaultTitle   = "reddit.com"
	userAgentKey   = "User-Agent"
	userAgentValue = "Mozilla/5.0 (X11; Linux x86_64; rv:134.0) Gecko/20100101 Firefox/134.0"
)

type RedditClient struct {
	postsClient    RedditPostsClient
	commentsClient RedditCommentsClient
}

func New() RedditClient {
	httpClient := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	postsClient := RedditPostsClient{httpClient}
	commentsClient := RedditCommentsClient{httpClient}
	return RedditClient{postsClient, commentsClient}
}

func (r RedditClient) GetHomePosts() (Posts, error) {
	return r.postsClient.GetHomePosts()
}

func (r RedditClient) GetSubredditPosts(subreddit string) (Posts, error) {
	return r.postsClient.GetSubredditPosts(subreddit)
}

func (r RedditClient) GetComments(url string) (Comments, error) {
	return r.commentsClient.GetComments(url)
}
