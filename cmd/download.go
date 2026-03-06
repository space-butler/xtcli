package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"xtcli/consts"
	"xtcli/xtream"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download <stream-id>",
	Short: "Download a stream by ID (e.g. a VOD movie found via search)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		streamID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		output, _ := cmd.Flags().GetString("output")
		streamType, _ := cmd.Flags().GetString("type")
		format, _ := cmd.Flags().GetString("format")
		quiet, _ := cmd.Flags().GetBool("quiet")
		return handleDownload(streamID, output, streamType, format, quiet)
	},
}

func init() {
	downloadCmd.Flags().StringP("output", "o", "", "Output file path (default: stream_<id>.<ext>)")
	downloadCmd.Flags().StringP("type", "t", "vod", "Stream type: live or vod")
	downloadCmd.Flags().StringP("format", "f", "mkv", "Format/extension for VOD (e.g. mkv, mp4)")
	downloadCmd.Flags().BoolP("quiet", "q", false, "Quiet mode (no progress output)")
	rootCmd.AddCommand(downloadCmd)
}

func handleDownload(streamID int64, output, streamType, format string, quiet bool) error {
	if cfg.VlcPath == "" {
		return fmt.Errorf("vlc_path not set in config (required for download); add it to ~/" + consts.CONFIG_FILE_NAME)
	}

	var streamURL string
	var err error

	switch streamType {
	case "vod":
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

	if output == "" {
		output = fmt.Sprintf("stream_%d.%s", streamID, format)
	}
	absOutput, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("output path: %w", err)
	}

	// Use VLC as the capture tool (same dependency as play)
	// -I dummy: headless; --sout: stream output to file
	sout := fmt.Sprintf("#std{access=file,mux=%s,dst=%s}", format, absOutput)

	if !quiet {
		fmt.Printf("Downloading stream %d to %s...\n", streamID, output)
	}

	vlcCmd := exec.Command(cfg.VlcPath,
		"-I", "dummy",
		"--no-video-title-show",
		"--quiet",
		"--play-and-exit",
		streamURL,
		"--sout", sout,
	)
	vlcCmd.Stdout = nil
	vlcCmd.Stderr = nil

	if err := vlcCmd.Start(); err != nil {
		return fmt.Errorf("failed to start VLC: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Progress goroutine: poll file size and update status line
	done := make(chan struct{})
	if !quiet {
		go func() {
			ticker := time.NewTicker(1000 * time.Millisecond)
			defer ticker.Stop()
			var lastSize int64
			var lastTime time.Time
			for {
				select {
				case <-ctx.Done():
					return
				case <-done:
					return
				case <-ticker.C:
					info, err := os.Stat(absOutput)
					if err != nil {
						continue
					}
					size := info.Size()
					now := time.Now()
					mb := float64(size) / (1024 * 1024)
					speed := ""
					if !lastTime.IsZero() && now.After(lastTime) {
						elapsed := now.Sub(lastTime).Seconds()
						if elapsed > 0 {
							deltaMB := float64(size-lastSize) / (1024 * 1024)
							mbPerSec := deltaMB / elapsed
							if mbPerSec > 0 {
								speed = fmt.Sprintf(" (%.1f MB/s)", mbPerSec)
							}
						}
					}
					lastSize = size
					lastTime = now
					fmt.Fprintf(os.Stderr, "\r  Downloaded: %.1f MB%s   ", mb, speed)
				}
			}
		}()
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- vlcCmd.Wait()
		close(done)
	}()

	<-done
	cancel()

	if !quiet {
		if _, err := os.Stat(absOutput); err == nil {
			info, _ := os.Stat(absOutput)
			sizeMB := float64(info.Size()) / (1024 * 1024)
			fmt.Fprintf(os.Stderr, "\r  Done: %.1f MB saved to %s\n", sizeMB, output)
		}
	}

	if err := <-errCh; err != nil {
		return fmt.Errorf("VLC exited with error: %w", err)
	}
	return nil
}
