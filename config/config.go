package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"xtream-dump/consts"
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
