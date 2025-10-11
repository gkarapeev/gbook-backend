package migrator

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func RunSingleMigration(db *sql.DB, migrationType, migrationName string) {
	fmt.Printf("Attempting to run '%s' migration for '%s'...\n", migrationType, migrationName)

	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("Failed to begin transaction: %v\n", err)
		return
	}

	var appliedName string
	err = tx.QueryRow(`SELECT name FROM applied_migrations WHERE name = $1`, migrationName).Scan(&appliedName)
	isApplied := err != sql.ErrNoRows

	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		fmt.Printf("Error checking migration status: %v\n", err)
		return
	}

	_ = godotenv.Load(".env")
	migrationsDir := os.Getenv("PG_MIGRATIONS_DIR")

	filename := fmt.Sprintf("%s.%s.sql", migrationName, migrationType)
	filePath := filepath.Join(migrationsDir, filename)

	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		tx.Rollback()
		fmt.Printf("Error reading migration file '%s': %v\n", filePath, err)
		return
	}

	sqlToExecute := string(sqlBytes)

	switch migrationType {
	case "up":
		if isApplied {
			fmt.Printf("Migration '%s' already applied. Skipping.\n", migrationName)
			tx.Rollback()
			return
		}

		if _, err := tx.Exec(sqlToExecute); err != nil {
			tx.Rollback()
			fmt.Printf("Error applying UP migration: %v\n", err)
			return
		}

		if _, err := tx.Exec(`INSERT INTO applied_migrations (name) VALUES ($1)`, migrationName); err != nil {
			tx.Rollback()
			fmt.Printf("Error inserting into applied_migrations table: %v\n", err)
			return
		}

		fmt.Printf("Successfully applied migration: %s\n", migrationName)

	case "down":
		if !isApplied {
			fmt.Printf("Migration '%s' not applied. Skipping rollback.\n", migrationName)
			tx.Rollback()
			return
		}

		if _, err := tx.Exec(sqlToExecute); err != nil {
			tx.Rollback()
			fmt.Printf("Error applying DOWN migration: %v\n", err)
			return
		}

		if _, err := tx.Exec(`DELETE FROM applied_migrations WHERE name = $1`, migrationName); err != nil {
			tx.Rollback()
			fmt.Printf("Error deleting from applied_migrations table: %v\n", err)
			return
		}

		fmt.Printf("Successfully rolled back migration: %s\n", migrationName)
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("Failed to commit transaction: %v\n", err)
		return
	}
}
