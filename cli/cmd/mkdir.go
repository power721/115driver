package cmd

import (
	"fmt"
	"path"
	"strings"

	"github.com/power721/115driver/cli/internal/output"
	"github.com/power721/115driver/cli/internal/resolver"
	"github.com/power721/115driver/pkg/driver"
	"github.com/spf13/cobra"
)

var mkdirParents bool

var mkdirCmd = &cobra.Command{
	Use:   "mkdir [-p] <remote_path>",
	Short: "Create a directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		remotePath := strings.TrimRight(args[0], "/")
		if remotePath == "" {
			return &exitError{code: output.ExitArgs, msg: "Cannot create root directory."}
		}
		dirName := path.Base(remotePath)
		parentPath := path.Dir(remotePath)
		if parentPath == "." {
			parentPath = "/"
		}

		if mkdirParents {
			return mkdirP(parentPath, dirName, remotePath)
		}

		parentID, err := resolver.ResolveDir(client, parentPath)
		if err != nil {
			return &exitError{code: output.ExitNotFound, msg: fmt.Sprintf("Parent directory not found: %s", parentPath)}
		}

		dirID, err := client.Mkdir(parentID, dirName)
		if err != nil {
			return &exitError{code: output.ExitError, msg: err.Error()}
		}

		printer.PrintSuccess(map[string]interface{}{
			"name":   dirName,
			"dir_id": dirID,
			"path":   remotePath,
		})
		if !jsonOutput {
			fmt.Printf("Created directory: %s (ID: %s)\n", remotePath, dirID)
		}
		return nil
	},
}

func mkdirP(parentPath, dirName, fullPath string) error {
	parts := strings.Split(strings.Trim(parentPath+"/"+dirName, "/"), "/")
	currentID := resolver.RootID
	createdPath := ""

	for _, part := range parts {
		if part == "" {
			continue
		}
		createdPath += "/" + part

		existingID, err := resolver.ResolveDir(client, createdPath)
		if err == nil && existingID != "" {
			currentID = existingID
			continue
		}

		newID, err := client.Mkdir(currentID, part)
		if err != nil {
			if err == driver.ErrExist {
				existingID, _ := resolver.ResolveDir(client, createdPath)
				if existingID != "" {
					currentID = existingID
					continue
				}
			}
			return &exitError{code: output.ExitError, msg: err.Error()}
		}
		currentID = newID
	}

	printer.PrintSuccess(map[string]interface{}{
		"name":   dirName,
		"dir_id": currentID,
		"path":   fullPath,
	})
	if !jsonOutput {
		fmt.Printf("Created directory: %s (ID: %s)\n", fullPath, currentID)
	}
	return nil
}

func init() {
	mkdirCmd.Flags().BoolVarP(&mkdirParents, "parents", "p", false, "Create parent directories as needed")
	rootCmd.AddCommand(mkdirCmd)
}
