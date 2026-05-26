package cmd

import (
	"path/filepath"
	"testing"
)

func TestResolveDownloadTargetPath_UsesExplicitFilePath(t *testing.T) {
	got := resolveDownloadTargetPath(filepath.Join("/tmp", "custom-name.mp4"), "remote-name.mp4")
	want := filepath.Join("/tmp", "custom-name.mp4")
	if got != want {
		t.Fatalf("unexpected target path: got %q want %q", got, want)
	}
}

func TestResolveDownloadTargetPath_AppendsNameForDirectoryHint(t *testing.T) {
	got := resolveDownloadTargetPath("/tmp/downloads/", "remote-name.mp4")
	want := filepath.Join("/tmp/downloads", "remote-name.mp4")
	if got != want {
		t.Fatalf("unexpected target path: got %q want %q", got, want)
	}
}
