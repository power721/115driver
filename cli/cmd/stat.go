package cmd

import (
	"github.com/power721/115driver/cli/internal/output"
	"github.com/power721/115driver/cli/internal/resolver"
	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stat <remote_path>",
	Short: "Show file or directory details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		remotePath := args[0]
		fileID, _, err := resolver.ResolvePath(client, remotePath)
		if err != nil {
			return &exitError{code: output.ExitNotFound, msg: err.Error()}
		}

		statInfo, err := client.Stat(fileID)
		if err != nil {
			return &exitError{code: output.ExitError, msg: err.Error()}
		}

		jsonStat := output.JSONStat{
			Name:       statInfo.Name,
			IsDir:      statInfo.IsDirectory,
			FileID:     fileID,
			Sha1:       statInfo.Sha1,
			PickCode:   statInfo.PickCode,
			CreateTime: statInfo.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime: statInfo.UpdateTime.Format("2006-01-02 15:04:05"),
			FileCount:  statInfo.FileCount,
			DirCount:   statInfo.DirCount,
			Parents:    make([]output.JSONDir, 0, len(statInfo.Parents)),
		}

		if !statInfo.IsDirectory {
			f, err := client.GetFile(fileID)
			if err != nil {
				return &exitError{code: output.ExitError, msg: "Failed to get file details: " + err.Error()}
			}
			jsonStat.Size = f.Size
		}

		for _, p := range statInfo.Parents {
			jsonStat.Parents = append(jsonStat.Parents, output.JSONDir{ID: p.ID, Name: p.Name})
		}

		printer.PrintStatTable(jsonStat)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statCmd)
}
