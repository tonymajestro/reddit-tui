package utils

import "fmt"

func NormalizeSubreddit(subreddit string) string {
	if len(subreddit) >= 2 && subreddit[:2] == "r/" {
		return subreddit
	}

	return fmt.Sprintf("r/%s", subreddit)
}

func TruncateString(s string, w int) string {
	if len(s) <= w {
		return s
	}

	return fmt.Sprintf("%s...", s[:w-3])
}

func Clamp(min, max, val int) int {
	if val < min {
		return min
	} else if val > max {
		return max
	}

	return val
}
