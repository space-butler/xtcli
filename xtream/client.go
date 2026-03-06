package xtream

import (
	"fmt"
	"strconv"
	"time"
	"xtcli/cache"
	"xtcli/consts"

	xtreamcodes "github.com/space-butler/go.xtream-codes"
)

type client struct {
	xtreamClient *xtreamcodes.XtreamClient
	username     string
	password     string
	serverURL    string
}

var cli *client = nil

func IsInitialized() bool {
	return cli != nil
}

func Initialize(username, password, serverURL string) error {
	return InitializeWithCacheTTL(username, password, serverURL, 24)
}

func InitializeWithCacheTTL(username, password, serverURL string, cacheTTLHours int) error {
	if IsInitialized() {
		return nil
	}

	c, err := xtreamcodes.NewClient(username, password, serverURL)
	if err != nil {
		cli = nil
		return err
	}

	cli = &client{
		xtreamClient: c,
		username:     username,
		password:     password,
		serverURL:    serverURL,
	}

	// Initialize cache
	cache.Initialize(cacheTTLHours)
	if err := cache.Load(); err != nil {
		// Log warning but don't fail initialization
		fmt.Printf("Warning: Failed to load cache: %v\n", err)
	}

	return nil
}

func GetCategories(catType consts.CategoryType) ([]Category, error) {
	if !IsInitialized() {
		return nil, ErrClientNotInitialized
	}

	// Try to get from cache first
	if cachedCategories, found := cache.GetCategories(catType); found {
		// Convert cache.Category to xtream.Category
		result := make([]Category, len(cachedCategories))
		for i, c := range cachedCategories {
			result[i] = Category{
				ID:     c.ID,
				Name:   c.Name,
				Parent: c.Parent,
			}
		}
		return result, nil
	}

	// Fetch from server
	var categories []xtreamcodes.Category
	var err error = nil
	switch catType {
	case consts.CATEGORY_TYPE_LIVE:
		categories, err = cli.xtreamClient.GetLiveCategories()
	case consts.CATEGORY_TYPE_VOD:
		categories, err = cli.xtreamClient.GetVideoOnDemandCategories()
	default:
		return nil, ErrUnsupportedCategoryType
	}
	if err != nil {
		return nil, err
	}

	var result []Category
	for _, c := range categories {
		result = append(result, Category{
			ID:     int64(c.ID),
			Name:   c.Name,
			Parent: int64(c.Parent),
		})
	}

	// Convert to cache.Category and store in cache
	cacheCategories := make([]cache.Category, len(result))
	for i, c := range result {
		cacheCategories[i] = cache.Category{
			ID:     c.ID,
			Name:   c.Name,
			Parent: c.Parent,
		}
	}
	cache.SetCategories(catType, cacheCategories)
	cache.Save()

	return result, nil
}

