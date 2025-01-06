package main

import (
	"log"
	"net/http"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/net/html"
)

const (
	URL       = "https://old.reddit.com"
	MAX_POSTS = 10
)

type (
	postsMsg    []post
	getPostsErr error
)

type htmlNode struct {
	*html.Node
}

func (n htmlNode) GetAttr(key string) string {
	for _, attr := range n.Attr {
		if attr.Key != key {
			continue
		}

		return attr.Val
	}

	return ""
}

func (n htmlNode) Classes() []string {
	classes := n.GetAttr("class")
	return strings.Fields(classes)
}

func (n htmlNode) ClassContains(c string) bool {
	classes := n.Classes()
	return slices.Contains(classes, c)
}

func (n htmlNode) Text() string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			return c.Data
		}
	}

	return ""
}

func (n htmlNode) TagEquals(tag string) bool {
	return n.Type == html.ElementNode && n.Data == tag
}

func (n htmlNode) NodeEquals(tag, class string) bool {
	return n.TagEquals(tag) && n.ClassContains(class)
}

func getPosts() tea.Msg {
	client := &http.Client{}
	res, err := client.Get(URL)
	if err != nil {
		return getPostsErr(err)
	}

	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Fatal("Could not html parse reddit home page")
	}

	var posts []post
	posts = getPostsHelper(doc, posts)
	return postsMsg(posts)
}

func getPostsHelper(n *html.Node, posts []post) []post {
	for c := range n.Descendants() {
		node := htmlNode{c}
		if node.NodeEquals("div", "entry") {
			p := createPost(node)
			posts = append(posts, p)

			if len(posts) >= MAX_POSTS {
				break
			}
		}
	}

	return posts
}

func createPost(n htmlNode) post {
	var p post
	for c := range n.Descendants() {
		cNode := htmlNode{c}

		if cNode.NodeEquals("a", "title") {
			p.title = cNode.Text()
			p.postUrl = cNode.GetAttr("href")
		} else if cNode.NodeEquals("a", "author") {
			p.author = cNode.Text()
		} else if cNode.NodeEquals("a", "subreddit") {
			p.subreddit = cNode.Text()
		} else if cNode.NodeEquals("a", "comments") {
			p.commentsUrl = cNode.GetAttr("href")
		} else if cNode.NodeEquals("time", "live-timestamp") {
			p.friendlyDate = cNode.Text()
		}
	}

	return p
}
