package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	return path
}

func TestReadProfileConfig_KeyExists(t *testing.T) {
	configPath := writeTestConfig(t, `
default_profile = "main"

[profiles.main]
cookie = "UID=123;CID=abc"
default_offline_save_dir = "123456"

[profiles.work]
cookie = "UID=456"
default_offline_save_dir = "789012"
`)
	got := ReadProfileConfig(configPath, "main", "default_offline_save_dir")
	if got != "123456" {
		t.Fatalf("expected '123456', got '%s'", got)
	}
}

func TestReadProfileConfig_KeyMissing(t *testing.T) {
	configPath := writeTestConfig(t, `
[profiles.main]
cookie = "UID=123;CID=abc"
`)
	got := ReadProfileConfig(configPath, "main", "nonexistent_key")
	if got != "" {
		t.Fatalf("expected empty string, got '%s'", got)
	}
}

func TestReadProfileConfig_ConfigNotExist(t *testing.T) {
	got := ReadProfileConfig("/nonexistent/path/config.toml", "main", "any_key")
	if got != "" {
		t.Fatalf("expected empty string for missing config, got '%s'", got)
	}
}

func TestReadProfileConfig_CustomProfile(t *testing.T) {
	configPath := writeTestConfig(t, `
[profiles.personal]
cookie = "UID=789"
default_offline_save_dir = "personal_dir"

[profiles.main]
default_offline_save_dir = "default_dir"
`)
	got := ReadProfileConfig(configPath, "personal", "default_offline_save_dir")
	if got != "personal_dir" {
		t.Fatalf("expected 'personal_dir', got '%s'", got)
	}
}

func TestReadProfileConfig_EmptyProfileFallsBackToHardcodedDefault(t *testing.T) {
	configPath := writeTestConfig(t, `
[profiles.main]
default_offline_save_dir = "main_dir"

[profiles.work]
default_offline_save_dir = "work_dir"
`)
	got := ReadProfileConfig(configPath, "", "default_offline_save_dir")
	if got != "main_dir" {
		t.Fatalf("expected 'main_dir' (hardcoded default), got '%s'", got)
	}
}

func TestReadProfileConfig_EmptyProfileNoDefault(t *testing.T) {
	configPath := writeTestConfig(t, `
[profiles.main]
default_offline_save_dir = "main_dir"
`)
	got := ReadProfileConfig(configPath, "", "default_offline_save_dir")
	if got != "main_dir" {
		t.Fatalf("expected 'main_dir' (hardcoded default), got '%s'", got)
	}
}
