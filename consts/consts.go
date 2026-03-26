package consts

// Exit codes
const (
	EXIT_NO_CONFIG = 1
)

const (
	CONFIG_DIR_NAME  = ".xtcli"
	CONFIG_FILE_NAME = "config.json"

	CACHE_DIR_NAME = "cache"
)

// Byte-unit conversion constants
const (
	BYTES_PER_KB = int64(1024)
	BYTES_PER_MB = int64(1024 * 1024)
	BYTES_PER_GB = int64(1024 * 1024 * 1024)
)

// Category Types
type CategoryType string

const (
	CATEGORY_TYPE_LIVE    CategoryType = "live"
	CATEGORY_TYPE_VOD     CategoryType = "vod"
	CATEGORY_TYPE_SERIES  CategoryType = "series"
	CATEGORY_TYPE_UNKNOWN CategoryType = "unknown"
)
