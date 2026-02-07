package storage

import "time"

// Project represents a collection of test suites
type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TestSuite represents a stored test suite
type TestSuite struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	Name        string    `json:"name"`
	YAMLContent string    `json:"yaml_content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TestRun represents an execution of a test suite
type TestRun struct {
	ID          int64      `json:"id"`
	SuiteID     int64      `json:"suite_id"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Status      string     `json:"status"`
	TotalTests  int        `json:"total_tests"`
	PassedTests int        `json:"passed_tests"`
	FailedTests int        `json:"failed_tests"`
}

// TestResult represents a single test result within a run
type TestResult struct {
	ID           int64     `json:"id"`
	RunID        int64     `json:"run_id"`
	TestName     string    `json:"test_name"`
	Passed       bool      `json:"passed"`
	StatusCode   int       `json:"status_code,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
	DurationMs   int64     `json:"duration_ms"`
	CreatedAt    time.Time `json:"created_at"`
}
