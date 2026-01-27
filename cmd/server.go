package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server-related operations",
	Long:  `Commands for interacting with and getting information about the Xtream Codes server.`,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
