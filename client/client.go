package client

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var debug = false

const (
	homeUrl        = "https://old.reddit.com"
	subredditUrl   = "https://old.reddit.com/r/"
	userAgentKey   = "User-Agent"
	userAgentValue = "Mozilla/5.0 (X11; Linux x86_64; rv:134.0) Gecko/20100101 Firefox/134.0"
)

type RedditClient struct {
	client *http.Client
}

func New() RedditClient {
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	return RedditClient{client}
}

func (r RedditClient) GetHomePosts() ([]Post, error) {
	return r.getPosts(homeUrl)
}

func (r RedditClient) GetSubredditPosts(subreddit string) ([]Post, error) {
	url := subredditUrl + subreddit
	return r.getPosts(url)
}

func (r RedditClient) GetComments(url string) ([]Comment, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(userAgentKey, userAgentValue)

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	var comments []Comment
	comments = createComments(doc, comments)
	return comments, nil
}

func (r RedditClient) getPosts(url string) ([]Post, error) {
	var reader io.Reader

	if url == homeUrl && debug {
		reader, _ = os.Open("samples/home.html")
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add(userAgentKey, userAgentValue)

		res, err := r.client.Do(req)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()
		reader = res.Body
	}

	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	var posts []Post
	posts = createPosts(doc, posts)
	return posts, nil
}

func createPosts(n *html.Node, posts []Post) []Post {
	for c := range n.Descendants() {
		node := HtmlNode{c}
		if node.NodeEquals("div", "thing") {
			post := createPost(node)
			posts = append(posts, post)
		}
	}

	return posts
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

func createComments(node *html.Node, comments []Comment) []Comment {
	for c := range node.Descendants() {
		node := HtmlNode{c}
		if node.NodeEquals("div", "comment") {
			comment := createComment(node)
			comments = append(comments, comment)
		}
	}

	return comments
}

func createComment(node HtmlNode) Comment {
	var comment Comment
	for c := range node.Descendants() {
		cNode := HtmlNode{c}

		if cNode.NodeEquals("a", "author") {
			comment.author = cNode.Text()
		} else if cNode.NodeEquals("div", "usertext-body") {
			comment.text = getCommentText(HtmlNode(cNode))
		} else if cNode.NodeEquals("span", "score", "likes") {
			comment.points = cNode.Text()
		} else if cNode.NodeEquals("time", "live-timestamp") {
			comment.friendlyDate = cNode.Text()
		}
	}

	return comment
}

func getCommentText(node HtmlNode) string {
	for c := range node.Descendants() {
		cNode := HtmlNode{c}
		if cNode.TagEquals("p") {
			return cNode.Text()
		}
	}

	return ""
}
