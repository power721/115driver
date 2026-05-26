package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/power721/115driver/internal/accountinfo"
	"github.com/power721/115driver/pkg/driver"
)

type AccountTools struct {
	client *driver.Pan115Client
}

func NewAccountTools(client *driver.Pan115Client) *AccountTools {
	return &AccountTools{
		client: client,
	}
}

func (at *AccountTools) RegisterTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "getAccountInfo",
		Description: "Get current account, storage space, and login device info",
	}, at.getAccountInfo)
}

func (at *AccountTools) getAccountInfo(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
	userInfo, err := at.client.GetUser()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get user info: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}
	info, err := at.client.GetInfo()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to get account info: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	responseJSON, err := marshalAccountInfoResult(accountinfo.FromDriverData(userInfo, info))
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to serialize account info: %v", err),
				},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: responseJSON,
			},
		},
	}, nil, nil
}

func marshalAccountInfoResult(info accountinfo.AccountInfo) (string, error) {
	responseJSON, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	return string(responseJSON), nil
}
