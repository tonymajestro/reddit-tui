package cache

import (
	"encoding/json"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"reddittui/client/common"
	"reddittui/model"
	"time"
)

type PostsCache interface {
	Get(path string) (model.Posts, error)
	Put(posts model.Posts, cacheFilePath string) error
}

type FilePostsCache struct {
	CacheBaseDir string
}

func NewFilePostsCache(cacheDir string) FilePostsCache {
	return FilePostsCache{CacheBaseDir: cacheDir}
}

// Get posts stored in cached file.
// Returns posts if they are present and not expired
func (f FilePostsCache) Get(filename string) (posts model.Posts, err error) {
	sanitizedFilename := url.QueryEscape(filename) + ".json"
	cacheFilePath := filepath.Join(f.CacheBaseDir, sanitizedFilename)

	cacheFile, err := os.Open(cacheFilePath)
	if os.IsNotExist(err) {
		return posts, common.ErrNotFound
	} else if err != nil {
		slog.Warn("Could not open cache file.", "error", err)
		return posts, common.ErrCannotOpenCacheFile
	}

	defer cacheFile.Close()

	decoder := json.NewDecoder(cacheFile)
	err = decoder.Decode(&posts)
	if err != nil {
		slog.Warn("Could not decode cached posts.", "error", err)
		return posts, common.ErrCannotDecodeCacheFile
	}

	if time.Now().After(posts.Expiry) {
		return posts, common.ErrCacheEntryExpired
	}

	return posts, nil
}

// Cache the posts, writing the contents to the given cache file
func (f FilePostsCache) Put(posts model.Posts, filename string) error {
	sanitizedFilename := url.QueryEscape(filename) + ".json"
	cacheFilePath := filepath.Join(f.CacheBaseDir, sanitizedFilename)

	cacheFile, err := os.OpenFile(cacheFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		slog.Warn("Could not open cache file for encoding", "error", err)
		return common.ErrCannotOpenCacheFile
	}

	defer cacheFile.Close()

	postsJson, err := json.MarshalIndent(posts, "", " ")
	if err != nil {
		slog.Warn("Could not encode posts for caching", "error", err)
		return common.ErrCannotEncodeCacheFile
	}

	_, err = cacheFile.Write(postsJson)
	if err != nil {
		slog.Warn("Could not encode posts for caching", "error", err)
		return common.ErrCannotEncodeCacheFile
	}

	return nil
}

type NoOpPostsCache struct{}

func NewNoOpPostsCache() NoOpPostsCache {
	return NoOpPostsCache{}
}

func (n NoOpPostsCache) Get(cacheFilePath string) (posts model.Posts, err error) {
	return posts, common.ErrNotFound
}

func (n NoOpPostsCache) Put(posts model.Posts, cacheFilePath string) error {
	return nil
}
