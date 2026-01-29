package config

type Config struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Host      string   `json:"host"`
	VlcPath   string   `json:"vlc_path,omitempty"`
	Favorites []string `json:"favorites,omitempty"`
}
