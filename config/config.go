package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"xtcli/consts"
)

// Exists checks if the config file exists and returns an error if it doesn't
func Exists() bool {
	configPath := GetConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, consts.CONFIG_FILE_NAME)
}

// Save writes the Config to a JSON file in the user's home directory
func Save(cfg *Config) error {
	configPath := GetConfigPath()

	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

// CreateDefault creates and saves a default config file with placeholder values
func CreateDefault() error {
	cfg := &Config{
		Username: "your_user_name",
		Password: "your_password",
		Host:     "https://path.to.your.xtream.iptv.server",
	}
	return Save(cfg)
}

// Load reads the Config from a JSON file in the user's home directory
func Load() (*Config, error) {
	configPath := GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// AddFavorites adds stream IDs to the favorites list
func AddFavorites(streamIDs []string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	// Create a map for quick lookup to avoid duplicates
	existing := make(map[string]bool)
	for _, id := range cfg.Favorites {
		existing[id] = true
	}

	// Add new stream IDs if not already present
	for _, id := range streamIDs {
		if !existing[id] {
			cfg.Favorites = append(cfg.Favorites, id)
			existing[id] = true
		}
	}

	return Save(cfg)
}

// RemoveFavorites removes stream IDs from the favorites list
func RemoveFavorites(streamIDs []string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	// Create a map of IDs to remove
	toRemove := make(map[string]bool)
	for _, id := range streamIDs {
		toRemove[id] = true
	}

	// Filter out the stream IDs to remove
	newFavorites := make([]string, 0)
	for _, id := range cfg.Favorites {
		if !toRemove[id] {
			newFavorites = append(newFavorites, id)
		}
	}

	cfg.Favorites = newFavorites
	return Save(cfg)
}

// GetFavorites returns the list of favorite stream IDs
func GetFavorites() ([]string, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}
	return cfg.Favorites, nil
}
