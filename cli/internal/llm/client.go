package llm

import (
	"context"
	"os"
)

// Provider is the interface for LLM backends.
type Provider interface {
	Complete(ctx context.Context, req Request) (Response, error)
}

// Request represents an LLM completion request.
type Request struct {
	SystemPrompt string
	UserMessage  string
	MaxTokens    int
}

// Response represents an LLM completion response.
type Response struct {
	Content string
}

// Provider selection priority:
//   1. WRAAS_PROVIDER=claude-code  → uses `claude` CLI (Max Plan)
//   2. WRAAS_PROVIDER=api          → uses Anthropic API (requires key)
//   3. WRAAS_API_KEY or ANTHROPIC_API_KEY set → uses Anthropic API
//   4. `claude` CLI available       → uses Claude Code automatically
//   5. fallback                     → MockProvider (demo mode)

// NewProvider creates the appropriate LLM provider based on environment.
func NewProvider() Provider {
	// Explicit provider override
	switch os.Getenv("WRAAS_PROVIDER") {
	case "claude-code":
		return &ClaudeCodeProvider{}
	case "api":
		return newAnthropicFromEnv()
	case "mock":
		return &MockProvider{}
	}

	// Auto-detect: API key takes priority
	apiKey := os.Getenv("WRAAS_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if apiKey != "" {
		model := os.Getenv("WRAAS_MODEL")
		if model == "" {
			model = "claude-sonnet-4-20250514"
		}
		return &AnthropicProvider{APIKey: apiKey, Model: model}
	}

	// Auto-detect: claude CLI available → use Max Plan
	if isClaudeCodeAvailable() {
		return &ClaudeCodeProvider{}
	}

	// Fallback: demo mode
	return &MockProvider{}
}

func newAnthropicFromEnv() Provider {
	apiKey := os.Getenv("WRAAS_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	model := os.Getenv("WRAAS_MODEL")
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}
	return &AnthropicProvider{APIKey: apiKey, Model: model}
}

// IsConfigured returns true if an LLM provider is explicitly selected or detected.
func IsConfigured() bool {
	if os.Getenv("WRAAS_PROVIDER") != "" {
		return true
	}
	if os.Getenv("WRAAS_API_KEY") != "" || os.Getenv("ANTHROPIC_API_KEY") != "" {
		return true
	}
	return isClaudeCodeAvailable()
}

// ProviderName returns a human-readable name for the active provider.
func ProviderName() string {
	switch os.Getenv("WRAAS_PROVIDER") {
	case "claude-code":
		return "Claude Code (Max Plan)"
	case "api":
		return "Anthropic API"
	case "mock":
		return "Demo mode"
	}
	if os.Getenv("WRAAS_API_KEY") != "" || os.Getenv("ANTHROPIC_API_KEY") != "" {
		return "Anthropic API"
	}
	if isClaudeCodeAvailable() {
		return "Claude Code (Max Plan)"
	}
	return "Demo mode"
}
