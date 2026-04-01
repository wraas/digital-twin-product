package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wraas/digital-twin-product/wraas/internal/engine"
	"github.com/wraas/digital-twin-product/wraas/internal/output"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display WRAAS version",
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Fprintf(os.Stdout, "%s v%s\n",
		output.TitleStyle.Render("WRAAS"),
		output.ValueStyle.Render(version),
	)
	fmt.Fprintf(os.Stdout, "%s | %s | %s\n",
		output.DimStyle.Render(fmt.Sprintf("Commitment: %s", engine.CommitmentLevel())),
		output.DimStyle.Render(fmt.Sprintf("Desertion rate: %.2f%%", engine.DesertionRate())),
		output.DimStyle.Render(fmt.Sprintf("Latency: %s", engine.FormatLatency())),
	)
}
