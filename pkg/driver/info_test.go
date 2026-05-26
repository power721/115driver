package driver

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfoResponseUnmarshalsDecimalSpaceSizes(t *testing.T) {
	const payload = `{
		"state": true,
		"data": {
			"space_info": {
				"all_total": {"size": 155555555555555.25, "size_format": "141.49TB"},
				"all_remain": {"size": "123.75", "size_format": "123B"},
				"all_use": {"size": 42, "size_format": "42B"}
			},
			"login_devices_info": {},
			"imei_info": false
		}
	}`

	var resp InfoResponse
	require.NoError(t, json.Unmarshal([]byte(payload), &resp))
	assert.Equal(t, int64(155555555555555), resp.Data.SpaceInfo.AllTotal.Size)
	assert.Equal(t, int64(123), resp.Data.SpaceInfo.AllRemain.Size)
	assert.Equal(t, int64(42), resp.Data.SpaceInfo.AllUse.Size)
}

func TestInfoResponseUnmarshalsExponentSpaceSizes(t *testing.T) {
	const payload = `{
		"state": true,
		"data": {
			"space_info": {
				"all_total": {"size": 1.5e6, "size_format": "1.43MB"},
				"all_remain": {"size": "2.5e3", "size_format": "2.44KB"},
				"all_use": {"size": 42, "size_format": "42B"}
			},
			"login_devices_info": {},
			"imei_info": false
		}
	}`

	var resp InfoResponse
	require.NoError(t, json.Unmarshal([]byte(payload), &resp))
	assert.Equal(t, int64(1500000), resp.Data.SpaceInfo.AllTotal.Size)
	assert.Equal(t, int64(2500), resp.Data.SpaceInfo.AllRemain.Size)
	assert.Equal(t, int64(42), resp.Data.SpaceInfo.AllUse.Size)
}

func TestInfoResponseRejectsOutOfRangeSpaceSizes(t *testing.T) {
	const payload = `{
		"state": true,
		"data": {
			"space_info": {
				"all_total": {"size": 1e19, "size_format": "9.09EB"},
				"all_remain": {"size": 0, "size_format": "0B"},
				"all_use": {"size": 42, "size_format": "42B"}
			},
			"login_devices_info": {},
			"imei_info": false
		}
	}`

	var resp InfoResponse
	err := json.Unmarshal([]byte(payload), &resp)
	require.Error(t, err)
	assert.ErrorIs(t, err, strconv.ErrRange)
}

func TestInfoResponseRejectsInfiniteSpaceSizes(t *testing.T) {
	const payload = `{
		"state": true,
		"data": {
			"space_info": {
				"all_total": {"size": "Inf", "size_format": "Inf"},
				"all_remain": {"size": 0, "size_format": "0B"},
				"all_use": {"size": 42, "size_format": "42B"}
			},
			"login_devices_info": {},
			"imei_info": false
		}
	}`

	var resp InfoResponse
	require.Error(t, json.Unmarshal([]byte(payload), &resp))
}
