package service

import (
	"fmt"
	"sync"

	"github.com/dawgdevv/apitestercli/internal/executor"
	"github.com/dawgdevv/apitestercli/pkg/models"
)

// Runner is responsible for executing test suites
type Runner struct {
	options RunOptions
}

// NewRunner creates a new test runner with the given options
func NewRunner(options RunOptions) *Runner {
	if options.MaxConcurrent <= 0 {
		options.MaxConcurrent = 10
	}
	return &Runner{options: options}
}

// RunSuite executes all tests in a suite and returns the results
func (r *Runner) RunSuite(suite *models.TestSuite) ([]executor.Result, error) {
	baseURL := suite.Env["base_url"]
	if baseURL == "" {
		return nil, fmt.Errorf("base_url not defined in env")
	}

	// Buffered channel to collect all results
	resultsChan := make(chan executor.Result, len(suite.Tests))

	// Adjust concurrency limit based on test count
	maxConcurrent := r.options.MaxConcurrent
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

			// Send to progress callback if configured
			if r.options.ProgressCallback != nil {
				r.options.ProgressCallback(result)
			}

			resultsChan <- result
		}(test)
	}

	// Close results channel when all goroutines finish
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect all results
	var results []executor.Result
	for result := range resultsChan {
		results = append(results, result)
	}

	return results, nil
}

// CountFailures returns the number of failed tests in the results
func CountFailures(results []executor.Result) int {
	failed := 0
	for _, result := range results {
		if !result.Passed {
			failed++
		}
	}
	return failed
}
