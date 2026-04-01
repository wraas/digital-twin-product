package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wraas/digital-twin-product/wraas/internal/config"
	"github.com/wraas/digital-twin-product/wraas/internal/engine"
	"github.com/wraas/digital-twin-product/wraas/internal/output"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long: `Read or set individual configuration values without editing wraas.yml directly.
Changes made via wraas config set are written to the config file and take
effect immediately. Most of them.`,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Read a config value",
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigSet,
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(engine.ExitConfigError)
	}

	value, err := config.GetValue(cfg, args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(engine.ExitConfigError)
	}

	fmt.Println(value)
	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	configPath := cfgFile
	if configPath == "" {
		configPath = config.FileName
	}

	if !config.Exists(configPath) {
		fmt.Fprintf(os.Stderr, "Error: %s not found. Run 'wraas init' first.\n", configPath)
		os.Exit(engine.ExitConfigError)
	}

	messages, err := config.SetValue(configPath, args[0], args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(engine.ExitConfigError)
	}

	for _, msg := range messages {
		if !quiet {
			output.Prompt(os.Stdout, output.WarnStyle.Render(msg))
		}
	}

	if !quiet {
		output.Prompt(os.Stdout, fmt.Sprintf("%s = %s",
			output.KeyStyle.Render(args[0]),
			output.ValueStyle.Render(args[1]),
		))
		output.Latency(os.Stdout)
	}

	return nil
}
