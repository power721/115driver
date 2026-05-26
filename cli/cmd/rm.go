package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/power721/115driver/cli/internal/output"
	"github.com/power721/115driver/cli/internal/resolver"
	"github.com/spf13/cobra"
)

var rmForce bool

var rmCmd = &cobra.Command{
	Use:   "rm <remote_path>",
	Short: "Delete file or directory (moves to recycle bin)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		remotePath := args[0]

		fileID, isDir, err := resolver.ResolvePath(client, remotePath)
		if err != nil {
			return &exitError{code: output.ExitNotFound, msg: err.Error()}
		}

		if isDir && !jsonOutput {
			fmt.Printf("Delete directory %s and all its contents? [y/N] ", remotePath)
			reader := bufio.NewReader(os.Stdin)
			resp, _ := reader.ReadString('\n')
			resp = strings.TrimSpace(strings.ToLower(resp))
			if resp != "y" && resp != "yes" {
				fmt.Println("Canceled.")
				return nil
			}
		}

		if err := client.Delete(fileID); err != nil {
			return &exitError{code: output.ExitError, msg: err.Error()}
		}

		printer.PrintSuccess(map[string]interface{}{
			"deleted":  []string{remotePath},
			"file_ids": []string{fileID},
		})
		if !jsonOutput {
			fmt.Printf("Deleted: %s\n", remotePath)
		}
		return nil
	},
}

func init() {
	rmCmd.Flags().BoolVarP(&rmForce, "force", "f", false, "Reserved for future permanent delete")
	rootCmd.AddCommand(rmCmd)
}
