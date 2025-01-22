package client

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type RedditCommentsClient struct {
	client *http.Client
}

type Comment struct {
	Author    string
	Text      string
	Points    string
	Timestamp string
	Children  []*Comment
	Hidden    bool
	Depth     int
}

type Comments struct {
	PostTitle string
	Author    string
	Subreddit string
	Text      string
	Timestamp string
	Comments  []Comment
}

func (c Comment) Title() string {
	return formatDepth(c.Text, c.Depth)
}

func (c Comment) Description() string {
	desc := fmt.Sprintf("%s  by %s  %s", c.Points, c.Author, c.Timestamp)
	return formatDepth(desc, c.Depth)
}

func (c Comment) FilterValue() string {
	return c.Author
}

func (r RedditCommentsClient) GetComments(url string) (Comments, error) {
	var comments Comments

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
	commentsData.Text = getPostText(root)
	commentsData.Subreddit = getSubreddit(root)
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

func getTitle(root HtmlNode) string {
	for n := range root.FindDescendants("meta") {
		if n.GetAttr("property") == "og:title" {
			return n.GetAttr("content")
		}
	}

	return ""
}

func getPostText(root HtmlNode) string {
	var postText strings.Builder

	if linkListingNode, ok := root.FindDescendant("div", "sitetable", "linklisting"); ok {
		if mdNode, ok := linkListingNode.FindDescendant("div", "md"); ok {
			for n := range mdNode.Descendants() {
				nNode := HtmlNode{n}

				if len(nNode.Text()) > 0 {
					postText.WriteString(nNode.Text())
				} else if nNode.TagEquals("hr") {
					postText.WriteString("\n")
				}

				switch nNode.Tag() {
				case "p", "div", "li":
					postText.WriteString("\n")
				}
			}
		}
	}

	return ""
}

func getSubreddit(root HtmlNode) string {
	if spanNode, ok := root.FindDescendant("span", "pagename", "redditname"); ok {
		if subredditNode, ok := spanNode.FindDescendant("a"); ok {
			return subredditNode.Text()
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
