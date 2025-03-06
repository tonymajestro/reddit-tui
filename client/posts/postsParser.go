package posts

import (
	"log/slog"
	"net/url"
	"reddittui/client/common"
	"reddittui/model"
	"strings"
)

type PostsParser interface {
	ParsePosts(common.HtmlNode) model.Posts
}

type OldRedditPostsParser struct{}

func (p OldRedditPostsParser) ParsePosts(root common.HtmlNode) model.Posts {
	var (
		posts       []model.Post
		description string
	)

	for d := range root.FindDescendants("div", "thing") {
		if d.ClassContains("promoted", "promotedlink") {
			// Skip ads and promotional content
			continue
		}

		post := p.parsePost(d)
		posts = append(posts, post)
	}

	for d := range root.FindDescendants("meta") {
		if d.GetAttr("name") == "description" {
			description = d.GetAttr("content")
		}
	}

	modelPosts := model.Posts{
		Posts:       posts,
		Description: description,
	}

	return modelPosts
}

func (p OldRedditPostsParser) parsePost(n common.HtmlNode) model.Post {
	var post model.Post
	for c := range n.Descendants() {
		cNode := common.HtmlNode{Node: c}

		if cNode.NodeEquals("a", "title") {
			post.PostTitle = cNode.Text()
			post.PostUrl = cNode.GetAttr("href")
		} else if cNode.NodeEquals("a", "author") {
			post.Author = cNode.Text()
		} else if cNode.NodeEquals("a", "subreddit") {
			post.Subreddit = cNode.Text()
		} else if cNode.NodeEquals("time", "live-timestamp") {
			post.FriendlyDate = cNode.Text()
		} else if cNode.NodeEquals("a", "comments") {
			post.CommentsUrl = cNode.GetAttr("href")
			post.TotalComments = strings.Fields(cNode.Text())[0]
		} else if cNode.NodeEquals("div", "likes") {
			post.TotalLikes = cNode.Text()
		}
	}

	return post
}

type RedlibParser struct {
	BaseUrl string
}

func (p RedlibParser) ParsePosts(root common.HtmlNode) model.Posts {
	var posts model.Posts

	for d := range root.FindDescendants("div", "post") {
		post := p.parsePost(d)
		posts.Posts = append(posts.Posts, post)
	}

	if descriptionNode, ok := root.FindDescendantById("p", "sub_description"); ok {
		posts.Description = descriptionNode.Text()
	}

	return posts
}

func (p RedlibParser) parsePost(n common.HtmlNode) model.Post {
	var post model.Post
	for c := range n.Descendants() {
		cNode := common.HtmlNode{Node: c}

		if cNode.NodeEquals("h2", "post_title") {
			for postTitleSubNode := range cNode.FindChildren("a") {
				post.PostTitle = postTitleSubNode.Text()
				commentsUrl, err := p.buildUrl(postTitleSubNode.GetAttr("href"))
				if err != nil {
					slog.Debug("Error parsing comments url", "error", err)
					continue
				}
				post.CommentsUrl = commentsUrl
			}
		} else if cNode.NodeEquals("a", "post_author") {
			post.Author = cNode.Text()
		} else if cNode.NodeEquals("a", "post_subreddit") {
			post.Subreddit = cNode.Text()
		} else if cNode.NodeEquals("span", "created") {
			post.FriendlyDate = cNode.Text()
		} else if cNode.NodeEquals("a", "post_comments") {
			commentsUrl, err := p.buildUrl(cNode.GetAttr("href"))
			if err != nil {
				slog.Debug("Error parsing comments url", "error", err)
				continue
			}

			post.CommentsUrl = commentsUrl
			post.TotalComments = cNode.GetAttr("title")
		} else if cNode.NodeEquals("div", "post_score") {
			post.TotalLikes = strings.TrimSpace(cNode.Text())
		}
	}

	return post
}

func (p RedlibParser) buildUrl(part string) (string, error) {
	return url.JoinPath(p.BaseUrl, part)
}
