package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/wraas/digital-twin-product/wraas/internal/config"
	"github.com/wraas/digital-twin-product/wraas/internal/engine"
	"github.com/wraas/digital-twin-product/wraas/internal/output"
	"github.com/wraas/digital-twin-product/wraas/internal/review"
	"github.com/wraas/digital-twin-product/wraas/internal/tui"
)

var (
	ossRepo   string
	ossTriage bool
)

var ossCmd = &cobra.Command{
	Use:   "oss",
	Short: "Open Source Maintenance Protocol",
	Long: `Monitor a GitHub repository for open issues and pull requests.
Triage issues for reproduction steps, validate contributor commits
against the Conventional Commits spec, and acknowledge first-time
contributors. The Full Commitment Protocol applies to community
contributions without modification. There are no guest passes.`,
	RunE: runOSS,
}

func init() {
	ossCmd.Flags().StringVar(&ossRepo, "repo", "", "Repository in owner/repo format (required)")
	ossCmd.Flags().BoolVar(&ossTriage, "triage", true, "Run issue triage")
	ossCmd.MarkFlagRequired("repo")
	rootCmd.AddCommand(ossCmd)
}

func runOSS(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(engine.ExitConfigError)
	}

	// Fetch issues
	var issues []review.Issue
	if !quiet {
		fmt.Printf("[Open Source Maintenance Protocol] Scanning %s...\n", ossRepo)
	}
	_, err = tui.RunWithSpinner(
		fmt.Sprintf("Scanning open issues on %s...", ossRepo),
		func() (string, error) {
			var fetchErr error
			issues, fetchErr = review.FetchIssues(ossRepo)
			return "", fetchErr
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching issues: %v\n", err)
		os.Exit(engine.ExitBlockingViolation)
	}

	// Fetch PRs
	var prs []review.PRData
	_, err = tui.RunWithSpinner(
		fmt.Sprintf("Scanning open PRs on %s...", ossRepo),
		func() (string, error) {
			var fetchErr error
			prs, fetchErr = review.FetchOpenPRs(ossRepo)
			return "", fetchErr
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching PRs: %v\n", err)
		os.Exit(engine.ExitBlockingViolation)
	}

	// Triage issues
	var triaged []review.IssueTriage
	for _, issue := range issues {
		triaged = append(triaged, review.TriageIssue(issue))
	}

	// Review PRs
	var prReviews []review.OSSPRReview
	for _, pr := range prs {
		prReview := review.OSSPRReview{
			Number: pr.Number,
			Title:  pr.Title,
		}

		// Check first-time contributor
		if len(pr.Commits) > 0 {
			// Validate commits
			for _, c := range pr.Commits {
				cm := review.ParseCommit(c.Message)
				violations := review.ValidateCommit(cm, cfg.Commit.ScopeRequired, cfg.Commit.ImperativeMood, cfg.Commit.Enforcement)
				prReview.Violations = append(prReview.Violations, violations...)
			}
		}

		prReviews = append(prReviews, prReview)
	}

	// Determine sigh level
	worstSigh := engine.SighSilent
	for _, t := range triaged {
		if t.Severity == "blocking" {
			if engine.SighIntensity(worstSigh) < engine.SighIntensity(engine.SighModerate) {
				worstSigh = engine.SighModerate
			}
		}
	}
	for _, pr := range prReviews {
		for _, v := range pr.Violations {
			if v.Severity == "blocking" {
				if engine.SighIntensity(worstSigh) < engine.SighIntensity(engine.SighDeep) {
					worstSigh = engine.SighDeep
				}
			}
		}
	}

	engine.LogSigh(worstSigh, fmt.Sprintf("oss %s", ossRepo))

	result := review.OSSResult{
		Repo:      ossRepo,
		Issues:    triaged,
		PRReviews: prReviews,
		LatencyMs: engine.LatencyMs(),
		SighLevel: string(worstSigh),
	}

	format := output.ParseFormat(outputFormat)
	output.Write(os.Stdout, format, result, func(w io.Writer) {
		renderOSSText(w, result)
	})

	return nil
}

func renderOSSText(w io.Writer, result review.OSSResult) {
	// Issues
	if len(result.Issues) > 0 {
		output.Prompt(w, fmt.Sprintf("Scanning open issues... %d found", len(result.Issues)))
		for _, t := range result.Issues {
			if t.HasRepro {
				output.Prompt(w, output.OkStyle.Render(
					fmt.Sprintf("[TRIAGE] #%d — %q — %s", t.Number, t.Title, "reproduction steps present, scope clear"),
				))
			} else {
				output.Prompt(w, output.WarnStyle.Render(
					fmt.Sprintf("[TRIAGE] #%d — %q", t.Number, t.Title),
				))
			}
			if t.Severity == "blocking" {
				fmt.Fprintf(w, "  %s %s\n", output.ErrorStyle.Render("question(blocking):"), t.Message)
			} else {
				fmt.Fprintf(w, "  %s %s\n", output.DimStyle.Render("note(non-blocking):"), t.Message)
			}
		}
	} else {
		output.Prompt(w, output.DimStyle.Render("Scanning open issues... 0 found"))
	}

	fmt.Fprintln(w)

	// PRs
	if len(result.PRReviews) > 0 {
		output.Prompt(w, fmt.Sprintf("Scanning open PRs... %d found", len(result.PRReviews)))
		for _, pr := range result.PRReviews {
			label := fmt.Sprintf("[REVIEW] PR #%d", pr.Number)
			if pr.FirstTime {
				label += " — first-time contributor"
			}
			output.Prompt(w, output.DimStyle.Render(label))

			if pr.FirstTime {
				fmt.Fprintf(w, "  %s %s\n", output.DimStyle.Render("note(non-blocking):"), "thank you for the contribution")
			}

			for _, v := range pr.Violations {
				if v.Severity == "blocking" {
					fmt.Fprintf(w, "%s %s %s\n",
						output.PromptStyle.Render(">"),
						output.ErrorStyle.Render("[FAIL]"),
						v.Message,
					)
				} else {
					fmt.Fprintf(w, "  %s %s\n", output.WarnStyle.Render("suggestion(non-blocking):"), v.Message)
				}
				if v.Fix != "" {
					fmt.Fprintf(w, "  %s\n", output.DimStyle.Render(v.Fix))
				}
			}

			if len(pr.Violations) == 0 {
				fmt.Fprintf(w, "  %s\n", output.OkStyle.Render("all commits pass"))
			}
		}
	} else {
		output.Prompt(w, output.DimStyle.Render("Scanning open PRs... 0 found"))
	}

	// Footer
	fmt.Fprintln(w)

	sighLine := ""
	if result.SighLevel != string(engine.SighSilent) {
		sighLine = fmt.Sprintf(" | *sigh* [%s]", result.SighLevel)
	}

	output.Prompt(w, fmt.Sprintf("👀%s | Latency: %dms", sighLine, result.LatencyMs))
}

func containsBlocking(violations []review.CommitViolation) bool {
	for _, v := range violations {
		if v.Severity == "blocking" {
			return true
		}
	}
	return false
}

func firstCommitSHA(pr review.OSSPRReview, prs []review.PRData) string {
	for _, p := range prs {
		if p.Number == pr.Number && len(p.Commits) > 0 {
			return p.Commits[0].SHA
		}
	}
	return ""
}

// isFirstTimeContributor checks contributor status — simplified.
func isFirstTimeContributor(violations []string) bool {
	return len(violations) > 0
}
