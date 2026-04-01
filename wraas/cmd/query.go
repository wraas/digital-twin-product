package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wraas/digital-twin-product/wraas/internal/config"
	"github.com/wraas/digital-twin-product/wraas/internal/engine"
	"github.com/wraas/digital-twin-product/wraas/internal/llm"
	"github.com/wraas/digital-twin-product/wraas/internal/output"
	"github.com/wraas/digital-twin-product/wraas/internal/tui"
)

var (
	queryInput        string
	queryContext      string
	includeWrongOpts  bool
	sighOverride      string
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Submit a decision query for evaluation",
	Long: `Submit a decision query for evaluation. WRAAS will generate the full option
space, evaluate each option, document all rejections, and return a
recommendation with rationale. The rationale is always included.
--verbose makes it longer.`,
	RunE: runQuery,
}

func init() {
	queryCmd.Flags().StringVar(&queryInput, "input", "", "The query text (required)")
	queryCmd.Flags().StringVar(&queryContext, "context", "", "Additional context as comma-separated key=value pairs")
	queryCmd.Flags().BoolVar(&includeWrongOpts, "include-wrong-options", true, "Include obviously wrong options in evaluation. Setting to false is itself an obviously wrong option.")
	queryCmd.Flags().StringVar(&sighOverride, "sigh", "auto", "Sigh calibration: auto, none, mild, moderate, deep, existential")
	queryCmd.MarkFlagRequired("input")
	rootCmd.AddCommand(queryCmd)
}

func runQuery(cmd *cobra.Command, args []string) error {
	// Parse context key=value pairs
	contextMap := make(map[string]string)
	if queryContext != "" {
		pairs := strings.Split(queryContext, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
			if len(kv) == 2 {
				contextMap[kv[0]] = kv[1]
			}
		}
	}

	input := engine.QueryInput{
		Input:               queryInput,
		Context:             contextMap,
		IncludeWrongOptions: includeWrongOpts,
		SighOverride:        sighOverride,
	}

	// Show DME-0001 notice before processing
	if !includeWrongOpts && !quiet {
		output.Prompt(os.Stdout, output.WarnStyle.Render("DME-0001: Setting include_wrong_options to false is itself an obviously wrong option. Value reset to true."))
		output.Prompt(os.Stdout, output.DimStyle.Render("Proceeding with full option space."))
		fmt.Println()
	}

	// Run query with spinner
	provider := llm.NewProvider()

	if !quiet {
		providerName := llm.ProviderName()
		if !llm.IsConfigured() {
			output.Prompt(os.Stdout, output.DimStyle.Render("LLM provider not configured. Running in demonstration mode."))
			output.Prompt(os.Stdout, output.DimStyle.Render("Set WRAAS_PROVIDER=claude-code to use your Max Plan, or set ANTHROPIC_API_KEY for API access."))
			fmt.Println()
		} else {
			output.Prompt(os.Stdout, output.DimStyle.Render("Provider: "+providerName))
		}
		_ = providerName
	}

	var result engine.QueryResult
	_, err := tui.RunWithSpinner(
		"Generating full option space (including obviously wrong ones)...",
		func() (string, error) {
			var err error
			result, err = engine.RunQuery(context.Background(), provider, input)
			return "", err
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(engine.ExitBlockingViolation)
	}

	// Output
	format := output.ParseFormat(outputFormat)

	output.Write(os.Stdout, format, result, func(w io.Writer) {
		renderQueryText(w, result)
	})

	return nil
}

func renderQueryText(w io.Writer, result engine.QueryResult) {
	// DME notices
	for _, notice := range result.DMENotices {
		output.Prompt(w, output.WarnStyle.Render(notice))
	}

	// Load config for max_width
	cfg, _ := config.Load(cfgFile)
	maxWidth := cfg.Output.MaxWidth

	// Render the LLM response with colors
	fmt.Fprintln(w)
	output.RenderResponse(w, result.Response, maxWidth)

	// Sigh
	if result.SighLevel != string(engine.SighSilent) {
		fmt.Fprintln(w)
		output.Sigh(w, result.SighLevel)
	}

	// Latency footer
	output.Latency(w)
}
