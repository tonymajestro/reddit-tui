package components

import "fmt"

func normalizeSubreddit(subreddit string) string {
	if len(subreddit) >= 2 && subreddit[:2] == "r/" {
		return subreddit
	}

	return fmt.Sprintf("r/%s", subreddit)
}
