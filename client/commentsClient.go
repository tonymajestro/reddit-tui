package client

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

const limitQueryParameter = "limit=500"

var postTextTrimRegex = regexp.MustCompile("\n\n\n+")

type RedditCommentsClient struct {
	client *http.Client
}

func (r RedditCommentsClient) GetComments(url string) (Comments, error) {
	var comments Comments

	url = addQueryParameter(url, limitQueryParameter)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return comments, err
	}
	req.Header.Add(userAgentKey, userAgentValue)

	res, err := r.client.Do(req)
	if err != nil {
		return comments, err
	}

	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		return comments, err
	}

	comments = createCommentsHelper(HtmlNode{doc})
	return comments, nil
}

func createCommentsHelper(root HtmlNode) Comments {
	var commentsData Comments
	var commentsList []Comment

	commentsData.PostTitle = getTitle(root)
	commentsData.PostText = getPostText(root)
	commentsData.PostAuthor = getPostAuthor(root)
	commentsData.PostTimestamp = getPostTimestamp(root)
	commentsData.Subreddit = getSubreddit(root)
	commentsData.PostPoints = getPostPoints(root)
	commentsData.Comments = createCommentsList(root, 0, commentsList)

	return commentsData
}

func createCommentsList(node HtmlNode, depth int, comments []Comment) []Comment {
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

		entryNode, ok := c.FindChild("div", "entry", "unvoted")
		if !ok {
			continue
		}

		comments = parseCommentNode(entryNode, depth, comments)

		if n, ok := c.FindChild("div", "child"); ok {
			comments = createCommentsList(n, depth+1, comments)
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
		comment.Text = strings.TrimSpace(getBlockHtml(usertextNode))
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

func getPostText(root HtmlNode) string {
	if linkListingNode, ok := root.FindDescendant("div", "sitetable", "linklisting"); ok {
		if mdNode, ok := linkListingNode.FindDescendant("div", "md"); ok {
			postText := getBlockHtml(mdNode)
			return postTextTrimRegex.ReplaceAllString(postText, "\n\n")
		}
	}

	return ""
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

func getBlockHtml(node HtmlNode) string {
	var content strings.Builder
	for child := range node.ChildNodes() {
		cNode := HtmlNode{child}

		var blockText strings.Builder
		collectBlockText(cNode, &blockText)
		content.WriteString(strings.TrimSpace(blockText.String()) + "\n")
	}

	return content.String()
}

func collectBlockText(blockNode HtmlNode, blockText *strings.Builder) {
	if blockNode.Type == html.TextNode {
		blockText.WriteString(blockNode.Data)
	} else if blockNode.Tag() == "li" || blockNode.Tag() == "ol" {
		blockText.WriteString("- ")
	}

	for child := range blockNode.ChildNodes() {
		collectBlockText(HtmlNode{child}, blockText)
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

func formatDepth(s string, depth int) string {
	var sb strings.Builder
	for range depth {
		sb.WriteString("  ")
	}
	sb.WriteString(s)

	return sb.String()
}

func addQueryParameter(url, query string) string {
	if strings.Contains(url, "?") {
		return fmt.Sprintf("%s&%s", url, query)
	}

	return fmt.Sprintf("%s?%s", url, query)
}
