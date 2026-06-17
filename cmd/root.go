// Package cmd contains all Cobra CLI commands for NagiosQL.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	devMode bool

	buildVersion   = "dev"
	buildDate      = "unknown"
)

// SetBuildInfo is called from main() with values injected via -ldflags.
func SetBuildInfo(version, date string) {
	buildVersion = version
	buildDate = date
}

// rootCmd is the base command; sub-commands are registered in their own files.
var rootCmd = &cobra.Command{
	Use:   "nagiosql",
	Short: "NagiosQL — Nagios Core configuration manager",
	Long: `NagiosQL manages Nagios Core monitoring configuration through a REST API
backed by a MariaDB database. It generates .cfg files, validates them with
nagios -v, and triggers graceful reloads via a file-based trigger watched by
reload-watcher.sh.`,
	// Serve is the default action when no sub-command is given.
	RunE: func(cmd *cobra.Command, args []string) error {
		return serveCmd.RunE(cmd, args)
	},
}

// Execute runs the root command. Called from main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.toml", "config file path")
	rootCmd.PersistentFlags().BoolVar(&devMode, "dev", false, "enable development mode (verbose logging)")
}
