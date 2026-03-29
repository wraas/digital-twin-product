package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is injected at build time via ldflags. Default for local dev builds.
var version = "1.3-dev"

var (
	cfgFile      string
	outputFormat string
	quiet        bool
	verbose      bool
)

var rootCmd = &cobra.Command{
	Use:     "wraas",
	Short:   "WRAAS — Weighted Romain Algorithmic Approximation Software",
	Version: version,
	Long: fmt.Sprintf(`WRAAS v%s — Weighted Romain Algorithmic Approximation Software

A high-fidelity neural network delivering the exact same brilliant insights
— and signature sighs — as Romain himself.

Desertion rate: 0.00%%. Latency: 113ms. Commitment: FULL.`, version),
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("WRAAS v%s\n", version))
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to wraas.yml (default: ./wraas.yml)")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "output", "text", "Output format: text, json, yaml")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "Suppress all output except final result. WRAAS will still sigh internally.")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Include full decision matrix in output. Long.")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
