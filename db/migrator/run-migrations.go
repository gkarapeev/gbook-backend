package migrator

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"
)

func RunAllMigrations(db *sql.DB) {
	fmt.Println("Attempting to run all unapplied migrations...")

	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		fmt.Printf("Error getting applied migrations: %v\n", err)
		return
	}

	_ = godotenv.Load(".env")
	migrationsDir := os.Getenv("PG_MIGRATIONS_DIR")

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		fmt.Printf("Error reading migrations directory: %v\n", err)
		return
	}

	var upFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			upFiles = append(upFiles, file.Name())
		}
	}
	sort.Strings(upFiles)

	for _, filename := range upFiles {
		migrationName := strings.TrimSuffix(filename, ".up.sql")
		if !appliedMigrations[migrationName] {
			RunSingleMigration(db, "up", migrationName)
		} else {
			fmt.Printf("Migration '%s' is already applied. Skipping.\n", migrationName)
		}
	}

	fmt.Println("All available migrations processed.")
}
