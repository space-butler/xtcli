package cmd

import (
	"fmt"
	"os"
	"xtcli/config"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Manage IPTV providers",
	Args:  cobra.NoArgs,
}

var providerAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add or update an IPTV provider",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name is required")
		}
		username, _ := cmd.Flags().GetString("username")
		if username == "" {
			return fmt.Errorf("--username is required")
		}
		password, _ := cmd.Flags().GetString("password")
		if password == "" {
			return fmt.Errorf("--password is required")
		}
		host, _ := cmd.Flags().GetString("host")
		if host == "" {
			return fmt.Errorf("--host is required")
		}

		p := config.Provider{
			Name:     name,
			Username: username,
			Password: password,
			Host:     host,
		}

		if err := config.AddProvider(p); err != nil {
			return err
		}

		fmt.Printf("Provider '%s' added successfully\n", name)
		return nil
	},
}

var providerDelCmd = &cobra.Command{
	Use:   "del <name>",
	Short: "Remove an IPTV provider",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		removed, err := config.RemoveProvider(args[0])
		if err != nil {
			return err
		}
		if !removed {
			fmt.Printf("Provider '%s' not found\n", args[0])
			return nil
		}
		fmt.Printf("Provider '%s' removed\n", args[0])
		return nil
	},
}

var providerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured IPTV providers",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		providers, defaultName, err := config.ListProviders()
		if err != nil {
			return err
		}

		if len(providers) == 0 {
			fmt.Println("No providers configured")
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header("Name", "Host", "Username", "Default")

		for _, p := range providers {
			isDefault := ""
			if p.Name == defaultName {
				isDefault = "*"
			}
			table.Append(p.Name, p.Host, p.Username, isDefault)
		}

		table.Render()
		return nil
	},
}

var providerSetDefaultCmd = &cobra.Command{
	Use:   "default <name>",
	Short: "Set the default IPTV provider",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.SetDefaultProvider(args[0]); err != nil {
			return err
		}
		fmt.Printf("Default provider set to '%s'\n", args[0])
		return nil
	},
}

func init() {
	providerAddCmd.Flags().StringP("name", "n", "", "Provider name (required)")
	providerAddCmd.Flags().StringP("username", "u", "", "Username (required)")
	providerAddCmd.Flags().StringP("password", "p", "", "Password (required)")
	providerAddCmd.Flags().String("host", "", "Server host URL (required)")

	providerCmd.AddCommand(providerAddCmd)
	providerCmd.AddCommand(providerDelCmd)
	providerCmd.AddCommand(providerListCmd)
	providerCmd.AddCommand(providerSetDefaultCmd)
	configCmd.AddCommand(providerCmd)
}
