package posts

import (
	"fmt"
	"log/slog"
	"net/http"
	"reddittui/client/cache"
	"reddittui/client/common"
	"reddittui/config"
	"reddittui/model"
	"reddittui/utils"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type RedditPostsClient struct {
	BaseUrl          string
	CacheTtl         time.Duration
	Client           *http.Client
	Cache            cache.PostsCache
	Parser           PostsParser
	KeywordFilters   []string
	SubredditFilters []string
}

func NewRedditPostsClient(
	baseUrl string,
	httpClient *http.Client,
	postsCache cache.PostsCache,
	configuration config.Config,
) RedditPostsClient {
	var parser PostsParser

	switch strings.ToLower(configuration.Server.Type) {
	case "old":
		parser = OldRedditPostsParser{}
	case "redlib":
		parser = RedlibParser{baseUrl}
	default:
		panic("Unrecognized server type in configuration: " + configuration.Server.Type)
	}

	return RedditPostsClient{
		BaseUrl:          baseUrl,
		CacheTtl:         time.Duration(configuration.Client.CacheTtlSeconds) * time.Second,
		Client:           httpClient,
		Cache:            postsCache,
		Parser:           parser,
		KeywordFilters:   configuration.Filter.Keywords,
		SubredditFilters: configuration.Filter.Subreddits,
	}
}

func (r RedditPostsClient) GetHomePosts() (model.Posts, error) {
	timer := utils.NewTimer("total time to retrieve home posts")
	defer timer.StopAndLog()

	posts, err := r.tryGetCachedPosts(r.BaseUrl)
	posts.IsHome = true

	return posts, err
}

func (r RedditPostsClient) GetSubredditPosts(subreddit string) (model.Posts, error) {
	timer := utils.NewTimer("total time to retrieve subreddit posts")
	defer timer.StopAndLog()

	postsUrl := r.GetSubredditUrl(subreddit)
	posts, err := r.tryGetCachedPosts(postsUrl)
	posts.Subreddit = subreddit

	return posts, err
}

// Try to get posts from cache. If they are not present, fetch them and cache the results
func (r RedditPostsClient) tryGetCachedPosts(postsUrl string) (posts model.Posts, err error) {
	timer := utils.NewTimer("fetching posts from cache")
	posts, err = r.Cache.Get(postsUrl)
	if err == nil {
		// return cached data
		timer.StopAndLog()
		return r.filterPosts(posts), nil
	}
	timer.StopAndLog()

	timer = utils.NewTimer("getting posts from server")
	posts, err = r.getPosts(postsUrl)
	if err != nil {
		timer.StopAndLog()
		return posts, err
	}
	timer.StopAndLog()

	timer = utils.NewTimer("filtering posts")
	posts = r.filterPosts(posts)
	timer.StopAndLog()

	timer = utils.NewTimer("putting posts in cache")
	r.Cache.Put(posts, postsUrl)
	timer.StopAndLog()
	return posts, nil
}

func (r RedditPostsClient) getPosts(url string) (posts model.Posts, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return posts, err
	}

	req.Header.Add(common.UserAgentHeaderKey, common.UserAgentHeaderValue)

	timer := utils.NewTimer("fetching posts from server")
	res, err := r.Client.Do(req)
	timer.StopAndLog("url", url)

	if err != nil {
		return posts, err
	} else if res.StatusCode != http.StatusOK {
		// Treat all non-200s as 404s
		return posts, common.ErrCannotLoadPosts
	}

	defer res.Body.Close()

	timer = utils.NewTimer("parsing posts html")
	doc, err := html.Parse(res.Body)
	timer.StopAndLog()
	if err != nil {
		return posts, err
	}

	timer = utils.NewTimer("converting posts html")
	posts = r.Parser.ParsePosts(common.HtmlNode{Node: doc})
	timer.StopAndLog()
	if len(posts.Posts) == 0 {
		// if there are no posts, assume 404.
		// reddit redirect invalid subreddits requests to some search page instead of doing 404
		slog.Warn("Subreddit not found")
		return posts, common.ErrNotFound
	}

	posts.Expiry = time.Now().Add(r.CacheTtl)
	return posts, nil
}

func (r RedditPostsClient) filterPosts(posts model.Posts) model.Posts {
	var filteredPosts []model.Post

outer:
	for _, post := range posts.Posts {
		for _, keyword := range r.KeywordFilters {
			if strings.Contains(strings.ToLower(post.PostTitle), strings.ToLower(keyword)) {
				slog.Debug("filtering post", "title", post.PostTitle)
				continue outer
			}
		}

		for _, subreddit := range r.SubredditFilters {
			subreddit = utils.NormalizeSubreddit(subreddit)
			if strings.EqualFold(post.Subreddit, subreddit) {
				slog.Debug("filtering post", "title", post.PostTitle)
				continue outer
			}
		}

		filteredPosts = append(filteredPosts, post)
	}

	posts.Posts = filteredPosts
	return posts
}

func (r RedditPostsClient) GetSubredditUrl(subreddit string) string {
	return fmt.Sprintf("%s/r/%s", r.BaseUrl, subreddit)
}
