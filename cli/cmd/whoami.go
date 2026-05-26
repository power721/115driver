package cmd

import (
	"fmt"

	"github.com/power721/115driver/cli/internal/output"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current authenticated user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		userInfo, err := client.GetUser()
		if err != nil {
			return &exitError{code: output.ExitAuth, msg: fmt.Sprintf("Failed to get user info: %v", err)}
		}

		if jsonOutput {
			printer.PrintSuccess(map[string]interface{}{
				"user_id":  userInfo.UserID,
				"username": userInfo.UserName,
				"vip":      userInfo.Vip,
				"expire":   userInfo.Expire,
			})
		} else {
			fmt.Printf("User ID:  %d\n", userInfo.UserID)
			fmt.Printf("Username: %s\n", userInfo.UserName)
			if userInfo.Vip > 0 {
				fmt.Printf("VIP:      %d (expires: %d)\n", userInfo.Vip, userInfo.Expire)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
