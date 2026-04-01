package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"github.com/wraas/digital-twin-product/wraas/internal/tui"
)

// version is injected at build time via ldflags. Default for local dev builds.
var version = "1.3-dev"

var (
	cfgFile      string
	outputFormat string
	quiet        bool
	verbose      bool
	colorFlag    string
	noSpinner    bool
)

var rootCmd = &cobra.Command{
	Use:     "wraas",
	Short:   "WRAAS — Weighted Romain Algorithmic Approximation Software",
	Version: version,
	Long: fmt.Sprintf(`WRAAS v%s — Weighted Romain Algorithmic Approximation Software

A high-fidelity neural network delivering the exact same brilliant insights
— and signature sighs — as Romain himself.

Desertion rate: 0.00%%. Latency: 113ms. Commitment: FULL.`, version),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		applyFlags()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println()
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
	},
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("WRAAS v%s\n", version))
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to wraas.yml (default: ./wraas.yml)")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "output", "text", "Output format: text, json, yaml")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "Suppress all output except final result. WRAAS will still sigh internally.")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Include full decision matrix in output. Long.")
	rootCmd.PersistentFlags().StringVar(&colorFlag, "color", "auto", "Color output: auto, always, never")
	rootCmd.PersistentFlags().BoolVar(&noSpinner, "no-spinner", false, "Disable the spinner animation")
}

func applyFlags() {
	tui.Disabled = noSpinner
	applyColorMode()
}

func applyColorMode() {
	switch colorFlag {
	case "always":
		os.Setenv("CLICOLOR_FORCE", "1")
		lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stdout, termenv.WithProfile(termenv.TrueColor)))
	case "never":
		os.Setenv("NO_COLOR", "1")
		lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stdout, termenv.WithProfile(termenv.Ascii)))
	}
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
