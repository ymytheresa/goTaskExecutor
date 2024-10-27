package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func startDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./completedTasks.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS completedTasks (
		id INTEGER PRIMARY KEY,
		task_id INTEGER NOT NULL
	);`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	printDB(db)
	return db, nil
}

func printDB(db *sql.DB) {
	rows, err := db.Query("SELECT id, task_id FROM completedTasks")
	if err != nil {
		log.Fatalf("Failed to query completedTasks: %v", err)
	}
	defer rows.Close()

	log.Println("Current Entries in completedTasks:")
	for rows.Next() {
		var id, taskID int
		if err := rows.Scan(&id, &taskID); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		log.Printf("ID: %d, Task ID: %d", id, taskID)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Row iteration error: %v", err)
	}
}
