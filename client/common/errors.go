package common

import "errors"

var (
	ErrCacheEntryExpired     = errors.New("entry is expired")
	ErrCannotLoadPosts       = errors.New("cannot load posts")
	ErrNotFound              = errors.New("not found")
	ErrCannotOpenCacheFile   = errors.New("cannot open cache file")
	ErrCannotEncodeCacheFile = errors.New("cannot encode cache file")
	ErrCannotDecodeCacheFile = errors.New("cannot decode cache file")
	ErrParsingCacheHeaders   = errors.New("could not parse cache-control header")
)
