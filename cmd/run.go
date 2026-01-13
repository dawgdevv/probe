package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/dawgdevv/apitestercli/internal/executor"
	"github.com/dawgdevv/apitestercli/internal/loader"
	"github.com/dawgdevv/apitestercli/pkg/models"
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
		baseURL := suite.Env["base_url"]
		if baseURL == "" {
			fmt.Println("base_url not defined in env")
			os.Exit(1)
		}

		// Buffered channel to collect all results
		results := make(chan executor.Result, len(suite.Tests))

		// Default concurrency limit (prevents overwhelming the API)
		maxConcurrent := 10
		if len(suite.Tests) < maxConcurrent {
			maxConcurrent = len(suite.Tests)
		}
		sem := make(chan struct{}, maxConcurrent)

		var wg sync.WaitGroup

		// Run tests in parallel with controlled concurrency
		for _, test := range suite.Tests {
			wg.Add(1)
			sem <- struct{}{} // Acquire semaphore

			go func(t models.TestCase) {
				defer wg.Done()
				defer func() { <-sem }() // Release semaphore

				result := executor.RunTest(baseURL, suite.Env, t)
				results <- result
			}(test)
		}

		// Close results channel when all goroutines finish
		go func() {
			wg.Wait()
			close(results)
		}()

		failed := 0

		for result := range results {
			if result.Passed {
				fmt.Printf("✔ %s (%d) [%v]\n", result.Name, result.StatusCode, result.Duration)
			} else {
				fmt.Printf("✖ %s (%v)\n", result.Name, result.Error)
				failed++
			}
		}

		fmt.Printf("\n%d tests , %d failed\n", len(suite.Tests), failed)

		if failed > 0 {
			os.Exit(1)
		}
	},
}
