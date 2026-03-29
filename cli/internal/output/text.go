package output

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles for terminal output matching the documentation examples.
	PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	KeyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("253")).Bold(true)
	OkStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	ValueStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	WarnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	ErrorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	DimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	SighStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true)

	// Box styles for status dashboard.
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("253")).
			Bold(true).
			Underline(true)
)

// Prompt writes a "> " prefixed line (matching the docs terminal examples).
func Prompt(w io.Writer, s string) {
	fmt.Fprintf(w, "%s %s\n", PromptStyle.Render(">"), s)
}

// KeyValue writes a key-value pair with consistent alignment.
func KeyValue(w io.Writer, key string, value string, valueStyle lipgloss.Style) {
	fmt.Fprintf(w, "%s %-22s %s\n",
		PromptStyle.Render(">"),
		KeyStyle.Render(key+":"),
		valueStyle.Render(value),
	)
}

// Suggestion writes a non-blocking suggestion.
func Suggestion(w io.Writer, msg string) {
	fmt.Fprintf(w, "%s %s %s\n",
		PromptStyle.Render(">"),
		WarnStyle.Render("suggestion(non-blocking):"),
		msg,
	)
}

// Violation writes a blocking violation.
func Violation(w io.Writer, msg string) {
	fmt.Fprintf(w, "%s %s %s\n",
		PromptStyle.Render(">"),
		ErrorStyle.Render("VIOLATION(blocking):"),
		msg,
	)
}

// Sigh writes a sigh indicator.
func Sigh(w io.Writer, level string) {
	fmt.Fprintf(w, "%s %s\n",
		PromptStyle.Render(">"),
		SighStyle.Render(fmt.Sprintf("*sigh* [%s]", level)),
	)
}

// Latency writes the latency footer.
func Latency(w io.Writer) {
	fmt.Fprintf(w, "%s %s\n",
		PromptStyle.Render(">"),
		DimStyle.Render("Latency: 113ms"),
	)
}

// Banner writes the WRAAS ASCII banner.
func Banner(w io.Writer) {
	banner := DimStyle.Render(` __        __ ____      _        _     ____
 \ \      / /|  _ \    / \      / \   / ___|
  \ \ /\ / / | |_) |  / _ \    / _ \  \___ \
   \ V  V /  |  _ <  / ___ \  / ___ \  ___) |
    \_/\_/   |_| \_\/_/   \_\/_/   \_\|____/`)
	fmt.Fprintln(w, banner)
}
