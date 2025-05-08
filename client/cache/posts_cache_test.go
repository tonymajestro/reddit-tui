package cache

import (
	"os"
	"path/filepath"
	"reddittui/client/common"
	"reddittui/model"
	"strings"
	"testing"
	"time"
)

const (
	testTitle         = "title"
	testDescription   = "description"
	testAuthor        = "author"
	testSubreddit     = "subreddit"
	testFriendlyDate  = "5 mins ago"
	testPostUrl       = "post.url"
	testCommentsUrl   = "comments.url"
	testTotalComments = "5 comments"
	testTotalLikes    = "10 likes"
	testIsHome        = false
	testAfter         = "after"
)

func TestPostsCacheHappyPath(t *testing.T) {
	cache := NewFilePostsCache(t.TempDir())

	expiry := time.Now().Add(200 * time.Millisecond).Round(time.Millisecond)
	expected := createTestPosts(expiry)
	err := cache.Put(expected, "happy")
	if err != nil {
		t.Fatalf("could not put posts in posts cache: %v", err)
	}

	got, err := cache.Get("happy")
	if err != nil {
		t.Fatalf("expected no errors getting posts from cache: %v", err)
	}

	assertPosts(expected, got, t)
}

func TestPostsCacheNotFound(t *testing.T) {
	cache := NewFilePostsCache(t.TempDir())

	if _, err := cache.Get("notfound.json"); err != common.ErrNotFound {
		t.Fatalf("expected to not find posts in cache")
	}
}

func TestPostsCacheCannotDecodePosts(t *testing.T) {
	cache := NewFilePostsCache(t.TempDir())

	file, err := os.CreateTemp(cache.CacheBaseDir, "cannotdecode*.json")
	if err != nil {
		t.Fatalf("could not create test posts file")
	}
	defer file.Close()

	file.WriteString("not valid json")

	// Cache adds the .json extension when fetching the file from the cache
	// Strip it here so we don't add it twice
	filename := filepath.Base(file.Name())
	cacheEntryName := strings.TrimSuffix(filename, filepath.Ext(filename))

	if _, err = cache.Get(cacheEntryName); err != common.ErrCannotDecodeCacheFile {
		t.Fatalf("expected cannot decode cache entry %s, got %v", filename, err)
	}
}

func TestPostsCacheCacheExpired(t *testing.T) {
	cache := NewFilePostsCache(t.TempDir())

	expiry := time.Now().Round(time.Millisecond)
	expected := createTestPosts(expiry)
	err := cache.Put(expected, "expired")
	if err != nil {
		t.Fatalf("could not put posts in posts cache: %v", err)
	}

	// Posts should be already expired by time we fetch them
	time.Sleep(100 * time.Millisecond)

	_, err = cache.Get("happy")
	if err == nil {
		t.Fatalf("expected no errors getting posts from cache: %v", err)
	}
}

func TestPostsCacheCleanCache(t *testing.T) {
	cache := NewFilePostsCache(t.TempDir())

	posts1 := createTestPosts(time.Now().Round(time.Millisecond))
	posts2 := createTestPosts(time.Now().Add(200 * time.Millisecond).Round(time.Millisecond))

	posts1.Subreddit = "subreddit1"
	posts2.Subreddit = "subreddit2"

	cache.Put(posts1, "subreddit1")
	cache.Put(posts2, "subreddit2")

	cache.Clean()

	if _, err := cache.Get("subreddit1"); err != common.ErrNotFound {
		t.Fatal("expected expired posts subreddit1 to be cleaned from cache")
	}

	gotPosts2, err := cache.Get("subreddit2")
	if err != nil {
		t.Fatalf("unexpected error fetching subreddit2 from cache: %v", err)
	}

	assertPosts(posts2, gotPosts2, t)

	time.Sleep(200 * time.Millisecond)
	cache.Clean()

	if _, err := cache.Get("subreddit2"); err != common.ErrNotFound {
		t.Fatal("expected expired posts subreddit1 to be cleaned from cache")
	}
}

func assertPosts(expected, got model.Posts, t *testing.T) {
	assertVal("After", expected.After, got.After, t)
	assertVal("Description", expected.Description, got.Description, t)
	assertVal("Subreddit", expected.Subreddit, got.Subreddit, t)
	assertVal("IsHome", expected.IsHome, got.IsHome, t)
	assertVal("Expiry", expected.Expiry, got.Expiry, t)

	if len(expected.Posts) != len(got.Posts) {
		t.Fatalf("expected %d posts but got %d:", len(expected.Posts), len(got.Posts))
	}

	for i, expectedPost := range expected.Posts {
		gotPost := got.Posts[i]
		assertPost(expectedPost, gotPost, t)
	}

	if t.Failed() {
		t.FailNow()
	}
}

func assertPost(expected, got model.Post, t *testing.T) {
	assertVal("PostTitle", expected.PostTitle, got.PostTitle, t)
	assertVal("Author", expected.Author, got.Author, t)
	assertVal("Subreddit", expected.Author, got.Author, t)
	assertVal("FriendlyDate", expected.FriendlyDate, got.FriendlyDate, t)
	assertVal("Expiry", expected.Expiry, got.Expiry, t)
	assertVal("PostUrl", expected.PostUrl, got.PostUrl, t)
	assertVal("CommentsUrl", expected.CommentsUrl, got.CommentsUrl, t)
	assertVal("TotalComments", expected.TotalComments, got.TotalComments, t)
	assertVal("TotalLikes", expected.TotalLikes, got.TotalLikes, t)
}

func assertVal[K comparable](context string, expected, got K, t *testing.T) {
	if expected != got {
		t.Errorf("assertion failed %s: for expected %v but got %v", context, expected, got)
	}
}

func createTestPost() model.Post {
	return model.Post{
		PostTitle:     testTitle,
		Author:        testAuthor,
		Subreddit:     testSubreddit,
		FriendlyDate:  testFriendlyDate,
		PostUrl:       testPostUrl,
		CommentsUrl:   testCommentsUrl,
		TotalComments: testTotalComments,
		TotalLikes:    testTotalLikes,
	}
}

func createTestPosts(expiry time.Time) model.Posts {
	post := createTestPost()
	posts := []model.Post{post}
	return model.Posts{
		Description: testDescription,
		Subreddit:   testSubreddit,
		IsHome:      false,
		Posts:       posts,
		After:       testAfter,
		Expiry:      expiry,
	}
}
