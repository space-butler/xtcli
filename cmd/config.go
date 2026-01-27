package cmd

import (
	"fmt"
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

func init() {
	configCmd.AddCommand(configCreateCmd)
	rootCmd.AddCommand(configCmd)
}
