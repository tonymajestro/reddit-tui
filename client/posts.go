package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type RedditPostsClient struct {
	client *http.Client
}

type Post struct {
	PostTitle     string
	Author        string
	Subreddit     string
	FriendlyDate  string
	PostUrl       string
	CommentsUrl   string
	TotalComments string
	TotalLikes    string
}

type Posts struct {
	Description string
	Subreddit   string
	IsHome      bool
	Posts       []Post
}

func (p Post) Title() string {
	return fmt.Sprintf(" %s  %s", p.TotalLikes, p.PostTitle)
}

func (p Post) Description() string {
	var sb strings.Builder
	if strings.TrimSpace(p.Subreddit) != "" {
		sb.WriteString(p.Subreddit)
		sb.WriteString("  ")
	}

	if strings.TrimSpace(p.TotalComments) == "" {
		fmt.Fprintf(&sb, "%d comments  ", 0)
	} else {
		fmt.Fprintf(&sb, "%s comments  ", p.TotalComments)
	}

	fmt.Fprintf(&sb, "submitted %s by %s", p.FriendlyDate, p.Author)
	return sb.String()
}

func (p Post) FilterValue() string {
	return p.PostTitle
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
	var (
		reader io.Reader
		posts  Posts
	)

	if url == homeUrl && debug {
		reader, _ = os.Open("samples/home.html")
	} else {
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
		reader = res.Body
	}

	doc, err := html.Parse(reader)
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
