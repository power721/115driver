package cmd

import (
	"fmt"
	"strings"

	"github.com/power721/115driver/cli/internal/output"
	"github.com/power721/115driver/internal/accountinfo"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show current account and storage info",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		userInfo, err := client.GetUser()
		if err != nil {
			return &exitError{code: output.ExitAuth, msg: fmt.Sprintf("Failed to get user info: %v", err)}
		}
		info, err := client.GetInfo()
		if err != nil {
			return &exitError{code: output.ExitError, msg: fmt.Sprintf("Failed to get account info: %v", err)}
		}

		account := accountinfo.FromDriverData(userInfo, info)
		if jsonOutput {
			printer.PrintSuccess(account)
		} else {
			fmt.Print(formatAccountInfoText(account))
		}
		return nil
	},
}

func formatAccountInfoText(info accountinfo.AccountInfo) string {
	var b strings.Builder
	fmt.Fprintf(&b, "User ID:  %d\n", info.User.UserID)
	fmt.Fprintf(&b, "Username: %s\n", info.User.Username)
	if info.User.VIP > 0 {
		fmt.Fprintf(&b, "VIP:      %d (expires: %d)\n", info.User.VIP, info.User.Expire)
	}
	fmt.Fprintln(&b, "Space:")
	fmt.Fprintf(&b, "  Total:    %s\n", formatSpaceSize(info.Space.Total))
	fmt.Fprintf(&b, "  Remain:   %s\n", formatSpaceSize(info.Space.Remain))
	fmt.Fprintf(&b, "  Used:     %s\n", formatSpaceSize(info.Space.Used))
	if info.LoginDevices.Last.Device != "" || info.LoginDevices.Last.IP != "" {
		fmt.Fprintf(&b, "Last Login: %s %s %s\n", info.LoginDevices.Last.Device, info.LoginDevices.Last.IP, info.LoginDevices.Last.City)
	}
	return b.String()
}

func formatSpaceSize(size accountinfo.SizeInfo) string {
	if size.SizeFormat != "" {
		return fmt.Sprintf("%s (%d bytes)", size.SizeFormat, size.Size)
	}
	return fmt.Sprintf("%s (%d bytes)", output.FormatFileSize(size.Size), size.Size)
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
