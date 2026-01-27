package cmd

import (
	"fmt"
	"xtream-dump/config"

	"github.com/spf13/cobra"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "xtream-dump",
	Short: "A tool to dump Xtream Codes IPTV server data",
	Long:  `xtream-dump is a command-line tool that allows users to extract and dump data from Xtream Codes IPTV servers.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for commands that don't need it (like 'config create')
		if cmd.Name() == "create" || cmd.Name() == "config" {
			return nil
		}

		if !config.Exists() {
			return fmt.Errorf("config file does not exist. Run 'xtream-dump config create' to create a default config file")
		}

		var err error
		cfg, err = config.Load()
		return err
	},
}

func Execute() error {
	return rootCmd.Execute()
}
