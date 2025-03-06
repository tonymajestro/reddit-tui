package comments

import (
	"fmt"
	"reddittui/client/common"
	"reddittui/model"
	"reddittui/utils"
	"strings"

	"golang.org/x/net/html"
)

type CommentsParser interface {
	ParseComments(common.HtmlNode, string) model.Comments
}

type OldRedditCommentsParser struct{}

func (p OldRedditCommentsParser) ParseComments(root common.HtmlNode, url string) model.Comments {
	var commentsData model.Comments
	var commentsList []model.Comment

	commentsData.PostTitle = p.getTitle(root)
	commentsData.PostAuthor = p.getPostAuthor(root)
	commentsData.PostTimestamp = p.getPostTimestamp(root)
	commentsData.Subreddit = p.getSubreddit(root)
	commentsData.PostPoints = p.getPostPoints(root)
	commentsData.Comments = p.parseCommentsList(root, 0, commentsList)

	postText, postUrl := p.getPostContent(root)
	if postUrl == "" {
		// Self post
		postUrl = url
	}
	commentsData.PostText = postText
	commentsData.PostUrl = postUrl

	return commentsData
}

func (p OldRedditCommentsParser) parseCommentsList(node common.HtmlNode, depth int, comments []model.Comment) []model.Comment {
	var commentsNode common.HtmlNode

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

		comment := p.parseCommentNode(entryNode, depth)
		comments = append(comments, comment)

		if n, ok := c.FindChild("div", "child"); ok {
			comments = p.parseCommentsList(n, depth+1, comments)
		}
	}

	return comments
}

func (p OldRedditCommentsParser) parseCommentNode(node common.HtmlNode, depth int) model.Comment {
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

	return comment
}

func (p OldRedditCommentsParser) getTitle(root common.HtmlNode) string {
	for n := range root.FindDescendants("meta") {
		if n.GetAttr("property") == "og:title" {
			return n.GetAttr("content")
		}
	}

	return ""
}

func (p OldRedditCommentsParser) getPostContent(root common.HtmlNode) (content, url string) {
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
			content := fmt.Sprintf("%s\n\n", common.HyperLinkStyle.Render(url))
			return content, url

		}
	}

	return "", ""
}

func (p OldRedditCommentsParser) getPostAuthor(root common.HtmlNode) string {
	if linkListingNode, ok := root.FindDescendant("div", "sitetable", "linklisting"); ok {
		if authorNode, ok := linkListingNode.FindDescendant("a", "author"); ok {
			return authorNode.Text()
		}
	}

	return ""
}

func (p OldRedditCommentsParser) getPostTimestamp(root common.HtmlNode) string {
	if linkListingNode, ok := root.FindDescendant("div", "sitetable", "linklisting"); ok {
		if timestampNode, ok := linkListingNode.FindDescendant("time", "live-timestamp"); ok {
			return timestampNode.Text()
		}
	}

	return ""
}

func (p OldRedditCommentsParser) getSubreddit(root common.HtmlNode) string {
	if spanNode, ok := root.FindDescendant("span", "pagename", "redditname"); ok {
		if subredditNode, ok := spanNode.FindDescendant("a"); ok {
			return subredditNode.Text()
		}
	}

	return ""
}

