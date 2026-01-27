package xtream

import (
	"fmt"
	"strconv"
	"time"
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

	return nil
}

func GetCategories(catType consts.CategoryType) ([]Category, error) {
	if !IsInitialized() {
		return nil, ErrClientNotInitialized
	}

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
	return result, nil
}

func GetStreamsByCategory(categoryID int64) ([]Stream, error) {
	if !IsInitialized() {
		return nil, ErrClientNotInitialized
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

	return result, nil
}

func GetStream(streamID int64) (*Stream, error) {
	if !IsInitialized() {
		return nil, ErrClientNotInitialized
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

func GetShortEPG(streamID int64, limit int) ([]EPG, error) {
	if !IsInitialized() {
		return nil, ErrClientNotInitialized
	}

	if limit <= 0 {
		limit = 4
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

func GetXMLTVFile() ([]byte, error) {
	if !IsInitialized() {
		return nil, ErrClientNotInitialized
	}
	return cli.xtreamClient.GetXMLTV()
}
