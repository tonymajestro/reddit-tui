package client

import (
	"slices"
	"strings"

	"golang.org/x/net/html"
)

type htmlNode struct {
	*html.Node
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
