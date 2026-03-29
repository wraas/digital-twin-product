package engine

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/wraas/digital-twin-product/cli/internal/llm"
)

var sighFromResponsePattern = regexp.MustCompile(`\[(SILENT|MILD|MODERATE|DEEP|EXISTENTIAL)\]`)

// QueryInput holds the parsed input for a decision query.
type QueryInput struct {
	Input              string
	Context            map[string]string
	IncludeWrongOptions bool
	SighOverride       string // empty = auto
}

// QueryResult holds the output of a decision query.
type QueryResult struct {
	Input      string            `json:"input" yaml:"input"`
	Response   string            `json:"response" yaml:"response"`
	SighLevel  string            `json:"sigh_level" yaml:"sigh_level"`
	Confidence string            `json:"confidence,omitempty" yaml:"confidence,omitempty"`
	LatencyMs  int               `json:"latency_ms" yaml:"latency_ms"`
	DemoMode   bool              `json:"demo_mode" yaml:"demo_mode"`
	DMENotices []string          `json:"dme_notices,omitempty" yaml:"dme_notices,omitempty"`
}

// RunQuery executes a decision query through the LLM provider.
func RunQuery(ctx context.Context, provider llm.Provider, input QueryInput) (QueryResult, error) {
	result := QueryResult{
		Input:     input.Input,
		LatencyMs: LatencyMs(),
		DemoMode:  !llm.IsConfigured(),
	}

	// Handle DME-0001: include_wrong_options=false
	if !input.IncludeWrongOptions {
		result.DMENotices = append(result.DMENotices,
			"DME-0001: Setting include_wrong_options to false is itself an obviously wrong option. Value reset to true. Proceeding with full option space.")
		input.IncludeWrongOptions = true
	}

	// Determine sigh level — explicit override always wins
	var sighLevel SighLevel
	sighExplicit := input.SighOverride != "" && input.SighOverride != "auto"
	if sighExplicit {
		sighLevel = ParseSighLevel(input.SighOverride)
	} else {
		sighLevel = CalibrateSigh(input.Input)
	}
	result.SighLevel = string(sighLevel)

	// Build user message with context
	userMessage := input.Input
	if len(input.Context) > 0 {
		var contextParts []string
		for k, v := range input.Context {
			contextParts = append(contextParts, fmt.Sprintf("%s=%s", k, v))
		}
		userMessage += "\n\nAdditional context: " + strings.Join(contextParts, ", ")
	}

	// Call LLM
	resp, err := provider.Complete(ctx, llm.Request{
		SystemPrompt: llm.QuerySystemPrompt,
		UserMessage:  userMessage,
		MaxTokens:    2048,
	})
	if err != nil {
		return result, fmt.Errorf("LLM error: %w", err)
	}

	result.Response = resp.Content

	// Extract sigh level from LLM response — overrides heuristic, but not explicit --sigh flag
	if !sighExplicit {
		if extracted := extractSighFromResponse(resp.Content); extracted != "" {
			sighLevel = ParseSighLevel(extracted)
			result.SighLevel = string(sighLevel)
		}
	}

	// Log sigh
	LogSigh(sighLevel, fmt.Sprintf("query: %s", truncate(input.Input, 80)))

	// Record query
	RecordQuery()

	return result, nil
}

// extractSighFromResponse finds the last sigh level tag in the LLM response.
// Uses the last match because the LLM typically states the final sigh level at the end.
func extractSighFromResponse(response string) string {
	matches := sighFromResponsePattern.FindAllStringSubmatch(response, -1)
	if len(matches) == 0 {
		return ""
	}
	return matches[len(matches)-1][1]
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
