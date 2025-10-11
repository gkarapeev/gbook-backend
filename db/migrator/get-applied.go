package migrator

import "database/sql"

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT name FROM applied_migrations")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	applied := make(map[string]bool)

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		applied[name] = true
	}

	return applied, nil
}
