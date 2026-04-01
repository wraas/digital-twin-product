package review

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Issue represents a GitHub issue.
type Issue struct {
	Number int
	Title  string
	Body   string
	Author string
}

// OSSResult holds the full triage output.
type OSSResult struct {
	Repo      string          `json:"repo" yaml:"repo"`
	Issues    []IssueTriage   `json:"issues" yaml:"issues"`
	PRReviews []OSSPRReview   `json:"pr_reviews" yaml:"pr_reviews"`
	LatencyMs int             `json:"latency_ms" yaml:"latency_ms"`
	SighLevel string          `json:"sigh_level" yaml:"sigh_level"`
}

// IssueTriage is the triage result for a single issue.
type IssueTriage struct {
	Number   int    `json:"number" yaml:"number"`
	Title    string `json:"title" yaml:"title"`
	HasRepro bool   `json:"has_repro" yaml:"has_repro"`
	Severity string `json:"severity" yaml:"severity"`
	Message  string `json:"message" yaml:"message"`
}

// OSSPRReview is the review result for a single PR.
type OSSPRReview struct {
	Number     int               `json:"number" yaml:"number"`
	Title      string            `json:"title" yaml:"title"`
	Author     string            `json:"author" yaml:"author"`
	FirstTime  bool              `json:"first_time" yaml:"first_time"`
	Violations []CommitViolation `json:"violations,omitempty" yaml:"violations,omitempty"`
}

// FetchIssues fetches open issues for a repository via gh CLI.
func FetchIssues(repo string) ([]Issue, error) {
	data, err := ghAPI(fmt.Sprintf("repos/%s/issues?state=open&per_page=30", repo))
	if err != nil {
		return nil, fmt.Errorf("fetching issues: %w", err)
	}

	var raw []struct {
		Number      int    `json:"number"`
		Title       string `json:"title"`
		Body        string `json:"body"`
		PullRequest *struct{} `json:"pull_request"`
		User        struct {
			Login string `json:"login"`
		} `json:"user"`
	}
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return nil, fmt.Errorf("parsing issues: %w", err)
	}

	var issues []Issue
	for _, r := range raw {
		// GitHub issues API includes PRs — skip them
		if r.PullRequest != nil {
			continue
		}
		issues = append(issues, Issue{
			Number: r.Number,
			Title:  r.Title,
			Body:   r.Body,
			Author: r.User.Login,
		})
	}
	return issues, nil
}

// FetchOpenPRs fetches open pull requests for a repository.
func FetchOpenPRs(repo string) ([]PRData, error) {
	data, err := ghAPI(fmt.Sprintf("repos/%s/pulls?state=open&per_page=30", repo))
	if err != nil {
		return nil, fmt.Errorf("fetching PRs: %w", err)
	}

	var raw []struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		Body   string `json:"body"`
		User   struct {
			Login string `json:"login"`
		} `json:"user"`
		AuthorAssociation string `json:"author_association"`
	}
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return nil, fmt.Errorf("parsing PRs: %w", err)
	}

	var prs []PRData
	for _, r := range raw {
		pr := PRData{
			Number: r.Number,
			Title:  r.Title,
			Body:   r.Body,
		}

		// Fetch commits for each PR
		commitsJSON, err := ghAPI(fmt.Sprintf("repos/%s/pulls/%d/commits", repo, r.Number))
		if err == nil {
			var commits []struct {
				SHA    string `json:"sha"`
				Commit struct {
					Message string `json:"message"`
				} `json:"commit"`
			}
			if json.Unmarshal([]byte(commitsJSON), &commits) == nil {
				for _, c := range commits {
					sha := c.SHA
					if len(sha) > 7 {
						sha = sha[:7]
					}
					pr.Commits = append(pr.Commits, PRCommit{
						SHA:     sha,
						Message: c.Commit.Message,
					})
				}
			}
		}

		prs = append(prs, pr)
	}
	return prs, nil
}

// IsFirstTimeContributor checks if the author has no previous merged PRs.
func IsFirstTimeContributor(repo string, author string) bool {
	data, err := ghAPI(fmt.Sprintf("repos/%s/pulls?state=closed&per_page=1&creator=%s", repo, author))
	if err != nil {
		return true // assume first-time if we can't check
	}
	return data == "[]" || strings.TrimSpace(data) == "[\n]" || strings.TrimSpace(data) == "[]"
}

// TriageIssue evaluates an issue for reproduction steps and clarity.
func TriageIssue(issue Issue) IssueTriage {
	body := strings.ToLower(issue.Body)

	hasRepro := hasReproductionSteps(body)
	hasVersion := hasVersionInfo(body)
	hasCodeBlock := strings.Contains(issue.Body, "```")

	triage := IssueTriage{
		Number: issue.Number,
		Title:  issue.Title,
	}

	switch {
	case len(strings.TrimSpace(issue.Body)) < 20:
		// Nearly empty body
		triage.HasRepro = false
		triage.Severity = "blocking"
		triage.Message = "what doesn't work? what version? what input? what expected output?"

	case !hasRepro && !hasCodeBlock:
		// No reproduction steps and no code examples
		triage.HasRepro = false
		triage.Severity = "blocking"
		triage.Message = "can you provide reproduction steps and the exact error output?"

	case hasRepro || (hasCodeBlock && hasVersion):
		// Good issue
		triage.HasRepro = true
		triage.Severity = "non-blocking"
		triage.Message = "reproduction steps present, scope clear. evaluated. added to decision matrix. response queued."

	default:
		// Partial info
		triage.HasRepro = false
		triage.Severity = "blocking"
		triage.Message = "partial information provided. can you include version, expected behavior, and actual behavior?"
	}

	return triage
}

func hasReproductionSteps(body string) bool {
	markers := []string{
		"steps to reproduce",
		"reproduction",
		"how to reproduce",
		"to reproduce",
		"1.",
		"step 1",
		"expected behavior",
		"expected result",
		"actual behavior",
		"actual result",
	}
	for _, m := range markers {
		if strings.Contains(body, m) {
			return true
		}
	}
	return false
}

func hasVersionInfo(body string) bool {
	markers := []string{
		"version",
		"v1.", "v2.", "v3.", "v4.", "v5.",
		"@v",
		"node ",
		"go ",
		"python ",
	}
	for _, m := range markers {
		if strings.Contains(body, m) {
			return true
		}
	}
	return false
}
