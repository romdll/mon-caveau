package database

import (
	"fmt"
)

type Migration struct {
	Version int
	SQL     string
}

func ApplyMigrations() error {
	migrations := []Migration{}

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version INT PRIMARY KEY)`)
	if err != nil {
		return fmt.Errorf("failed to ensure migrations table exists: %w", err)
	}

	var currentVersion int
	err = db.QueryRow(`SELECT IFNULL(MAX(version), 0) FROM schema_migrations`).Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("failed to fetch current schema version: %w", err)
	}

	for _, migration := range migrations {
		if migration.Version > currentVersion {
			logger.Printf("Applying migration %d...\n", migration.Version)

			if _, err := db.Exec(migration.SQL); err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
			}

			if _, err := db.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, migration.Version); err != nil {
				return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
			}

			logger.Printf("Migration %d applied successfully\n", migration.Version)
		}
	}

	logger.Println("All migrations applied")
	return nil
}
