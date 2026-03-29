package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wraas/digital-twin-product/cli/internal/config"
	"github.com/wraas/digital-twin-product/cli/internal/engine"
	"github.com/wraas/digital-twin-product/cli/internal/output"
	"github.com/wraas/digital-twin-product/cli/internal/review"
	"github.com/wraas/digital-twin-product/cli/internal/tui"
)

var (
	prNumber int
	prRepo   string
	strict   bool
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review a pull request",
	Long: `Submit a pull request for WRAAS review. The review covers conventional commit
compliance, scope consistency, breaking change documentation, and whether the
PR description answers the question "why" with sufficient conviction.`,
	RunE: runReview,
}

func init() {
	reviewCmd.Flags().IntVar(&prNumber, "pr", 0, "Pull request number (required)")
	reviewCmd.Flags().StringVar(&prRepo, "repo", "", "Repository in owner/repo format (required)")
	reviewCmd.Flags().BoolVar(&strict, "strict", true, "Enforce all conventional commit rules. Default: true. There is no false.")
	reviewCmd.MarkFlagRequired("pr")
	reviewCmd.MarkFlagRequired("repo")
	rootCmd.AddCommand(reviewCmd)
}

// ReviewOutput is the structured output for json/yaml.
type ReviewOutput struct {
	PR         int                    `json:"pr" yaml:"pr"`
	Repo       string                 `json:"repo" yaml:"repo"`
	Commits    int                    `json:"commits" yaml:"commits"`
	Files      int                    `json:"files" yaml:"files"`
	Violations []ReviewViolationOut   `json:"violations" yaml:"violations"`
	Compliance []ReviewComplianceOut  `json:"compliance" yaml:"compliance"`
	Actions    []ReviewActionOut      `json:"actions,omitempty" yaml:"actions,omitempty"`
	SighLevel  string                 `json:"sigh_level" yaml:"sigh_level"`
	LatencyMs  int                    `json:"latency_ms" yaml:"latency_ms"`
}

type ReviewViolationOut struct {
	Commit   string `json:"commit" yaml:"commit"`
	Severity string `json:"severity" yaml:"severity"`
	Message  string `json:"message" yaml:"message"`
	Fix      string `json:"fix,omitempty" yaml:"fix,omitempty"`
}

type ReviewComplianceOut struct {
	Module   string `json:"module" yaml:"module"`
	Severity string `json:"severity" yaml:"severity"`
	Message  string `json:"message" yaml:"message"`
	DME      string `json:"dme,omitempty" yaml:"dme,omitempty"`
}

type ReviewActionOut struct {
	Action string `json:"action" yaml:"action"`
	Level  string `json:"level" yaml:"level"`
	Reason string `json:"reason" yaml:"reason"`
}

