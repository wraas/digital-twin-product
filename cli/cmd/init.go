package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/wraas/digital-twin-product/cli/internal/config"
	"github.com/wraas/digital-twin-product/cli/internal/engine"
	"github.com/wraas/digital-twin-product/cli/internal/output"
	"github.com/wraas/digital-twin-product/cli/internal/tui"
)

var forceInit bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise WRAAS in the current directory",
	Long: `Initialises WRAAS in the current directory. Creates wraas.yml with default
configuration and registers the directory with the engine.

WRAAS was already running. init creates the config file and makes the
relationship official.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVar(&forceInit, "force", false, "Overwrite existing config")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	configPath := cfgFile
	if configPath == "" {
		configPath = config.FileName
	}

	if config.Exists(configPath) && !forceInit {
		fmt.Fprintln(os.Stderr, output.ErrorStyle.Render("Config already exists. Use --force to overwrite."))
		os.Exit(engine.ExitBlockingViolation)
	}

	if config.Exists(configPath) && forceInit {
		// Archive existing config
		backupPath := fmt.Sprintf("%s.bak.%d", configPath, time.Now().Unix())
		if err := os.Rename(configPath, backupPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error archiving config: %v\n", err)
			os.Exit(engine.ExitConfigError)
		}
		if !quiet {
			output.Prompt(os.Stdout, fmt.Sprintf("Previous config archived to %s", output.DimStyle.Render(backupPath)))
		}
	}

	// Write with spinner
	var path string
	if !quiet {
		var err error
		path, err = writeWithSpinner()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(engine.ExitConfigError)
		}
	} else {
		dir, _ := os.Getwd()
		var err error
		path, err = config.WriteDefault(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(engine.ExitConfigError)
		}
	}

	if !quiet {
		fmt.Println()
		output.Prompt(os.Stdout, output.OkStyle.Render("Config written to ")+output.ValueStyle.Render(path))
		output.Prompt(os.Stdout, "WRAAS was already running. "+output.DimStyle.Render("init")+" creates the config file and makes the relationship official.")
		output.Latency(os.Stdout)
	}

	return nil
}

func writeWithSpinner() (string, error) {
	var path string
	_, err := tui.RunWithSpinner("Initialising WRAAS...", func() (string, error) {
		time.Sleep(400 * time.Millisecond) // Brief delay for effect
		dir, _ := os.Getwd()
		var err error
		path, err = config.WriteDefault(dir)
		return path, err
	})
	return path, err
}
