package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/power721/115driver/pkg/driver"
)

// ShareTools holds share-related MCP tools
type ShareTools struct {
	client *driver.Pan115Client
}

// NewShareTools creates a new ShareTools instance
func NewShareTools(client *driver.Pan115Client) *ShareTools {
	return &ShareTools{
		client: client,
	}
}

// GetShareSnapArgs defines arguments for get share snap tool
type GetShareSnapArgs struct {
	ShareCode   string `json:"share_code" jsonschema:"required,share code"`
	ReceiveCode string `json:"receive_code" jsonschema:"required,receive code"`
	DirID       string `json:"dir_id" jsonschema:"directory ID to list, default is root directory"`
	Offset      int    `json:"offset,omitempty" jsonschema:"offset for pagination, default is 0"`
	Limit       int    `json:"limit,omitempty" jsonschema:"number of items to return, default is 20"`
}

// RegisterTools registers share-related tools with the MCP server
func (st *ShareTools) RegisterTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "getShareSnap",
		Description: "Get shared files and directories snapshot information",
	}, st.getShareSnap)
}

func (st *ShareTools) getShareSnap(ctx context.Context, req *mcp.CallToolRequest, args GetShareSnapArgs) (*mcp.CallToolResult, any, error) {
	var (
		result *driver.ShareSnapResp
		err    error
	)

	// Prepare queries
	queries := make([]driver.Query, 0)
	if args.Limit > 0 {
		queries = append(queries, driver.QueryLimit(args.Limit))
	}
	if args.Offset > 0 {
		queries = append(queries, driver.QueryOffset(args.Offset))
	}

	result, err = st.client.GetShareSnap(args.ShareCode, args.ReceiveCode, args.DirID, queries...)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get share snap: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	// Serialize result to JSON
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
