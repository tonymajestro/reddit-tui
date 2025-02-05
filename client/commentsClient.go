package client

import (
	"fmt"
	"log/slog"
	"net/http"
	"reddittui/client/cache"
	"reddittui/model"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var postTextTrimRegex = regexp.MustCompile("\n\n\n+")

type RedditCommentsClient struct {
	Client *http.Client
	Cache  cache.CommentsCache
}

func (r RedditCommentsClient) GetComments(url string) (comments model.Comments, err error) {
	comments, err = r.Cache.Get(url)
	if err == nil {
		// return cached data
		slog.Info("Cache hit")
		return comments, nil
	}

	urlWithLimit := addQueryParameter(url, limitQueryParameter)
	req, err := http.NewRequest("GET", urlWithLimit, nil)
	if err != nil {
		return comments, err
	}
	req.Header.Add(userAgentKey, userAgentValue)

	res, err := r.Client.Do(req)
	if err != nil {
		return comments, err
	}

	defer res.Body.Close()
	maxAge, err := getMaxAge(res)
	if err != nil {
		slog.Error("Error getting cache headers from response", "error", err.Error(), "url", url)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return comments, err
	}

	comments = parseComments(HtmlNode{doc}, url)
	comments.Expiry = time.Now().Add(maxAge)

	r.Cache.Put(comments, url)
	return comments, nil
}

func parseComments(root HtmlNode, url string) model.Comments {
	var commentsData model.Comments
	var commentsList []model.Comment

	commentsData.PostTitle = getTitle(root)
	commentsData.PostAuthor = getPostAuthor(root)
	commentsData.PostTimestamp = getPostTimestamp(root)
	commentsData.Subreddit = getSubreddit(root)
	commentsData.PostPoints = getPostPoints(root)
	commentsData.Comments = parseCommentsList(root, 0, commentsList)

	postText, postUrl := getPostContent(root)
	if postUrl == "" {
		// Self post
		postUrl = url
	}
	commentsData.PostText = postText
	commentsData.PostUrl = postUrl

	return commentsData
}

func parseCommentsList(node HtmlNode, depth int, comments []model.Comment) []model.Comment {
	var commentsNode HtmlNode

	commentsNode, ok := node.FindDescendant("div", "sitetable", "nestedlisting")
	if !ok {
		commentsNode, ok = node.FindDescendant("div", "sitetable", "listing")
		if !ok {
			return comments
		}
	}

	for c := range commentsNode.FindChildren("div", "thing", "comment") {
		if c.ClassContains("deleted") {
			// Skip deleted comments and their children
			// todo: figure out how to render these properly
			continue
		}

		entryNode, ok := c.FindChild("div", "entry")
		if !ok {
			continue
		}

		comments = parseCommentNode(entryNode, depth, comments)

		if n, ok := c.FindChild("div", "child"); ok {
			comments = parseCommentsList(n, depth+1, comments)
		}
	}

	return comments
}

func parseCommentNode(node HtmlNode, depth int, comments []model.Comment) []model.Comment {
	var comment model.Comment
	comment.Depth = depth

	if taglineNode, ok := node.FindChild("p", "tagline"); ok {
		if authorNode, ok := taglineNode.FindChild("a", "author"); ok {
			comment.Author = authorNode.Text()
		}

		// Default to 1 point if the comment is too new to show points
		points := "1 point"
		if likesNode, ok := taglineNode.FindChild("span", "score", "likes"); ok {
			points = likesNode.Text()
		}
		comment.Points = points

		if timestampNode, ok := taglineNode.FindChild("time", "live-timestamp"); ok {
			comment.Timestamp = timestampNode.Text()
		}
	}

	if usertextNode, ok := node.FindChild("form", "usertext"); ok {
		comment.Text = strings.TrimSpace(renderHtmlNode(usertextNode))
	}

	comments = append(comments, comment)
	return comments
}

func getTitle(root HtmlNode) string {
	for n := range root.FindDescendants("meta") {
		if n.GetAttr("property") == "og:title" {
			return n.GetAttr("content")
		}
	}

	return ""
}

func getPostContent(root HtmlNode) (content, url string) {
	if linkListingNode, ok := root.FindDescendant("div", "sitetable", "linklisting"); ok {
		// self post
		if mdNode, ok := linkListingNode.FindDescendant("div", "md"); ok {
			postText := renderHtmlNode(mdNode)
			content = postTextTrimRegex.ReplaceAllString(postText, "\n\n")
			return content, ""
		}
	}

	if entry, ok := root.FindDescendant("div", "entry", "unvoted"); ok {
		// link post
		if linkNode, ok := entry.FindDescendant("a", "title"); ok {
			url = linkNode.GetAttr("href")
			content := fmt.Sprintf("%s\n\n", hyperLinkStyle.Render(url))
			return content, url

		}
	}

	return "", ""
}

func getPostAuthor(root HtmlNode) string {
	if linkListingNode, ok := root.FindDescendant("div", "sitetable", "linklisting"); ok {
		if authorNode, ok := linkListingNode.FindDescendant("a", "author"); ok {
			return authorNode.Text()
		}
	}

	return ""
}

func getPostTimestamp(root HtmlNode) string {
	if linkListingNode, ok := root.FindDescendant("div", "sitetable", "linklisting"); ok {
		if timestampNode, ok := linkListingNode.FindDescendant("time", "live-timestamp"); ok {
			return timestampNode.Text()
		}
	}

	return ""
}

func renderHtmlNode(node HtmlNode) string {
	var content strings.Builder
	for child := range node.ChildNodes() {
		cNode := HtmlNode{child}

		var nodeResults strings.Builder
		renderHtmlNodeHelper(cNode, &nodeResults)
		content.WriteString(nodeResults.String())
		content.WriteString("\n")
	}

	return content.String()
}

func renderHtmlNodeHelper(node HtmlNode, results *strings.Builder) {
	if node.Type == html.TextNode {
		results.WriteString(node.Data)
	} else if node.Tag() == "a" {
		results.WriteString(renderAnchor(node))
		return
	} else if node.Tag() == "li" {
		results.WriteString(node.Text())
		return
	}

	for child := range node.ChildNodes() {
		renderHtmlNodeHelper(HtmlNode{child}, results)
	}
}

func getSubreddit(root HtmlNode) string {
	if spanNode, ok := root.FindDescendant("span", "pagename", "redditname"); ok {
		if subredditNode, ok := spanNode.FindDescendant("a"); ok {
			return subredditNode.Text()
		}
	}

	return ""
}

func getPostPoints(root HtmlNode) string {
	if linkListingNode, ok := root.FindDescendant("div", "sitetable", "linklisting"); ok {
		if likesNode, ok := linkListingNode.FindDescendant("div", "score", "likes"); ok {
			return likesNode.Text()
		}

		if unvotedNode, ok := linkListingNode.FindDescendant("div", "score", "unvoted"); ok {
			return unvotedNode.Text()
		}

		// Fallback to any score node
		if pointsNode, ok := linkListingNode.FindDescendant("div", "score"); ok {
			return pointsNode.Text()
		}
	}

	return ""
}
