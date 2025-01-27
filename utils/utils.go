package utils

import "fmt"

func NormalizeSubreddit(subreddit string) string {
	if len(subreddit) >= 2 && subreddit[:2] == "r/" {
		return subreddit
	}

	return fmt.Sprintf("r/%s", subreddit)
}

func Clamp(min, max, val int) int {
	if val < min {
		return min
	} else if val > max {
		return max
	}

	return val
}
