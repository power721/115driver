package driver

import (
	"strconv"
)

// SearchOption defines options for search
type SearchOption struct {
	// Offset for pagination
	Offset int
	// Limit number of results
	Limit int
	// SearchValue search keyword
	SearchValue string
	// Date filter
	Date string
	// Aid area ID
	Aid string
	// Cid category ID
	Cid string
	// PickCode pickcode
	PickCode string
	// Type file type filter 0:all 1:folder 2:document 3:image 4:video 5:audio 6:archive
	Type int
	// CountFolders whether to count folders
	CountFolders int
	// Source source filter
	Source string
	// Star star file only
	Star string
	// Suffix file suffix filter
	Suffix string
	// Order sort field
	Order string
	// Asc ascending order 0:descending 1:ascending
	Asc int
}

// SearchFile represents a file in search results
type SearchFile struct {
	// File ID
	FileID string `json:"fid"`
	// Category ID
	CategoryID IntString `json:"cid"`
	// File name
	Name string `json:"n"`
	// File size
	Size StringInt64 `json:"s"`
	// SHA1 hash
	Sha1 string `json:"sha"`
	// PickCode
	PickCode string `json:"pc"`
	// Is directory
	IsDirectory int `json:"fc"`
	// Is starred file
	IsStar StringInt `json:"m"`
	// Update time
	UpdateTime string `json:"t"`
	// Create time
	CreateTime StringInt64 `json:"tp"`
	// File type icon
	Icon string `json:"ico"`
	// Highlighted file name
	HighlightName string `json:"ns"`
	// File labels
	Labels []*LabelInfo `json:"fl"`
	// Thumbnail URL
	ThumbURL string `json:"u"`
}

// SearchResult represents search results
type SearchResult struct {
	// File list
	Files []File `json:"data"`
	// Total count
	Count int `json:"count"`
	// File count
	FileCount int `json:"file_count"`
	// Folder count
	FolderCount int `json:"folder_count"`
	// Page size
	PageSize int `json:"page_size"`
	// Offset
	Offset int `json:"offset"`
	// Sort field
	Order string `json:"order"`
	// Ascending order
	IsAsc int `json:"is_asc"`
}

// Search searches for files using given options
func (c *Pan115Client) Search(opts *SearchOption) (*SearchResult, error) {
	result := FileListResp{}
	params := map[string]string{
		"aid":           "7",
		"cid":           "0",
		"format":        "json",
		"offset":        "0",
		"limit":         "30",
		"search_value":  "",
		"type":          "0",
		"count_folders": "1",
		"o":             "file_name",
		"asc":           "1",
	}

	// Set search parameters
	if opts != nil {
		if opts.Offset >= 0 {
			params["offset"] = strconv.Itoa(opts.Offset)
		}
		if opts.Limit > 0 {
			params["limit"] = strconv.Itoa(opts.Limit)
		}
		if opts.SearchValue != "" {
			params["search_value"] = opts.SearchValue
		}
		if opts.Date != "" {
			params["date"] = opts.Date
		}
		if opts.Aid != "" {
			params["aid"] = opts.Aid
		}
		if opts.Cid != "" {
			params["cid"] = opts.Cid
		}
		if opts.PickCode != "" {
			params["pick_code"] = opts.PickCode
		}
		if opts.Type > 0 {
			params["type"] = strconv.Itoa(opts.Type)
		}
		if opts.CountFolders > 0 {
			params["count_folders"] = strconv.Itoa(opts.CountFolders)
		}
		if opts.Source != "" {
			params["source"] = opts.Source
		}
		if opts.Star != "" {
			params["star"] = opts.Star
		}
		if opts.Suffix != "" {
			params["suffix"] = opts.Suffix
		}
		if opts.Order != "" {
			params["o"] = opts.Order
		}
		params["asc"] = strconv.Itoa(opts.Asc)
	}

	req := c.NewRequest().
		SetQueryParams(params).
		SetResult(&result).
		ForceContentType("application/json;charset=UTF-8")

	resp, err := req.Get(ApiFileSearch)
	if err = CheckErr(err, &result, resp); err != nil {
		return nil, err
	}

	// Convert results
	searchResult := &SearchResult{
		Count:       result.Count,
		FileCount:   0, // Not available in FileListResp
		FolderCount: 0, // Not available in FileListResp
		PageSize:    result.PageSize,
		Offset:      result.Offset,
		Order:       result.Order,
		IsAsc:       result.IsAsc,
		Files:       make([]File, 0, len(result.Files)),
	}

	for _, fileInfo := range result.Files {
		searchResult.Files = append(searchResult.Files, *(&File{}).from(&fileInfo))
	}

	return searchResult, nil
}