package driver

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type Query func(query *map[string]string)

// QueryLimit set query limit
func QueryLimit(limit int) Query {
	return func(query *map[string]string) {
		(*query)["limit"] = strconv.FormatInt(int64(limit), 10)
	}
}

// QueryOffset set query offset
func QueryOffset(offset int) Query {
	return func(query *map[string]string) {
		(*query)["offset"] = strconv.FormatInt(int64(offset), 10)
	}
}

// GetShareSnapWithUA get share snap info with user agent
func (c *Pan115Client) GetShareSnapWithUA(ua, shareCode, receiveCode, dirID string, Queries ...Query) (*ShareSnapResp, error) {
	if isCalledByAlistV3() {
		return nil, ErrorNotSupportAlist
	}
	result := ShareSnapResp{}
	query := map[string]string{
		"share_code":   shareCode,
		"receive_code": receiveCode,
		"cid":          dirID,
		"limit":        "20",
		"asc":          "0",
		"offset":       "0",
		"format":       "json",
	}

	for _, q := range Queries {
		q(&query)
	}

	req := c.NewRequest().
		SetQueryParams(query).
		SetHeader("referer", BuildShareReferer(shareCode, receiveCode)).
		SetHeader("User-Agent", ua).
		ForceContentType("application/json;charset=UTF-8").
		SetResult(&result)
	resp, err := req.Get(ApiShareSnap)
	if err := CheckErr(err, &result, resp); err != nil {
		return nil, err
	}

	return &result, nil
}

func BuildShareReferer(shareCode, receiveCode string) string {
	return fmt.Sprintf("https://115cdn.com/s/%s?password=%s&", shareCode, receiveCode)
}

// GetShareSnap get share snap info
func (c *Pan115Client) GetShareSnap(shareCode, receiveCode, dirID string, Queries ...Query) (*ShareSnapResp, error) {
	return c.GetShareSnapWithUA("", shareCode, receiveCode, dirID, Queries...)
}

// ShareSendResp is the response of share/send (create share).
// The share content is an immutable snapshot of the directory at creation time:
// deleting the source files afterwards does not affect the share, and adding
// files requires creating a new share.
type ShareSendResp struct {
	BasicResp
	Data struct {
		TotalSize       int64  `json:"total_size"`
		ShareTitle      string `json:"share_title"`
		FileCount       int    `json:"file_count"`
		FolderCount     int    `json:"folder_count"`
		FileCategory    int    `json:"file_category"`
		ReceiveCode     string `json:"receive_code"`
		ShareExDuration string `json:"share_ex_duration"`
		ShareExTime     int64  `json:"share_ex_time"`
		ShareCode       string `json:"share_code"`
		ShareURL        string `json:"share_url"`
		ShareCommand    string `json:"share_command"`
	} `json:"data"`
}

// CreateShare creates a share of fileIDs (comma separated file/directory ids),
// always with ignore_warn=1 (skip the risk-warning pre-check and create
// directly — this is an unattended-flow library). NOTE: a freshly created
// share only lasts 15 days by default — call UpdateShareDuration(code, -1)
// for a permanent one, or use CreatePermanentShare.
func (c *Pan115Client) CreateShare(fileIDs string) (*ShareSendResp, error) {
	if isCalledByAlistV3() {
		return nil, ErrorNotSupportAlist
	}
	if c.UserID == 0 {
		return nil, errors.New("user id unknown, LoginCheck required before creating share")
	}
	result := ShareSendResp{}
	req := c.NewRequest().
		SetFormData(map[string]string{
			"user_id":     strconv.FormatInt(c.UserID, 10),
			"file_ids":    fileIDs,
			"ignore_warn": "1",
			"is_asc":      "0",
			"order":       "file_name",
		}).
		SetHeader("Referer", "https://115.com/").
		ForceContentType("application/json;charset=UTF-8").
		SetResult(&result)
	resp, err := req.Post(ApiShareSend)
	if err := CheckErr(err, &result, resp); err != nil {
		return nil, err
	}
	if result.Data.ShareCode == "" {
		return nil, errors.New("share/send succeeded but no share_code returned")
	}
	return &result, nil
}

// UpdateShareDuration changes the expiry of an existing share.
// duration=-1 means permanent (长期); positive values are days.
func (c *Pan115Client) UpdateShareDuration(shareCode string, duration int) error {
	if isCalledByAlistV3() {
		return ErrorNotSupportAlist
	}
	result := BasicResp{}
	req := c.NewRequest().
		SetFormData(map[string]string{
			"share_code":     shareCode,
			"share_duration": strconv.Itoa(duration),
		}).
		SetHeader("Referer", "https://cdnres.115.com/").
		ForceContentType("application/json;charset=UTF-8").
		SetResult(&result)
	resp, err := req.Post(ApiShareUpdate)
	if err := CheckErr(err, &result, resp); err != nil {
		return err
	}
	return nil
}

// CreatePermanentShare creates a share and immediately extends it to
// permanent (share/send defaults to 15 days, updateshare -1 removes expiry).
// Cookie-only web API: no open-platform equivalent exists.
func (c *Pan115Client) CreatePermanentShare(fileIDs string) (*ShareSendResp, error) {
	result, err := c.CreateShare(fileIDs)
	if err != nil {
		return nil, err
	}
	if err := c.UpdateShareDuration(result.Data.ShareCode, -1); err != nil {
		return nil, errors.Wrapf(err, "extend share %s to permanent failed", result.Data.ShareCode)
	}
	return result, nil
}
