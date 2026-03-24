package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"xtcli/config"
	"xtcli/consts"
	"xtcli/xtream"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:          "download <stream-id>",
	Short:        "Download a stream by ID (e.g. a VOD movie found via search)",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		streamType, _ := cmd.Flags().GetString("type")
		format, _ := cmd.Flags().GetString("format")
		quiet, _ := cmd.Flags().GetBool("quiet")
		favArg, _ := cmd.Flags().GetString("fav")

		if favArg != "" {
			fav, err := config.GetFavorite(favArg, activeProvider.Name)
			if err != nil {
				return err
			}
			return handleDownload(fav.StreamID, output, fav.Type, format, quiet)
		}

		if len(args) == 0 {
			return fmt.Errorf("either a stream-id argument or --fav flag is required")
		}

		streamID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		return handleDownload(streamID, output, streamType, format, quiet)
	},
}

func init() {
	downloadCmd.Flags().StringP("output", "o", "", "Output file path (default: stream_<id>.<ext>)")
	downloadCmd.Flags().StringP("type", "t", "vod", "Stream type: live or vod")
	downloadCmd.Flags().StringP("format", "f", "mkv", "Format/extension for VOD (e.g. mkv, mp4)")
	downloadCmd.Flags().BoolP("quiet", "q", false, "Quiet mode (no progress output)")
	downloadCmd.Flags().StringP("fav", "", "", "Favorite number or name to download")
	rootCmd.AddCommand(downloadCmd)
}

func handleDownload(streamID int64, output, streamType, format string, quiet bool) error {
	if cfg.VlcPath == "" {
		return fmt.Errorf("vlc_path not set in config (required for download); add it to ~/" + consts.CONFIG_DIR_NAME + "/" + consts.CONFIG_FILE_NAME)
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
					speed := ""
					if !lastTime.IsZero() && now.After(lastTime) {
						elapsed := now.Sub(lastTime).Seconds()
						if elapsed > 0 {
							deltaBytes := size - lastSize
							if deltaBytes > 0 {
								bytesPerSec := float64(deltaBytes) / elapsed
								if bytesPerSec >= float64(consts.BYTES_PER_MB) {
									speed = fmt.Sprintf("(%.1f MB/s)", bytesPerSec/float64(consts.BYTES_PER_MB))
								} else {
									speed = fmt.Sprintf("(%.1f KB/s)", bytesPerSec/float64(consts.BYTES_PER_KB))
								}
							}
						}
					}
					lastSize = size
					lastTime = now
					fmt.Fprintf(os.Stderr, "\r  Downloaded: %.1f MB %s\033[K", float64(size)/float64(consts.BYTES_PER_MB), speed)
				}
			}
		}()
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- vlcCmd.Wait()
		close(done)
	}()

	interrupted := false
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	go func() {
		select {
		case <-sigCh:
			interrupted = true
			vlcCmd.Process.Kill()
		case <-done:
		}
	}()

	<-done
	cancel()

	if interrupted {
		<-errCh
		if !quiet {
			fmt.Fprintf(os.Stderr, "\r  Cancelled.\033[K\n")
		}
		return nil
	}

	if !quiet {
		if _, err := os.Stat(absOutput); err == nil {
			info, _ := os.Stat(absOutput)
			sizeMB := float64(info.Size()) / float64(consts.BYTES_PER_MB)
			fmt.Fprintf(os.Stderr, "\r  Done: %.1f MB saved to %s\n", sizeMB, output)
		}
	}

	if err := <-errCh; err != nil {
		return fmt.Errorf("VLC exited with error: %w", err)
	}
	return nil
}
