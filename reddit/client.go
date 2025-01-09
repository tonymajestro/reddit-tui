package reddit

import (
	"log"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

const (
	homeUrl      = "https://old.reddit.com/r/all"
	subredditUrl = "https://old.reddit.com/r/"
)

type (
	postsMsg    []post
	getPostsErr error
)

type htmlNode struct {
	*html.Node
}

func GetHomePosts() ([]post, error) {
    return getPosts(homeUrl)
}

func GetSubredditPosts(subreddit string) ([]post, error) {
    url := subredditUrl + subreddit
    return getPosts(url)
}

func (n htmlNode) getAttr(key string) string {
	for _, attr := range n.Attr {
		if attr.Key != key {
			continue
		}

		return attr.Val
	}

	return ""
}

func (n htmlNode) classes() []string {
	classes := n.getAttr("class")
	return strings.Fields(classes)
}

func (n htmlNode) classContains(c string) bool {
	classes := n.classes()
	return slices.Contains(classes, c)
}

func (n htmlNode) text() string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			return c.Data
		}
	}

	return ""
}

func (n htmlNode) tagEquals(tag string) bool {
	return n.Type == html.ElementNode && n.Data == tag
}

func (n htmlNode) nodeEquals(tag, class string) bool {
	return n.tagEquals(tag) && n.classContains(class)
}

func getPosts(url string) ([]post, error) {
	client := &http.Client{}
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Fatal("Could not html parse reddit home page")
	}

	var posts []post
	posts = populatePosts(doc, posts)
	return posts, nil
}

func populatePosts(n *html.Node, posts []post) []post {
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

func createPost(n htmlNode) post {
	var p post
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
