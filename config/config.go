package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	return filepath.Join(homeDir, consts.CONFIG_DIR_NAME, consts.CONFIG_FILE_NAME)
}

// Save writes the Config to a JSON file
func Save(cfg *Config) error {
	configPath := GetConfigPath()

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

// CreateDefault creates and saves a default config file with placeholder values
func CreateDefault() error {
	cfg := &Config{
		DefaultProvider: "default",
		Providers: []Provider{
			{
				Name:     "default",
				Username: "your_user_name",
				Password: "your_password",
				Host:     "https://path.to.your.xtream.iptv.server",
			},
		},
		VlcPath: "/path/to/vlc",
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

// AddFavorite adds a favorite to the favorites list, keyed by name
func AddFavorite(fav Favorite) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	// Replace if name already exists, preserving its number
	for i, f := range cfg.Favorites {
		if strings.EqualFold(f.Name, fav.Name) {
			fav.Number = f.Number
			cfg.Favorites[i] = fav
			return Save(cfg)
		}
	}

	// Assign the next number
	maxNum := 0
	for _, f := range cfg.Favorites {
		if f.Number > maxNum {
			maxNum = f.Number
		}
	}
	fav.Number = maxNum + 1

	cfg.Favorites = append(cfg.Favorites, fav)
	return Save(cfg)
}

// RemoveFavorites removes favorites by number or name (case-insensitive).
// Returns the number of favorites actually removed.
func RemoveFavorites(args []string) (int, error) {
	cfg, err := Load()
	if err != nil {
		return 0, err
	}

	newFavorites := make([]Favorite, 0)
	for _, f := range cfg.Favorites {
		remove := false
		for _, arg := range args {
			// Try matching by number first
			if n, err := strconv.Atoi(arg); err == nil {
				if f.Number == n {
					remove = true
					break
				}
			} else if strings.EqualFold(f.Name, arg) {
				remove = true
				break
			}
		}
		if !remove {
			newFavorites = append(newFavorites, f)
		}
	}

	removed := len(cfg.Favorites) - len(newFavorites)
	cfg.Favorites = newFavorites
	return removed, Save(cfg)
}

// GetFavorites returns the list of favorites
func GetFavorites() ([]Favorite, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}
	return cfg.Favorites, nil
}

// GetFavorite looks up a single favorite by number (if arg is an integer) or name (case-insensitive).
func GetFavorite(arg string) (*Favorite, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	for _, f := range cfg.Favorites {
		if n, err := strconv.Atoi(arg); err == nil {
			if f.Number == n {
				result := f
				return &result, nil
			}
		} else if strings.EqualFold(f.Name, arg) {
			result := f
			return &result, nil
		}
	}

	return nil, fmt.Errorf("favorite %q not found", arg)
}

// GetCacheTTL returns the cache time-to-live in hours (default: 24)
func GetCacheTTL() (int, error) {
	cfg, err := Load()
	if err != nil {
		return 24, err
	}
	if cfg.CacheTTL <= 0 {
		return 24, nil
	}
	return cfg.CacheTTL, nil
}

// GetProvider returns the provider matching the given name (case-insensitive).
// If name is empty, it returns the default provider.
func GetProvider(name string) (*Provider, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	if name == "" {
		name = cfg.DefaultProvider
	}

	if name == "" {
		if len(cfg.Providers) > 0 {
			p := cfg.Providers[0]
			return &p, nil
		}
		return nil, fmt.Errorf("no providers configured; run 'xtcli config provider add' to add one")
	}

	for _, p := range cfg.Providers {
		if strings.EqualFold(p.Name, name) {
			result := p
			return &result, nil
		}
	}

	return nil, fmt.Errorf("provider %q not found", name)
}

// AddProvider adds or updates a provider by name.
func AddProvider(p Provider) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	for i, existing := range cfg.Providers {
		if strings.EqualFold(existing.Name, p.Name) {
			cfg.Providers[i] = p
			return Save(cfg)
		}
	}

	cfg.Providers = append(cfg.Providers, p)

	// If this is the first provider, set it as default
	if cfg.DefaultProvider == "" {
		cfg.DefaultProvider = p.Name
	}

	return Save(cfg)
}

// RemoveProvider removes a provider by name (case-insensitive).
// Returns true if a provider was actually removed.
func RemoveProvider(name string) (bool, error) {
	cfg, err := Load()
	if err != nil {
		return false, err
	}

	newProviders := make([]Provider, 0)
	removed := false
	for _, p := range cfg.Providers {
		if strings.EqualFold(p.Name, name) {
			removed = true
		} else {
			newProviders = append(newProviders, p)
		}
	}

	if !removed {
		return false, nil
	}

	cfg.Providers = newProviders

	// If the default was removed, clear it or set to first remaining
	if strings.EqualFold(cfg.DefaultProvider, name) {
		if len(cfg.Providers) > 0 {
			cfg.DefaultProvider = cfg.Providers[0].Name
		} else {
			cfg.DefaultProvider = ""
		}
	}

	return true, Save(cfg)
}

// ListProviders returns all configured providers.
func ListProviders() ([]Provider, string, error) {
	cfg, err := Load()
	if err != nil {
		return nil, "", err
	}
	return cfg.Providers, cfg.DefaultProvider, nil
}

// SetDefaultProvider sets the default provider by name.
func SetDefaultProvider(name string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	found := false
	for _, p := range cfg.Providers {
		if strings.EqualFold(p.Name, name) {
			cfg.DefaultProvider = p.Name
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("provider %q not found", name)
	}

	return Save(cfg)
}
