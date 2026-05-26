package driver

import (
	"net/http"
	"testing"
)

func TestBuildDownloadHeaders_MergesResponseCookiesIntoCookieHeader(t *testing.T) {
	requestHeaders := http.Header{}
	requestHeaders.Set("User-Agent", UA115Browser)
	requestHeaders.Set("Cookie", "UID=1; CID=2")

	responseCookies := []*http.Cookie{
		{
			Name:  "download_token",
			Value: "abc123",
		},
	}

	got := buildDownloadHeaders(requestHeaders, responseCookies)

	if got.Get("User-Agent") != UA115Browser {
		t.Fatalf("unexpected user agent: %q", got.Get("User-Agent"))
	}
	if got.Get("Cookie") != "UID=1; CID=2; download_token=abc123" {
		t.Fatalf("unexpected cookie header: %q", got.Get("Cookie"))
	}
}
