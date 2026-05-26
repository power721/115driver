package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/power721/115driver/pkg/driver"
)

// FileTools holds file-related MCP tools
type FileTools struct {
	client *driver.Pan115Client
}

// NewFileTools creates a new FileTools instance
func NewFileTools(client *driver.Pan115Client) *FileTools {
	return &FileTools{
		client: client,
	}
}

// MkdirArgs defines arguments for mkdir tool
type MkdirArgs struct {
	ParentID string `json:"parent_id" jsonschema:"parent directory ID"`
	Name     string `json:"name" jsonschema:"name of the new directory"`
}

// DeleteArgs defines arguments for delete tool
type DeleteArgs struct {
	FileIDs []string `json:"file_ids" jsonschema:"IDs of files or directories to delete"`
}

// RenameArgs defines arguments for rename tool
type RenameArgs struct {
	FileID  string `json:"file_id" jsonschema:"ID of file or directory to rename"`
	NewName string `json:"new_name" jsonschema:"new name for the file or directory"`
}

// MoveArgs defines arguments for move tool
type MoveArgs struct {
	DirID   string   `json:"dir_id" jsonschema:"target directory ID"`
	FileIDs []string `json:"file_ids" jsonschema:"IDs of files or directories to move"`
}

// CopyArgs defines arguments for copy tool
type CopyArgs struct {
	DirID   string   `json:"dir_id" jsonschema:"target directory ID"`
	FileIDs []string `json:"file_ids" jsonschema:"IDs of files or directories to copy"`
}

// StatArgs defines arguments for stat tool
type StatArgs struct {
	FileID string `json:"file_id" jsonschema:"ID of file or directory to get info"`
}

// UploadFromURLArgs defines arguments for uploading from URL
type UploadFromURLArgs struct {
	URL      string `json:"url" jsonschema:"URL of the file to download and upload"`
	DirID    string `json:"dir_id" jsonschema:"target directory ID in 115 cloud for saving the file"`
	FileName string `json:"file_name,omitempty" jsonschema:"optional filename for the uploaded file, defaults to original filename"`
}

// UploadFromLocalArgs defines arguments for uploading from local file
type UploadFromLocalArgs struct {
	LocalPath string `json:"local_path" jsonschema:"absolute path to the local file to upload"`
	DirID     string `json:"dir_id" jsonschema:"target directory ID in 115 cloud"`
	FileName  string `json:"file_name,omitempty" jsonschema:"optional filename for the uploaded file, defaults to original filename"`
}

// DownloadFileArgs defines arguments for downloading a file
type DownloadFileArgs struct {
	PickCode  string `json:"pick_code" jsonschema:"pick code of the file to download"`
	LocalPath string `json:"local_path" jsonschema:"local path where the downloaded file will be saved"`
	UserAgent string `json:"user_agent,omitempty" jsonschema:"optional user agent for the download request, uses 115 browser UA if not specified"`
}

// GetDownloadInfoArgs defines arguments for getting download information
type GetDownloadInfoArgs struct {
	PickCode  string `json:"pick_code" jsonschema:"pick code of the file to get info for"`
	UserAgent string `json:"user_agent,omitempty" jsonschema:"optional user agent for the download request, uses 115 browser UA if not specified"`
}

// GetDownloadInfoResult defines the result for getting download information
type GetDownloadInfoResult struct {
	URL      string `json:"url" jsonschema:"download URL"`
	FileName string `json:"file_name" jsonschema:"file name"`
	Size     int64  `json:"size" jsonschema:"file size in bytes"`
}

// DownloadFileResult defines the result for downloading a file
type DownloadFileResult struct {
	Message string `json:"message" jsonschema:"result message"`
}

// RegisterTools registers file-related tools with the MCP server
func (ft *FileTools) RegisterTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "mkdir",
		Description: "Create a new directory",
	}, ft.mkdir)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete",
		Description: "Delete files or directories",
	}, ft.delete)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "rename",
		Description: "Rename a file or directory",
	}, ft.rename)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "move",
		Description: "Move files or directories to another directory",
	}, ft.move)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "copy",
		Description: "Copy files or directories to another directory",
	}, ft.copy)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "stat",
		Description: "Get detailed information about a file or directory",
	}, ft.stat)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "upload_from_url",
		Description: "Upload a file to 115 cloud storage from a URL",
	}, ft.uploadFromURL)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "upload_from_local",
		Description: "Upload a local file to 115 cloud storage",
	}, ft.uploadFromLocal)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "download_file",
		Description: "Download a file from 115 cloud storage to local path",
	}, ft.downloadFile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_download_info",
		Description: "Get download information for a file including URL, file name, and size",
	}, ft.getDownloadInfo)
}

func (ft *FileTools) mkdir(ctx context.Context, req *mcp.CallToolRequest, args MkdirArgs) (*mcp.CallToolResult, any, error) {
	dirID, err := ft.client.Mkdir(args.ParentID, args.Name)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to create directory: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	result := map[string]string{
		"directory_id": dirID,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to serialize result: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil, nil
}

func (ft *FileTools) delete(ctx context.Context, req *mcp.CallToolRequest, args DeleteArgs) (*mcp.CallToolResult, any, error) {
	if len(args.FileIDs) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "No file IDs provided",
				},
			},
			IsError: true,
		}, nil, nil
	}

	err := ft.client.Delete(args.FileIDs...)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to delete files: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "Files deleted successfully",
			},
		},
	}, nil, nil
}

func (ft *FileTools) rename(ctx context.Context, req *mcp.CallToolRequest, args RenameArgs) (*mcp.CallToolResult, any, error) {
	err := ft.client.Rename(args.FileID, args.NewName)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to rename file: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "File renamed successfully",
			},
		},
	}, nil, nil
}

