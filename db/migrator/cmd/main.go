package main

import (
	"fmt"
	"os"
	"this_project_id_285410/db"
	"this_project_id_285410/db/migrator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/main.go [up|down|run-all] [migration_name]")
		os.Exit(1)
	}

	command := os.Args[1]

	db, err := db.OpenDB()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := migrator.EnsureMigrationsTable(db); err != nil {
		fmt.Printf("Error ensuring migrations table: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "up", "down":
		if len(os.Args) != 3 {
			fmt.Println("Usage: go run cmd/main.go [up|down] [migration_name]")
			os.Exit(1)
		}
		migrator.RunSingleMigration(db, command, os.Args[2])
	case "run-all":
		migrator.RunAllMigrations(db)
	default:
		fmt.Println("Invalid command. Use 'up', 'down', or 'run-all'.")
		os.Exit(1)
	}
}
