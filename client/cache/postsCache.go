package cache

import (
	"encoding/json"
	"io/fs"
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
	Clean()
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

func (f FilePostsCache) Clean() {
	filepath.WalkDir(f.CacheBaseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Error("cache error", "error", err)
			return nil
		}

		// Only clean up posts which are stored in root cache directory, skip everything else.
		if d.IsDir() {
			if path == f.CacheBaseDir {
				return nil
			} else {
				return filepath.SkipDir
			}
		} else if filepath.Ext(path) != ".json" {
			return nil
		}

		cacheFile, err := os.Open(path)
		if err != nil {
			slog.Debug("Could not open cache file.", "error", err)
			return nil
		}

		defer cacheFile.Close()

		var posts model.Posts
		decoder := json.NewDecoder(cacheFile)
		err = decoder.Decode(&posts)
		if err != nil {
			slog.Debug("Could not decode cached posts.", "error", err)
			return nil
		}

		// Delete cached posts file if it is expired
		if time.Now().After(posts.Expiry) {
			slog.Debug("Removing expired cache posts", "path", path)
			err = os.Remove(path)
			if err != nil {
				slog.Debug("Could not delete expired cache file", "error", err)
				return nil
			}
		}

		return nil
	})
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

func (f NoOpPostsCache) Clean() {
}
