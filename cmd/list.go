package cmd

import (
	"fmt"
	"os"
	"strconv"
	"xtcli/consts"
	"xtcli/xtream"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List categories from the Xtream Codes IPTV server",
	Args:  cobra.NoArgs,
}

var listCategoriesCmd = &cobra.Command{
	Use:     "categories [type]",
	Aliases: []string{"c", "cat"},
	Short:   "List categories (type: live or vod, default: live)",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		categoryType := "live"
		if len(args) > 0 {
			categoryType = args[0]
		}

		var catType consts.CategoryType
		switch categoryType {
		case "live":
			catType = consts.CATEGORY_TYPE_LIVE
		case "vod":
			catType = consts.CATEGORY_TYPE_VOD
		default:
			catType = consts.CATEGORY_TYPE_LIVE
		}
		return handleListCategories(catType)
	},
}

var listStreamsCmd = &cobra.Command{
	Use:     "streams <category-id>",
	Aliases: []string{"s"},
	Short:   "List streams for a given category ID",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		categoryID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		return handleListStreams(categoryID)
	},
}

var listStreamCmd = &cobra.Command{
	Use:   "stream <stream-id>",
	Short: "Show details for a single stream ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		streamID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		return handleListStream(streamID)
	},
}

var listEPGCmd = &cobra.Command{
	Use:     "epg <stream-id>",
	Aliases: []string{"e"},
	Short:   "List short EPG data for a given stream ID",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		streamID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		limit, _ := cmd.Flags().GetInt("limit")
		return handleListEPG(streamID, limit)
	},
}

var listURLCmd = &cobra.Command{
	Use:     "url <stream-id>",
	Aliases: []string{"u"},
	Short:   "Get the stream URL for a given stream ID",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		streamID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		format, _ := cmd.Flags().GetString("format")
		return handleGetURL(streamID, format)
	},
}

func init() {
	listURLCmd.Flags().StringP("format", "f", "ts", "Stream format (ts, m3u8, etc.)")
	listEPGCmd.Flags().IntP("limit", "l", 4, "Number of EPG entries to retrieve")

	listCmd.AddCommand(listCategoriesCmd)
	listCmd.AddCommand(listStreamsCmd)
	listCmd.AddCommand(listStreamCmd)
	listCmd.AddCommand(listEPGCmd)
	listCmd.AddCommand(listURLCmd)
	rootCmd.AddCommand(listCmd)
}

func handleListCategories(catType consts.CategoryType) error {
	categories, err := xtream.GetCategories(catType)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ID", "Name")

	for _, category := range categories {
		table.Append(
			strconv.FormatInt(category.ID, 10),
			category.Name,
		)
	}

	table.Render()
	return nil
}

func handleListStreams(categoryID int64) error {
	streams, err := xtream.GetStreamsByCategory(categoryID)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ID", "Name", "Type")

	for _, stream := range streams {
		table.Append(
			strconv.FormatInt(stream.ID, 10),
			stream.Name,
			stream.Type,
		)
	}

	table.Render()
	return nil
}

func handleListStream(streamID int64) error {
	stream, err := xtream.GetStream(streamID)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Field", "Value")

	table.Append("ID", strconv.FormatInt(stream.ID, 10))
	table.Append("Name", stream.Name)
	table.Append("Category ID", strconv.FormatInt(stream.CategoryID, 10))
	table.Append("Category Name", stream.CategoryName)
	table.Append("Type", stream.Type)
	table.Append("Extension", stream.ContainerExtension)
	table.Append("Number", strconv.FormatInt(stream.Number, 10))
	table.Append("EPG Channel ID", stream.EPGChannelID)
	if !stream.Added.IsZero() {
		table.Append("Added", stream.Added.Format("2006-01-02 15:04:05"))
	}
	if stream.Icon != "" {
		table.Append("Icon", stream.Icon)
	}

	table.Render()
	return nil
}

func handleListEPG(streamID int64, limit int) error {
	epgData, err := xtream.GetShortEPG(streamID, limit)
	if err != nil {
		return err
	}

	if len(epgData) == 0 {
		fmt.Println("No EPG data found for this stream ID")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Title", "Start", "End")

	for _, epg := range epgData {
		table.Append(
			epg.Title,
			epg.StartTimestamp.Format("01/02 3:04 PM"),
			epg.StopTimestamp.Format("01/02 3:04 PM"),
		)
	}

	table.Render()
	return nil
}

func handleGetURL(streamID int64, format string) error {
	url, err := xtream.GetStreamURL(streamID, format)
	if err != nil {
		return err
	}

	fmt.Println(url)
	return nil
}
