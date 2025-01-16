package client

import (
	"iter"
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
	var classes []string

	class := n.GetAttr("class")
	for _, c := range strings.Fields(class) {
		classes = append(classes, strings.TrimSpace(c))
	}

	return classes
}

func (n HtmlNode) Class() string {
	return n.GetAttr("class")
}

func (n HtmlNode) ClassContains(classesToFind ...string) bool {
	for _, c := range classesToFind {
		if !slices.Contains(n.Classes(), strings.TrimSpace(c)) {
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

func (n HtmlNode) Tag() string {
	return n.Data
}

func (n HtmlNode) TagEquals(tag string) bool {
	return n.Type == html.ElementNode && n.Data == tag
}

func (n HtmlNode) NodeEquals(tag string, classes ...string) bool {
	return n.TagEquals(tag) && n.ClassContains(classes...)
}

func (n HtmlNode) FindDescendant(tag string, classes ...string) (HtmlNode, bool) {
	var descendant HtmlNode

	for c := range n.Descendants() {
		descendant = HtmlNode{c}
		if len(classes) == 0 && descendant.TagEquals(tag) {
			return descendant, true
		} else if descendant.NodeEquals(tag, classes...) {
			return descendant, true
		}
	}

	return descendant, false
}

func (n HtmlNode) FindDescendants(tag string, classes ...string) iter.Seq[HtmlNode] {
	return func(yield func(HtmlNode) bool) {
		for c := range n.Descendants() {
			childNode := HtmlNode{c}

			if len(classes) == 0 && childNode.TagEquals(tag) {
				if !yield(childNode) {
					return
				}
			} else if childNode.NodeEquals(tag, classes...) {
				if !yield(childNode) {
					return
				}
			}
		}
	}
}

func (n HtmlNode) FindChild(tag string, classes ...string) (HtmlNode, bool) {
	var child HtmlNode

	for c := range n.ChildNodes() {
		child = HtmlNode{c}
		if len(classes) == 0 && child.TagEquals(tag) {
			return child, true
		} else if child.NodeEquals(tag, classes...) {
			return child, true
		}
	}

	return child, false
}

func (n HtmlNode) FindChildren(tag string, classes ...string) iter.Seq[HtmlNode] {
	return func(yield func(HtmlNode) bool) {
		for c := range n.ChildNodes() {
			childNode := HtmlNode{c}

			if len(classes) == 0 && childNode.TagEquals(tag) {
				if !yield(childNode) {
					return
				}
			} else if childNode.NodeEquals(tag, classes...) {
				if !yield(childNode) {
					return
				}
			}
		}
	}
}
