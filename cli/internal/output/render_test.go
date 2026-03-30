package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestVisibleLen(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"hello", 5},
		{"\x1b[38;5;39mhello\x1b[0m", 5},
		{"\x1b[1;38;5;196m✗ REJECTED\x1b[0m", 10},
		{"", 0},
		{"no ansi here", 12},
	}
	for _, tt := range tests {
		got := visibleLen(tt.input)
		if got != tt.want {
			t.Errorf("visibleLen(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestWrapRendered_ShortLine(t *testing.T) {
	raw := "Short line."
	rendered := renderLine(raw)
	result := wrapRendered(raw, rendered, 120)
	if len(result) != 1 {
		t.Errorf("expected 1 line, got %d", len(result))
	}
}

func TestWrapRendered_LongLine(t *testing.T) {
	raw := "This is a very long line that should be wrapped because it exceeds the maximum width that we have configured for the output rendering"
	rendered := renderLine(raw)
	result := wrapRendered(raw, rendered, 60)
	if len(result) < 2 {
		t.Errorf("expected multiple lines, got %d", len(result))
	}
	// Each visible line should be at most 60 chars
	for i, line := range result {
		vl := visibleLen(line)
		if vl > 60 {
			t.Errorf("line %d visible length %d exceeds max 60", i, vl)
		}
	}
}

func TestWrapRendered_PreservesFormatting(t *testing.T) {
	raw := "This mentions DME-0404 and has a **bold section** in the text that should still render after wrapping"
	rendered := renderLine(raw)
	result := wrapRendered(raw, rendered, 50)
	joined := strings.Join(result, "\n")
	// Should contain ANSI codes for DME reference (orange)
	if !strings.Contains(joined, "DME-0404") {
		t.Error("wrapped output should contain DME-0404")
	}
}

func TestRenderResponse_MaxWidthZero(t *testing.T) {
	var buf bytes.Buffer
	input := "Short line.\n\nAnother line."
	RenderResponse(&buf, input, 0)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) < 2 {
		t.Errorf("expected at least 2 lines, got %d", len(lines))
	}
}

func TestRenderResponse_MaxWidthApplied(t *testing.T) {
	var buf bytes.Buffer
	input := "This is a very long line that definitely exceeds forty characters and should be wrapped into multiple lines by the renderer"
	RenderResponse(&buf, input, 50)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) < 2 {
		t.Errorf("expected wrapping into multiple lines, got %d lines", len(lines))
	}
}
