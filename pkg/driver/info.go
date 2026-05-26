package driver

import (
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
)

// GetInfo get space info and login device info.
func (c *Pan115Client) GetInfo() (InfoData, error) {
	result := InfoResponse{}
	req := c.NewRequest().
		SetResult(&result).
		ForceContentType("application/json;charset=UTF-8")

	resp, err := req.Get(ApiFileIndexInfo)

	if err = CheckErr(err, &result, resp); err != nil {
		return InfoData{}, err
	}
	return result.Data, nil
}

type InfoResponse struct {
	BasicResp
	Data InfoData `json:"data"`
}

type InfoData struct {
	SpaceInfo        SpaceInfo        `json:"space_info"`
	LoginDevicesInfo LoginDevicesInfo `json:"login_devices_info"`
	ImeiInfo         bool             `json:"imei_info"`
}

type TotalSize struct {
	Size       int64  `json:"size"`
	SizeFormat string `json:"size_format"`
}

func (s *TotalSize) UnmarshalJSON(b []byte) error {
	size, sizeFormat, err := unmarshalSpaceSize(b)
	if err != nil {
		return err
	}
	s.Size = size
	s.SizeFormat = sizeFormat
	return nil
}

type RemainSize struct {
	Size       int64  `json:"size"`
	SizeFormat string `json:"size_format"`
}

func (s *RemainSize) UnmarshalJSON(b []byte) error {
	size, sizeFormat, err := unmarshalSpaceSize(b)
	if err != nil {
		return err
	}
	s.Size = size
	s.SizeFormat = sizeFormat
	return nil
}

type UseSize struct {
	Size       int64  `json:"size"`
	SizeFormat string `json:"size_format"`
}

func (s *UseSize) UnmarshalJSON(b []byte) error {
	size, sizeFormat, err := unmarshalSpaceSize(b)
	if err != nil {
		return err
	}
	s.Size = size
	s.SizeFormat = sizeFormat
	return nil
}

func unmarshalSpaceSize(b []byte) (int64, string, error) {
	var raw struct {
		Size       json.RawMessage `json:"size"`
		SizeFormat string          `json:"size_format"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return 0, "", err
	}
	size, err := parseSpaceSize(raw.Size)
	if err != nil {
		return 0, "", err
	}
	return size, raw.SizeFormat, nil
}

func parseSpaceSize(b []byte) (int64, error) {
	if len(b) == 0 || string(b) == "null" {
		return 0, nil
	}
	var value string
	if b[0] == '"' {
		if err := json.Unmarshal(b, &value); err != nil {
			return 0, err
		}
	} else {
		value = string(b)
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}
	size, _, err := big.ParseFloat(value, 10, 256, big.ToZero)
	if err != nil {
		return 0, err
	}
	sizeInt, _ := size.Int(nil)
	if sizeInt == nil {
		return 0, strconv.ErrRange
	}
	if !sizeInt.IsInt64() {
		return 0, strconv.ErrRange
	}
	return sizeInt.Int64(), nil
}

type SpaceInfo struct {
	AllTotal  TotalSize  `json:"all_total"`
	AllRemain RemainSize `json:"all_remain"`
	AllUse    UseSize    `json:"all_use"`
}

type LastDevice struct {
	IP       string `json:"ip"`
	Device   string `json:"device"`
	DeviceID string `json:"device_id"`
	Network  string `json:"network"`
	Os       string `json:"os"`
	City     string `json:"city"`
	Utime    int    `json:"utime"`
}

type Device struct {
	IsCurrent int    `json:"is_current"`
	Ssoent    string `json:"ssoent"`
	Utime     int    `json:"utime"`
	Device    string `json:"device"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Desc      string `json:"desc"`
	IP        string `json:"ip"`
	City      string `json:"city"`
	IsUnusual int    `json:"is_unusual"`
}

type LoginDevicesInfo struct {
	Last LastDevice `json:"last"`
	List []Device   `json:"list"`
}
