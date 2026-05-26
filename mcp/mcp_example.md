# 115driver MCP Server Example

This document provides examples of how to use the 115driver MCP server.

## Starting the Server

To start the MCP server, run:

```bash
go build -o mcp-server mcp/main.go
```

Then run the server with your 115 cookies:

```bash
./mcp-server --cookie="UID=your_uid;CID=your_cid;SEID=your_seid"
```

The server will listen on stdin/stdout for MCP requests.

## Available Tools

### Account Tools

1. `getAccountInfo`: Get current account, storage space, and login device info
   - Parameters: none
   - Returns: `user`, `space.total`, `space.remain`, `space.used`, `login_devices`, and `imei_info`

### Directory Tools

1. `listDirectory`: List files and directories in a specific directory
   - Parameters:
     - `dir_id` (string): Directory ID to list, default is root directory (0)
     - `offset` (int64): Offset for pagination, default is 0
     - `limit` (int64): Number of items to return, default is all items

2. `mkdir`: Create a new directory
   - Parameters:
     - `parent_id` (string): Parent directory ID
     - `name` (string): Name of the new directory

### File Tools

1. `delete`: Delete files or directories
   - Parameters:
     - `file_ids` (array of strings): IDs of files or directories to delete

2. `rename`: Rename a file or directory
   - Parameters:
     - `file_id` (string): ID of file or directory to rename
     - `new_name` (string): New name for the file or directory

3. `move`: Move files or directories to another directory
   - Parameters:
     - `dir_id` (string): Target directory ID
     - `file_ids` (array of strings): IDs of files or directories to move

4. `copy`: Copy files or directories to another directory
   - Parameters:
     - `dir_id` (string): Target directory ID
     - `file_ids` (array of strings): IDs of files or directories to copy

5. `stat`: Get detailed information about a file or directory
   - Parameters:
     - `file_id` (string): ID of file or directory to get info

### Recycle Bin Tools

1. `listRecycleBin`: List items in the recycle bin
   - Parameters:
     - `offset` (string): Offset for pagination, default is 0
     - `limit` (string): Number of items to return, default is 40

2. `revertRecycleBin`: Revert items from the recycle bin
   - Parameters:
     - `item_ids` (array of strings): IDs of items to revert

3. `cleanRecycleBin`: Clean items from the recycle bin
   - Parameters:
     - `password` (string): Password for cleaning recycle bin
     - `item_ids` (array of strings): IDs of items to clean

### Share Tools

1. `getShareSnap`: Get shared files and directories snapshot information
   - Parameters:
     - `share_code` (string): Share code
     - `receive_code` (string): Receive code
     - `dir_id` (string): Directory ID to list, default is root directory

### Search Tools

1. `search`: Search for files and directories in the 115 cloud storage
   - Parameters:
     - `search_value` (string): Search keyword
     - `offset` (int): Offset for pagination, default is 0
     - `limit` (int): Limit number of results, default is 30
     - `type` (int): File type filter, 0:all 1:folder 2:document 3:image 4:video 5:audio 6:archive
     - `order` (string): Sort field, e.g. file_name, user_ptime
     - `asc` (int): Ascending order, 0:descending 1:ascending

### Offline Download Tools

1. `listOfflineTasks`: List offline download tasks
   - Parameters:
     - `page` (int64): Page number for pagination, default is 1

2. `addOfflineTaskURIs`: Add offline tasks by download URIs, supports http, ed2k, magnet
   - Parameters:
     - `uris` (array of strings): Download URIs, supports http, ed2k, magnet
     - `save_dir_id` (string): Directory ID to save downloaded files

3. `deleteOfflineTasks`: Delete offline tasks
   - Parameters:
     - `hashes` (array of strings): Task hashes to delete
     - `delete_files` (bool): Whether to delete associated files, default is false

4. `clearOfflineTasks`: Clear offline tasks
   - Parameters:
     - `clear_flag` (int64): Clear flag, 0: clear completed tasks, 1: clear all tasks

## Example Request/Response

### Basic Directory Listing Request

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "tools/call",
  "params": {
    "name": "listDirectory",
    "arguments": {
      "dir_id": "0"
    }
  }
}
```

### Basic Directory Listing Response

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "[{\"Id\":\"12345\",\"Name\":\"Documents\",\"Size\":0,\"Type\":1,\"CreateTime\":1234567890},{\"Id\":\"67890\",\"Name\":\"image.jpg\",\"Size\":1024,\"Type\":0,\"CreateTime\":1234567891}]"
      }
    ]
  }
}
```

### Search Request

Search for documents containing the word "report":

```json
{
  "jsonrpc": "2.0",
  "id": "2",
  "method": "tools/call",
  "params": {
    "name": "search",
    "arguments": {
      "search_value": "report",
      "limit": 10,
      "type": 2
    }
  }
}
```

### Search Response

```json
{
  "jsonrpc": "2.0",
  "id": "2",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "{\"count\":5,\"files\":[{\"file_id\":\"12345\",\"name\":\"report.pdf\",\"size\":1024,\"is_directory\":false},{\"file_id\":\"12346\",\"name\":\"annual-report.docx\",\"size\":2048,\"is_directory\":false}],\"offset\":0,\"page_size\":10}"
      }
    ]
  }
}
```

## Notes

1. Valid cookies must be provided for authentication
2. The file list in the response is returned as a text content in JSON string format
3. The Type field indicates the file type: typically 1 for directories and 0 for files
4. All tools follow the standard MCP tool calling conventions
