package llm

import (
	"context"
	"strings"
)

// MockProvider returns deterministic demo responses for offline use.
type MockProvider struct{}

func (p *MockProvider) Complete(_ context.Context, req Request) (Response, error) {
	// Detect if this is a review or query based on system prompt
	if strings.Contains(req.SystemPrompt, "code review") {
		return Response{Content: mockReviewResponse(req.UserMessage)}, nil
	}
	return Response{Content: mockQueryResponse(req.UserMessage)}, nil
}

func mockQueryResponse(input string) string {
	return `## Decision Matrix Evaluation

**Query analyzed.** Generating full option space (including obviously wrong ones).

### Options Evaluated

1. **Option A: Direct approach** [VIABLE]
   Straightforward implementation that addresses the stated requirements.
   Confidence: 87.3%

2. **Option B: Over-engineered solution** [REJECTED]
   Introduces unnecessary abstraction layers. Solves problems that do not exist yet.
   Rejection rationale: Architectural complexity without demonstrated need.
   See DME-0047.

3. **Option C: "Just ship it"** [REJECTED]
   No explanation needed.

### Recommendation

Option A. It solves the stated problem without introducing the problems described in Options B and C. The evaluation is complete. The wrong options have been documented. This is how it works.

**Confidence:** 87.3%
**Sigh level:** [MILD]

> Note: LLM provider not configured. This is demonstration output. For full evaluation, set WRAAS_API_KEY or ANTHROPIC_API_KEY.`
}

func mockReviewResponse(input string) string {
	return `## PR Analysis

The PR description is present but could more explicitly address "why" this change is being made. The "what" is clear from the diff. The "why" requires inference.

Scope consistency across commits is acceptable. No architectural concerns detected.

**Sigh level:** [MILD]

> Note: LLM provider not configured. This is demonstration output.`
}
