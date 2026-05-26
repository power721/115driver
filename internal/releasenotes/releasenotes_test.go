package releasenotes

import (
	"errors"
	"testing"
	"time"
)

func TestRenderMarkdown_GroupsConventionalCommitsAndFiltersNoise(t *testing.T) {
	releaseDate := time.Date(2026, time.May, 6, 22, 59, 0, 0, time.FixedZone("CST", 8*60*60))

	got, err := RenderMarkdown(Release{
		Tag:           "v1.3.0",
		PreviousTag:   "v1.2.3",
		RepositoryURL: "https://github.com/power721/115driver",
		Date:          releaseDate,
		Commits: []Commit{
			{Subject: "Merge branch 'feature/cli' into develop"},
			{Subject: "build: add project Makefile"},
			{Subject: "ci: add GoReleaser for multi-platform builds, fix --json=true early flag scan"},
			{Subject: "feat(cli): implement full CLI tool for 115 cloud storage"},
			{Subject: "feat(cli): add config and version commands"},
			{Subject: "fix(cli): pre-scan --json flag before Cobra parses args"},
			{Subject: "fix(cli): use plain text for QR login progress on stderr in --json mode"},
			{Subject: "refactor(cli): deduplicate code from simplify review"},
			{Subject: "docs: add CLI section to README"},
			{Subject: "chore: remove accidentally committed binary and QWEN.md"},
		},
	})
	if err != nil {
		t.Fatalf("RenderMarkdown returned error: %v", err)
	}

	want := `## [v1.3.0](https://github.com/power721/115driver/compare/v1.2.3...v1.3.0) (2026-05-06)

### 🚀 Features

* **cli:** implement full CLI tool for 115 cloud storage
* **cli:** add config and version commands

### 🐞 Bug Fixes

* **cli:** pre-scan --json flag before Cobra parses args
* **cli:** use plain text for QR login progress on stderr in --json mode

### 🔨 Refactoring

* **cli:** deduplicate code from simplify review
`
	if got != want {
		t.Fatalf("unexpected markdown:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestRenderMarkdown_WithoutPreviousTagOmitsCompareLink(t *testing.T) {
	releaseDate := time.Date(2025, time.December, 8, 18, 39, 0, 0, time.UTC)

	got, err := RenderMarkdown(Release{
		Tag:           "v1.2.0",
		RepositoryURL: "https://github.com/power721/115driver",
		Date:          releaseDate,
		Commits: []Commit{
			{Subject: "feat(mcp): add Model Context Protocol server implementation"},
		},
	})
	if err != nil {
		t.Fatalf("RenderMarkdown returned error: %v", err)
	}

	want := `## v1.2.0 (2025-12-08)

### 🚀 Features

* **mcp:** add Model Context Protocol server implementation
`
	if got != want {
		t.Fatalf("unexpected markdown:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestRenderMarkdown_IncludesRevertsSection(t *testing.T) {
	releaseDate := time.Date(2026, time.May, 6, 22, 59, 0, 0, time.UTC)

	got, err := RenderMarkdown(Release{
		Tag:           "v1.3.1",
		PreviousTag:   "v1.3.0",
		RepositoryURL: "https://github.com/power721/115driver",
		Date:          releaseDate,
		Commits: []Commit{
			{Subject: `Revert "feat(cli): add config and version commands"`},
		},
	})
	if err != nil {
		t.Fatalf("RenderMarkdown returned error: %v", err)
	}

	want := `## [v1.3.1](https://github.com/power721/115driver/compare/v1.3.0...v1.3.1) (2026-05-06)

### Reverts

* feat(cli): add config and version commands
`
	if got != want {
		t.Fatalf("unexpected markdown:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestLoadRelease_UsesPreviousTagDateAndCommitsFromGit(t *testing.T) {
	runner := fakeGitRunner{
		results: map[string]string{
			"describe\x00--tags\x00--abbrev=0\x00v1.3.0^": "v1.2.3\n",
			"log\x00-1\x00--format=%cI\x00v1.3.0":         "2026-05-06T22:59:00+08:00\n",
			"log\x00--format=%s\x00v1.2.3..v1.3.0":        "feat(cli): add config and version commands\nfix(cli): pre-scan --json flag before Cobra parses args\n",
		},
	}

	got, err := LoadRelease(runner, "https://github.com/power721/115driver", "v1.3.0")
	if err != nil {
		t.Fatalf("LoadRelease returned error: %v", err)
	}

	if got.Tag != "v1.3.0" {
		t.Fatalf("unexpected tag: %s", got.Tag)
	}
	if got.PreviousTag != "v1.2.3" {
		t.Fatalf("unexpected previous tag: %s", got.PreviousTag)
	}
	if got.Date.Format(time.RFC3339) != "2026-05-06T22:59:00+08:00" {
		t.Fatalf("unexpected date: %s", got.Date.Format(time.RFC3339))
	}
	if len(got.Commits) != 2 {
		t.Fatalf("unexpected commit count: %d", len(got.Commits))
	}
	if got.Commits[0].Subject != "feat(cli): add config and version commands" {
		t.Fatalf("unexpected first commit: %q", got.Commits[0].Subject)
	}
}

func TestLoadRelease_WithoutPreviousTagFallsBackToTaggedHistory(t *testing.T) {
	runner := fakeGitRunner{
		results: map[string]string{
			"log\x00-1\x00--format=%cI\x00v1.2.0": "2025-12-08T18:39:00Z\n",
			"log\x00--format=%s\x00v1.2.0":        "feat(mcp): add Model Context Protocol server implementation\n",
		},
		errors: map[string]error{
			"describe\x00--tags\x00--abbrev=0\x00v1.2.0^": errors.New("no previous tag"),
		},
	}

	got, err := LoadRelease(runner, "https://github.com/power721/115driver", "v1.2.0")
	if err != nil {
		t.Fatalf("LoadRelease returned error: %v", err)
	}

	if got.PreviousTag != "" {
		t.Fatalf("expected empty previous tag, got %q", got.PreviousTag)
	}
	if len(got.Commits) != 1 {
		t.Fatalf("unexpected commit count: %d", len(got.Commits))
	}
}

type fakeGitRunner struct {
	results map[string]string
	errors  map[string]error
}

func (f fakeGitRunner) Run(args ...string) (string, error) {
	key := joinArgs(args...)
	if err, ok := f.errors[key]; ok {
		return "", err
	}
	if result, ok := f.results[key]; ok {
		return result, nil
	}
	return "", errors.New("unexpected command: " + key)
}
