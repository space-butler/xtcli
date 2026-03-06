package cmd

import (
	"fmt"
	"xtcli/config"
	"xtcli/xtream"

	"github.com/spf13/cobra"
)

var cfg *config.Config
var activeProvider *config.Provider

var rootCmd = &cobra.Command{
	Use:   "xtcli",
	Short: "A tool to dump Xtream Codes IPTV server data",
	Long:  `xtcli is a command-line tool that allows users to extract and dump data from Xtream Codes IPTV servers.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for commands that don't need it (like 'config create')
		if cmd.Name() == "create" || cmd.Name() == "config" {
			return nil
		}

		// Allow provider management commands without a valid provider
		if cmd.Name() == "provider" || cmd.Name() == "add" || cmd.Name() == "del" ||
			cmd.Name() == "list" || cmd.Name() == "default" {
			// Check if we're under the provider subcommand tree
			for p := cmd.Parent(); p != nil; p = p.Parent() {
				if p.Name() == "provider" || p.Name() == "config" {
					return nil
				}
			}
		}

		if !config.Exists() {
			return fmt.Errorf("config file does not exist. Run 'xtcli config create' to create a default config file")
		}

		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}

		// Resolve the active provider
		providerName, _ := cmd.Flags().GetString("provider")
		activeProvider, err = config.GetProvider(providerName)
		if err != nil {
			return err
		}

		// Initialize xtream client with cache TTL from config
		cacheTTL, _ := config.GetCacheTTL()
		if err := xtream.InitializeWithCacheTTL(activeProvider.Username, activeProvider.Password, activeProvider.Host, cacheTTL); err != nil {
			return err
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().String("provider", "", "IPTV provider name (default: from config)")
}
