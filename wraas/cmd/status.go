package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/wraas/digital-twin-product/wraas/internal/config"
	"github.com/wraas/digital-twin-product/wraas/internal/engine"
	"github.com/wraas/digital-twin-product/wraas/internal/output"
	"github.com/wraas/digital-twin-product/wraas/internal/tui"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display WRAAS operational status",
	Long: `Displays current WRAAS operational status, including engine state, sigh
calibration, and desertion rate. The desertion rate is included for
completeness. It has always been 0.00%.`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

// StatusOutput is the structured output for json/yaml formats.
type StatusOutput struct {
	Engine          string  `json:"engine" yaml:"engine"`
	CommitmentLevel string  `json:"commitment_level" yaml:"commitment_level"`
	SighCalibration string  `json:"sigh_calibration" yaml:"sigh_calibration"`
	DesertionRate   float64 `json:"desertion_rate" yaml:"desertion_rate"`
	LatencyMs       int     `json:"latency_ms" yaml:"latency_ms"`
	LastQuery       string  `json:"last_query" yaml:"last_query"`
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(engine.ExitConfigError)
	}

	state, _ := engine.LoadState()

	lastQuery := "no queries yet"
	if !state.LastQuery.IsZero() {
		lastQuery = state.LastQuery.Format("2006-01-02 15:04:05")
	}

	var messages []string
	if cfg.Commitment.Level != "FULL" {
		messages = append(messages, "Noted. Proceeding at FULL commitment as configured by design.")
	}

	format := output.ParseFormat(outputFormat)

	data := StatusOutput{
		Engine:          "RUNNING",
		CommitmentLevel: engine.CommitmentLevel(),
		SighCalibration: "ACTIVE",
		DesertionRate:   engine.DesertionRate(),
		LatencyMs:       engine.LatencyMs(),
		LastQuery:       lastQuery,
	}

	output.Write(os.Stdout, format, data, func(_ io.Writer) {
		fmt.Print(tui.RenderStatus(tui.StatusData{
			EngineStatus:    "RUNNING",
			CommitmentLevel: engine.CommitmentLevel(),
			SighCalibration: "ACTIVE",
			DesertionRate:   fmt.Sprintf("%.2f%%", engine.DesertionRate()),
			Latency:         engine.FormatLatency(),
			LastQuery:       lastQuery,
			Messages:        messages,
		}))
	})

	return nil
}