// GetStreamsByCategory returns live streams for the given category ID.
// It uses the caching layer: reads from cache when valid, and on miss fetches from the provider
// then updates the cache (SetStreams + Save) before returning.
func GetStreamsByCategory(categoryID int64) ([]Stream, error) {
	if !IsInitialized() {
		return nil, ErrClientNotInitialized
	}

	// Try to get from cache first
	if cachedStreams, found := cache.GetStreams(categoryID); found {
		// Convert cache.Stream to xtream.Stream
		result := make([]Stream, len(cachedStreams))
		for i, s := range cachedStreams {
			result[i] = Stream{
				Added:              s.Added,
				CategoryID:         s.CategoryID,
				CategoryName:       s.CategoryName,
				ContainerExtension: s.ContainerExtension,
				CustomSid:          s.CustomSid,
				DirectSource:       s.DirectSource,
				EPGChannelID:       s.EPGChannelID,
				Icon:               s.Icon,
				ID:                 s.ID,
				Name:               s.Name,
				Number:             s.Number,
				Rating:             s.Rating,
				Type:               s.Type,
			}
		}
		return result, nil
	}

	// Get live streams for the category
	streams, err := cli.xtreamClient.GetLiveStreams(strconv.FormatInt(categoryID, 10))
	if err != nil {
		return nil, err
	}

	var result []Stream
	for _, s := range streams {
		var added time.Time
		if s.Added != nil {
			added = s.Added.Time
		}

		result = append(result, Stream{
			Added:              added,
			CategoryID:         int64(s.CategoryID),
			CategoryName:       s.CategoryName,
			ContainerExtension: s.ContainerExtension,
			CustomSid:          s.CustomSid,
			DirectSource:       s.DirectSource,
			EPGChannelID:       s.EPGChannelID,
			Icon:               s.Icon,
			ID:                 int64(s.ID),
			Name:               s.Name,
			Number:             int64(s.Number),
			Rating:             float32(s.Rating),
			Type:               s.Type,
		})
	}

	// Convert to cache.Stream and store in cache
	cacheStreams := make([]cache.Stream, len(result))
	for i, s := range result {
		cacheStreams[i] = cache.Stream{
			Added:              s.Added,
			CategoryID:         s.CategoryID,
			CategoryName:       s.CategoryName,
			ContainerExtension: s.ContainerExtension,
			CustomSid:          s.CustomSid,
			DirectSource:       s.DirectSource,
			EPGChannelID:       s.EPGChannelID,
			Icon:               s.Icon,
			ID:                 s.ID,
			Name:               s.Name,
			Number:             s.Number,
			Rating:             s.Rating,
			Type:               s.Type,
		}
	}
	cache.SetStreams(categoryID, cacheStreams)
	cache.Save()

	return result, nil
}

// GetVodStreamsByCategory returns VOD (video on demand) streams for the given category ID.
// It uses the caching layer: reads from cache when valid, and on miss fetches from the provider
// then updates the cache (SetVODStreams + Save) before returning.
func GetVodStreamsByCategory(categoryID int64) ([]Stream, error) {
	if !IsInitialized() {
		return nil, ErrClientNotInitialized
	}

	// Try to get from cache first
	if cachedStreams, found := cache.GetVODStreams(categoryID); found {
		result := make([]Stream, len(cachedStreams))
		for i, s := range cachedStreams {
			result[i] = Stream{
				Added:              s.Added,
				CategoryID:         s.CategoryID,
				CategoryName:       s.CategoryName,
				ContainerExtension: s.ContainerExtension,
				CustomSid:          s.CustomSid,
				DirectSource:       s.DirectSource,
				EPGChannelID:       s.EPGChannelID,
				Icon:               s.Icon,
				ID:                 s.ID,
				Name:               s.Name,
				Number:             s.Number,
				Rating:             s.Rating,
				Type:               s.Type,
			}
		}
		return result, nil
	}

	// Fetch VOD streams from server
	streams, err := cli.xtreamClient.GetVideoOnDemandStreams(strconv.FormatInt(categoryID, 10))
	if err != nil {
		return nil, err
	}

	var result []Stream
	for _, s := range streams {
		var added time.Time
		if s.Added != nil {
			added = s.Added.Time
		}

		result = append(result, Stream{
			Added:              added,
			CategoryID:         int64(s.CategoryID),
			CategoryName:       s.CategoryName,
			ContainerExtension: s.ContainerExtension,
			CustomSid:          s.CustomSid,
			DirectSource:       s.DirectSource,
			EPGChannelID:       s.EPGChannelID,
			Icon:               s.Icon,
			ID:                 int64(s.ID),
			Name:               s.Name,
			Number:             int64(s.Number),
			Rating:             float32(s.Rating),
			Type:               s.Type,
		})
	}

	// Store in cache
	cacheStreams := make([]cache.Stream, len(result))
	for i, s := range result {
		cacheStreams[i] = cache.Stream{
			Added:              s.Added,
			CategoryID:         s.CategoryID,
			CategoryName:       s.CategoryName,
			ContainerExtension: s.ContainerExtension,
			CustomSid:          s.CustomSid,
			DirectSource:       s.DirectSource,
			EPGChannelID:       s.EPGChannelID,
			Icon:               s.Icon,
			ID:                 s.ID,
			Name:               s.Name,
			Number:             s.Number,
			Rating:             s.Rating,
			Type:               s.Type,
		}
	}
	cache.SetVODStreams(categoryID, cacheStreams)
	cache.Save()

	return result, nil
}

