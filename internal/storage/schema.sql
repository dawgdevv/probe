-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Test suites table
CREATE TABLE IF NOT EXISTS test_suites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    yaml_content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Test runs table
CREATE TABLE IF NOT EXISTS test_runs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    suite_id INTEGER NOT NULL,
    started_at DATETIME NOT NULL,
    completed_at DATETIME,
    status TEXT NOT NULL, -- 'running', 'passed', 'failed', 'error'
    total_tests INTEGER DEFAULT 0,
    passed_tests INTEGER DEFAULT 0,
    failed_tests INTEGER DEFAULT 0,
    FOREIGN KEY (suite_id) REFERENCES test_suites(id) ON DELETE CASCADE
);

-- Test results table (individual test results within a run)
CREATE TABLE IF NOT EXISTS test_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id INTEGER NOT NULL,
    test_name TEXT NOT NULL,
    passed BOOLEAN NOT NULL,
    status_code INTEGER,
    error_message TEXT,
    duration_ms INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (run_id) REFERENCES test_runs(id) ON DELETE CASCADE
);

-- Schema migrations tracking
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_test_suites_project ON test_suites(project_id);
CREATE INDEX IF NOT EXISTS idx_test_runs_suite ON test_runs(suite_id);
CREATE INDEX IF NOT EXISTS idx_test_results_run ON test_results(run_id);
CREATE INDEX IF NOT EXISTS idx_test_runs_status ON test_runs(status);
