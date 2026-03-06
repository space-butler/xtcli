package config

type Favorite struct {
	Number   int    `json:"number"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	StreamID int64  `json:"stream_id"`
}

type Provider struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

type Config struct {
	DefaultProvider string     `json:"default_provider,omitempty"`
	Providers       []Provider `json:"providers,omitempty"`
	VlcPath         string     `json:"vlc_path,omitempty"`
	Favorites       []Favorite `json:"favorites,omitempty"`
	CacheTTL        int        `json:"cache_ttl,omitempty"` // Cache time-to-live in hours (default: 24)
}
