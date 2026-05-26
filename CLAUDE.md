# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

115driver is a Go driver package for 115 cloud storage (a Chinese cloud service). It provides a comprehensive API for file operations, directory management, search, offline downloads, and includes an MCP (Model Context Protocol) server for AI application integration.

## Development Commands

### Build
```bash
# Build the main MCP server binary
go build -o mcp/115driver-mcp-server ./mcp/main.go

# Install dependencies
go mod tidy
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./pkg/driver/

# Run tests with verbose output
go test -v ./pkg/driver/

# Run a specific test
go test -v ./pkg/driver/ -run TestFunctionName
```

### Running the MCP Server
```bash
# The MCP server requires cookie authentication
./mcp/115driver-mcp-server --cookie="UID=xxx;CID=xxx;SEID=xxx;KID=xxx"

# The server communicates via stdin/stdout using JSON-RPC 2.0
```

## Architecture

### Core Package Structure

**`pkg/driver/`** - Main 115 cloud storage driver
- `client.go` - Core `Pan115Client` using resty/v2 for HTTP requests
- `options.go` - Options pattern for flexible client configuration
- `consts.go` - Constants including cookie names, user agents, API endpoints
- `file.go`, `dir.go`, `upload.go`, `download.go` - File operations
- `search.go` - Search functionality with filters and pagination
- `offline.go` - Offline download management
- `share.go` - Share link creation and management
- `recycle.go` - Recycle bin operations
- `api.go` - API endpoint definitions

**`pkg/crypto/`** - Cryptographic implementations
- 115-specific encryption protocols
- ECDH key exchange, AES encryption
- Required for secure communication with 115 services

**`mcp/`** - Model Context Protocol server
- `main.go` - Entry point with cookie flag authentication
- `server/server.go` - MCP server core with tool registration
- `server/tools/` - Individual MCP tool implementations (one file per tool category)

### Authentication

The driver uses cookie-based authentication with four required components:
- `UID` - User ID
- `CID` - Credential ID
- `SEID` - Session ID
- `KID` - Key ID

These can be imported from a cookie string using `FromCookie()` or set directly via the `Credential` struct.

### MCP Tool Architecture

Each tool category is implemented as a separate file in `mcp/server/tools/`:
- `dir.go` - Directory listing and creation
- `file.go` - File operations (rename, move, delete)
- `search.go` - File search with filters
- `offline.go` - Offline download management
- `share.go` - Share link operations
- `recycle.go` - Recycle bin management
- `upload_download.go` - Upload and download utilities

Tools follow a consistent pattern:
1. Create a struct holding the `Pan115Client`
2. Define argument structs with JSON schema tags
3. Implement the handler function
4. Register via `RegisterTools()` method

### Key Design Patterns

**Options Pattern**: Used extensively for client configuration
```go
client := driver.New(
    driver.UA(driver.UA115Browser),
    driver.WithDebug(),
).ImportCredential(cr)
```

**Error Handling**: Custom error types in `error.go` with 115 API error code mapping

**Response Structures**: All API responses have corresponding structs with JSON field mappings

### Known Typename (for backward compatibility)

Note: The codebase contains intentional typos maintained for backward compatibility:
- `Defalut()` function in `client.go` (should be `Default()`) - DEPRECATED: use `Default()` instead
- `UADefalut` constant in `consts.go` (should be `UA115Default`) - DEPRECATED: use `UADefault` instead
- `DefalutUploadMultipartOptions()` function in `option.go` (should be `DefaultUploadMultipartOptions()`) - DEPRECATED: use `DefaultUploadMultipartOptions()` instead

These functions are marked as deprecated and will be removed in a future version.

### MCP Server Details

The MCP server uses stdin/stdout transport for JSON-RPC 2.0 communication. All tool responses are returned as JSON strings in text content blocks. The server registers tools automatically on startup and handles authentication via the required `--cookie` flag.

### File Operations Reference

- **Rapid Upload**: SHA1-based deduplication to avoid re-uploading existing files
- **Multipart Upload**: Large file support via Aliyun OSS integration
- **Search**: Supports filtering by type (folder, document, image, video, audio, archive), date range, and sorting
- **Offline Downloads**: Supports HTTP, ED2K, and magnet links
