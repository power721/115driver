package accountinfo

import (
	"encoding/json"

	"github.com/power721/115driver/pkg/driver"
)

type AccountInfo struct {
	User         User                    `json:"user"`
	Space        Space                   `json:"space"`
	LoginDevices driver.LoginDevicesInfo `json:"login_devices"`
	ImeiInfo     json.RawMessage         `json:"imei_info"`
}

type User struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	VIP      int    `json:"vip"`
	Expire   int    `json:"expire"`
}

type Space struct {
	Total  SizeInfo `json:"total"`
	Remain SizeInfo `json:"remain"`
	Used   SizeInfo `json:"used"`
}

type SizeInfo struct {
	Size       int64  `json:"size"`
	SizeFormat string `json:"size_format"`
}

func FromDriverData(user *driver.UserInfo, info driver.InfoData) AccountInfo {
	account := AccountInfo{
		Space: Space{
			Total: SizeInfo{
				Size:       info.SpaceInfo.AllTotal.Size,
				SizeFormat: info.SpaceInfo.AllTotal.SizeFormat,
			},
			Remain: SizeInfo{
				Size:       info.SpaceInfo.AllRemain.Size,
				SizeFormat: info.SpaceInfo.AllRemain.SizeFormat,
			},
			Used: SizeInfo{
				Size:       info.SpaceInfo.AllUse.Size,
				SizeFormat: info.SpaceInfo.AllUse.SizeFormat,
			},
		},
		LoginDevices: info.LoginDevicesInfo,
		ImeiInfo:     info.ImeiInfo,
	}
	if user != nil {
		account.User = User{
			UserID:   user.UserID,
			Username: user.UserName,
			VIP:      user.Vip,
			Expire:   user.Expire,
		}
	}
	return account
}
