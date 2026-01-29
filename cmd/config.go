package cmd

import (
	"fmt"
	"xtcli/cache"
	"xtcli/config"

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

func init() {
	cacheCmd.AddCommand(cacheClearCmd)
	cacheCmd.AddCommand(cacheInfoCmd)
	favCmd.AddCommand(favAddCmd)
	favCmd.AddCommand(favDelCmd)
	favCmd.AddCommand(favListCmd)
	configCmd.AddCommand(configCreateCmd)
	configCmd.AddCommand(favCmd)
	configCmd.AddCommand(cacheCmd)
	rootCmd.AddCommand(configCmd)
}