func (p OldRedditCommentsParser) getPostPoints(root common.HtmlNode) string {
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

type RedlibCommentsParser struct{}

func (p RedlibCommentsParser) ParseComments(root common.HtmlNode, url string) model.Comments {
	var (
		commentsData model.Comments
		commentsList []model.Comment
	)

	mainNode, ok := root.FindDescendant("main")
	if !ok {
		return commentsData
	}

	if headerNode, ok := mainNode.FindDescendant("div", "post", "highlighted"); ok {
		commentsData.PostAuthor = p.getPostAuthor(headerNode)
		commentsData.PostTimestamp = p.getPostTimestamp(headerNode)
		commentsData.Subreddit = p.getSubreddit(headerNode)
	}

	commentsData.PostTitle = p.getTitle(root)
	commentsData.PostPoints = p.getPostPoints(mainNode)
	commentsData.Comments = p.parseCommentsList(mainNode, 0, commentsList)

	postText, postUrl := p.getPostContent(mainNode)
	if postUrl == "" {
		// Self post
		postUrl = url
	}
	commentsData.PostText = postText
	commentsData.PostUrl = postUrl

	return commentsData
}

func (p RedlibCommentsParser) getTitle(root common.HtmlNode) string {
	titleNode, ok := root.FindDescendant("title")
	if !ok {
		return ""
	}

	// Strip subreddit from title
	title := titleNode.Text()
	index := strings.Index(title, "- r/")
	if index < 0 {
		return title
	}
	return strings.TrimSpace(title[:index])
}

func (p RedlibCommentsParser) getPostAuthor(root common.HtmlNode) string {
	authorNode, ok := root.FindDescendant("a", "post_author")
	if !ok {
		return ""
	}

	author := authorNode.Text()
	if len(author) > 2 && author[:2] == "u/" {
		author = author[2:]
	}

	return author
}

func (p RedlibCommentsParser) getPostTimestamp(root common.HtmlNode) string {
	timestampNode, ok := root.FindDescendant("span", "created")
	if !ok {
		return ""
	}

	return timestampNode.Text()
}

func (p RedlibCommentsParser) getSubreddit(root common.HtmlNode) string {
	subredditNode, ok := root.FindDescendant("a", "post_subreddit")
	if !ok {
		return ""
	}

	return subredditNode.Text()
}

func (p RedlibCommentsParser) getPostPoints(root common.HtmlNode) string {
	pointsNode, ok := root.FindDescendant("div", "post_score")
	if !ok {
		return ""
	}

	return strings.TrimSpace(pointsNode.Text())
}

func (p RedlibCommentsParser) getPostContent(root common.HtmlNode) (content, url string) {
	// self post
	if postBodyNode, ok := root.FindDescendant("div", "post_body"); ok {
		if mdNode, ok := postBodyNode.FindDescendant("div", "md"); ok {
			postText := renderHtmlNode(mdNode)
			content = postTextTrimRegex.ReplaceAllString(postText, "\n\n")
			return content, ""
		}
	}

	// link post
	for linkNode := range root.FindChildren("a") {
		if linkNode.GetAttr("id") == "post_url" {
			url = linkNode.GetAttr("href")
			content := fmt.Sprintf("%s\n\n", common.HyperLinkStyle.Render(url))
			return content, url
		}
	}

	return "", ""
}

func (p RedlibCommentsParser) parseCommentsList(root common.HtmlNode, depth int, comments []model.Comment) []model.Comment {
	for threadNode := range root.FindDescendants("div", "thread") {
		comments = p.parseThread(threadNode, depth, comments)
	}

	return comments
}

func (p RedlibCommentsParser) parseThread(root common.HtmlNode, depth int, comments []model.Comment) []model.Comment {
	commentNode, ok := root.FindChild("div", "comment")
	if !ok {
		return comments
	}

	comment := p.parseCommentNode(commentNode, depth)
	comments = append(comments, comment)

	if n, ok := commentNode.FindDescendant("blockquote", "replies"); ok {
		comments = p.parseThread(n, depth+1, comments)
	}

	return comments
}

func (p RedlibCommentsParser) parseCommentNode(node common.HtmlNode, depth int) model.Comment {
	var comment model.Comment
	comment.Depth = depth

	if leftNode, ok := node.FindDescendant("div", "comment_left"); ok {
		if scoreNode, ok := leftNode.FindChild("p", "comment_score"); ok {
			points := "1 point"
			if scoreNode.GetAttr("title") != "Hidden" {
				points = utils.GetSingularPlural(strings.TrimSpace(scoreNode.Text()), "point", "points")
			}
			comment.Points = strings.TrimSpace(points)
		}
	}

	if rightNode, ok := node.FindDescendant("details", "comment_right"); ok {
		if authorNode, ok := rightNode.FindDescendant("a", "comment_author"); ok {
			author := authorNode.Text()
			if len(author) > 2 && author[:2] == "u/" {
				author = author[2:]
			}
			comment.Author = author
		}

		if timestampNode, ok := rightNode.FindDescendant("a", "created"); ok {
			comment.Timestamp = timestampNode.Text()
		}

		if commentBodyNode, ok := node.FindDescendant("div", "md"); ok {
			commentText := strings.TrimSpace(renderHtmlNode(commentBodyNode))
			comment.Text = postTextTrimRegex.ReplaceAllString(commentText, "\n\n")
		}
	}

	return comment
}

func renderHtmlNode(node common.HtmlNode) string {
	var content strings.Builder
	for child := range node.ChildNodes() {
		cNode := common.HtmlNode{Node: child}

		var nodeResults strings.Builder
		renderHtmlNodeHelper(cNode, &nodeResults)
		content.WriteString(nodeResults.String())
		content.WriteString("\n")
	}

	return content.String()
}

func renderHtmlNodeHelper(node common.HtmlNode, results *strings.Builder) {
	if node.Type == html.TextNode {
		results.WriteString(node.Data)
	} else if node.Tag() == "a" {
		results.WriteString(common.RenderAnchor(node))
		return
	} else if node.Tag() == "li" {
		results.WriteString(node.Text())
		return
	}

	for child := range node.ChildNodes() {
		renderHtmlNodeHelper(common.HtmlNode{Node: child}, results)
	}
}
