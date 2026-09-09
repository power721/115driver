package cmd

import (
	"fmt"

	"github.com/power721/115driver/cli/internal/auth"
	"github.com/power721/115driver/cli/internal/output"
	"github.com/power721/115driver/cli/internal/resolver"
	"github.com/power721/115driver/pkg/driver"
	"github.com/spf13/cobra"
)

var offlineSaveDir string

var offlineCmd = &cobra.Command{
	Use:   "offline",
	Short: "Manage offline downloads",
}

var offlineAddCmd = &cobra.Command{
	Use:   "add <url>",
	Short: "Add an offline download task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		saveDirID := resolver.RootID
		saveDirName := ""
		if offlineSaveDir != "" {
			// -d flag takes priority
			id, err := resolver.ResolveDir(client, offlineSaveDir)
			if err != nil {
				return &exitError{code: output.ExitNotFound, msg: fmt.Sprintf("Save directory not found: %s", offlineSaveDir)}
			}
			saveDirID = id
			saveDirName = offlineSaveDir
		} else {
			// Try config default
			if cfgDir := auth.ReadProfileConfig(configPath, profile, "default_offline_save_dir"); cfgDir != "" {
				id, err := resolver.ResolveDir(client, cfgDir)
				if err != nil {
					return &exitError{code: output.ExitNotFound, msg: fmt.Sprintf("Save directory not found: %s (from config default_offline_save_dir)", cfgDir)}
				}
				saveDirID = id
				saveDirName = cfgDir
			}
		}

		hashes, err := client.AddOfflineTaskURIs([]string{url}, saveDirID)
		if err != nil {
			return &exitError{code: output.ExitError, msg: err.Error()}
		}

		printer.PrintSuccess(map[string]interface{}{
			"url":      url,
			"hashes":   hashes,
			"save_dir": saveDirName,
		})
		if !jsonOutput {
			fmt.Printf("Offline task added: %s\n", url)
		}
		return nil
	},
}

var offlineListCmd = &cobra.Command{
	Use:   "list",
	Short: "List offline download tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		var allTasks []*driver.OfflineTask
		var total int64

		for page := int64(1); ; page++ {
			result, err := client.ListOfflineTask(page)
			if err != nil {
				return &exitError{code: output.ExitError, msg: err.Error()}
			}
			allTasks = append(allTasks, result.Tasks...)
			total = result.Total
			if page >= result.PageCount {
				break
			}
		}

		tasks := make([]map[string]interface{}, 0, len(allTasks))
		for _, t := range allTasks {
			tasks = append(tasks, map[string]interface{}{
				"name":    t.Name,
				"hash":    t.InfoHash,
				"status":  t.GetStatus(),
				"percent": t.Percent,
				"size":    t.Size,
			})
		}

		if jsonOutput {
			printer.PrintSuccess(map[string]interface{}{
				"total": total,
				"tasks": tasks,
			})
		} else {
			fmt.Printf("Offline tasks (%d total):\n\n", total)
			printer.PrintOfflineTable(tasks)
		}
		return nil
	},
}

var offlineRmCmd = &cobra.Command{
	Use:   "rm <hash>",
	Short: "Remove an offline download task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash := args[0]

		if err := client.DeleteOfflineTasks([]string{hash}, false); err != nil {
			return &exitError{code: output.ExitError, msg: err.Error()}
		}

		printer.PrintSuccess(map[string]interface{}{
			"deleted_hashes": []string{hash},
		})
		if !jsonOutput {
			fmt.Printf("Removed offline task: %s\n", hash)
		}
		return nil
	},
}

func init() {
	offlineAddCmd.Flags().StringVarP(&offlineSaveDir, "dir", "d", "", "Remote directory to save downloaded files")
	offlineCmd.AddCommand(offlineAddCmd)
	offlineCmd.AddCommand(offlineListCmd)
	offlineCmd.AddCommand(offlineRmCmd)
	rootCmd.AddCommand(offlineCmd)
}
