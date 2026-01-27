package consts

// Exit codes
const (
	EXIT_NO_CONFIG = 1
)

// File Names
const (
	CONFIG_FILE_NAME = ".xtream-dump"
)

// Category Types
type CategoryType string

const (
	CATEGORY_TYPE_LIVE    CategoryType = "live"
	CATEGORY_TYPE_VOD     CategoryType = "vod"
	CATEGORY_TYPE_SERIES  CategoryType = "series"
	CATEGORY_TYPE_UNKNOWN CategoryType = "unknown"
)
