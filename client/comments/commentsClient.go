package comments

import (
	"log/slog"
	"net/http"
	"reddittui/client/cache"
	"reddittui/client/common"
	"reddittui/model"
	"reddittui/utils"
	"regexp"
	"time"

	"golang.org/x/net/html"
)

const defaultTtl = 1 * time.Hour

var postTextTrimRegex = regexp.MustCompile("\n\n\n+")

type RedditCommentsClient struct {
	BaseUrl string
	Client  *http.Client
	Cache   cache.CommentsCache
	Parser  CommentsParser
}

func NewRedditCommentsClient(baseUrl, serverType string, httpClient *http.Client, commentsCache cache.CommentsCache) RedditCommentsClient {
	var parser CommentsParser

	switch serverType {
	case "old":
		parser = OldRedditCommentsParser{}
	case "redlib":
		parser = RedlibCommentsParser{}
	default:
		panic("Unrecognized server type in configuration: " + serverType)
	}

	return RedditCommentsClient{
		BaseUrl: baseUrl,
		Client:  httpClient,
		Cache:   commentsCache,
		Parser:  parser,
	}
}

func (r RedditCommentsClient) GetComments(url string) (comments model.Comments, err error) {
	totalTimer := utils.NewTimer("total time to retrieve comments")
	defer totalTimer.StopAndLog()

	timer := utils.NewTimer("fetching comments from cache")
	comments, err = r.Cache.Get(url)
	if err == nil {
		// return cached data
		timer.StopAndLog()
		return comments, nil
	}
	timer.StopAndLog()

	urlWithLimit := common.AddQueryParameter(url, common.LimitQueryParameter)
	req, err := http.NewRequest("GET", urlWithLimit, nil)
	if err != nil {
		return comments, err
	}
	req.Header.Add(common.UserAgentHeaderKey, common.UserAgentHeaderValue)

	timer = utils.NewTimer("fetching comments from server")
	res, err := r.Client.Do(req)
	timer.StopAndLog("url", url)
	if err != nil {
		return comments, err
	}

	if res.StatusCode != http.StatusOK {
		slog.Error("Error fetching comments from server", "StatusCode", res.StatusCode)
		return comments, common.ErrNotFound
	}

	defer res.Body.Close()

	timer = utils.NewTimer("parsing comments html")
	doc, err := html.Parse(res.Body)
	timer.StopAndLog()
	if err != nil {
		return comments, err
	}

	timer = utils.NewTimer("converting comments html")
	comments = r.Parser.ParseComments(common.HtmlNode{Node: doc}, url)
	comments.Expiry = time.Now().Add(defaultTtl)
	timer.StopAndLog()

	timer = utils.NewTimer("putting comments in cache")
	r.Cache.Put(comments, url)
	timer.StopAndLog()

	return comments, nil
}
