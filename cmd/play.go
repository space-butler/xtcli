package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"xtcli/config"
	"xtcli/xtream"

	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play <stream-id>",
	Short: "Play stream",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		streamType, _ := cmd.Flags().GetString("type")
		format, _ := cmd.Flags().GetString("format")
		favArg, _ := cmd.Flags().GetString("fav")

		if favArg != "" {
			fav, err := config.GetFavorite(favArg, activeProvider.Name)
			if err != nil {
				return err
			}
			return handlePlayStream(fav.StreamID, fav.Type, format, fav.Name)
		}

		if len(args) == 0 {
			return fmt.Errorf("either a stream-id argument or --fav flag is required")
		}

		streamID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		return handlePlayStream(streamID, streamType, format, "")
	},
}

func init() {
	playCmd.Flags().StringP("type", "t", "live", "Stream type: live or vod")
	playCmd.Flags().StringP("format", "f", "ts", "Format/extension (e.g. ts, m3u8, mp4). Actual supported formats may depend on your provider.")
	playCmd.Flags().StringP("fav", "", "", "Favorite number or name to play")
	rootCmd.AddCommand(playCmd)
}

func handlePlayStream(streamID int64, streamType, format, label string) error {
	var streamURL string
	var err error

	switch streamType {
	case "vod":
		if format == "" {
			format = "mkv"
		}
		streamURL, err = xtream.GetVodStreamURL(streamID, format)
	default:
		if format == "" {
			format = "ts"
		}
		streamURL, err = xtream.GetStreamURL(streamID, format)
	}
	if err != nil {
		return err
	}

	if label != "" {
		fmt.Printf("Playing '%s' (%d)...\n", label, streamID)
	} else {
		fmt.Printf("Playing stream ID %d...\n", streamID)
	}

	cmd := exec.Command(cfg.VlcPath, streamURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start VLC: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case <-sigCh:
		cmd.Process.Kill()
		<-done
		return nil
	case err := <-done:
		if err != nil {
			return fmt.Errorf("VLC exited with error: %w", err)
		}
		return nil
	}
}
