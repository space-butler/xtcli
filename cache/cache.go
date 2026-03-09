package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"xtcli/consts"
)

var providerName string
var cacheTTL time.Duration

// Initialize sets up the cache with the specified provider and TTL in hours
func Initialize(provider string, ttlHours int) {
	providerName = provider
	cacheTTL = time.Duration(ttlHours) * time.Hour
}

// GetCachePath returns the root cache directory for the active provider
func GetCachePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, consts.CONFIG_DIR_NAME, consts.CACHE_DIR_NAME, providerName)
}

// filePath builds and ensures a cache file path: <root>/<subdir>/<filename>.json
func cacheFilePath(subdir, filename string) string {
	return filepath.Join(GetCachePath(), subdir, filename+".json")
}

// readCacheFile reads a cached JSON file and unmarshals the data payload into dest.
// Returns false if the file doesn't exist or the data is stale.
func readCacheFile(path string, dest interface{}) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var cf cachedFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return false
	}

	if cf.Timestamp.IsZero() || time.Since(cf.Timestamp) > cacheTTL {
		return false
	}

	if err := json.Unmarshal(cf.Data, dest); err != nil {
		return false
	}

	return true
}

// writeCacheFile marshals payload, wraps it with a timestamp, and writes to path.
func writeCacheFile(path string, payload interface{}) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	cf := cachedFile{
		Timestamp: time.Now(),
		Data:      raw,
	}

	data, err := json.MarshalIndent(cf, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// Clear removes all cached data for the active provider
func Clear() error {
	cachePath := GetCachePath()
	if _, err := os.Stat(cachePath); err == nil {
		return os.RemoveAll(cachePath)
	}
	return nil
}

// --- Categories ---

// GetCategories returns cached categories or nil if not available
func GetCategories(catType consts.CategoryType) ([]Category, bool) {
	path := cacheFilePath("categories", string(catType))
	var categories []Category
	if readCacheFile(path, &categories) && len(categories) > 0 {
		return categories, true
	}
	return nil, false
}

// SetCategories stores categories in the cache
func SetCategories(catType consts.CategoryType, categories []Category) {
	path := cacheFilePath("categories", string(catType))
	_ = writeCacheFile(path, categories)
}

// --- Live Streams ---

// GetStreams returns cached live streams for a category or nil if not available
func GetStreams(categoryID int64) ([]Stream, bool) {
	path := cacheFilePath("live", strconv.FormatInt(categoryID, 10))
	var streams []Stream
	if readCacheFile(path, &streams) {
		return streams, true
	}
	return nil, false
}

// SetStreams stores live streams in the cache for a category
func SetStreams(categoryID int64, streams []Stream) {
	path := cacheFilePath("live", strconv.FormatInt(categoryID, 10))
	_ = writeCacheFile(path, streams)
}

// --- VOD Streams ---

// GetVODStreams returns cached VOD streams for a category or nil if not available
func GetVODStreams(categoryID int64) ([]Stream, bool) {
	path := cacheFilePath("vod", strconv.FormatInt(categoryID, 10))
	var streams []Stream
	if readCacheFile(path, &streams) {
		return streams, true
	}
	return nil, false
}

// SetVODStreams stores VOD streams in the cache for a category
func SetVODStreams(categoryID int64, streams []Stream) {
	path := cacheFilePath("vod", strconv.FormatInt(categoryID, 10))
	_ = writeCacheFile(path, streams)
}

// --- EPG ---

// GetEPG returns cached EPG data for a stream or nil if not available
func GetEPG(streamID int64) ([]EPG, bool) {
	path := cacheFilePath("epg", strconv.FormatInt(streamID, 10))
	var epg []EPG
	if readCacheFile(path, &epg) {
		return epg, true
	}
	return nil, false
}

// SetEPG stores EPG data in the cache for a stream
func SetEPG(streamID int64, epg []EPG) {
	path := cacheFilePath("epg", strconv.FormatInt(streamID, 10))
	_ = writeCacheFile(path, epg)
}

// --- All Streams ---

// GetAllStreams returns all cached live streams from all category files
func GetAllStreams() []Stream {
	liveDir := filepath.Join(GetCachePath(), "live")
	return readAllStreamsFromDir(liveDir)
}

// GetAllVODStreams returns all cached VOD streams from all category files
func GetAllVODStreams() []Stream {
	vodDir := filepath.Join(GetCachePath(), "vod")
	return readAllStreamsFromDir(vodDir)
}

func readAllStreamsFromDir(dir string) []Stream {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var allStreams []Stream
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		var streams []Stream
		if readCacheFile(path, &streams) {
			allStreams = append(allStreams, streams...)
		}
	}
	return allStreams
}

// Save is a no-op in the new cache design (each Set* writes immediately).
// Kept for API compatibility during transition.
func Save() error {
	return nil
}

// Load is a no-op in the new cache design (reads happen on demand per file).
// Kept for API compatibility during transition.
func Load() error {
	return nil
}

// Info returns cache information for display
func Info() string {
	path := GetCachePath()
	info := fmt.Sprintf("Cache directory: %s\n", path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		info += "Status: empty (no cached data)\n"
		return info
	}

	var totalFiles int
	var totalSize int64
	_ = filepath.Walk(path, func(_ string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !fi.IsDir() {
			totalFiles++
			totalSize += fi.Size()
		}
		return nil
	})

	info += fmt.Sprintf("Files: %d\n", totalFiles)
	info += fmt.Sprintf("Total size: %.1f KB\n", float64(totalSize)/1024)
	return info
}
