package resolver

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/power721/115driver/pkg/driver"
)

const RootID = "0"

// TODO: add LRU cache for path resolution

type pathResolverClient interface {
	DirName2CID(dir string) (*driver.APIGetDirIDResp, error)
	List(dirID string, opts ...driver.ListOption) (*[]driver.File, error)
}

func ResolveDir(client pathResolverClient, remotePath string) (string, error) {
	if remotePath == "" || remotePath == "/" {
		return RootID, nil
	}

	cleaned := strings.TrimPrefix(remotePath, "/")
	cleaned = strings.TrimSuffix(cleaned, "/")

	if cleaned == "" {
		return RootID, nil
	}

	resp, err := client.DirName2CID(cleaned)
	if err != nil {
		return "", fmt.Errorf("directory not found: %s (%w)", remotePath, err)
	}
	if string(resp.CategoryID) == RootID {
		return "", fmt.Errorf("directory not found: %s", remotePath)
	}
	return string(resp.CategoryID), nil
}

func ResolveFile(client pathResolverClient, remotePath string) (string, error) {
	cleaned := strings.TrimPrefix(remotePath, "/")
	cleaned = strings.TrimSuffix(cleaned, "/")

	dir := path.Dir(cleaned)
	fileName := path.Base(cleaned)

	var dirID string
	if dir == "." || dir == "" {
		dirID = RootID
	} else {
		var err error
		dirID, err = ResolveDir(client, dir)
		if err != nil {
			return "", err
		}
	}

	files, err := client.List(dirID)
	if err != nil {
		return "", fmt.Errorf("failed to list directory: %w", err)
	}

	for _, f := range *files {
		if f.Name == fileName && !f.IsDirectory {
			return f.FileID, nil
		}
	}
	return "", fmt.Errorf("file not found: %s", remotePath)
}

func ResolvePath(client pathResolverClient, remotePath string) (string, bool, error) {
	if remotePath == "" || remotePath == "/" {
		return RootID, true, nil
	}

	// Try as directory first
	dirID, err := ResolveDir(client, remotePath)
	if err == nil && dirID != "" {
		return dirID, true, nil
	}

	// Try as file
	fileID, fileErr := ResolveFile(client, remotePath)
	if fileErr != nil {
		return "", false, fmt.Errorf("%w; also tried as directory: %v", fileErr, err)
	}
	return fileID, false, nil
}

func ResolveLocalDownloadPath(localTarget, fileName string) string {
	if localTarget == "" {
		return fileName
	}

	if fi, err := osStat(localTarget); err == nil && fi.IsDir() {
		return filepath.Join(localTarget, fileName)
	}

	if strings.HasSuffix(localTarget, string(filepath.Separator)) {
		return filepath.Join(strings.TrimSuffix(localTarget, string(filepath.Separator)), fileName)
	}

	if filepath.Ext(filepath.Base(localTarget)) == "" {
		return filepath.Join(localTarget, fileName)
	}

	return localTarget
}

var osStat = func(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
