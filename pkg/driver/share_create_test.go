package driver

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// newShareTestEnv routes both share endpoints of a Pan115Client to a mock
// server (the endpoints are absolute URLs, so rewriting happens at the
// transport layer, same trick as the UA tests). The mock dispatches by path
// and records every request's form values and Referer.
type shareTestEnv struct {
	server  *httptest.Server
	client  *Pan115Client
	request []shareTestRequest
}

type shareTestRequest struct {
	path    string
	form    url.Values
	referer string
}

func newShareTestEnv(t *testing.T, sendResp, updateResp string) *shareTestEnv {
	t.Helper()
	env := &shareTestEnv{}
	env.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		env.request = append(env.request, shareTestRequest{
			path:    r.URL.Path,
			form:    r.PostForm,
			referer: r.Header.Get("Referer"),
		})
		switch r.URL.Path {
		case "/share/send":
			_, _ = io.WriteString(w, sendResp)
		case "/share/updateshare":
			_, _ = io.WriteString(w, updateResp)
		default:
			_, _ = fmt.Fprintf(w, `{"state":false,"error":"unexpected path %s"}`, r.URL.Path)
		}
	}))
	t.Cleanup(env.server.Close)
	tr := &recordingTransport{base: &http.Transport{}, mockURL: mustParseURL(t, env.server.URL)}
	env.client = New(WithClient(&http.Client{Transport: tr}))
	env.client.UserID = 6338615
	return env
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

const shareSendOK = `{"state":true,"error":"","errno":0,"data":{"total_size":150671436782,` +
	`"share_title":"测试剧 (2024)","file_count":0,"folder_count":1,"file_category":0,` +
	`"receive_code":"6666","bind_receive_code":1,"share_ex_duration":"15天","share_ex_time":1790477955,` +
	`"skip_login_state":0,"share_code":"swsexqo3hjs","share_url":"https://115cdn.com/s/swsexqo3hjs",` +
	`"share_command":"/swsexqo3hjs-6666/","define_receive_code":1}}`

const shareUpdateOK = `{"state":true,"error":"","errno":0}`

func TestCreateShareForm(t *testing.T) {
	env := newShareTestEnv(t, shareSendOK, shareUpdateOK)
	resp, err := env.client.CreateShare("3189661663297662199", false)
	if err != nil {
		t.Fatal(err)
	}
	req := env.request[0]
	if req.path != "/share/send" {
		t.Fatalf("path = %s, want /share/send", req.path)
	}
	for key, want := range map[string]string{
		"user_id":     "6338615",
		"file_ids":    "3189661663297662199",
		"ignore_warn": "0",
		"order":       "file_name",
	} {
		if got := req.form.Get(key); got != want {
			t.Fatalf("form[%s] = %q, want %q", key, got, want)
		}
	}
	if req.referer != "https://115.com/" {
		t.Fatalf("referer = %q", req.referer)
	}
	if resp.Data.ShareCode != "swsexqo3hjs" || resp.Data.ReceiveCode != "6666" {
		t.Fatalf("share code/receive code = %q/%q", resp.Data.ShareCode, resp.Data.ReceiveCode)
	}
	if resp.Data.ShareExDuration != "15天" {
		t.Fatalf("default duration = %q, want 15天 (caller must extend to permanent)", resp.Data.ShareExDuration)
	}
}

func TestCreateShareIgnoreWarn(t *testing.T) {
	env := newShareTestEnv(t, shareSendOK, shareUpdateOK)
	if _, err := env.client.CreateShare("123", true); err != nil {
		t.Fatal(err)
	}
	if got := env.request[0].form.Get("ignore_warn"); got != "1" {
		t.Fatalf("ignore_warn = %q, want 1", got)
	}
}

func TestCreatePermanentShareTwoSteps(t *testing.T) {
	env := newShareTestEnv(t, shareSendOK, shareUpdateOK)
	resp, err := env.client.CreatePermanentShare("3189661663297662199")
	if err != nil {
		t.Fatal(err)
	}
	if len(env.request) != 2 {
		t.Fatalf("requests = %d, want 2 (send + updateshare)", len(env.request))
	}
	upd := env.request[1]
	if upd.path != "/share/updateshare" {
		t.Fatalf("second path = %s, want /share/updateshare", upd.path)
	}
	if got := upd.form.Get("share_code"); got != "swsexqo3hjs" {
		t.Fatalf("updateshare share_code = %q", got)
	}
	if got := upd.form.Get("share_duration"); got != "-1" {
		t.Fatalf("share_duration = %q, want -1 (permanent)", got)
	}
	if upd.referer != "https://cdnres.115.com/" {
		t.Fatalf("updateshare referer = %q", upd.referer)
	}
	if resp.Data.ShareURL != "https://115cdn.com/s/swsexqo3hjs" {
		t.Fatalf("share url = %q", resp.Data.ShareURL)
	}
}

func TestCreateShareRequiresUserID(t *testing.T) {
	env := newShareTestEnv(t, shareSendOK, shareUpdateOK)
	env.client.UserID = 0
	_, err := env.client.CreateShare("123", true)
	if err == nil || !strings.Contains(err.Error(), "LoginCheck") {
		t.Fatalf("err = %v, want LoginCheck required", err)
	}
	if len(env.request) != 0 {
		t.Fatalf("no request should be sent, got %d", len(env.request))
	}
}

func TestCreateShareApiError(t *testing.T) {
	env := newShareTestEnv(t, `{"state":false,"error":"分享次数已达上限","errno":990001}`, shareUpdateOK)
	_, err := env.client.CreateShare("123", true)
	if err == nil || !strings.Contains(err.Error(), "分享次数已达上限") {
		t.Fatalf("err = %v, want API error passthrough", err)
	}
}

func TestCreatePermanentShareUpdateFailure(t *testing.T) {
	env := newShareTestEnv(t, shareSendOK, `{"state":false,"error":"更新失败","errno":990002}`)
	_, err := env.client.CreatePermanentShare("123")
	if err == nil || !strings.Contains(err.Error(), "swsexqo3hjs") {
		t.Fatalf("err = %v, want share code context for manual cleanup", err)
	}
}
