package common

import "errors"

var (
	ErrCacheEntryExpired     = errors.New("entry is expired")
	ErrNotFound              = errors.New("entry not found")
	ErrCannotOpenCacheFile   = errors.New("cannot open cache file")
	ErrCannotEncodeCacheFile = errors.New("cannot encode cache file")
	ErrCannotDecodeCacheFile = errors.New("cannot decode cache file")
)
