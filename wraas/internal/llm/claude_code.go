package llm

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// ClaudeCodeProvider implements Provider by shelling out to the claude CLI.
// This allows using a Claude Max Plan subscription instead of API credits.
type ClaudeCodeProvider struct{}

func (p *ClaudeCodeProvider) Complete(ctx context.Context, req Request) (Response, error) {
	// Build the prompt: system prompt + user message combined
	prompt := req.SystemPrompt + "\n\n---\n\nUser query:\n" + req.UserMessage

	args := []string{
		"--print",           // non-interactive, output only
		"--model", "sonnet", // use sonnet model
		prompt,
	}

	if req.MaxTokens > 0 {
		args = append([]string{"--max-turns", "1"}, args...)
	}

	cmd := exec.CommandContext(ctx, "claude", args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return Response{}, fmt.Errorf("claude CLI error: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return Response{}, fmt.Errorf("claude CLI not found or failed: %w", err)
	}

	return Response{Content: strings.TrimSpace(string(out))}, nil
}

// isClaudeCodeAvailable checks if the claude CLI is installed and accessible.
func isClaudeCodeAvailable() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}
