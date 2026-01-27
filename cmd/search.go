package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"xtream-dump/consts"
	"xtream-dump/xtream"

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
	Short: "Search for streams by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		searchTerm := args[0]
		return handleSearchStream(searchTerm)
	},
}

func init() {
	searchCmd.AddCommand(searchStreamCmd)
	rootCmd.AddCommand(searchCmd)
}

func handleSearchStream(searchTerm string) error {
	err := xtream.Initialize(cfg.Username, cfg.Password, cfg.Host)
	if err != nil {
		return err
	}

	// Get all live categories
	categories, err := xtream.GetCategories(consts.CATEGORY_TYPE_LIVE)
	if err != nil {
		return err
	}

	// Search results
	type searchResult struct {
		Category string
		ID       int64
		Title    string
	}
	var results []searchResult

	// Search through all categories
	searchLower := strings.ToLower(searchTerm)
	for _, category := range categories {
		streams, err := xtream.GetStreamsByCategory(category.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to get streams for category %s: %v\n", category.Name, err)
			continue
		}

		// Filter streams by search term
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

	// Display results
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
