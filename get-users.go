package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func getUsersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("["))
	first := true
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		if !first {
			w.Write([]byte(","))
		}
		first = false
		w.Write([]byte(`{"id":` +
			fmt.Sprintf("%d", id) + `,"name":"` + name + `"}`))
	}
	w.Write([]byte("]"))
}
