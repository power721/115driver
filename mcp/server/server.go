package server

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/power721/115driver/mcp/server/tools"
	"github.com/power721/115driver/pkg/driver"
)

// Server represents the 115driver MCP server
type Server struct {
	mcpServer *mcp.Server
	client    *driver.Pan115Client
}

// NewServer creates a new 115driver MCP server
func NewServer() *Server {
	return &Server{
		mcpServer: mcp.NewServer(&mcp.Implementation{
			Name:    "115driver-mcp-server",
			Version: "1.0.0",
		}, nil),
	}
}

// WithClient sets the 115 driver client for the server
func (s *Server) WithClient(client *driver.Pan115Client) *Server {
	s.client = client
	return s
}

// Start runs the MCP server
func (s *Server) Start(ctx context.Context) error {
	// Register all tools
	s.registerTools()

	// Run the server on the stdio transport
	if err := s.mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Printf("Server failed: %v", err)
		return err
	}
	return nil
}

// registerTools registers all available tools with the MCP server
func (s *Server) registerTools() {
	// Register account tools
	accountTools := tools.NewAccountTools(s.client)
	accountTools.RegisterTools(s.mcpServer)

	// Register directory tools
	dirTools := tools.NewDirTools(s.client)
	dirTools.RegisterTools(s.mcpServer)

	// Register file tools
	fileTools := tools.NewFileTools(s.client)
	fileTools.RegisterTools(s.mcpServer)

	// Register recycle tools
	recycleTools := tools.NewRecycleTools(s.client)
	recycleTools.RegisterTools(s.mcpServer)

	// Register share tools
	shareTools := tools.NewShareTools(s.client)
	shareTools.RegisterTools(s.mcpServer)

	// Register search tools
	searchTools := tools.NewSearchTools(s.client)
	searchTools.RegisterTools(s.mcpServer)

	// Register offline tools
	offlineTools := tools.NewOfflineTools(s.client)
	offlineTools.RegisterTools(s.mcpServer)
}
