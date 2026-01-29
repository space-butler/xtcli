package cmd

import (
	"fmt"
	"xtcli/cache"
	"xtcli/config"
	"xtcli/consts"
	"xtcli/xtream"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Args:  cobra.NoArgs,
}

var configCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a default config file",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.CreateDefault(); err != nil {
			return err
		}

		configPath := config.GetConfigPath()

		fmt.Printf("Default config file created at: %s\n", configPath)
		fmt.Println("Please edit the file with your Xtream Codes server details.")
		return nil
	},
}

var favCmd = &cobra.Command{
	Use:   "fav",
	Short: "Manage favorite streams",
	Args:  cobra.NoArgs,
}

var favAddCmd = &cobra.Command{
	Use:   "add <stream_id> [stream_id]...",
	Short: "Add streams to favorites",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.AddFavorites(args); err != nil {
			return err
		}

		fmt.Printf("Added %d stream(s) to favorites\n", len(args))
		return nil
	},
}

var favDelCmd = &cobra.Command{
	Use:   "del <stream_id> [stream_id]...",
	Short: "Remove streams from favorites",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.RemoveFavorites(args); err != nil {
			return err
		}

		fmt.Printf("Removed %d stream(s) from favorites\n", len(args))
		return nil
	},
}

var favListCmd = &cobra.Command{
	Use:   "list",
	Short: "List favorite streams",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		favorites, err := config.GetFavorites()
		if err != nil {
			return err
		}

		if len(favorites) == 0 {
			fmt.Println("No favorite streams")
			return nil
		}

		fmt.Println("Favorite streams:")
		for _, id := range favorites {
			fmt.Printf("  - %s\n", id)
		}
		return nil
	},
}

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage cache",
	Args:  cobra.NoArgs,
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the local cache",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cache.Clear(); err != nil {
			return err
		}

		fmt.Println("Cache cleared successfully")
		return nil
	},
}

var cacheUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the local cache with fresh data from server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleCacheUpdate()
	},
}

var cacheInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show cache information",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cachePath := cache.GetCachePath()
		fmt.Printf("Cache file: %s\n", cachePath)

		cacheTTL, _ := config.GetCacheTTL()
		fmt.Printf("Cache TTL: %d hours\n", cacheTTL)

		return nil
	},
}

func handleCacheUpdate() error {
	fmt.Println("Updating cache with fresh data from server...")

	// Clear existing cache to force refresh
	if err := cache.Clear(); err != nil {
		return fmt.Errorf("failed to clear cache: %v", err)
	}

	// Fetch live categories
	fmt.Print("Fetching live categories... ")
	liveCategories, err := xtream.GetCategories(consts.CATEGORY_TYPE_LIVE)
	if err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("failed to fetch live categories: %v", err)
	}
	fmt.Printf("OK (%d categories)\n", len(liveCategories))

	// Fetch VOD categories
	fmt.Print("Fetching VOD categories... ")
	vodCategories, err := xtream.GetCategories(consts.CATEGORY_TYPE_VOD)
	if err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("failed to fetch VOD categories: %v", err)
	}
	fmt.Printf("OK (%d categories)\n", len(vodCategories))

	// Fetch streams for each live category
	totalStreams := 0
	fmt.Printf("Fetching streams for %d live categories...\n", len(liveCategories))
	for i, category := range liveCategories {
		fmt.Printf("  [%d/%d] %s... ", i+1, len(liveCategories), category.Name)
		streams, err := xtream.GetStreamsByCategory(category.ID)
		if err != nil {
			fmt.Printf("FAILED (%v)\n", err)
			continue
		}
		fmt.Printf("OK (%d streams)\n", len(streams))
		totalStreams += len(streams)
	}

	fmt.Printf("\nCache update completed successfully!\n")
	fmt.Printf("Cached %d live categories, %d VOD categories, and %d streams\n",
		len(liveCategories), len(vodCategories), totalStreams)

	return nil
}

func init() {
	cacheCmd.AddCommand(cacheClearCmd)
	cacheCmd.AddCommand(cacheUpdateCmd)
	cacheCmd.AddCommand(cacheInfoCmd)
	favCmd.AddCommand(favAddCmd)
	favCmd.AddCommand(favDelCmd)
	favCmd.AddCommand(favListCmd)
	configCmd.AddCommand(configCreateCmd)
	configCmd.AddCommand(favCmd)
	configCmd.AddCommand(cacheCmd)
	rootCmd.AddCommand(configCmd)
}
