package storage

import (
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed schema.sql
var schemaSQL string

const currentSchemaVersion = 1

// runMigrations initializes the database schema and applies any migrations
func runMigrations(db *sql.DB) error {
	// Check if migrations table exists and get current version
	var version int
	err := db.QueryRow("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version)
	if err != nil && err != sql.ErrNoRows {
		// Migrations table might not exist yet, create schema from scratch
		if err := initSchema(db); err != nil {
			return fmt.Errorf("failed to initialize schema: %w", err)
		}
		return nil
	}

	// If we're at current version, nothing to do
	if version >= currentSchemaVersion {
		return nil
	}

	// Future: Add migration logic here for schema updates
	// For now, we only have version 1 (initial schema)

	return nil
}

// initSchema creates all tables and indexes from schema.sql
func initSchema(db *sql.DB) error {
	// Execute schema SQL
	if _, err := db.Exec(schemaSQL); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	// Record migration
	if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", currentSchemaVersion); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return nil
}
