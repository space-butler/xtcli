package xtream

import "errors"

var (
	ErrClientNotInitialized    = errors.New("xtream client not initialized")
	ErrUnsupportedCategoryType = errors.New("unsupported category type")
)
