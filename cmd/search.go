package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"xtcli/cache"
	"xtcli/consts"
	"xtcli/xtream"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for streams on the Xtream Codes IPTV server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		streamType, _ := cmd.Flags().GetString("type")
		return handleSearchStream(args[0], 0, streamType)
	},
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

var searchEPGCmd = &cobra.Command{
	Use:   "epg <search-string>",
	Short: "Search EPG program titles. Use -c to fetch live from a category, otherwise searches cached EPG.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		categoryID, _ := cmd.Flags().GetInt64("category")
		return handleSearchEPG(args[0], categoryID)
	},
}

func init() {
	searchCmd.Flags().StringP("type", "t", "live", "Stream type: live or vod")
	searchStreamCmd.Flags().Int64P("category", "c", 0, "Limit search to this category ID (0 = all categories)")
	searchStreamCmd.Flags().StringP("type", "t", "live", "Stream type: live or vod")
	searchEPGCmd.Flags().Int64P("category", "c", 0, "Fetch EPG on demand for streams in this category ID")
	searchCmd.AddCommand(searchStreamCmd)
	searchCmd.AddCommand(searchEPGCmd)
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
						Category: category.Name,
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
	table.Header("ID", "Category", "Title")

	for _, result := range results {
		table.Append(
			strconv.FormatInt(result.ID, 10),
			result.Category,
			result.Title,
		)
	}

	table.Render()
	fmt.Printf("\nFound %d stream(s) matching '%s'\n", len(results), searchTerm)
	return nil
}

func handleSearchEPG(searchTerm string, categoryID int64) error {
	type epgResult struct {
		StreamID int64
		Category string
		Title    string
	}

	searchLower := strings.ToLower(searchTerm)

	// Build a category ID -> name lookup from cached categories
	categoryNames := make(map[int64]string)
	if cats, ok := cache.GetCategories(consts.CATEGORY_TYPE_LIVE); ok {
		for _, c := range cats {
			categoryNames[c.ID] = c.Name
		}
	}

	var results []epgResult
	seen := make(map[int64]bool)

	if categoryID != 0 {
		// Fetch streams and EPG on demand for the given category
		streams, err := xtream.GetStreamsByCategory(categoryID)
		if err != nil {
			return err
		}
		catName := categoryNames[categoryID]
		for _, stream := range streams {
			epgEntries, err := xtream.GetShortEPG(stream.ID, 10)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to get EPG for stream %d: %v\n", stream.ID, err)
				continue
			}
			for _, epg := range epgEntries {
				if strings.Contains(strings.ToLower(epg.Title), searchLower) {
					if !seen[stream.ID] {
						seen[stream.ID] = true
						results = append(results, epgResult{
							StreamID: stream.ID,
							Category: catName,
							Title:    stream.Name,
						})
					}
					break
				}
			}
		}
	} else {
		// Search all cached EPG files directly (not dependent on stream cache)
		epgStreamIDs := cache.GetCachedEPGStreamIDs()

		if len(epgStreamIDs) == 0 {
			// No cached EPG — fetch on demand across all categories
			categories, err := xtream.GetCategories(consts.CATEGORY_TYPE_LIVE)
			if err != nil {
				return err
			}
			total := len(categories)
			for i, category := range categories {
				fmt.Fprintf(os.Stderr, "\rSearching EPG: category %d/%d (%s)...", i+1, total, category.Name)
				streams, err := xtream.GetStreamsByCategory(category.ID)
				if err != nil {
					continue
				}
				catName := category.Name
				for _, stream := range streams {
					epgEntries, err := xtream.GetShortEPG(stream.ID, 10)
					if err != nil {
						continue
					}
					for _, epg := range epgEntries {
						if strings.Contains(strings.ToLower(epg.Title), searchLower) {
							if !seen[stream.ID] {
								seen[stream.ID] = true
								results = append(results, epgResult{
									StreamID: stream.ID,
									Category: catName,
									Title:    stream.Name,
								})
							}
							break
						}
					}
				}
			}
			fmt.Fprintln(os.Stderr)
		} else {
			// Build a stream ID -> stream lookup for name/category resolution
			streamIndex := make(map[int64]cache.Stream)
			for _, s := range cache.GetAllStreamsAny() {
				streamIndex[s.ID] = s
			}

			for _, streamID := range epgStreamIDs {
				epgEntries, found := cache.GetEPGAny(streamID)
				if !found {
					continue
				}
				for _, epg := range epgEntries {
					if strings.Contains(strings.ToLower(epg.Title), searchLower) {
						if !seen[streamID] {
							seen[streamID] = true
							streamName := strconv.FormatInt(streamID, 10)
							catName := ""
							if s, ok := streamIndex[streamID]; ok {
								streamName = s.Name
								catName = categoryNames[s.CategoryID]
							}
							results = append(results, epgResult{
								StreamID: streamID,
								Category: catName,
								Title:    streamName,
							})
						}
						break
					}
				}
			}
		}
	}

	if len(results) == 0 {
		fmt.Printf("No EPG entries found matching '%s'\n", searchTerm)
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ID", "Category", "Title")

	for _, result := range results {
		table.Append(
			strconv.FormatInt(result.StreamID, 10),
			result.Category,
			result.Title,
		)
	}

	table.Render()
	fmt.Printf("\nFound %d EPG entry/entries matching '%s'\n", len(results), searchTerm)
	return nil
}
