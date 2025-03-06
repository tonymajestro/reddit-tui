package client

import (
	"fmt"
	"net/url"
	"strings"
)

func NormalizeBaseUrl(baseUrl string) (string, error) {
	parsed, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	if parsed.Scheme != "https" {
		parsed.Scheme = "https"
	}

	url := parsed.String()
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}

	return url, nil
}

func GetPostUrl(baseUrl, post string) string {
	if strings.Contains(post, "http") || strings.Contains(post, baseUrl) {
		index := strings.Index(post, baseUrl)
		if index == -1 {
			return post
		}

		rest := post[index+len(baseUrl):]
		return baseUrl + rest
	} else {
		// User passed in post ID, build URL from ID
		return fmt.Sprintf("%s/%s", baseUrl, post)
	}
}
