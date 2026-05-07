package images

import "errors"

var (
	ErrUnsupportedImageURL = errors.New("url does not reference an image")
	ErrImageTooLarge       = errors.New("image exceeds maximum size")
	ErrCannotDecodeImage   = errors.New("cannot decode image")
	ErrCannotFetchImage    = errors.New("cannot fetch image")
)
