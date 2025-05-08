package cache

import (
	"fmt"
	"reddittui/client/common"
	"reddittui/model"
	"testing"
	"time"
)

const (
	testPostPoints       = "5 points"
	testPostText         = "text"
	testPostTimestamp    = "5 mins ago"
	testCommentAuthor    = "author"
	testCommentText      = "comments"
	testCommentPoints    = "5 points"
	testCommentTimestamp = "5 mins ago"
	testBaseUrl          = "old.reddit.com"
	testCommentDepth     = 0
)

func TestCommentsCacheHappyPath(t *testing.T) {
	cache := NewFileCommentsCache(testBaseUrl, t.TempDir())

	expiry := time.Now().Add(200 * time.Millisecond).Round(time.Millisecond)
	expected := createTestComments(expiry)

	commentsUrl := generateCommentsFileUrl(testSubreddit, "happy")
	err := cache.Put(expected, commentsUrl)
	if err != nil {
		t.Fatalf("could not put comments in comments cache: %v", err)
	}

	got, err := cache.Get(commentsUrl)
	if err != nil {
		t.Fatalf("expected no errors getting comments from cache: %v", err)
	}

	assertComments(expected, got, t)
}

func TestCommentsCacheCacheNotFound(t *testing.T) {
	cache := NewFileCommentsCache(testBaseUrl, t.TempDir())

	if _, err := cache.Get("notfound.json"); err != common.ErrNotFound {
		t.Fatalf("expected to not find comments in cache")
	}
}

func assertComments(expected, got model.Comments, t *testing.T) {
	assertVal("PostTitle", expected.PostTitle, got.PostTitle, t)
	assertVal("PostAuthor", expected.PostAuthor, got.PostAuthor, t)
	assertVal("Subreddit", expected.Subreddit, got.Subreddit, t)
	assertVal("PostPoints", expected.PostPoints, got.PostPoints, t)
	assertVal("PostText", expected.PostText, got.PostText, t)
	assertVal("PostTimestamp", expected.PostTimestamp, got.PostTimestamp, t)
	assertVal("Expiry", expected.Expiry, got.Expiry, t)

	if len(expected.Comments) != len(got.Comments) {
		t.Fatalf("expected %d comments but got %d:", len(expected.Comments), len(got.Comments))
	}

	for i, expectedComment := range expected.Comments {
		gotComment := got.Comments[i]
		assertComment(expectedComment, gotComment, t)
	}

	if t.Failed() {
		t.FailNow()
	}
}

func assertComment(expectedComment, gotComment model.Comment, t *testing.T) {
	assertVal("Author", expectedComment.Author, gotComment.Author, t)
	assertVal("Text", expectedComment.Text, gotComment.Text, t)
	assertVal("Points", expectedComment.Points, gotComment.Points, t)
	assertVal("Timestamp", expectedComment.Timestamp, gotComment.Timestamp, t)
	assertVal("Depth", expectedComment.Depth, gotComment.Depth, t)
}

func createTestComment() model.Comment {
	return model.Comment{
		Author:    testCommentAuthor,
		Text:      testCommentText,
		Points:    testCommentPoints,
		Timestamp: testCommentTimestamp,
		Depth:     testCommentDepth,
	}
}

func createTestComments(expiry time.Time) model.Comments {
	return model.Comments{
		PostTitle:     testTitle,
		PostAuthor:    testAuthor,
		Subreddit:     testSubreddit,
		PostPoints:    testPostPoints,
		PostText:      testPostUrl,
		PostTimestamp: testPostTimestamp,
		Expiry:        expiry,
		Comments:      []model.Comment{createTestComment()},
	}
}

func generateCommentsFileUrl(subreddit, filename string) string {
	return fmt.Sprintf("%s/r/%s/%s", testBaseUrl, subreddit, filename)
}
