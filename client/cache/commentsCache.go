package cache

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"reddittui/client/common"
	"reddittui/model"
	"strings"
	"time"
)

type CommentsCache interface {
	Get(path string) (model.Comments, error)
	Put(comments model.Comments, path string) error
	Clean()
}

type FileCommentsCache struct {
	BaseUrl      string
	CacheBaseDir string
}

func NewFileCommentsCache(baseUrl, cacheDir string) FileCommentsCache {
	return FileCommentsCache{
		BaseUrl:      baseUrl,
		CacheBaseDir: cacheDir,
	}
}

// Get comments stored in cached file.
// Returns comments if they are present and not expired
func (f FileCommentsCache) Get(filename string) (comments model.Comments, err error) {
	subreddit := f.GetSubredditFromUrl(filename)
	if len(subreddit) == 0 {
		return comments, common.ErrNotFound
	}

	sanitizedFilename := url.QueryEscape(filename) + ".json"
	cacheFilePath := filepath.Join(f.CacheBaseDir, subreddit, sanitizedFilename)

	cacheFile, err := os.Open(cacheFilePath)
	if os.IsNotExist(err) {
		slog.Info("not found: " + cacheFilePath)
		return comments, common.ErrNotFound
	} else if err != nil {
		slog.Warn("Could not open cache file.", "error", err)
		return comments, common.ErrCannotOpenCacheFile
	}

	defer cacheFile.Close()

	decoder := json.NewDecoder(cacheFile)
	err = decoder.Decode(&comments)
	if err != nil {
		slog.Warn("Could not decode cached comments.", "error", err)
		return comments, common.ErrCannotDecodeCacheFile
	}

	if time.Now().After(comments.Expiry) {
		return comments, common.ErrCacheEntryExpired
	}

	return comments, nil
}

// Cache the comments, writing the contents to the given cache file
func (f FileCommentsCache) Put(comments model.Comments, filename string) error {
	subreddit := f.GetSubredditFromUrl(filename)

	cacheDir := filepath.Join(f.CacheBaseDir, subreddit)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		slog.Warn("Could not create subreddit comments cache directory", "error", err)
		return err
	}

	sanitizedFilename := url.QueryEscape(filename) + ".json"
	cacheFilePath := filepath.Join(cacheDir, sanitizedFilename)
	cacheFile, err := os.OpenFile(cacheFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		slog.Warn("Could not open cache file for encoding", "error", err)
		return common.ErrCannotOpenCacheFile
	}

	defer cacheFile.Close()

	commentsJson, err := json.MarshalIndent(comments, "", " ")
	if err != nil {
		slog.Warn("Could not encode comments for caching", "error", err)
		return common.ErrCannotEncodeCacheFile
	}

	_, err = cacheFile.Write(commentsJson)
	if err != nil {
		slog.Warn("Could not encode comments for caching", "error", err)
		return common.ErrCannotEncodeCacheFile
	}

	return nil
}

func (f FileCommentsCache) GetSubredditFromUrl(commentsUrl string) string {
	part := fmt.Sprintf("%s/r/", f.BaseUrl)
	if !strings.Contains(commentsUrl, part) {
		return ""
	}

	subreddit := commentsUrl[len(part):]
	if strings.Contains(subreddit, "/") {
		subreddit = subreddit[:strings.Index(subreddit, "/")]
	}

	return subreddit
}

func (f FileCommentsCache) Clean() {
	// First delete expired comments files in each subreddit directory
	filepath.WalkDir(f.CacheBaseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Ignore directories, only look at cache comment files
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		cacheFile, err := os.Open(path)
		if err != nil {
			slog.Debug("Could not open cache file.", "error", err)
			return nil
		}

		defer cacheFile.Close()

		var comments model.Comments
		decoder := json.NewDecoder(cacheFile)
		err = decoder.Decode(&comments)
		if err != nil {
			slog.Debug("Could not decode cached comments.", "error", err)
			return nil
		}

		// Delete cached comments file if it is expired
		if time.Now().After(comments.Expiry) {
			err = os.Remove(path)
			if err != nil {
				slog.Warn("Could not delete expired cache file")
				return nil
			}
		}

		return nil
	})

	// Wait for OS to process deleted files
	time.Sleep(500 * time.Millisecond)

	// Delete subreddit directories that are empty
	filepath.WalkDir(f.CacheBaseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Error("cache error", "error", err)
			return nil
		}

		if !d.IsDir() || path == f.CacheBaseDir {
			return nil
		}

		dir, err := os.Open(path)
		if err != nil {
			slog.Warn("Could not open cache dir", "error", err)
			return filepath.SkipDir
		}

		// Get list of filenames in subreddit directory
		files, err := dir.Readdirnames(0)
		if err != nil {
			slog.Warn("Could not read contests of subreddit cache dir", "error", err)
			return filepath.SkipDir
		}

		// Delete subreddit directory if it is empty
		if len(files) == 0 {
			err = os.Remove(path)
			if err != nil {
				slog.Warn("Could not delete empty cache dir")
				return filepath.SkipDir
			}
		}

		return nil
	})
}

type NoOpCommentsCache struct{}

func NewNoOpCommentsCache() NoOpCommentsCache {
	return NoOpCommentsCache{}
}

func (n NoOpCommentsCache) Get(cacheFilePath string) (comments model.Comments, err error) {
	return comments, common.ErrNotFound
}

func (n NoOpCommentsCache) Put(comments model.Comments, cacheFilePath string) error {
	return nil
}

func (n NoOpCommentsCache) Clean() {
}
