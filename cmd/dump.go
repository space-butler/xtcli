package cmd

import (
	"fmt"
	"os"
	"strconv"
	"xtream-dump/xtream"

	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump data from the Xtream Codes IPTV server",
	Args:  cobra.NoArgs,
}

var dumpXMLTVCmd = &cobra.Command{
	Use:   "xmltv <file-path>",
	Short: "Dump XMLTV EPG data to a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		return handleDumpXMLTV(filePath)
	},
}

var dumpM3UCmd = &cobra.Command{
	Use:   "m3u <file-path> <stream-id>...",
	Short: "Dump M3U8 playlist for one or more streams to a file",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// First argument is the file path
		filePath := args[0]

		// All other arguments are stream IDs
		streamIDs := make([]int64, 0, len(args)-1)
		for i := 1; i < len(args); i++ {
			streamID, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid stream ID '%s': %w", args[i], err)
			}
			streamIDs = append(streamIDs, streamID)
		}

		return handleDumpM3U(streamIDs, filePath)
	},
}

func init() {
	dumpCmd.AddCommand(dumpXMLTVCmd)
	dumpCmd.AddCommand(dumpM3UCmd)
	rootCmd.AddCommand(dumpCmd)
}

func handleDumpXMLTV(filePath string) error {
	err := xtream.Initialize(cfg.Username, cfg.Password, cfg.Host)
	if err != nil {
		return err
	}

	xmltvData, err := xtream.GetXMLTVFile()
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, xmltvData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Successfully wrote XMLTV data to %s (%d bytes)\n", filePath, len(xmltvData))
	return nil
}

func handleDumpM3U(streamIDs []int64, filePath string) error {
	err := xtream.Initialize(cfg.Username, cfg.Password, cfg.Host)
	if err != nil {
		return err
	}

	// Create M3U8 content
	m3uContent := "#EXTM3U\n"

	// Process each stream ID
	for _, streamID := range streamIDs {
		// Get stream details
		stream, err := xtream.GetStream(streamID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to get stream %d: %v\n", streamID, err)
			continue
		}

		// Get stream URL
		streamURL, err := xtream.GetStreamURL(streamID, "m3u8")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to get URL for stream %d: %v\n", streamID, err)
			continue
		}

		// Add stream entry to M3U content
		m3uContent += fmt.Sprintf("#EXTINF:-1 tvg-id=\"%s\" tvg-name=\"%s\" tvg-logo=\"%s\" group-title=\"%s\",%s\n",
			stream.EPGChannelID,
			stream.Name,
			stream.Icon,
			stream.CategoryName,
			stream.Name,
		)
		m3uContent += streamURL + "\n"

		fmt.Printf("Added stream: %s (Category: %s)\n", stream.Name, stream.CategoryName)
	}

	// Write to file
	err = os.WriteFile(filePath, []byte(m3uContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("\nSuccessfully wrote M3U playlist with %d stream(s) to %s\n", len(streamIDs), filePath)
	return nil
}
