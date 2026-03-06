package cmd

import (
	"fmt"
	"os"
	"strconv"
	"xtcli/config"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var favCmd = &cobra.Command{
	Use:   "fav",
	Short: "Manage favorite streams",
	Args:  cobra.NoArgs,
}

var favAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a stream to favorites",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name is required and cannot be empty")
		}

		streamType, _ := cmd.Flags().GetString("type")
		streamID, _ := cmd.Flags().GetInt64("id")
		if streamID == 0 {
			return fmt.Errorf("--id is required")
		}

		fav := config.Favorite{
			Name:     name,
			Type:     streamType,
			StreamID: streamID,
		}

		if err := config.AddFavorite(fav); err != nil {
			return err
		}

		fmt.Printf("Added favorite '%s' (type: %s, stream ID: %d)\n", name, streamType, streamID)
		return nil
	},
}

var favDelCmd = &cobra.Command{
	Use:   "del <number|name> [number|name]...",
	Short: "Remove favorites by number or name",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		removed, err := config.RemoveFavorites(args)
		if err != nil {
			return err
		}

		if removed == 0 {
			fmt.Println("No matching favorites found")
			return nil
		}

		fmt.Printf("Removed %d favorite(s)\n", removed)
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

		table := tablewriter.NewWriter(os.Stdout)
		table.Header("#", "Name", "Type", "Stream ID")

		for _, fav := range favorites {
			table.Append(strconv.Itoa(fav.Number), fav.Name, fav.Type, strconv.FormatInt(fav.StreamID, 10))
		}

		table.Render()
		return nil
	},
}

func init() {
	favAddCmd.Flags().StringP("name", "n", "", "Shortcut name for the favorite (required)")
	favAddCmd.Flags().StringP("type", "t", "live", "Stream type: live or vod")
	favAddCmd.Flags().Int64P("id", "i", 0, "Stream ID (required)")

	favCmd.AddCommand(favAddCmd)
	favCmd.AddCommand(favDelCmd)
	favCmd.AddCommand(favListCmd)
	rootCmd.AddCommand(favCmd)
}
