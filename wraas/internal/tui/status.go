package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wraas/digital-twin-product/wraas/internal/engine"
)

var (
	statusBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(0, 2).
			MarginTop(1).
			MarginBottom(1)

	statusTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	statusKey = lipgloss.NewStyle().
			Foreground(lipgloss.Color("253")).
			Bold(true).
			Width(24)

	statusOk = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	statusValue = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))

	statusDim = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

// StatusData holds the data for the status dashboard.
type StatusData struct {
	EngineStatus    string
	CommitmentLevel string
	SighCalibration string
	DesertionRate   string
	Latency         string
	LastQuery       string
	Messages        []string
}

// RenderStatus renders the status dashboard as a styled box.
func RenderStatus(data StatusData) string {
	var b strings.Builder

	b.WriteString(statusTitle.Render("WRAAS Engine Status"))
	b.WriteString("\n\n")

	rows := []struct {
		key   string
		value string
		style lipgloss.Style
	}{
		{"Engine", data.EngineStatus, statusOk},
		{"Commitment level", data.CommitmentLevel, statusOk},
		{"Sigh calibration", data.SighCalibration, statusOk},
		{"Desertion rate", data.DesertionRate, statusValue},
		{"Latency", data.Latency, statusValue},
		{"Last query", data.LastQuery, statusDim},
	}

	for _, row := range rows {
		b.WriteString(fmt.Sprintf("%s %s\n",
			statusKey.Render(row.key+":"),
			row.style.Render(row.value),
		))
	}

	result := statusBorder.Render(b.String())

	if len(data.Messages) > 0 {
		result += "\n"
		for _, msg := range data.Messages {
			result += statusDim.Render(fmt.Sprintf("  ℹ %s", msg)) + "\n"
		}
	}

	_ = engine.LatencyMs() // ensure package is referenced
	return result
}
