package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/power721/115driver/pkg/driver"
)

// RecycleTools holds recycle bin-related MCP tools
type RecycleTools struct {
	client *driver.Pan115Client
}

// NewRecycleTools creates a new RecycleTools instance
func NewRecycleTools(client *driver.Pan115Client) *RecycleTools {
	return &RecycleTools{
		client: client,
	}
}

// ListRecycleArgs defines arguments for listing recycle bin items
type ListRecycleArgs struct {
	Offset string `json:"offset" jsonschema:"offset for pagination, default is 0"`
	Limit  string `json:"limit" jsonschema:"number of items to return, default is 40"`
}

// RevertRecycleArgs defines arguments for reverting recycle bin items
type RevertRecycleArgs struct {
	ItemIDs []string `json:"item_ids" jsonschema:"IDs of items to revert"`
}

// CleanRecycleArgs defines arguments for cleaning recycle bin items
type CleanRecycleArgs struct {
	Password string   `json:"password" jsonschema:"password for cleaning recycle bin"`
	ItemIDs  []string `json:"item_ids" jsonschema:"IDs of items to clean"`
}

// RegisterTools registers recycle bin-related tools with the MCP server
func (rt *RecycleTools) RegisterTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "listRecycleBin",
		Description: "List items in the recycle bin",
	}, rt.listRecycleBin)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "revertRecycleBin",
		Description: "Revert items from the recycle bin",
	}, rt.revertRecycleBin)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "cleanRecycleBin",
		Description: "Clean items from the recycle bin",
	}, rt.cleanRecycleBin)
}

func (rt *RecycleTools) listRecycleBin(ctx context.Context, req *mcp.CallToolRequest, args ListRecycleArgs) (*mcp.CallToolResult, any, error) {
	offset, err := strconv.Atoi(args.Offset)
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(args.Limit)
	if err != nil {
		limit = 40
	}

	items, err := rt.client.ListRecycleBin(offset, limit)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to list recycle bin: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	resultJSON, err := json.Marshal(items)
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

func (rt *RecycleTools) revertRecycleBin(ctx context.Context, req *mcp.CallToolRequest, args RevertRecycleArgs) (*mcp.CallToolResult, any, error) {
	if len(args.ItemIDs) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "No item IDs provided",
				},
			},
			IsError: true,
		}, nil, nil
	}

	err := rt.client.RevertRecycleBin(args.ItemIDs...)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to revert recycle bin items: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "Items reverted successfully",
			},
		},
	}, nil, nil
}

func (rt *RecycleTools) cleanRecycleBin(ctx context.Context, req *mcp.CallToolRequest, args CleanRecycleArgs) (*mcp.CallToolResult, any, error) {
	if len(args.ItemIDs) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "No item IDs provided",
				},
			},
			IsError: true,
		}, nil, nil
	}

	err := rt.client.CleanRecycleBin(args.Password, args.ItemIDs...)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to clean recycle bin items: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "Items cleaned successfully",
			},
		},
	}, nil, nil
}
