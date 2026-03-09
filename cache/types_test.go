package cache

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Category ---

func TestCategory_RoundTrip(t *testing.T) {
	original := Category{ID: 42, Name: "Sports", Parent: 5}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var got Category
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, original, got)
}

func TestCategory_JSONKeys(t *testing.T) {
	cat := Category{ID: 1, Name: "Movies", Parent: 0}

	data, err := json.Marshal(cat)
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))

	assert.Equal(t, float64(1), raw["category_id"])
	assert.Equal(t, "Movies", raw["category_name"])
	assert.Equal(t, float64(0), raw["parent_id"])
	assert.NotContains(t, raw, "ID")
	assert.NotContains(t, raw, "Name")
}

func TestCategory_ZeroValues(t *testing.T) {
	var cat Category
	require.NoError(t, json.Unmarshal([]byte(`{"category_id":0,"category_name":"","parent_id":0}`), &cat))
	assert.Equal(t, Category{}, cat)
}

func TestCategory_UnknownFieldsIgnored(t *testing.T) {
	data := `{"category_id":7,"category_name":"News","parent_id":0,"extra_field":"ignored"}`
	var cat Category
	require.NoError(t, json.Unmarshal([]byte(data), &cat))
	assert.Equal(t, int64(7), cat.ID)
	assert.Equal(t, "News", cat.Name)
}

// --- Stream ---

func TestStream_RoundTrip(t *testing.T) {
	added := time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC)
	original := Stream{
		Added:              added,
		CategoryID:         10,
		CategoryName:       "News",
		ContainerExtension: "ts",
		CustomSid:          "sid123",
		DirectSource:       "http://example.com/direct",
		EPGChannelID:       "epg.news.1",
		Icon:               "http://example.com/icon.png",
		ID:                 99,
		Name:               "CNN",
		Number:             1,
		Rating:             8.5,
		Type:               "live",
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var got Stream
	require.NoError(t, json.Unmarshal(data, &got))

	assert.Equal(t, original.CategoryID, got.CategoryID)
	assert.Equal(t, original.CategoryName, got.CategoryName)
	assert.Equal(t, original.ContainerExtension, got.ContainerExtension)
	assert.Equal(t, original.CustomSid, got.CustomSid)
	assert.Equal(t, original.DirectSource, got.DirectSource)
	assert.Equal(t, original.EPGChannelID, got.EPGChannelID)
	assert.Equal(t, original.Icon, got.Icon)
	assert.Equal(t, original.ID, got.ID)
	assert.Equal(t, original.Name, got.Name)
	assert.Equal(t, original.Number, got.Number)
	assert.Equal(t, original.Rating, got.Rating)
	assert.Equal(t, original.Type, got.Type)
	assert.True(t, original.Added.Equal(got.Added))
}

func TestStream_JSONKeys(t *testing.T) {
	s := Stream{ID: 5, Name: "Test", Icon: "icon.png", Type: "live", Number: 3, CategoryID: 2}

	data, err := json.Marshal(s)
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))

	assert.Contains(t, raw, "stream_id")
	assert.Contains(t, raw, "name")
	assert.Contains(t, raw, "stream_icon")
	assert.Contains(t, raw, "stream_type")
	assert.Contains(t, raw, "num")
	assert.Contains(t, raw, "category_id")
	assert.Contains(t, raw, "category_name")
	assert.Contains(t, raw, "epg_channel_id")
	assert.Contains(t, raw, "custom_sid")
	assert.Contains(t, raw, "container_extension")
	assert.Contains(t, raw, "rating")
	assert.Contains(t, raw, "added")
}

func TestStream_DirectSource_OmittedWhenEmpty(t *testing.T) {
	s := Stream{ID: 1, Name: "Test", DirectSource: ""}

	data, err := json.Marshal(s)
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))

	assert.NotContains(t, raw, "direct_source", "direct_source should be omitted when empty")
}

func TestStream_DirectSource_PresentWhenSet(t *testing.T) {
	s := Stream{ID: 1, Name: "Test", DirectSource: "http://direct.example.com/stream"}

	data, err := json.Marshal(s)
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))

	assert.Equal(t, "http://direct.example.com/stream", raw["direct_source"])
}

func TestStream_ZeroValues(t *testing.T) {
	data := `{"stream_id":0,"name":"","num":0,"rating":0,"stream_type":"","stream_icon":"","category_id":0,"category_name":"","epg_channel_id":"","custom_sid":"","container_extension":"","added":"0001-01-01T00:00:00Z"}`
	var s Stream
	require.NoError(t, json.Unmarshal([]byte(data), &s))
	assert.Equal(t, int64(0), s.ID)
	assert.Equal(t, "", s.Name)
	assert.Equal(t, float32(0), s.Rating)
}

