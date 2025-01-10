package client

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

const (
	homeUrl      = "https://old.reddit.com"
	subredditUrl = "https://old.reddit.com/r/"
)

type RedditClient struct {
	client *http.Client
}

func New() RedditClient {
	client := &http.Client{}
	return RedditClient{client}
}

func (r RedditClient) GetHomePosts() ([]Post, error) {
	return r.getPosts(homeUrl)
}

func (r RedditClient) GetSubredditPosts(subreddit string) ([]Post, error) {
	url := subredditUrl + subreddit
	return r.getPosts(url)
}

func (r RedditClient) getPosts(url string) ([]Post, error) {
	var reader io.Reader

	if url == homeUrl {
		reader, _ = os.Open("samples/home.html")
	} else {
		res, err := r.client.Get(url)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()
		reader = res.Body
	}

	doc, err := html.Parse(reader)
	if err != nil {
		log.Fatal("Could not html parse reddit home page")
	}

	var posts []Post
	posts = populatePosts(doc, posts)
	return posts, nil
}

func populatePosts(n *html.Node, posts []Post) []Post {
	if n == nil {
		return posts
	}

	for c := range n.Descendants() {
		node := htmlNode{c}
		if node.nodeEquals("div", "thing") {
			p := createPost(node)
			posts = append(posts, p)
		}
	}

	return posts
}

func createPost(n htmlNode) Post {
	var p Post
	for c := range n.Descendants() {
		cNode := htmlNode{c}

		if cNode.nodeEquals("a", "title") {
			p.title = cNode.text()
			p.postUrl = cNode.getAttr("href")
		} else if cNode.nodeEquals("a", "author") {
			p.author = cNode.text()
		} else if cNode.nodeEquals("a", "subreddit") {
			p.subreddit = cNode.text()
		} else if cNode.nodeEquals("time", "live-timestamp") {
			p.friendlyDate = cNode.text()
		} else if cNode.nodeEquals("a", "comments") {
			p.commentsUrl = cNode.getAttr("href")
			p.totalComments = strings.Fields(cNode.text())[0]
		} else if cNode.nodeEquals("div", "likes") {
			p.totalLikes = cNode.text()
		}
	}

	return p
}
