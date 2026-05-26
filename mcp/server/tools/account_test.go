package tools

import (
	"encoding/json"
	"testing"

	"github.com/power721/115driver/internal/accountinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalAccountInfoResult(t *testing.T) {
	info := accountinfo.AccountInfo{
		User: accountinfo.User{
			UserID:   12345,
			Username: "alice",
		},
		Space: accountinfo.Space{
			Total:  accountinfo.SizeInfo{Size: 1000, SizeFormat: "1000B"},
			Remain: accountinfo.SizeInfo{Size: 250, SizeFormat: "250B"},
			Used:   accountinfo.SizeInfo{Size: 750, SizeFormat: "750B"},
		},
	}

	got, err := marshalAccountInfoResult(info)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal([]byte(got), &decoded))
	assert.Contains(t, decoded, "user")
	assert.Contains(t, decoded, "space")
	assert.Contains(t, decoded, "login_devices")
	assert.Contains(t, decoded, "imei_info")
}
