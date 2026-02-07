package formatter

import (
	"encoding/json"
	"time"

	"github.com/dawgdevv/apitestercli/internal/executor"
)

// JSONFormatter formats test results as JSON
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// TestResultJSON represents a test result in JSON format
type TestResultJSON struct {
	Name       string `json:"name"`
	Passed     bool   `json:"passed"`
	StatusCode int    `json:"status_code,omitempty"`
	Error      string `json:"error,omitempty"`
	Duration   string `json:"duration"`
}

// SuiteResultJSON represents the complete suite results
type SuiteResultJSON struct {
	TotalTests  int              `json:"total_tests"`
	PassedTests int              `json:"passed_tests"`
	FailedTests int              `json:"failed_tests"`
	Results     []TestResultJSON `json:"results"`
	Timestamp   time.Time        `json:"timestamp"`
}

// Format converts results to JSON structure
func (f *JSONFormatter) Format(results []executor.Result) SuiteResultJSON {
	passed := 0
	failed := 0
	jsonResults := make([]TestResultJSON, len(results))

	for i, result := range results {
		errorMsg := ""
		if result.Error != nil {
			errorMsg = result.Error.Error()
			failed++
		} else {
			passed++
		}

		jsonResults[i] = TestResultJSON{
			Name:       result.Name,
			Passed:     result.Passed,
			StatusCode: result.StatusCode,
			Error:      errorMsg,
			Duration:   result.Duration.String(),
		}
	}

	return SuiteResultJSON{
		TotalTests:  len(results),
		PassedTests: passed,
		FailedTests: failed,
		Results:     jsonResults,
		Timestamp:   time.Now(),
	}
}

// Marshal converts results to JSON bytes
func (f *JSONFormatter) Marshal(results []executor.Result) ([]byte, error) {
	formatted := f.Format(results)
	return json.MarshalIndent(formatted, "", "  ")
}
