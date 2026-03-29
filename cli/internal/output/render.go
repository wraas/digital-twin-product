package output

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles for rendered LLM response.
var (
	h1Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true).
		Underline(true)

	h2Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)

	h3Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("75")).
		Bold(true)

	viableStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	rejectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	dmeRefStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	confidencePctStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")).
				Bold(true)

	renderBoldStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("253"))

	hrRenderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	blockquoteRenderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Italic(true)

	codeInlineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	listBulletStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)
)

// Patterns for matching response elements.
var (
	h1Pattern         = regexp.MustCompile(`^##\s+(.+)$`)
	h2Pattern         = regexp.MustCompile(`^###\s+(.+)$`)
	h3Pattern         = regexp.MustCompile(`^####\s+(.+)$`)
	hrPattern         = regexp.MustCompile(`^---+$`)
	blockquotePattern = regexp.MustCompile(`^>\s*(.*)$`)
	boldPattern       = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	codeInlinePattern = regexp.MustCompile("`([^`]+)`")
	viablePattern     = regexp.MustCompile(`\[VIABLE[^\]]*\]`)
	rejectedPattern   = regexp.MustCompile(`\[REJECTED[^\]]*\]`)
	dmePattern        = regexp.MustCompile(`(DME-\d{4})`)
	confidencePattern = regexp.MustCompile(`(\d+\.?\d*%)`)
	sighTagPattern    = regexp.MustCompile(`\[(SILENT|MILD|MODERATE|DEEP|EXISTENTIAL)\]`)
	listItemPattern   = regexp.MustCompile(`^(\d+\.)\s`)
)

// RenderResponse colorizes an LLM response for terminal display.
func RenderResponse(w io.Writer, response string) {
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		rendered := renderLine(line)
		fmt.Fprintln(w, rendered)
	}
}

func renderLine(line string) string {
	trimmed := strings.TrimSpace(line)

	// Empty line
	if trimmed == "" {
		return ""
	}

	// Horizontal rule
	if hrPattern.MatchString(trimmed) {
		return hrRenderStyle.Render("  ─────────────────────────────────────────────────")
	}

	// Headers — check h2/h3 before h1 since ### starts with ##
	if m := h3Pattern.FindStringSubmatch(trimmed); m != nil {
		return "\n  " + h3Style.Render(renderInline(m[1]))
	}
	if m := h2Pattern.FindStringSubmatch(trimmed); m != nil {
		return "\n  " + h2Style.Render(renderInline(m[1]))
	}
	if m := h1Pattern.FindStringSubmatch(trimmed); m != nil {
		return "\n  " + h1Style.Render(renderInline(m[1]))
	}

	// Blockquotes
	if m := blockquotePattern.FindStringSubmatch(trimmed); m != nil {
		inner := renderInline(m[1])
		return "  " + blockquoteRenderStyle.Render("│ "+inner)
	}

	// Process inline formatting for regular lines
	result := renderInline(trimmed)

	// Numbered list items — colorize the number
	if m := listItemPattern.FindStringSubmatch(trimmed); m != nil {
		result = listItemPattern.ReplaceAllStringFunc(result, func(s string) string {
			parts := listItemPattern.FindStringSubmatch(s)
			return listBulletStyle.Render(parts[1]) + " "
		})
	}

	return "  " + result
}

// renderInline applies inline styling (bold, code, VIABLE/REJECTED, DME, etc.)
func renderInline(text string) string {
	// Order matters: replace tokens that could be inside bold first,
	// then bold, then remaining patterns.

	// Inline code `...` — do first so content inside won't be double-styled
	text = codeInlinePattern.ReplaceAllStringFunc(text, func(s string) string {
		inner := s[1 : len(s)-1]
		return codeInlineStyle.Render(inner)
	})

	// Bold **...**
	text = boldPattern.ReplaceAllStringFunc(text, func(s string) string {
		inner := s[2 : len(s)-2]
		// Apply inner styling within bold text
		inner = styleTokens(inner)
		return renderBoldStyle.Render(inner)
	})

	// Style remaining tokens outside bold
	text = styleTokens(text)

	return text
}

// styleTokens applies WRAAS-specific token styling.
func styleTokens(text string) string {
	// [VIABLE ...] → green bold (preserves trailing text like "— noted with mild appreciation")
	text = viablePattern.ReplaceAllStringFunc(text, func(s string) string {
		inner := s[1 : len(s)-1] // strip brackets
		return viableStyle.Render("✓ " + inner)
	})

	// [REJECTED ...] → red bold (preserves trailing text like "— see DME-0047")
	text = rejectedPattern.ReplaceAllStringFunc(text, func(s string) string {
		inner := s[1 : len(s)-1] // strip brackets
		return rejectedStyle.Render("✗ " + inner)
	})

	// DME-XXXX → orange bold
	text = dmePattern.ReplaceAllStringFunc(text, func(s string) string {
		return dmeRefStyle.Render(s)
	})

	// Sigh level tags → italic gray
	text = sighTagPattern.ReplaceAllStringFunc(text, func(s string) string {
		return SighStyle.Render(s)
	})

	// Confidence percentages on lines containing "confidence" (case-insensitive)
	if strings.Contains(strings.ToLower(text), "confidence") {
		text = confidencePattern.ReplaceAllStringFunc(text, func(s string) string {
			return confidencePctStyle.Render(s)
		})
	}

	return text
}
