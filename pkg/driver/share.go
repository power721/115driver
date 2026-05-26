package driver

import (
	"fmt"
	"strconv"
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
		ForceContentType("application/json;charset=UTF-8").
		SetResult(&result)
	if ua != "" {
		req.SetHeader("User-Agent", ua)
	}
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
