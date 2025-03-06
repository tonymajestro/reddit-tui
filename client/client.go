package client

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"reddittui/client/cache"
	"reddittui/client/comments"
	"reddittui/client/common"
	"reddittui/client/posts"
	"reddittui/config"
	"reddittui/model"
	"reddittui/utils"
	"time"
)

type RedditClient struct {
	BaseUrl        string
	postsClient    posts.RedditPostsClient
	commentsClient comments.RedditCommentsClient
}

func NewRedditClient(configuration config.Config) RedditClient {
	baseUrl, err := NormalizeBaseUrl(configuration.Server.Domain)
	if err != nil {
		log.Fatalf("Could not parse reddit server url: %s", configuration.Server.Domain)
	}

	// Support legacy core.ClientTimeout configuration value, use the greater of the two
	timeoutSeconds := max(configuration.Core.ClientTimeout, configuration.Client.TimeoutSeconds)
	httpClient := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}

	postsCache, commentsCache := InitializeCaches(baseUrl, configuration.Core.BypassCache)
	postsClient := posts.NewRedditPostsClient(baseUrl, httpClient, postsCache, configuration)
	commentsClient := comments.NewRedditCommentsClient(baseUrl, configuration.Server.Type, httpClient, commentsCache)

	return RedditClient{
		baseUrl,
		postsClient,
		commentsClient,
	}
}

func (r RedditClient) GetHomePosts() (model.Posts, error) {
	return r.postsClient.GetHomePosts()
}

func (r RedditClient) GetSubredditPosts(subreddit string) (model.Posts, error) {
	return r.postsClient.GetSubredditPosts(subreddit)
}

func (r RedditClient) GetComments(url string) (model.Comments, error) {
	return r.commentsClient.GetComments(url)
}

func InitializeCaches(baseUrl string, bypassCache bool) (cache.PostsCache, cache.CommentsCache) {
	if bypassCache {
		return cache.NewNoOpPostsCache(), cache.NewNoOpCommentsCache()
	}

	// read cache dir from env var
	cacheDir, err := utils.GetCacheDir()
	if err != nil {
		slog.Warn("Cannot open cache dir, skipping cache")
		return cache.NewNoOpPostsCache(), cache.NewNoOpCommentsCache()
	}

	// ensure root cache dir exists
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		slog.Warn("Cannot create root cache dir, skipping cache")
		return cache.NewNoOpPostsCache(), cache.NewNoOpCommentsCache()
	}

	// use root cache dir for posts
	postsCache := cache.NewFilePostsCache(cacheDir)

	// ensure comments cache dir exists
	commentsCacheDir := filepath.Join(cacheDir, common.CommentsCacheDirName)
	err = os.MkdirAll(commentsCacheDir, 0755)
	if err != nil {
		slog.Warn("Cannot create comments cache dir, skipping comments cache")
		return postsCache, cache.NewNoOpCommentsCache()
	}

	commentsCache := cache.NewFileCommentsCache(baseUrl, commentsCacheDir)
	return postsCache, commentsCache
}