func runReview(cmd *cobra.Command, args []string) error {
	// --strict has no false. If someone managed to set it to false, note it.
	if cmd.Flags().Changed("strict") && !strict {
		if !quiet {
			output.Prompt(os.Stdout, output.WarnStyle.Render("--strict=false is not available. There is no false. Proceeding with strict enforcement."))
		}
		strict = true
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(engine.ExitConfigError)
	}

	// Fetch PR with spinner
	if !quiet {
		fmt.Printf("[Code Review Simulation] Fetching PR #%d from %s...\n", prNumber, prRepo)
	}

	var pr review.PRData
	_, fetchErr := tui.RunWithSpinner(
		fmt.Sprintf("Fetching diff from %s#%d...", prRepo, prNumber),
		func() (string, error) {
			var err error
			pr, err = review.FetchPR(prRepo, prNumber)
			return "", err
		},
	)
	if fetchErr != nil {
		fmt.Fprintf(os.Stderr, "Error fetching PR: %v\n", fetchErr)
		os.Exit(engine.ExitBlockingViolation)
	}

	if !quiet {
		output.Prompt(os.Stdout, fmt.Sprintf("%d commits, %d files changed", len(pr.Commits), len(pr.Files)))
	}

	// Run all checks
	var allViolations []ReviewViolationOut
	var allCompliance []ReviewComplianceOut
	var allActions []ReviewActionOut
	hasBlocking := false
	worstSigh := engine.SighSilent

	// 1. Conventional commit validation
	var parsedCommits []review.CommitMessage
	for _, c := range pr.Commits {
		cm := review.ParseCommit(c.Message)
		parsedCommits = append(parsedCommits, cm)

		violations := review.ValidateCommit(cm, cfg.Commit.ScopeRequired, cfg.Commit.ImperativeMood, cfg.Commit.Enforcement)
		for _, v := range violations {
			allViolations = append(allViolations, ReviewViolationOut{
				Commit:   c.SHA,
				Severity: v.Severity,
				Message:  v.Message,
				Fix:      v.Fix,
			})
			if v.Severity == "blocking" {
				hasBlocking = true
			}
			// Calibrate sigh based on violation type
			if strings.Contains(v.Message, "wip") {
				worstSigh = engine.SighExistential
			} else if strings.Contains(v.Message, "Missing scope") {
				if engine.SighIntensity(worstSigh) < engine.SighIntensity(engine.SighModerate) {
					worstSigh = engine.SighModerate
				}
			}
		}
	}

	// 2. Scope consistency
	scopeViolations := review.ScopeConsistency(parsedCommits)
	for _, v := range scopeViolations {
		allViolations = append(allViolations, ReviewViolationOut{
			Severity: v.Severity,
			Message:  v.Message,
		})
		if engine.SighIntensity(worstSigh) < engine.SighIntensity(engine.SighDeep) {
			worstSigh = engine.SighDeep
		}
	}

	// 3. Compliance checks
	complianceChecks := review.CheckCompliance(pr.Files, cfg.Compliance.Asciidoc, cfg.Compliance.AntoraStructure)
	for _, c := range complianceChecks {
		allCompliance = append(allCompliance, ReviewComplianceOut{
			Module:   c.Module,
			Severity: c.Severity,
			Message:  c.Message,
			DME:      c.DME,
		})
	}

	// 4. GitHub Actions suggestions
	actionSuggestions := review.SuggestActions(pr.Files, cfg.GithubActions.SlugAction, cfg.GithubActions.DrawioExport)
	for _, a := range actionSuggestions {
		allActions = append(allActions, ReviewActionOut{
			Action: a.Action,
			Level:  a.Level,
			Reason: a.Reason,
		})
	}

	// 5. Check PR description quality
	if strings.TrimSpace(pr.Body) == "" {
		allViolations = append(allViolations, ReviewViolationOut{
			Severity: "non-blocking",
			Message:  "PR description is empty. Does not answer \"why\".",
		})
	}

	// Default sigh if nothing triggered
	if worstSigh == engine.SighSilent && len(allViolations) == 0 {
		worstSigh = engine.SighSilent
	} else if worstSigh == engine.SighSilent {
		worstSigh = engine.SighMild
	}

	// Log sigh
	engine.LogSigh(worstSigh, fmt.Sprintf("review %s#%d", prRepo, prNumber))

	// Output
	format := output.ParseFormat(outputFormat)

	result := ReviewOutput{
		PR:         prNumber,
		Repo:       prRepo,
		Commits:    len(pr.Commits),
		Files:      len(pr.Files),
		Violations: allViolations,
		Compliance: allCompliance,
		Actions:    allActions,
		SighLevel:  string(worstSigh),
		LatencyMs:  engine.LatencyMs(),
	}

	output.Write(os.Stdout, format, result, func(w io.Writer) {
		renderReviewText(w, result, hasBlocking)
	})

	if hasBlocking {
		os.Exit(engine.ExitBlockingViolation)
	}
	return nil
}

func renderReviewText(w io.Writer, result ReviewOutput, hasBlocking bool) {
	// Violations
	if len(result.Violations) > 0 {
		output.Prompt(w, output.HeaderStyle.Render("Conventional Commit Validation"))
		for _, v := range result.Violations {
			if v.Severity == "blocking" {
				output.Violation(w, v.Message)
			} else {
				output.Suggestion(w, v.Message)
			}
			if v.Fix != "" {
				output.Prompt(w, "  "+output.DimStyle.Render(v.Fix))
			}
		}
	} else {
		output.Prompt(w, output.OkStyle.Render("Conventional Commit Validation: all commits pass"))
	}

	// Compliance
	if len(result.Compliance) > 0 {
		output.Prompt(w, output.HeaderStyle.Render("Compliance Checks"))
		for _, c := range result.Compliance {
			prefix := fmt.Sprintf("[%s]", c.Module)
			msg := c.Message
			if c.DME != "" {
				msg += fmt.Sprintf(" (see %s)", c.DME)
			}
			output.Suggestion(w, prefix+" "+msg)
		}
	}

	// Actions
	if len(result.Actions) > 0 {
		output.Prompt(w, output.HeaderStyle.Render("GitHub Actions"))
		for _, a := range result.Actions {
			output.Suggestion(w, fmt.Sprintf("have you considered %s? %s", a.Action, a.Reason))
		}
	}

	// Sigh
	if result.SighLevel != string(engine.SighSilent) {
		output.Sigh(w, result.SighLevel)
	}

	// Summary
	fmt.Fprintln(w)
	if hasBlocking {
		output.Prompt(w, output.ErrorStyle.Render(fmt.Sprintf("Review complete | %d violation(s) | Latency: %dms", len(result.Violations), result.LatencyMs)))
	} else {
		output.Prompt(w, output.OkStyle.Render(fmt.Sprintf("Review complete | %d suggestion(s) | Latency: %dms", len(result.Violations)+len(result.Compliance)+len(result.Actions), result.LatencyMs)))
	}
}
