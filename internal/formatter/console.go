package formatter

import (
	"fmt"

	"github.com/dawgdevv/apitestercli/internal/executor"
)

// ConsoleFormatter formats test results for terminal output
type ConsoleFormatter struct{}

// NewConsoleFormatter creates a new console formatter
func NewConsoleFormatter() *ConsoleFormatter {
	return &ConsoleFormatter{}
}

// FormatResult formats a single test result for console output
func (f *ConsoleFormatter) FormatResult(result executor.Result) string {
	if result.Passed {
		return fmt.Sprintf("✔ %s (%d) [%v]", result.Name, result.StatusCode, result.Duration)
	}
	return fmt.Sprintf("✖ %s (%v)", result.Name, result.Error)
}

// FormatSummary formats the test suite summary
func (f *ConsoleFormatter) FormatSummary(total, failed int) string {
	return fmt.Sprintf("\n%d tests , %d failed\n", total, failed)
}

// PrintResult prints a single result to stdout
func (f *ConsoleFormatter) PrintResult(result executor.Result) {
	fmt.Println(f.FormatResult(result))
}

// PrintSummary prints the summary to stdout
func (f *ConsoleFormatter) PrintSummary(total, failed int) {
	fmt.Print(f.FormatSummary(total, failed))
}
