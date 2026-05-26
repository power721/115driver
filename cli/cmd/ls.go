package cmd

import (
	"github.com/power721/115driver/cli/internal/output"
	"github.com/power721/115driver/cli/internal/resolver"
	"github.com/spf13/cobra"
)

var lsLong bool

var lsCmd = &cobra.Command{
	Use:   "ls [remote_path]",
	Short: "List directory contents",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		remotePath := "/"
		if len(args) > 0 {
			remotePath = args[0]
		}

		dirID, err := resolver.ResolveDir(client, remotePath)
		if err != nil {
			return &exitError{code: output.ExitNotFound, msg: err.Error()}
		}

		files, err := client.List(dirID)
		if err != nil {
			return &exitError{code: output.ExitError, msg: err.Error()}
		}

		jsonFiles := make([]output.JSONFile, 0, len(*files))
		for _, f := range *files {
			jsonFiles = append(jsonFiles, output.FileToJSON(&f))
		}

		if lsLong {
			printer.PrintFileTable(remotePath, jsonFiles)
		} else {
			printer.PrintFileList(remotePath, jsonFiles)
		}
		return nil
	},
}

func init() {
	lsCmd.Flags().BoolVarP(&lsLong, "long", "l", false, "Show detailed listing")
	rootCmd.AddCommand(lsCmd)
}
