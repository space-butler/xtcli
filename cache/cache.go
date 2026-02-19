package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
	"xtcli/consts"
)

const CACHE_FILE_NAME = ".xtcli-cache"

var cacheData *CacheData = nil
var cacheTTL time.Duration

// Initialize sets up the cache with the specified TTL in hours
func Initialize(ttlHours int) {
	cacheTTL = time.Duration(ttlHours) * time.Hour
}

// GetCachePath returns the full path to the cache file
func GetCachePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, CACHE_FILE_NAME)
}

// Load reads the cache from disk
func Load() error {
	cachePath := GetCachePath()

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Cache doesn't exist, initialize empty cache
			cacheData = &CacheData{
				Timestamp:      time.Time{},
				LiveCategories: []Category{},
				VODCategories:  []Category{},
				Streams:        make(map[int64][]Stream),
				VODStreams:     make(map[int64][]Stream),
				EPGData:        make(map[int64][]EPG),
			}
			return nil
		}
		return err
	}

	var cache CacheData
	if err := json.Unmarshal(data, &cache); err != nil {
		// If unmarshal fails, start with empty cache
		cacheData = &CacheData{
			Timestamp:      time.Time{},
			LiveCategories: []Category{},
			VODCategories:  []Category{},
			Streams:        make(map[int64][]Stream),
			VODStreams:     make(map[int64][]Stream),
			EPGData:        make(map[int64][]EPG),
		}
		return nil
	}

	cacheData = &cache
	if cacheData.VODStreams == nil {
		cacheData.VODStreams = make(map[int64][]Stream)
	}
	return nil
}

// Save writes the cache to disk
func Save() error {
	if cacheData == nil {
		return nil
	}

	cacheData.Timestamp = time.Now()
	cachePath := GetCachePath()

	data, err := json.MarshalIndent(cacheData, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0600)
}

// IsStale checks if the cache has expired
func IsStale() bool {
	if cacheData == nil {
		return true
	}

	if cacheData.Timestamp.IsZero() {
		return true
	}

	return time.Since(cacheData.Timestamp) > cacheTTL
}

// Clear removes all cached data
func Clear() error {
	cacheData = &CacheData{
		Timestamp:      time.Time{},
		LiveCategories: []Category{},
		VODCategories:  []Category{},
		Streams:        make(map[int64][]Stream),
		VODStreams:     make(map[int64][]Stream),
		EPGData:        make(map[int64][]EPG),
	}

	cachePath := GetCachePath()
	if _, err := os.Stat(cachePath); err == nil {
		return os.Remove(cachePath)
	}
	return nil
}

// GetCategories returns cached categories or nil if not available
func GetCategories(catType consts.CategoryType) ([]Category, bool) {
	if cacheData == nil || IsStale() {
		return nil, false
	}

	switch catType {
	case consts.CATEGORY_TYPE_LIVE:
		if len(cacheData.LiveCategories) > 0 {
			return cacheData.LiveCategories, true
		}
	case consts.CATEGORY_TYPE_VOD:
		if len(cacheData.VODCategories) > 0 {
			return cacheData.VODCategories, true
		}
	}

	return nil, false
}

// SetCategories stores categories in the cache
func SetCategories(catType consts.CategoryType, categories []Category) {
	if cacheData == nil {
		cacheData = &CacheData{
			LiveCategories: []Category{},
			VODCategories:  []Category{},
			Streams:        make(map[int64][]Stream),
			VODStreams:     make(map[int64][]Stream),
			EPGData:        make(map[int64][]EPG),
		}
	}

	switch catType {
	case consts.CATEGORY_TYPE_LIVE:
		cacheData.LiveCategories = categories
	case consts.CATEGORY_TYPE_VOD:
		cacheData.VODCategories = categories
	}
}

// GetStreams returns cached streams for a category or nil if not available
func GetStreams(categoryID int64) ([]Stream, bool) {
	if cacheData == nil || IsStale() {
		return nil, false
	}

	if streams, ok := cacheData.Streams[categoryID]; ok {
		return streams, true
	}

	return nil, false
}

// SetStreams stores streams in the cache for a category
func SetStreams(categoryID int64, streams []Stream) {
	if cacheData == nil {
		cacheData = &CacheData{
			LiveCategories: []Category{},
			VODCategories:  []Category{},
			Streams:        make(map[int64][]Stream),
			VODStreams:     make(map[int64][]Stream),
			EPGData:        make(map[int64][]EPG),
		}
	}

	if cacheData.Streams == nil {
		cacheData.Streams = make(map[int64][]Stream)
	}

	cacheData.Streams[categoryID] = streams
}

// GetVODStreams returns cached VOD streams for a category or nil if not available
func GetVODStreams(categoryID int64) ([]Stream, bool) {
	if cacheData == nil || IsStale() {
		return nil, false
	}

	if cacheData.VODStreams == nil {
		return nil, false
	}
	if streams, ok := cacheData.VODStreams[categoryID]; ok {
		return streams, true
	}

	return nil, false
}

// SetVODStreams stores VOD streams in the cache for a category
func SetVODStreams(categoryID int64, streams []Stream) {
	if cacheData == nil {
		cacheData = &CacheData{
			LiveCategories: []Category{},
			VODCategories:  []Category{},
			Streams:        make(map[int64][]Stream),
			VODStreams:     make(map[int64][]Stream),
			EPGData:        make(map[int64][]EPG),
		}
	}

	if cacheData.VODStreams == nil {
		cacheData.VODStreams = make(map[int64][]Stream)
	}

	cacheData.VODStreams[categoryID] = streams
}

// GetEPG returns cached EPG data for a stream or nil if not available
func GetEPG(streamID int64) ([]EPG, bool) {
	if cacheData == nil || IsStale() {
		return nil, false
	}

	if epg, ok := cacheData.EPGData[streamID]; ok {
		return epg, true
	}

	return nil, false
}

// SetEPG stores EPG data in the cache for a stream
func SetEPG(streamID int64, epg []EPG) {
	if cacheData == nil {
		cacheData = &CacheData{
			LiveCategories: []Category{},
			VODCategories:  []Category{},
			Streams:        make(map[int64][]Stream),
			VODStreams:     make(map[int64][]Stream),
			EPGData:        make(map[int64][]EPG),
		}
	}

	if cacheData.EPGData == nil {
		cacheData.EPGData = make(map[int64][]EPG)
	}

	cacheData.EPGData[streamID] = epg
}

// GetAllStreams returns all cached streams from all categories
func GetAllStreams() []Stream {
	if cacheData == nil || IsStale() {
		return nil
	}

	var allStreams []Stream
	for _, streams := range cacheData.Streams {
		allStreams = append(allStreams, streams...)
	}

	return allStreams
}
