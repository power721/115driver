package cmd

import (
	"fmt"
	"path"

	"github.com/power721/115driver/cli/internal/output"
	"github.com/power721/115driver/cli/internal/resolver"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename <remote_path> <new_name>",
	Short: "Rename a file or directory",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		remotePath := args[0]
		newName := args[1]

		fileID, _, err := resolver.ResolvePath(client, remotePath)
		if err != nil {
			return &exitError{code: output.ExitNotFound, msg: err.Error()}
		}

		if err := client.Rename(fileID, newName); err != nil {
			return &exitError{code: output.ExitError, msg: err.Error()}
		}

		printer.PrintSuccess(map[string]interface{}{
			"old_name": path.Base(remotePath),
			"new_name": newName,
			"file_id":  fileID,
		})
		if !jsonOutput {
			fmt.Printf("Renamed to: %s\n", newName)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
