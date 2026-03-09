package cache

import (
	"encoding/json"
	"time"
)

// Category represents an IPTV category
type Category struct {
	ID     int64  `json:"category_id"`
	Name   string `json:"category_name"`
	Parent int64  `json:"parent_id"`
}

// Stream represents an IPTV stream
type Stream struct {
	Added              time.Time `json:"added"`
	CategoryID         int64     `json:"category_id"`
	CategoryName       string    `json:"category_name"`
	ContainerExtension string    `json:"container_extension"`
	CustomSid          string    `json:"custom_sid"`
	DirectSource       string    `json:"direct_source,omitempty"`
	EPGChannelID       string    `json:"epg_channel_id"`
	Icon               string    `json:"stream_icon"`
	ID                 int64     `json:"stream_id"`
	Name               string    `json:"name"`
	Number             int64     `json:"num"`
	Rating             float32   `json:"rating"`
	Type               string    `json:"stream_type"`
}

// EPG represents electronic program guide data
type EPG struct {
	ChannelID      string    `json:"channel_id"`
	Description    string    `json:"description"`
	End            string    `json:"end"`
	EPGID          int64     `json:"epg_id"`
	HasArchive     bool      `json:"has_archive"`
	ID             int64     `json:"id"`
	Lang           string    `json:"lang"`
	NowPlaying     bool      `json:"now_playing"`
	Start          string    `json:"start"`
	StartTimestamp time.Time `json:"start_timestamp"`
	StopTimestamp  time.Time `json:"stop_timestamp"`
	Title          string    `json:"title"`
}

// cachedFile wraps any cached payload with a timestamp for TTL checks
type cachedFile struct {
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}