// --- EPG ---

func TestEPG_RoundTrip(t *testing.T) {
	start := time.Date(2024, 6, 1, 20, 0, 0, 0, time.UTC)
	stop := time.Date(2024, 6, 1, 21, 30, 0, 0, time.UTC)
	original := EPG{
		ChannelID:      "cnn.us",
		Description:    "Evening news programme",
		End:            "2024-06-01 21:30:00",
		EPGID:          1001,
		HasArchive:     true,
		ID:             500,
		Lang:           "en",
		NowPlaying:     false,
		Start:          "2024-06-01 20:00:00",
		StartTimestamp: start,
		StopTimestamp:  stop,
		Title:          "Evening News",
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var got EPG
	require.NoError(t, json.Unmarshal(data, &got))

	assert.Equal(t, original.ChannelID, got.ChannelID)
	assert.Equal(t, original.Description, got.Description)
	assert.Equal(t, original.End, got.End)
	assert.Equal(t, original.EPGID, got.EPGID)
	assert.Equal(t, original.HasArchive, got.HasArchive)
	assert.Equal(t, original.ID, got.ID)
	assert.Equal(t, original.Lang, got.Lang)
	assert.Equal(t, original.NowPlaying, got.NowPlaying)
	assert.Equal(t, original.Start, got.Start)
	assert.Equal(t, original.Title, got.Title)
	assert.True(t, original.StartTimestamp.Equal(got.StartTimestamp))
	assert.True(t, original.StopTimestamp.Equal(got.StopTimestamp))
}

func TestEPG_JSONKeys(t *testing.T) {
	epg := EPG{ID: 1, Title: "Show", ChannelID: "ch1", HasArchive: true, NowPlaying: false}

	data, err := json.Marshal(epg)
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))

	assert.Contains(t, raw, "channel_id")
	assert.Contains(t, raw, "description")
	assert.Contains(t, raw, "end")
	assert.Contains(t, raw, "epg_id")
	assert.Contains(t, raw, "has_archive")
	assert.Contains(t, raw, "id")
	assert.Contains(t, raw, "lang")
	assert.Contains(t, raw, "now_playing")
	assert.Contains(t, raw, "start")
	assert.Contains(t, raw, "start_timestamp")
	assert.Contains(t, raw, "stop_timestamp")
	assert.Contains(t, raw, "title")
}

func TestEPG_HasArchive_Bool(t *testing.T) {
	epgTrue := EPG{ID: 1, HasArchive: true}
	epgFalse := EPG{ID: 2, HasArchive: false}

	dataTrue, err := json.Marshal(epgTrue)
	require.NoError(t, err)
	dataFalse, err := json.Marshal(epgFalse)
	require.NoError(t, err)

	var gotTrue, gotFalse EPG
	require.NoError(t, json.Unmarshal(dataTrue, &gotTrue))
	require.NoError(t, json.Unmarshal(dataFalse, &gotFalse))

	assert.True(t, gotTrue.HasArchive)
	assert.False(t, gotFalse.HasArchive)
}

// --- cachedFile ---

func TestCachedFile_RoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	cat := Category{ID: 7, Name: "Cached", Parent: 0}

	raw, err := json.Marshal(cat)
	require.NoError(t, err)

	original := cachedFile{Timestamp: now, Data: json.RawMessage(raw)}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var got cachedFile
	require.NoError(t, json.Unmarshal(data, &got))

	assert.True(t, original.Timestamp.Equal(got.Timestamp))

	var gotCat Category
	require.NoError(t, json.Unmarshal(got.Data, &gotCat))
	assert.Equal(t, cat, gotCat)
}

func TestCachedFile_JSONKeys(t *testing.T) {
	cf := cachedFile{Timestamp: time.Now(), Data: json.RawMessage(`{}`)}

	data, err := json.Marshal(cf)
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))

	assert.Contains(t, raw, "timestamp")
	assert.Contains(t, raw, "data")
}

func TestCachedFile_PreservesSlicePayload(t *testing.T) {
	streams := []Stream{
		{ID: 1, Name: "Alpha", Type: "live"},
		{ID: 2, Name: "Beta", Type: "live"},
	}

	raw, err := json.Marshal(streams)
	require.NoError(t, err)

	cf := cachedFile{Timestamp: time.Now(), Data: json.RawMessage(raw)}

	data, err := json.Marshal(cf)
	require.NoError(t, err)

	var got cachedFile
	require.NoError(t, json.Unmarshal(data, &got))

	var gotStreams []Stream
	require.NoError(t, json.Unmarshal(got.Data, &gotStreams))
	require.Len(t, gotStreams, 2)
	assert.Equal(t, int64(1), gotStreams[0].ID)
	assert.Equal(t, "Beta", gotStreams[1].Name)
}
