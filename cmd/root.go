package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "apitestercli",
	Short: "API testing CLI for CI/CD",
	Long:  "apitester runs api tests defined in YAML files and exits with CI-Friendly status codes.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
