package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Store provides database operations for API Tester CLI
type Store struct {
	db *sql.DB
}

// NewStore creates a new storage instance and initializes the database
func NewStore(dataDir string) (*Store, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open database connection
	dbPath := filepath.Join(dataDir, "apitester.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(1) // SQLite works best with single connection
	db.SetMaxIdleConns(1)

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	store := &Store{db: db}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return store, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// --- Project operations ---

// CreateProject creates a new project
func (s *Store) CreateProject(name, description string) (*Project, error) {
	result, err := s.db.Exec(
		"INSERT INTO projects (name, description) VALUES (?, ?)",
		name, description,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get project ID: %w", err)
	}

	return s.GetProject(id)
}

// GetProject retrieves a project by ID
func (s *Store) GetProject(id int64) (*Project, error) {
	var project Project
	err := s.db.QueryRow(
		"SELECT id, name, description, created_at, updated_at FROM projects WHERE id = ?",
		id,
	).Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &project, nil
}

// ListProjects retrieves all projects
func (s *Store) ListProjects() ([]Project, error) {
	rows, err := s.db.Query(
		"SELECT id, name, description, created_at, updated_at FROM projects ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// --- Test Suite operations ---

// CreateTestSuite creates a new test suite
func (s *Store) CreateTestSuite(projectID int64, name, yamlContent string) (*TestSuite, error) {
	result, err := s.db.Exec(
		"INSERT INTO test_suites (project_id, name, yaml_content) VALUES (?, ?, ?)",
		projectID, name, yamlContent,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create test suite: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get suite ID: %w", err)
	}

	return s.GetTestSuite(id)
}

// GetTestSuite retrieves a test suite by ID
func (s *Store) GetTestSuite(id int64) (*TestSuite, error) {
	var suite TestSuite
	err := s.db.QueryRow(
		"SELECT id, project_id, name, yaml_content, created_at, updated_at FROM test_suites WHERE id = ?",
		id,
	).Scan(&suite.ID, &suite.ProjectID, &suite.Name, &suite.YAMLContent, &suite.CreatedAt, &suite.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get test suite: %w", err)
	}

	return &suite, nil
}

// ListTestSuites retrieves all test suites for a project
func (s *Store) ListTestSuites(projectID int64) ([]TestSuite, error) {
	rows, err := s.db.Query(
		"SELECT id, project_id, name, yaml_content, created_at, updated_at FROM test_suites WHERE project_id = ? ORDER BY created_at DESC",
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list test suites: %w", err)
	}
	defer rows.Close()

	var suites []TestSuite
	for rows.Next() {
		var suite TestSuite
		if err := rows.Scan(&suite.ID, &suite.ProjectID, &suite.Name, &suite.YAMLContent, &suite.CreatedAt, &suite.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan test suite: %w", err)
		}
		suites = append(suites, suite)
	}

	return suites, nil
}

// --- Test Run operations ---

// CreateTestRun creates a new test run
func (s *Store) CreateTestRun(suiteID int64, totalTests int) (*TestRun, error) {
	result, err := s.db.Exec(
		"INSERT INTO test_runs (suite_id, started_at, status, total_tests) VALUES (?, ?, ?, ?)",
		suiteID, time.Now(), "running", totalTests,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create test run: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get run ID: %w", err)
	}

	return s.GetTestRun(id)
}

// GetTestRun retrieves a test run by ID
func (s *Store) GetTestRun(id int64) (*TestRun, error) {
	var run TestRun
	var completedAt sql.NullTime

	err := s.db.QueryRow(
		"SELECT id, suite_id, started_at, completed_at, status, total_tests, passed_tests, failed_tests FROM test_runs WHERE id = ?",
		id,
	).Scan(&run.ID, &run.SuiteID, &run.StartedAt, &completedAt, &run.Status, &run.TotalTests, &run.PassedTests, &run.FailedTests)

	if err != nil {
		return nil, fmt.Errorf("failed to get test run: %w", err)
	}

	if completedAt.Valid {
		run.CompletedAt = &completedAt.Time
	}

	return &run, nil
}

// CompleteTestRun marks a test run as completed
func (s *Store) CompleteTestRun(id int64, status string, passed, failed int) error {
	_, err := s.db.Exec(
		"UPDATE test_runs SET completed_at = ?, status = ?, passed_tests = ?, failed_tests = ? WHERE id = ?",
		time.Now(), status, passed, failed, id,
	)
	if err != nil {
		return fmt.Errorf("failed to complete test run: %w", err)
	}
	return nil
}

// ListTestRuns retrieves all test runs for a suite
func (s *Store) ListTestRuns(suiteID int64, limit int) ([]TestRun, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := s.db.Query(
		"SELECT id, suite_id, started_at, completed_at, status, total_tests, passed_tests, failed_tests FROM test_runs WHERE suite_id = ? ORDER BY started_at DESC LIMIT ?",
		suiteID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list test runs: %w", err)
	}
	defer rows.Close()

	var runs []TestRun
	for rows.Next() {
		var run TestRun
		var completedAt sql.NullTime
		if err := rows.Scan(&run.ID, &run.SuiteID, &run.StartedAt, &completedAt, &run.Status, &run.TotalTests, &run.PassedTests, &run.FailedTests); err != nil {
			return nil, fmt.Errorf("failed to scan test run: %w", err)
		}
		if completedAt.Valid {
			run.CompletedAt = &completedAt.Time
		}
		runs = append(runs, run)
	}

	return runs, nil
}

// --- Test Result operations ---

// SaveTestResult saves a single test result
func (s *Store) SaveTestResult(runID int64, testName string, passed bool, statusCode int, errorMessage string, durationMs int64) error {
	_, err := s.db.Exec(
		"INSERT INTO test_results (run_id, test_name, passed, status_code, error_message, duration_ms) VALUES (?, ?, ?, ?, ?, ?)",
		runID, testName, passed, statusCode, errorMessage, durationMs,
	)
	if err != nil {
		return fmt.Errorf("failed to save test result: %w", err)
	}
	return nil
}

// GetTestResults retrieves all results for a test run
func (s *Store) GetTestResults(runID int64) ([]TestResult, error) {
	rows, err := s.db.Query(
		"SELECT id, run_id, test_name, passed, status_code, error_message, duration_ms, created_at FROM test_results WHERE run_id = ? ORDER BY created_at",
		runID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get test results: %w", err)
	}
	defer rows.Close()

	var results []TestResult
	for rows.Next() {
		var result TestResult
		if err := rows.Scan(&result.ID, &result.RunID, &result.TestName, &result.Passed, &result.StatusCode, &result.ErrorMessage, &result.DurationMs, &result.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan test result: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}
