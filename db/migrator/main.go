package migrator

import (
	"database/sql"
)

func EnsureMigrationsTable(db *sql.DB) error {
	sql := `
		CREATE TABLE IF NOT EXISTS applied_migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		);
	`
	_, err := db.Exec(sql)
	return err
}
