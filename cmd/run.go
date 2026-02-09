package cmd

import (
	"fmt"
	"os"

	"github.com/dawgdevv/probe/internal/executor"
	"github.com/dawgdevv/probe/internal/formatter"
	"github.com/dawgdevv/probe/internal/loader"
	"github.com/dawgdevv/probe/internal/service"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run <file>",
	Short: "Run API tests from a YAML file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		suite, err := loader.LoadSuite(args[0])
		if err != nil {
			fmt.Println("Errors:", err)
			os.Exit(1)
		}

		// Create formatter for console output
		consoleFormatter := formatter.NewConsoleFormatter()

		// Create runner with progress callback for real-time output
		runner := service.NewRunner(service.RunOptions{
			MaxConcurrent: 10,
			ProgressCallback: func(result executor.Result) {
				consoleFormatter.PrintResult(result)
			},
		})

		// Execute test suite
		results, err := runner.RunSuite(suite)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		// Print summary
		failed := service.CountFailures(results)
		consoleFormatter.PrintSummary(len(suite.Tests), failed)

		if failed > 0 {
			os.Exit(1)
		}
	},
}
