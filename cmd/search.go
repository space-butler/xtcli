package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"xtcli/consts"
	"xtcli/xtream"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for streams on the Xtream Codes IPTV server",
	Args:  cobra.NoArgs,
}

var searchStreamCmd = &cobra.Command{
	Use:   "stream <search-string>",
	Short: "Search for streams by name (optionally within a category)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		searchTerm := args[0]
		categoryID, _ := cmd.Flags().GetInt64("category")
		streamType, _ := cmd.Flags().GetString("type")
		return handleSearchStream(searchTerm, categoryID, streamType)
	},
}

func init() {
	searchStreamCmd.Flags().Int64P("category", "c", 0, "Limit search to this category ID (0 = all categories)")
	searchStreamCmd.Flags().StringP("type", "t", "live", "Stream type: live or vod")
	searchCmd.AddCommand(searchStreamCmd)
	rootCmd.AddCommand(searchCmd)
}

func handleSearchStream(searchTerm string, categoryID int64, streamType string) error {
	type searchResult struct {
		Category string
		ID       int64
		Title    string
	}
	var results []searchResult
	searchLower := strings.ToLower(searchTerm)

	if categoryID != 0 {
		// Search inside a single category
		var streams []xtream.Stream
		var err error
		switch streamType {
		case "vod":
			streams, err = xtream.GetVodStreamsByCategory(categoryID)
		default:
			streams, err = xtream.GetStreamsByCategory(categoryID)
		}
		if err != nil {
			return err
		}
		for _, stream := range streams {
			if strings.Contains(strings.ToLower(stream.Name), searchLower) {
				results = append(results, searchResult{
					Category: stream.CategoryName,
					ID:       stream.ID,
					Title:    stream.Name,
				})
			}
		}
	} else {
		// Search across all categories
		catType := consts.CATEGORY_TYPE_LIVE
		if streamType == "vod" {
			catType = consts.CATEGORY_TYPE_VOD
		}
		categories, err := xtream.GetCategories(catType)
		if err != nil {
			return err
		}

		for _, category := range categories {
			var streams []xtream.Stream
			switch streamType {
			case "vod":
				streams, err = xtream.GetVodStreamsByCategory(category.ID)
			default:
				streams, err = xtream.GetStreamsByCategory(category.ID)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to get streams for category %s: %v\n", category.Name, err)
				continue
			}

			for _, stream := range streams {
				if strings.Contains(strings.ToLower(stream.Name), searchLower) {
					results = append(results, searchResult{
						Category: stream.CategoryName,
						ID:       stream.ID,
						Title:    stream.Name,
					})
				}
			}
		}
	}

	if len(results) == 0 {
		fmt.Printf("No streams found matching '%s'\n", searchTerm)
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Category", "ID", "Title")

	for _, result := range results {
		table.Append(
			result.Category,
			strconv.FormatInt(result.ID, 10),
			result.Title,
		)
	}

	table.Render()
	fmt.Printf("\nFound %d stream(s) matching '%s'\n", len(results), searchTerm)
	return nil
}
