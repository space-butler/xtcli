package cmd

import (
	"fmt"
	"xtcli/config"
	"xtcli/xtream"

	"github.com/spf13/cobra"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "xtcli",
	Short: "A tool to dump Xtream Codes IPTV server data",
	Long:  `xtcli is a command-line tool that allows users to extract and dump data from Xtream Codes IPTV servers.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for commands that don't need it (like 'config create')
		if cmd.Name() == "create" || cmd.Name() == "config" {
			return nil
		}

		if !config.Exists() {
			return fmt.Errorf("config file does not exist. Run 'xtcli config create' to create a default config file")
		}

		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}

		// Initialize xtream client with cache TTL from config
		cacheTTL, _ := config.GetCacheTTL()
		if err := xtream.InitializeWithCacheTTL(cfg.Username, cfg.Password, cfg.Host, cacheTTL); err != nil {
			return err
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
