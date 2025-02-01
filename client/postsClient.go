package client

import (
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type RedditPostsClient struct {
	client *http.Client
}

func (r RedditPostsClient) GetHomePosts() (Posts, error) {
	posts, err := r.getPosts(homeUrl)
	posts.IsHome = true

	return posts, err
}

func (r RedditPostsClient) GetSubredditPosts(subreddit string) (Posts, error) {
	url := subredditUrl + subreddit
	posts, err := r.getPosts(url)

	posts.Subreddit = subreddit
	posts.IsHome = false

	return posts, err
}

func (r RedditPostsClient) getPosts(url string) (Posts, error) {
	var posts Posts

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return posts, err
	}

	req.Header.Add(userAgentKey, userAgentValue)

	res, err := r.client.Do(req)
	if err != nil {
		return posts, err
	}

	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		return posts, err
	}

	return createPosts(HtmlNode{doc}), nil
}

func createPosts(root HtmlNode) Posts {
	var (
		posts       []Post
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

	return Posts{
		Posts:       posts,
		Description: description,
	}
}

func createPost(n HtmlNode) Post {
	var p Post
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
