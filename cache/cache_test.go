package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"xtcli/consts"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTest initialises the cache with a test-specific provider name (derived
// from the test name) and a 24-hour TTL. It returns a cleanup function that
// clears the cache directory and restores the original package-level state.
func setupTest(t *testing.T) {
	t.Helper()

	// Save previous state so tests are hermetic even when run in sequence.
	prevProvider := providerName
	prevTTL := cacheTTL

	// Derive a filesystem-safe provider name from the test name.
	safe := strings.NewReplacer("/", "_", " ", "_", ":", "_").Replace(t.Name())
	Initialize(safe, 24)

	t.Cleanup(func() {
		_ = Clear()
		providerName = prevProvider
		cacheTTL = prevTTL
	})
}

// --- Initialize / GetCachePath ---

func TestInitialize_SetsProviderAndTTL(t *testing.T) {
	defer func(p string, d time.Duration) { providerName = p; cacheTTL = d }(providerName, cacheTTL)

	Initialize("myprovider", 6)

	assert.Equal(t, "myprovider", providerName)
	assert.Equal(t, 6*time.Hour, cacheTTL)
}

func TestGetCachePath_ContainsProviderName(t *testing.T) {
	setupTest(t)

	p := GetCachePath()
	assert.True(t, filepath.IsAbs(p))
	assert.Contains(t, p, providerName)
	assert.Contains(t, p, consts.CACHE_DIR_NAME)
	assert.Contains(t, p, consts.CONFIG_DIR_NAME)
}

// --- writeCacheFile / readCacheFile ---

func TestWriteReadCacheFile_RoundTrip(t *testing.T) {
	setupTest(t)

	path := cacheFilePath("test", "roundtrip")
	payload := []Category{{ID: 1, Name: "Alpha"}, {ID: 2, Name: "Beta"}}

	require.NoError(t, writeCacheFile(path, payload))

	var got []Category
	assert.True(t, readCacheFile(path, &got))
	require.Len(t, got, 2)
	assert.Equal(t, payload[0], got[0])
	assert.Equal(t, payload[1], got[1])
}

func TestReadCacheFile_FileNotFound(t *testing.T) {
	setupTest(t)

	path := cacheFilePath("test", "nonexistent")
	var dest []Category
	assert.False(t, readCacheFile(path, &dest))
}

func TestReadCacheFile_CorruptJSON(t *testing.T) {
	setupTest(t)

	path := cacheFilePath("test", "corrupt")
	dir := filepath.Dir(path)
	require.NoError(t, os.MkdirAll(dir, 0755))
	require.NoError(t, os.WriteFile(path, []byte(`not valid json`), 0600))

	var dest []Category
	assert.False(t, readCacheFile(path, &dest))
}

func TestReadCacheFile_ExpiredTTL(t *testing.T) {
	setupTest(t)

	path := cacheFilePath("test", "expired")

	// Write a cache file whose timestamp is older than the current TTL.
	payload, err := json.Marshal([]Category{{ID: 99, Name: "Old"}})
	require.NoError(t, err)

	oldFile := cachedFile{
		Timestamp: time.Now().Add(-25 * time.Hour),
		Data:      json.RawMessage(payload),
	}
	data, err := json.MarshalIndent(oldFile, "", "    ")
	require.NoError(t, err)

	dir := filepath.Dir(path)
	require.NoError(t, os.MkdirAll(dir, 0755))
	require.NoError(t, os.WriteFile(path, data, 0600))

	var dest []Category
	assert.False(t, readCacheFile(path, &dest), "expired cache entry should not be returned")
}

func TestReadCacheFile_ZeroTimestamp(t *testing.T) {
	setupTest(t)

	path := cacheFilePath("test", "zero_ts")

	payload, err := json.Marshal([]Category{{ID: 1, Name: "Test"}})
	require.NoError(t, err)

	// A zero Timestamp must be treated as invalid/stale.
	zeroFile := cachedFile{
		Timestamp: time.Time{},
		Data:      json.RawMessage(payload),
	}
	data, err := json.MarshalIndent(zeroFile, "", "    ")
	require.NoError(t, err)

	dir := filepath.Dir(path)
	require.NoError(t, os.MkdirAll(dir, 0755))
	require.NoError(t, os.WriteFile(path, data, 0600))

	var dest []Category
	assert.False(t, readCacheFile(path, &dest), "zero timestamp should be treated as stale")
}

func TestReadCacheFile_PayloadTypeMismatch(t *testing.T) {
	setupTest(t)

	path := cacheFilePath("test", "mismatch")

	// Write a numeric array – incompatible with []Category.
	payload, err := json.Marshal([]int{1, 2, 3})
	require.NoError(t, err)

	cf := cachedFile{
		Timestamp: time.Now(),
		Data:      json.RawMessage(payload),
	}
	data, err := json.MarshalIndent(cf, "", "    ")
	require.NoError(t, err)

	dir := filepath.Dir(path)
	require.NoError(t, os.MkdirAll(dir, 0755))
	require.NoError(t, os.WriteFile(path, data, 0600))

	// Unmarshal into a struct slice; JSON number arrays decode into structs as
	// zero values rather than erroring, so the read itself may succeed. What
	// matters is we don't panic and we get a valid bool back.
	var dest []Category
	_ = readCacheFile(path, &dest) // must not panic
}

// --- Clear ---

func TestClear_RemovesCacheDir(t *testing.T) {
	setupTest(t)

	// Seed some data so the directory is created.
	SetCategories(consts.CATEGORY_TYPE_LIVE, []Category{{ID: 1, Name: "Test"}})

	cachePath := GetCachePath()
	_, err := os.Stat(cachePath)
	require.NoError(t, err, "cache directory should exist after write")

	require.NoError(t, Clear())

	_, err = os.Stat(cachePath)
	assert.True(t, os.IsNotExist(err), "cache directory should be gone after Clear")
}

func TestClear_NonExistentDir_NoError(t *testing.T) {
	setupTest(t)
	// No data written; directory does not exist.
	assert.NoError(t, Clear())
}

// --- GetCategories / SetCategories ---

func TestGetSetCategories_LiveRoundTrip(t *testing.T) {
	setupTest(t)

	cats := []Category{
		{ID: 1, Name: "Action"},
		{ID: 2, Name: "Comedy", Parent: 1},
	}
	SetCategories(consts.CATEGORY_TYPE_LIVE, cats)

	got, ok := GetCategories(consts.CATEGORY_TYPE_LIVE)
	require.True(t, ok)
	require.Len(t, got, 2)
	assert.Equal(t, cats[0], got[0])
	assert.Equal(t, cats[1], got[1])
}

func TestGetSetCategories_AllTypes(t *testing.T) {
	setupTest(t)

	for _, ct := range []consts.CategoryType{
		consts.CATEGORY_TYPE_LIVE,
		consts.CATEGORY_TYPE_VOD,
		consts.CATEGORY_TYPE_SERIES,
	} {
		cats := []Category{{ID: 10, Name: string(ct)}}
		SetCategories(ct, cats)

		got, ok := GetCategories(ct)
		require.Truef(t, ok, "expected cache hit for category type %s", ct)
		require.Lenf(t, got, 1, "unexpected length for category type %s", ct)
		assert.Equal(t, string(ct), got[0].Name)
	}
}

func TestGetCategories_Miss(t *testing.T) {
	setupTest(t)

	got, ok := GetCategories(consts.CATEGORY_TYPE_LIVE)
	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestGetCategories_EmptySlice_ReturnsMiss(t *testing.T) {
	setupTest(t)

	// An empty slice is cached but GetCategories requires len > 0 to be a hit.
	SetCategories(consts.CATEGORY_TYPE_LIVE, []Category{})

	got, ok := GetCategories(consts.CATEGORY_TYPE_LIVE)
	assert.False(t, ok)
	assert.Nil(t, got)
}

// --- GetStreams / SetStreams ---

func TestGetSetStreams_RoundTrip(t *testing.T) {
	setupTest(t)

	streams := []Stream{
		{ID: 1, Name: "BBC News", CategoryID: 5, Type: "live"},
		{ID: 2, Name: "CNN", CategoryID: 5, Type: "live"},
	}
	SetStreams(5, streams)

	got, ok := GetStreams(5)
	require.True(t, ok)
	require.Len(t, got, 2)
	assert.Equal(t, streams[0].ID, got[0].ID)
	assert.Equal(t, streams[1].Name, got[1].Name)
}

func TestGetStreams_Miss(t *testing.T) {
	setupTest(t)

	got, ok := GetStreams(999)
	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestGetSetStreams_MultipleCategoriesIsolated(t *testing.T) {
	setupTest(t)

	SetStreams(1, []Stream{{ID: 10, Name: "Stream A", CategoryID: 1}})
	SetStreams(2, []Stream{{ID: 20, Name: "Stream B", CategoryID: 2}})

	gotCat1, ok1 := GetStreams(1)
	gotCat2, ok2 := GetStreams(2)

	require.True(t, ok1)
	require.True(t, ok2)
	assert.Equal(t, int64(10), gotCat1[0].ID)
	assert.Equal(t, int64(20), gotCat2[0].ID)
}

// --- GetVODStreams / SetVODStreams ---

func TestGetSetVODStreams_RoundTrip(t *testing.T) {
	setupTest(t)

	streams := []Stream{
		{ID: 100, Name: "Inception", CategoryID: 3, Type: "movie", ContainerExtension: "mkv"},
	}
	SetVODStreams(3, streams)

	got, ok := GetVODStreams(3)
	require.True(t, ok)
	require.Len(t, got, 1)
	assert.Equal(t, "Inception", got[0].Name)
	assert.Equal(t, "mkv", got[0].ContainerExtension)
}

func TestGetVODStreams_Miss(t *testing.T) {
	setupTest(t)

	got, ok := GetVODStreams(404)
	assert.False(t, ok)
	assert.Nil(t, got)
}

// --- GetEPG / SetEPG ---

func TestGetSetEPG_RoundTrip(t *testing.T) {
	setupTest(t)

	start := time.Date(2024, 9, 1, 19, 0, 0, 0, time.UTC)
	stop := time.Date(2024, 9, 1, 20, 0, 0, 0, time.UTC)
	epgs := []EPG{
		{
			ID:             1,
			EPGID:          42,
			ChannelID:      "bbc.uk",
			Title:          "Panorama",
			Description:    "Investigation programme",
			Lang:           "en",
			HasArchive:     true,
			NowPlaying:     false,
			Start:          "2024-09-01 19:00:00",
			End:            "2024-09-01 20:00:00",
			StartTimestamp: start,
			StopTimestamp:  stop,
		},
	}
	SetEPG(77, epgs)

	got, ok := GetEPG(77)
	require.True(t, ok)
	require.Len(t, got, 1)
	assert.Equal(t, "Panorama", got[0].Title)
	assert.Equal(t, "bbc.uk", got[0].ChannelID)
	assert.True(t, got[0].HasArchive)
	assert.True(t, start.Equal(got[0].StartTimestamp))
	assert.True(t, stop.Equal(got[0].StopTimestamp))
}

func TestGetEPG_Miss(t *testing.T) {
	setupTest(t)

	got, ok := GetEPG(9999)
	assert.False(t, ok)
	assert.Nil(t, got)
}

// --- GetAllStreams / GetAllVODStreams ---

func TestGetAllStreams_AggregatesCategories(t *testing.T) {
	setupTest(t)

	SetStreams(1, []Stream{{ID: 1, Name: "A", CategoryID: 1}})
	SetStreams(2, []Stream{{ID: 2, Name: "B", CategoryID: 2}, {ID: 3, Name: "C", CategoryID: 2}})

	all := GetAllStreams()
	require.Len(t, all, 3)

	ids := make(map[int64]bool)
	for _, s := range all {
		ids[s.ID] = true
	}
	assert.True(t, ids[1])
	assert.True(t, ids[2])
	assert.True(t, ids[3])
}

func TestGetAllStreams_NoCacheDir_ReturnsNil(t *testing.T) {
	setupTest(t)
	// Nothing written, live dir doesn't exist.
	all := GetAllStreams()
	assert.Nil(t, all)
}

func TestGetAllVODStreams_AggregatesCategories(t *testing.T) {
	setupTest(t)

	SetVODStreams(10, []Stream{{ID: 10, Name: "Movie A", CategoryID: 10}})
	SetVODStreams(20, []Stream{{ID: 20, Name: "Movie B", CategoryID: 20}})

	all := GetAllVODStreams()
	require.Len(t, all, 2)

	ids := make(map[int64]bool)
	for _, s := range all {
		ids[s.ID] = true
	}
	assert.True(t, ids[10])
	assert.True(t, ids[20])
}

func TestGetAllVODStreams_NoCacheDir_ReturnsNil(t *testing.T) {
	setupTest(t)
	all := GetAllVODStreams()
	assert.Nil(t, all)
}

// --- Info ---

func TestInfo_EmptyCache(t *testing.T) {
	setupTest(t)

	// No data written; cache directory does not exist.
	info := Info()
	assert.Contains(t, info, GetCachePath())
	assert.Contains(t, info, "empty")
}

func TestInfo_WithData(t *testing.T) {
	setupTest(t)

	SetCategories(consts.CATEGORY_TYPE_LIVE, []Category{{ID: 1, Name: "Live"}})
	SetStreams(1, []Stream{{ID: 1, Name: "Stream"}})

	info := Info()
	assert.Contains(t, info, GetCachePath())
	assert.Contains(t, info, "Files:")
	assert.Contains(t, info, "Total size:")
}

// --- Save / Load (no-ops) ---

func TestSave_IsNoOp(t *testing.T) {
	assert.NoError(t, Save())
}

func TestLoad_IsNoOp(t *testing.T) {
	assert.NoError(t, Load())
}
