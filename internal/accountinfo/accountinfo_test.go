package accountinfo

import (
	"testing"

	"github.com/power721/115driver/pkg/driver"
	"github.com/stretchr/testify/assert"
)

func TestFromDriverDataMapsAccountInfo(t *testing.T) {
	user := &driver.UserInfo{
		UserID:   12345,
		UserName: "alice",
		Vip:      1,
		Expire:   1770000000,
	}
	info := driver.InfoData{
		SpaceInfo: driver.SpaceInfo{
			AllTotal:  driver.TotalSize{Size: 1000, SizeFormat: "1000B"},
			AllRemain: driver.RemainSize{Size: 250, SizeFormat: "250B"},
			AllUse:    driver.UseSize{Size: 750, SizeFormat: "750B"},
		},
		LoginDevicesInfo: driver.LoginDevicesInfo{
			Last: driver.LastDevice{IP: "127.0.0.1", Device: "Browser"},
			List: []driver.Device{
				{Device: "Browser", IsCurrent: 1},
			},
		},
		ImeiInfo: true,
	}

	got := FromDriverData(user, info)

	assert.Equal(t, int64(12345), got.User.UserID)
	assert.Equal(t, "alice", got.User.Username)
	assert.Equal(t, 1, got.User.VIP)
	assert.Equal(t, 1770000000, got.User.Expire)
	assert.Equal(t, int64(1000), got.Space.Total.Size)
	assert.Equal(t, "1000B", got.Space.Total.SizeFormat)
	assert.Equal(t, int64(250), got.Space.Remain.Size)
	assert.Equal(t, int64(750), got.Space.Used.Size)
	assert.Equal(t, "127.0.0.1", got.LoginDevices.Last.IP)
	assert.Len(t, got.LoginDevices.List, 1)
	assert.True(t, got.ImeiInfo)
}