func (ft *FileTools) move(ctx context.Context, req *mcp.CallToolRequest, args MoveArgs) (*mcp.CallToolResult, any, error) {
	if len(args.FileIDs) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "No file IDs provided",
				},
			},
			IsError: true,
		}, nil, nil
	}

	err := ft.client.Move(args.DirID, args.FileIDs...)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to move files: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "Files moved successfully",
			},
		},
	}, nil, nil
}

func (ft *FileTools) copy(ctx context.Context, req *mcp.CallToolRequest, args CopyArgs) (*mcp.CallToolResult, any, error) {
	if len(args.FileIDs) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "No file IDs provided",
				},
			},
			IsError: true,
		}, nil, nil
	}

	err := ft.client.Copy(args.DirID, args.FileIDs...)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to copy files: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "Files copied successfully",
			},
		},
	}, nil, nil
}

func (ft *FileTools) stat(ctx context.Context, req *mcp.CallToolRequest, args StatArgs) (*mcp.CallToolResult, any, error) {
	info, err := ft.client.Stat(args.FileID)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get file info: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	result := map[string]interface{}{
		"name":         info.Name,
		"pick_code":    info.PickCode,
		"sha1":         info.Sha1,
		"is_directory": info.IsDirectory,
		"file_count":   info.FileCount,
		"dir_count":    info.DirCount,
		"create_time":  info.CreateTime,
		"update_time":  info.UpdateTime,
		"parents":      info.Parents,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to serialize result: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil, nil
}

func (ft *FileTools) uploadFromURL(ctx context.Context, req *mcp.CallToolRequest, args UploadFromURLArgs) (*mcp.CallToolResult, any, error) {
	// Download the file from the URL
	resp, err := ft.client.Client.R().SetDoNotParseResponse(true).Get(args.URL)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to download file from URL: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}
	defer resp.RawBody().Close()

	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to download file from URL, status code: %d", resp.StatusCode()),
				},
			},
			IsError: true,
		}, nil, nil
	}

	// If fileName is empty, try to extract it from the URL
	fileName := args.FileName
	if fileName == "" {
		parts := strings.Split(args.URL, "/")
		fileName = parts[len(parts)-1]
		if fileName == "" {
			fileName = "downloaded_file"
		}
	}

	// Create a temporary file to store the downloaded content
	tempFile, err := os.CreateTemp("", "115_mcp_upload_*")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to create temporary file: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}
	defer os.Remove(tempFile.Name()) // Clean up the temp file afterwards
	defer tempFile.Close()

	// Copy the response body to the temporary file
	_, err = io.Copy(tempFile, resp.RawBody())
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to save downloaded content to temporary file: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	// Get the file size
	fileInfo, err := tempFile.Stat()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get file info: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}
	fileSize := fileInfo.Size()

	// Seek back to the beginning of the file
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to seek to beginning of file: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	// Upload the downloaded content to 115 using the existing method
	err = ft.client.RapidUploadOrByOSS(args.DirID, fileName, fileSize, tempFile)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to upload file to 115: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	result := map[string]string{
		"message": "File uploaded successfully from URL",
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to serialize result: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil, nil
}

func (ft *FileTools) uploadFromLocal(ctx context.Context, req *mcp.CallToolRequest, args UploadFromLocalArgs) (*mcp.CallToolResult, any, error) {
	// Open the local file
	file, err := os.Open(args.LocalPath)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to open local file: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}
	defer file.Close()

	// Get file info to determine file size
	fileInfo, err := file.Stat()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get file info: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	// If fileName is empty, use the basename of the local file
	fileName := args.FileName
	if fileName == "" {
		fileName = fileInfo.Name()
	}

	// Get file size
	fileSize := fileInfo.Size()

	// Seek to the beginning of the file
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to seek to beginning of file: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	// Upload the file using the existing method
	err = ft.client.RapidUploadOrByOSS(args.DirID, fileName, fileSize, file)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to upload file to 115: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	result := map[string]string{
		"message": "Local file uploaded successfully",
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to serialize result: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil, nil
}

func (ft *FileTools) downloadFile(ctx context.Context, req *mcp.CallToolRequest, args DownloadFileArgs) (*mcp.CallToolResult, any, error) {
	// Get download info with the specified User-Agent
	downloadInfo, err := ft.client.DownloadWithUA(args.PickCode, args.UserAgent)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get download info: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	// Perform the actual download using the same User-Agent
	reqDownload := ft.client.Client.R()
	if args.UserAgent != "" {
		reqDownload = reqDownload.SetHeader("User-Agent", args.UserAgent)
	}

	resp, err := reqDownload.SetDoNotParseResponse(true).Get(downloadInfo.Url.Url)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to download file: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}
	defer resp.RawBody().Close()

	// Create the local file
	localFile, err := os.Create(args.LocalPath)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to create local file: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}
	defer localFile.Close()

	// Copy the downloaded content to the local file
	_, err = io.Copy(localFile, resp.RawBody())
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to save file locally: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	result := DownloadFileResult{
		Message: "File downloaded successfully",
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to serialize result: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil, nil
}

func (ft *FileTools) getDownloadInfo(ctx context.Context, req *mcp.CallToolRequest, args GetDownloadInfoArgs) (*mcp.CallToolResult, any, error) {
	// Get download info with the specified User-Agent
	downloadInfo, err := ft.client.DownloadWithUA(args.PickCode, args.UserAgent)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get download info: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	result := GetDownloadInfoResult{
		URL:      downloadInfo.Url.Url,
		FileName: downloadInfo.FileName,
		Size:     int64(downloadInfo.FileSize),
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to serialize result: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil, nil
}