func GetStream(streamID int64) (*Stream, error) {
	if !IsInitialized() {
		return nil, ErrClientNotInitialized
	}

	// Try to find in all cached streams first
	allStreams := cache.GetAllStreams()
	for _, s := range allStreams {
		if s.ID == streamID {
			return &Stream{
				Added:              s.Added,
				CategoryID:         s.CategoryID,
				CategoryName:       s.CategoryName,
				ContainerExtension: s.ContainerExtension,
				CustomSid:          s.CustomSid,
				DirectSource:       s.DirectSource,
				EPGChannelID:       s.EPGChannelID,
				Icon:               s.Icon,
				ID:                 s.ID,
				Name:               s.Name,
				Number:             s.Number,
				Rating:             s.Rating,
				Type:               s.Type,
			}, nil
		}
	}

	// Get all live streams to populate cache and find the stream
	streams, err := cli.xtreamClient.GetLiveStreams("")
	if err != nil {
		return nil, err
	}

	for _, s := range streams {
		if int64(s.ID) == streamID {
			var added time.Time
			if s.Added != nil {
				added = s.Added.Time
			}

			result := &Stream{
				Added:              added,
				CategoryID:         int64(s.CategoryID),
				CategoryName:       s.CategoryName,
				ContainerExtension: s.ContainerExtension,
				CustomSid:          s.CustomSid,
				DirectSource:       s.DirectSource,
				EPGChannelID:       s.EPGChannelID,
				Icon:               s.Icon,
				ID:                 int64(s.ID),
				Name:               s.Name,
				Number:             int64(s.Number),
				Rating:             float32(s.Rating),
				Type:               s.Type,
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("stream with ID %d not found", streamID)
}

func GetVodStream(streamID int64) (*Stream, error) {
	if !IsInitialized() {
		return nil, ErrClientNotInitialized
	}

	// Try to find in all cached VOD streams first
	allStreams := cache.GetAllStreams()
	for _, s := range allStreams {
		if s.ID == streamID && s.Type == "movie" {
			return &Stream{
				Added:              s.Added,
				CategoryID:         s.CategoryID,
				CategoryName:       s.CategoryName,
				ContainerExtension: s.ContainerExtension,
				CustomSid:          s.CustomSid,
				DirectSource:       s.DirectSource,
				EPGChannelID:       s.EPGChannelID,
				Icon:               s.Icon,
				ID:                 s.ID,
				Name:               s.Name,
				Number:             s.Number,
				Rating:             s.Rating,
				Type:               s.Type,
			}, nil
		}
	}

	// Fetch all VOD streams from server (no category filter)
	streams, err := cli.xtreamClient.GetVideoOnDemandStreams("")
	if err != nil {
		return nil, err
	}

	for _, s := range streams {
		if int64(s.ID) == streamID {
			var added time.Time
			if s.Added != nil {
				added = s.Added.Time
			}

			return &Stream{
				Added:              added,
				CategoryID:         int64(s.CategoryID),
				CategoryName:       s.CategoryName,
				ContainerExtension: s.ContainerExtension,
				CustomSid:          s.CustomSid,
				DirectSource:       s.DirectSource,
				EPGChannelID:       s.EPGChannelID,
				Icon:               s.Icon,
				ID:                 int64(s.ID),
				Name:               s.Name,
				Number:             int64(s.Number),
				Rating:             float32(s.Rating),
				Type:               s.Type,
			}, nil
		}
	}

	return nil, fmt.Errorf("VOD stream with ID %d not found", streamID)
}

func GetShortEPG(streamID int64, limit int) ([]EPG, error) {
	if !IsInitialized() {
		return nil, ErrClientNotInitialized
	}

	if limit <= 0 {
		limit = 4
	}

	// Try to get from cache first
	if cachedEPG, found := cache.GetEPG(streamID); found {
		// Convert cache.EPG to xtream.EPG
		result := make([]EPG, 0, len(cachedEPG))
		for i, e := range cachedEPG {
			if i >= limit {
				break
			}
			result = append(result, EPG{
				ChannelID:      e.ChannelID,
				Description:    e.Description,
				End:            e.End,
				EPGID:          e.EPGID,
				HasArchive:     e.HasArchive,
				ID:             e.ID,
				Lang:           e.Lang,
				NowPlaying:     e.NowPlaying,
				Start:          e.Start,
				StartTimestamp: e.StartTimestamp,
				StopTimestamp:  e.StopTimestamp,
				Title:          e.Title,
			})
		}
		return result, nil
	}

	epgData, err := cli.xtreamClient.GetShortEPG(strconv.FormatInt(streamID, 10), limit)
	if err != nil {
		return nil, err
	}

	var result []EPG
	for _, e := range epgData {
		// ConvertibleBoolean has unexported bool field, so we check the JSON representation
		hasArchive := false
		nowPlaying := false

		// The internal bool is not exported, but we can use type assertion or check the value indirectly
		// Since ConvertibleBoolean doesn't export the bool, we'll parse the JSON representation
		hasArchiveJSON, _ := e.HasArchive.MarshalJSON()
		nowPlayingJSON, _ := e.NowPlaying.MarshalJSON()

		if string(hasArchiveJSON) == "1" || string(hasArchiveJSON) == "\"1\"" {
			hasArchive = true
		}
		if string(nowPlayingJSON) == "1" || string(nowPlayingJSON) == "\"1\"" {
			nowPlaying = true
		}

		result = append(result, EPG{
			ChannelID:      e.ChannelID,
			Description:    string(e.Description),
			End:            strconv.FormatInt(int64(e.End), 10),
			EPGID:          int64(e.EPGID),
			HasArchive:     hasArchive,
			ID:             int64(e.ID),
			Lang:           e.Lang,
			NowPlaying:     nowPlaying,
			Start:          e.Start,
			StartTimestamp: e.StartTimestamp.Time,
			StopTimestamp:  e.StopTimestamp.Time,
			Title:          string(e.Title),
		})
	}

	// Convert to cache.EPG and store in cache
	cacheEPG := make([]cache.EPG, len(result))
	for i, e := range result {
		cacheEPG[i] = cache.EPG{
			ChannelID:      e.ChannelID,
			Description:    e.Description,
			End:            e.End,
			EPGID:          e.EPGID,
			HasArchive:     e.HasArchive,
			ID:             e.ID,
			Lang:           e.Lang,
			NowPlaying:     e.NowPlaying,
			Start:          e.Start,
			StartTimestamp: e.StartTimestamp,
			StopTimestamp:  e.StopTimestamp,
			Title:          e.Title,
		}
	}
	cache.SetEPG(streamID, cacheEPG)
	cache.Save()

	return result, nil
}

func GetStreamURL(streamID int64, format string) (string, error) {
	if !IsInitialized() {
		return "", ErrClientNotInitialized
	}

	// First, we need to fetch all live streams to populate the internal cache
	// The xtream-codes library requires this before GetStreamURL can work
	_, err := cli.xtreamClient.GetLiveStreams("")
	if err != nil {
		return "", err
	}

	url, err := cli.xtreamClient.GetStreamURL(int(streamID), format)
	if err != nil {
		return "", err
	}

	return url, nil
}

// GetVodStreamURL returns the direct stream URL for a VOD stream (movie).
// Format is typically "mp4" or "mkv". The Xtream API uses baseURL/movie/username/password/streamID.ext
func GetVodStreamURL(streamID int64, format string) (string, error) {
	if !IsInitialized() {
		return "", ErrClientNotInitialized
	}
	if format == "" {
		format = "mp4"
	}
	return fmt.Sprintf("%s/movie/%s/%s/%d.%s", cli.serverURL, cli.username, cli.password, streamID, format), nil
}

func GetXMLTVFile() ([]byte, error) {
	if !IsInitialized() {
		return nil, ErrClientNotInitialized
	}
	return cli.xtreamClient.GetXMLTV()
}

func GetServerInfo() (*xtreamcodes.ServerInfo, *xtreamcodes.UserInfo, error) {
	if !IsInitialized() {
		return nil, nil, ErrClientNotInitialized
	}

	// The server and user info are populated when the client is created
	return &cli.xtreamClient.ServerInfo, &cli.xtreamClient.UserInfo, nil
}
