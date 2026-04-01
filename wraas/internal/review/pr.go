package review

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// PRData holds the relevant data from a pull request.
type PRData struct {
	Number      int
	Title       string
	Body        string
	Commits     []PRCommit
	Files       []PRFile
	BaseBranch  string
	HeadBranch  string
}

// PRCommit represents a commit in a PR.
type PRCommit struct {
	SHA     string
	Message string
}

// PRFile represents a changed file in a PR.
type PRFile struct {
	Filename string
	Status   string // "added", "modified", "removed"
}

// FetchPR fetches PR data using the gh CLI.
func FetchPR(repo string, prNumber int) (PRData, error) {
	// Fetch PR metadata
	prJSON, err := ghAPI(fmt.Sprintf("repos/%s/pulls/%d", repo, prNumber))
	if err != nil {
		return PRData{}, fmt.Errorf("fetching PR: %w", err)
	}

	var pr struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		Body   string `json:"body"`
		Base   struct {
			Ref string `json:"ref"`
		} `json:"base"`
		Head struct {
			Ref string `json:"ref"`
		} `json:"head"`
	}
	if err := json.Unmarshal([]byte(prJSON), &pr); err != nil {
		return PRData{}, fmt.Errorf("parsing PR: %w", err)
	}

	data := PRData{
		Number:     pr.Number,
		Title:      pr.Title,
		Body:       pr.Body,
		BaseBranch: pr.Base.Ref,
		HeadBranch: pr.Head.Ref,
	}

	// Fetch commits
	commitsJSON, err := ghAPI(fmt.Sprintf("repos/%s/pulls/%d/commits", repo, prNumber))
	if err != nil {
		return data, fmt.Errorf("fetching commits: %w", err)
	}

	var commits []struct {
		SHA    string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
		} `json:"commit"`
	}
	if err := json.Unmarshal([]byte(commitsJSON), &commits); err != nil {
		return data, fmt.Errorf("parsing commits: %w", err)
	}

	for _, c := range commits {
		data.Commits = append(data.Commits, PRCommit{
			SHA:     c.SHA[:7],
			Message: c.Commit.Message,
		})
	}

	// Fetch files
	filesJSON, err := ghAPI(fmt.Sprintf("repos/%s/pulls/%d/files", repo, prNumber))
	if err != nil {
		return data, fmt.Errorf("fetching files: %w", err)
	}

	var files []struct {
		Filename string `json:"filename"`
		Status   string `json:"status"`
	}
	if err := json.Unmarshal([]byte(filesJSON), &files); err != nil {
		return data, fmt.Errorf("parsing files: %w", err)
	}

	for _, f := range files {
		data.Files = append(data.Files, PRFile{
			Filename: f.Filename,
			Status:   f.Status,
		})
	}

	return data, nil
}

func ghAPI(endpoint string) (string, error) {
	cmd := exec.Command("gh", "api", endpoint)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("gh api %s: %s", endpoint, strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", err
	}
	return string(out), nil
}
