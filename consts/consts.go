package consts

// Exit codes
const (
	EXIT_NO_CONFIG = 1
)

// File Names
const (
	CONFIG_DIR_NAME  = ".xtcli"
	CONFIG_FILE_NAME = "config.json"
)

// Category Types
type CategoryType string

const (
	CATEGORY_TYPE_LIVE    CategoryType = "live"
	CATEGORY_TYPE_VOD     CategoryType = "vod"
	CATEGORY_TYPE_SERIES  CategoryType = "series"
	CATEGORY_TYPE_UNKNOWN CategoryType = "unknown"
)
