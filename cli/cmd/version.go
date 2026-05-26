package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("115driver version " + resolveVersion())
	},
}

func resolveVersion() string {
	if version != "dev" {
		return version
	}
	// Fallback: read commit hash from go build info (works with go install)
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			return "dev-" + s.Value[:7]
		}
	}
	return "dev"
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
