package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"xtcli/xtream"

	"github.com/olekukonko/tablewriter"
	xtreamcodes "github.com/space-butler/go.xtream-codes"
	"github.com/spf13/cobra"
)

var serverInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display server and account information",
	Long:  `Display comprehensive information about the Xtream Codes server and your account status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleServerInfo()
	},
}

func handleServerInfo() error {
	err := xtream.Initialize(cfg.Username, cfg.Password, cfg.Host)
	if err != nil {
		return fmt.Errorf("failed to initialize xtream client: %w", err)
	}

	serverInfo, userInfo, err := xtream.GetServerInfo()
	if err != nil {
		return fmt.Errorf("failed to get server info: %w", err)
	}

	// Display Server Information
	fmt.Println("=== Server Information ===")
	serverTable := tablewriter.NewWriter(os.Stdout)
	serverTable.Header("Property", "Value")

	serverTable.Append("URL", serverInfo.URL)
	serverTable.Append("Protocol", serverInfo.Protocol)
	serverTable.Append("Port", strconv.FormatInt(int64(serverInfo.Port), 10))
	serverTable.Append("HTTPS Port", strconv.FormatInt(int64(serverInfo.HTTPSPort), 10))
	serverTable.Append("RTMP Port", strconv.FormatInt(int64(serverInfo.RTMPPort), 10))
	serverTable.Append("Timezone", serverInfo.Timezone)
	serverTable.Append("Current Time", serverInfo.TimeNow)
	serverTable.Append("Process Status", formatBool(serverInfo.Process))

	serverTable.Render()

	fmt.Println("\n=== Account Information ===")
	userTable := tablewriter.NewWriter(os.Stdout)
	userTable.Header("Property", "Value")

	userTable.Append("Username", userInfo.Username)
	userTable.Append("Status", userInfo.Status)
	userTable.Append("Active Connections", strconv.FormatInt(int64(userInfo.ActiveConnections), 10))
	userTable.Append("Max Connections", strconv.FormatInt(int64(userInfo.MaxConnections), 10))
	userTable.Append("Is Trial", formatConvertibleBool(userInfo.IsTrial))
	userTable.Append("Authenticated", formatConvertibleBool(userInfo.Auth))

	if userInfo.ExpDate != nil {
		expDate := userInfo.ExpDate.Time
		userTable.Append("Expiration Date", expDate.Format("2006-01-02 15:04:05"))

		// Calculate remaining days
		now := time.Now()
		remainingDays := int(expDate.Sub(now).Hours() / 24)
		if remainingDays > 0 {
			userTable.Append("Days Remaining", strconv.Itoa(remainingDays))
		} else {
			userTable.Append("Days Remaining", "EXPIRED")
		}
	} else {
		userTable.Append("Expiration Date", "N/A")
	}

	createdAt := userInfo.CreatedAt.Time
	userTable.Append("Created At", createdAt.Format("2006-01-02 15:04:05"))

	if len(userInfo.AllowedOutputFormats) > 0 {
		formats := ""
		for i, format := range userInfo.AllowedOutputFormats {
			if i > 0 {
				formats += ", "
			}
			formats += format
		}
		userTable.Append("Allowed Formats", formats)
	}

	if userInfo.Message != "" {
		userTable.Append("Message", userInfo.Message)
	}

	userTable.Render()

	return nil
}

func formatBool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func formatConvertibleBool(cb xtreamcodes.ConvertibleBoolean) string {
	// ConvertibleBoolean doesn't have a direct method to get bool value
	// We can marshal it and check the result
	data, err := cb.MarshalJSON()
	if err != nil {
		return "Unknown"
	}
	// MarshalJSON returns "1" for true, "0" for false
	if string(data) == "1" {
		return "Yes"
	}
	return "No"
}

func init() {
	serverCmd.AddCommand(serverInfoCmd)
}
