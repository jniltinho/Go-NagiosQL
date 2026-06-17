package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version, build date, and Go runtime version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("nagiosql %s\n", buildVersion)
		fmt.Printf("  built:  %s\n", buildDate)
		fmt.Printf("  go:     %s\n", runtime.Version())
		fmt.Printf("  os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
