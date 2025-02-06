package client

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"reddittui/client/cache"
	"reddittui/model"
	"reddittui/utils"
	"strconv"
	"strings"
	"time"
)

const (
	homeUrl              = "https://old.reddit.com"
	subredditUrl         = "https://old.reddit.com/r/"
	defaultTitle         = "reddit.com"
	userAgentKey         = "User-Agent"
	userAgentValue       = "Mozilla/5.0 (X11; Linux x86_64; rv:134.0) Gecko/20100101 Firefox/134.0"
	cacheControlHeader   = "Cache-Control"
	maxAge               = "max-age"
	commentsCacheDirName = "comments"
)

type RedditClient struct {
	postsClient    RedditPostsClient
	commentsClient RedditCommentsClient
	bypassCache    bool
}

func NewRedditClient(bypassCache bool) RedditClient {
	httpClient := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	postsCache, commentsCache := InitializeCaches(bypassCache)
	postsClient := RedditPostsClient{
		Client: httpClient,
		Cache:  postsCache,
	}
	commentsClient := RedditCommentsClient{
		Client: httpClient,
		Cache:  commentsCache,
	}

	return RedditClient{
		postsClient,
		commentsClient,
		bypassCache,
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

func InitializeCaches(bypassCache bool) (cache.PostsCache, cache.CommentsCache) {
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
	commentsCacheDir := filepath.Join(cacheDir, commentsCacheDirName)
	err = os.MkdirAll(commentsCacheDir, 0755)
	if err != nil {
		slog.Warn("Cannot create comments cache dir, skipping comments cache")
		return postsCache, cache.NewNoOpCommentsCache()
	}

	commentsCache := cache.NewFileCommentsCache(commentsCacheDir)
	return postsCache, commentsCache
}

func getMaxAge(res *http.Response) (maxAge time.Duration, err error) {
	cacheHeader := strings.Join(res.Header[cacheControlHeader], ";")
	parts := strings.Split(cacheHeader, "=")
	if len(parts) != 2 {
		return maxAge, errParsingCacheHeaders
	}

	maxAgeSeconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return maxAge, errParsingCacheHeaders
	}

	maxAge = time.Duration(maxAgeSeconds) * time.Second
	return maxAge, nil
}
