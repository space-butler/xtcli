package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"xtream-dump/xtream"

	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play <stream-id>",
	Short: "Play stream",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var streamID int64 = 0
		streamID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		return handlePlayStream(streamID)
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
}

func handlePlayStream(streamID int64) error {
	err := xtream.Initialize(cfg.Username, cfg.Password, cfg.Host)
	if err != nil {
		return err
	}
	format := "ts"
	streamURL, err := xtream.GetStreamURL(streamID, format)
	if err != nil {
		return err
	}
	fmt.Printf("Playing stream ID %d with VLC...\n", streamID)
	cmd := exec.Command(cfg.VlcPath, streamURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
