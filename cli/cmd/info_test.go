package cmd

import (
	"strings"
	"testing"

	"github.com/power721/115driver/internal/accountinfo"
	"github.com/stretchr/testify/assert"
)

func TestFormatAccountInfoText(t *testing.T) {
	info := accountinfo.AccountInfo{
		User: accountinfo.User{
			UserID:   12345,
			Username: "alice",
			VIP:      1,
			Expire:   1770000000,
		},
		Space: accountinfo.Space{
			Total:  accountinfo.SizeInfo{Size: 1000, SizeFormat: "1000B"},
			Remain: accountinfo.SizeInfo{Size: 250, SizeFormat: "250B"},
			Used:   accountinfo.SizeInfo{Size: 750, SizeFormat: "750B"},
		},
	}

	got := formatAccountInfoText(info)

	assert.True(t, strings.Contains(got, "User ID:  12345"))
	assert.True(t, strings.Contains(got, "Username: alice"))
	assert.True(t, strings.Contains(got, "VIP:      1 (expires: 1770000000)"))
	assert.True(t, strings.Contains(got, "Total:    1000B (1000 bytes)"))
	assert.True(t, strings.Contains(got, "Remain:   250B (250 bytes)"))
	assert.True(t, strings.Contains(got, "Used:     750B (750 bytes)"))
}
