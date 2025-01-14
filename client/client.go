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
	comments = createComments(HtmlNode{doc}, comments)
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
	posts = createPosts(HtmlNode{doc}, posts)
	return posts, nil
}

func createPosts(root HtmlNode, posts []Post) []Post {
	for d := range root.FindDescendants("div", "thing") {
		post := createPost(d)
		posts = append(posts, post)
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

func createComments(root HtmlNode, comments []Comment) []Comment {
	file, _ := os.Create("debug.log")
	defer file.Close()

	commentsRootNode, ok := root.FindDescendant("div", "sitetable", "nestedlisting")
	if !ok {
		return comments
	}

	for c := range commentsRootNode.FindChildren("div", "thing", "comment") {
		if n, ok := c.FindChild("div", "entry", "unvoted"); ok {
			comments = parseRootCommentNode(n, comments)
		}
	}

	return comments
}

func parseRootCommentNode(node HtmlNode, comments []Comment) []Comment {
	var comment Comment

	if taglineNode, ok := node.FindChild("p", "tagline"); ok {
		if authorNode, ok := taglineNode.FindChild("a", "author"); ok {
			comment.Author = authorNode.Text()
		}

		if likesNode, ok := taglineNode.FindChild("span", "score", "likes"); ok {
			comment.Points = likesNode.Text()
		}

		if timestampNode, ok := taglineNode.FindChild("time", "live-timestamp"); ok {
			comment.Timestamp = timestampNode.Text()
		}
	}

	if usertextNode, ok := node.FindChild("form", "usertext"); ok {
		var usertext strings.Builder
		for n := range usertextNode.FindDescendants("p") {
			usertext.WriteString(n.Text())
			usertext.WriteRune(' ')
		}

		comment.Text = usertext.String()
	}

	comments = append(comments, comment)
	return comments
}

func createComment(node HtmlNode) Comment {
	var comment Comment
	for c := range node.Descendants() {
		cNode := HtmlNode{c}

		if cNode.NodeEquals("a", "author") {
			comment.Author = cNode.Text()
		} else if cNode.NodeEquals("div", "usertext-body") {
			comment.Text = getCommentText(cNode)
		} else if cNode.NodeEquals("span", "score", "likes") {
			comment.Points = cNode.Text()
		} else if cNode.NodeEquals("time", "live-timestamp") {
			comment.Timestamp = cNode.Text()
		}
	}

	return comment
}

func getCommentText(node HtmlNode) string {
	if c, ok := node.FindDescendant("p"); ok {
		return c.Text()
	} else {
		return ""
	}
}
