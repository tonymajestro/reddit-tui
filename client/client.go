package client

import (
	"fmt"
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
	defaultTitle   = "reddit.com"
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

func (r RedditClient) GetHomePosts() (Posts, error) {
	posts, err := r.getPosts(homeUrl)
	posts.Title = defaultTitle
	posts.IsHome = true

	return posts, err
}

func (r RedditClient) GetSubredditPosts(subreddit string) (Posts, error) {
	url := subredditUrl + subreddit
	posts, err := r.getPosts(url)

	posts.Title = fmt.Sprintf("r/%s", subreddit)
	posts.Subreddit = subreddit
	posts.IsHome = false

	return posts, err
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
	comments = createComments(HtmlNode{doc}, 0, comments)
	return comments, nil
}

func (r RedditClient) getPosts(url string) (Posts, error) {
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

func createComments(root HtmlNode, depth int, comments []Comment) []Comment {
	var commentsNode HtmlNode

	commentsNode, ok := root.FindDescendant("div", "sitetable", "nestedlisting")
	if !ok {
		commentsNode, ok = root.FindDescendant("div", "sitetable", "listing")
		if !ok {
			return comments
		}
	}

	for c := range commentsNode.FindChildren("div", "thing", "comment") {
		entryNode, ok := c.FindChild("div", "entry", "unvoted")
		if !ok {
			continue
		}

		comments = parseCommentNode(entryNode, depth, comments)

		if n, ok := c.FindChild("div", "child"); ok {
			comments = createComments(n, depth+1, comments)
		}
	}

	return comments
}

func parseCommentNode(node HtmlNode, depth int, comments []Comment) []Comment {
	var comment Comment
	comment.Depth = depth

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
