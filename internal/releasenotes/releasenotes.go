package releasenotes

import (
	"fmt"
	"strings"
	"time"
)

type Commit struct {
	Subject string
}

type Release struct {
	Tag           string
	PreviousTag   string
	RepositoryURL string
	Date          time.Time
	Commits       []Commit
}

type commitGroup struct {
	Title string
	Items []string
}

type GitRunner interface {
	Run(args ...string) (string, error)
}

var groupOrder = []struct {
	Type  string
	Title string
}{
	{Type: "feat", Title: "🚀 Features"},
	{Type: "fix", Title: "🐞 Bug Fixes"},
	{Type: "perf", Title: "🏎 Performance"},
	{Type: "refactor", Title: "🔨 Refactoring"},
}

func RenderMarkdown(release Release) (string, error) {
	if release.Tag == "" {
		return "", fmt.Errorf("tag is required")
	}
	if release.Date.IsZero() {
		return "", fmt.Errorf("date is required")
	}

	grouped := map[string][]string{}
	reverts := make([]string, 0)
	for _, commit := range release.Commits {
		if reverted, ok := parseRevert(commit.Subject); ok {
			reverts = append(reverts, reverted)
			continue
		}

		parsed, ok := parseCommit(commit.Subject)
		if !ok {
			continue
		}
		grouped[parsed.Type] = append(grouped[parsed.Type], parsed.Rendered())
	}

	var b strings.Builder
	if release.PreviousTag != "" {
		fmt.Fprintf(&b, "## [%s](%s/compare/%s...%s) (%s)\n", release.Tag, strings.TrimSuffix(release.RepositoryURL, "/"), release.PreviousTag, release.Tag, release.Date.Format("2006-01-02"))
	} else {
		fmt.Fprintf(&b, "## %s (%s)\n", release.Tag, release.Date.Format("2006-01-02"))
	}

	for _, group := range groupOrder {
		items := grouped[group.Type]
		if len(items) == 0 {
			continue
		}
		fmt.Fprintf(&b, "\n### %s\n\n", group.Title)
		for _, item := range items {
			fmt.Fprintf(&b, "* %s\n", item)
		}
	}

	if len(reverts) > 0 {
		fmt.Fprintf(&b, "\n### Reverts\n\n")
		for _, item := range reverts {
			fmt.Fprintf(&b, "* %s\n", item)
		}
	}

	return b.String(), nil
}

func LoadRelease(runner GitRunner, repositoryURL, tag string) (Release, error) {
	if tag == "" {
		return Release{}, fmt.Errorf("tag is required")
	}

	previousTag, err := trimOutput(runner.Run("describe", "--tags", "--abbrev=0", tag+"^"))
	if err != nil {
		previousTag = ""
	}

	dateText, err := trimOutput(runner.Run("log", "-1", "--format=%cI", tag))
	if err != nil {
		return Release{}, fmt.Errorf("read tag date: %w", err)
	}

	releaseDate, err := time.Parse(time.RFC3339, dateText)
	if err != nil {
		return Release{}, fmt.Errorf("parse tag date: %w", err)
	}

	rangeSpec := tag
	if previousTag != "" {
		rangeSpec = previousTag + ".." + tag
	}

	commitText, err := trimOutput(runner.Run("log", "--format=%s", rangeSpec))
	if err != nil {
		return Release{}, fmt.Errorf("read commits: %w", err)
	}

	commits := make([]Commit, 0)
	if commitText != "" {
		for _, line := range strings.Split(commitText, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			commits = append(commits, Commit{Subject: line})
		}
	}

	return Release{
		Tag:           tag,
		PreviousTag:   previousTag,
		RepositoryURL: repositoryURL,
		Date:          releaseDate,
		Commits:       commits,
	}, nil
}

type parsedCommit struct {
	Type    string
	Scope   string
	Subject string
}

func (c parsedCommit) Rendered() string {
	if c.Scope == "" {
		return c.Subject
	}
	return fmt.Sprintf("**%s:** %s", c.Scope, c.Subject)
}

func parseCommit(subject string) (parsedCommit, bool) {
	if strings.HasPrefix(subject, "Merge ") {
		return parsedCommit{}, false
	}

	header, body, found := strings.Cut(subject, ":")
	if !found {
		return parsedCommit{}, false
	}

	body = strings.TrimSpace(body)
	if body == "" {
		return parsedCommit{}, false
	}

	commitType := strings.TrimSpace(header)
	scope := ""
	if open := strings.IndexByte(header, '('); open >= 0 && strings.HasSuffix(header, ")") {
		commitType = strings.TrimSpace(header[:open])
		scope = strings.TrimSpace(header[open+1 : len(header)-1])
	}

	for _, group := range groupOrder {
		if commitType == group.Type {
			return parsedCommit{
				Type:    commitType,
				Scope:   scope,
				Subject: body,
			}, true
		}
	}

	return parsedCommit{}, false
}

func parseRevert(subject string) (string, bool) {
	if !strings.HasPrefix(subject, `Revert "`) || !strings.HasSuffix(subject, `"`) {
		return "", false
	}
	return strings.TrimSuffix(strings.TrimPrefix(subject, `Revert "`), `"`), true
}

func trimOutput(output string, err error) (string, error) {
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func joinArgs(args ...string) string {
	return strings.Join(args, "\x00")
}
