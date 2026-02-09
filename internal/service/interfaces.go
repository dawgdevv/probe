package service

import "github.com/dawgdevv/probe/internal/executor"

// ResultCollector defines the interface for collecting test results
type ResultCollector interface {
	Collect(result executor.Result)
	Results() []executor.Result
}

// ProgressCallback is called during test execution for real-time updates
type ProgressCallback func(result executor.Result)

// RunOptions configures test execution behavior
type RunOptions struct {
	MaxConcurrent    int
	ProgressCallback ProgressCallback
}
