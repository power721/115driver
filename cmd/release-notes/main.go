package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/power721/115driver/internal/releasenotes"
)

const defaultRepositoryURL = "https://github.com/power721/115driver"

func main() {
	var (
		tag  = flag.String("tag", envOrDefault("GITHUB_REF_NAME", ""), "git tag to render release notes for")
		repo = flag.String("repo-url", detectRepositoryURL(), "repository URL used in compare links")
	)
	flag.Parse()

	if *tag == "" {
		log.Fatal("missing -tag (or GITHUB_REF_NAME)")
	}

	release, err := releasenotes.LoadRelease(gitRunner{}, *repo, *tag)
	if err != nil {
		log.Fatal(err)
	}

	markdown, err := releasenotes.RenderMarkdown(release)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(markdown)
}

type gitRunner struct{}

func (gitRunner) Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func detectRepositoryURL() string {
	serverURL := strings.TrimSuffix(os.Getenv("GITHUB_SERVER_URL"), "/")
	repository := strings.TrimPrefix(os.Getenv("GITHUB_REPOSITORY"), "/")
	if serverURL != "" && repository != "" {
		return serverURL + "/" + repository
	}
	return defaultRepositoryURL
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
