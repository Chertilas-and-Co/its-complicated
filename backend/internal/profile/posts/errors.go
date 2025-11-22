package profile

import "errors"

// Custom errors for profile operations
var (
	ErrPostNotFound = errors.New("post not found")
	ErrAlreadyLiked = errors.New("already liked")
	ErrNotLiked     = errors.New("not liked")
)
