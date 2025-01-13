package client

import (
	"slices"
	"strings"

	"golang.org/x/net/html"
)

type HtmlNode struct {
	*html.Node
}

func (n HtmlNode) GetAttr(key string) string {
	for _, attr := range n.Attr {
		if attr.Key != key {
			continue
		}

		return attr.Val
	}

	return ""
}

func (n HtmlNode) Classes() []string {
	classes := n.GetAttr("class")
	return strings.Fields(classes)
}

func (n HtmlNode) ClassContains(classesToFind ...string) bool {
	for _, c := range classesToFind {
		if !slices.Contains(n.Classes(), c) {
			return false
		}
	}

	return true
}

func (n HtmlNode) Text() string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			return c.Data
		}
	}

	return ""
}

func (n HtmlNode) TagEquals(tag string) bool {
	return n.Type == html.ElementNode && n.Data == tag
}

func (n HtmlNode) NodeEquals(tag string, classes ...string) bool {
	return n.TagEquals(tag) && n.ClassContains(classes...)
}
