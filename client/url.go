package client

import (
	"net/url"
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

func GetPostUrl(baseUrl, post string) (string, error) {
	parsed, err := url.Parse(post)
	if err != nil {
		// User passed in post ID, build URL from ID
		return url.JoinPath(baseUrl, post)
	}

	// User passed in url, use base URL instead of the one passed in
	return url.JoinPath(baseUrl, parsed.Path)
}
