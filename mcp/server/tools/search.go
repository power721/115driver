package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/power721/115driver/pkg/driver"
)

// SearchTools holds search-related MCP tools
type SearchTools struct {
	client *driver.Pan115Client
}

// NewSearchTools creates a new SearchTools instance
func NewSearchTools(client *driver.Pan115Client) *SearchTools {
	return &SearchTools{
		client: client,
	}
}

// SearchArgs defines arguments for search tool
type SearchArgs struct {
	SearchValue string `json:"search_value" jsonschema:"search keyword"`
	Offset      int    `json:"offset" jsonschema:"offset for pagination, default is 0"`
	Limit       int    `json:"limit" jsonschema:"limit number of results, default is 30"`
	Type        int    `json:"type" jsonschema:"file type filter, 0:all 1:folder 2:document 3:image 4:video 5:audio 6:archive"`
	Order       string `json:"order" jsonschema:"sort field, e.g. file_name, user_ptime"`
	Asc         int    `json:"asc" jsonschema:"ascending order, 0:descending 1:ascending"`
}

// RegisterTools registers search-related tools with the MCP server
func (st *SearchTools) RegisterTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search",
		Description: "Search for files and directories in the 115 cloud storage",
	}, st.search)
}

func (st *SearchTools) search(ctx context.Context, req *mcp.CallToolRequest, args SearchArgs) (*mcp.CallToolResult, any, error) {
	opts := &driver.SearchOption{
		SearchValue: args.SearchValue,
		Offset:      args.Offset,
		Limit:       args.Limit,
		Type:        args.Type,
		Order:       args.Order,
		Asc:         args.Asc,
	}

	result, err := st.client.Search(opts)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to search files: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	// Convert files to a serializable format
	files := make([]map[string]interface{}, len(result.Files))
	for i, file := range result.Files {
		files[i] = map[string]interface{}{
			"file_id":      file.FileID,
			"parent_id":    file.ParentID,
			"name":         file.Name,
			"size":         file.Size,
			"pick_code":    file.PickCode,
			"sha1":         file.Sha1,
			"is_directory": file.IsDirectory,
			"star":         file.Star,
			"create_time":  file.CreateTime,
			"update_time":  file.UpdateTime,
		}
	}

	response := map[string]interface{}{
		"count":     result.Count,
		"files":     files,
		"offset":    result.Offset,
		"page_size": result.PageSize,
		"order":     result.Order,
		"is_asc":    result.IsAsc,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to serialize search results: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(responseJSON),
			},
		},
	}, nil, nil
}
