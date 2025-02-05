package client

import (
	"errors"
	"log/slog"
	"net/http"
	"reddittui/client/cache"
	"reddittui/model"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var errParsingCacheHeaders = errors.New("could not parse cache-control header")

type RedditPostsClient struct {
	Client *http.Client
	Cache  cache.PostsCache
}

func (r RedditPostsClient) GetHomePosts() (model.Posts, error) {
	posts, err := r.tryGetCachedPosts(homeUrl)
	posts.IsHome = true
	return posts, err
}

func (r RedditPostsClient) GetSubredditPosts(subreddit string) (model.Posts, error) {
	postsUrl := subredditUrl + subreddit
	posts, err := r.tryGetCachedPosts(postsUrl)
	posts.Subreddit = subreddit

	return posts, err
}

// Try to get posts from cache. If they are not present, fetch them from reddit.com and
// cache the results
func (r RedditPostsClient) tryGetCachedPosts(postsUrl string) (posts model.Posts, err error) {
	posts, err = r.Cache.Get(postsUrl)
	if err == nil {
		// return cached data
		return posts, nil
	}

	posts, err = r.getPosts(postsUrl)
	if err != nil {
		return posts, err
	}

	r.Cache.Put(posts, postsUrl)
	return posts, nil
}

func (r RedditPostsClient) getPosts(url string) (posts model.Posts, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return posts, err
	}

	req.Header.Add(userAgentKey, userAgentValue)

	res, err := r.Client.Do(req)
	if err != nil {
		return posts, err
	}

	defer res.Body.Close()

	maxAge, err := getMaxAge(res)
	if err != nil {
		slog.Error("Error getting cache headers from response", "error", err.Error(), "url", url)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return posts, err
	}

	posts = createPosts(HtmlNode{doc})
	posts.Expiry = time.Now().Add(maxAge)
	return posts, nil
}

func createPosts(root HtmlNode) model.Posts {
	var (
		posts       []model.Post
		description string
	)

	for d := range root.FindDescendants("div", "thing") {
		if d.ClassContains("promoted", "promotedlink") {
			// Skip ads and promotional content
			continue
		}

		post := createPost(d)
		posts = append(posts, post)
	}

	for d := range root.FindDescendants("meta") {
		if d.GetAttr("name") == "description" {
			description = d.GetAttr("content")
		}
	}

	return model.Posts{
		Posts:       posts,
		Description: description,
	}
}

func createPost(n HtmlNode) model.Post {
	var p model.Post
	for c := range n.Descendants() {
		cNode := HtmlNode{c}

		if cNode.NodeEquals("a", "title") {
			p.PostTitle = cNode.Text()
			p.PostUrl = cNode.GetAttr("href")
		} else if cNode.NodeEquals("a", "author") {
			p.Author = cNode.Text()
		} else if cNode.NodeEquals("a", "subreddit") {
			p.Subreddit = cNode.Text()
		} else if cNode.NodeEquals("time", "live-timestamp") {
			p.FriendlyDate = cNode.Text()
		} else if cNode.NodeEquals("a", "comments") {
			p.CommentsUrl = cNode.GetAttr("href")
			p.TotalComments = strings.Fields(cNode.Text())[0]
		} else if cNode.NodeEquals("div", "likes") {
			p.TotalLikes = cNode.Text()
		}
	}

	return p
}
