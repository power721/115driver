package cmd

import (
	"github.com/spf13/cobra"
)

var cpCmd = &cobra.Command{
	Use:   "cp <source_path> <destination_dir>",
	Short: "Copy files into a destination directory",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return moveOrCopy(args[0], args[1], client.Copy)
	},
}

func init() {
	rootCmd.AddCommand(cpCmd)
}
