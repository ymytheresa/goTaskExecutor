package service

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

	createTableSQL := `CREATE TABLE IF NOT EXISTS completedTasks (
		id INTEGER PRIMARY KEY,
		task_id INTEGER NOT NULL
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_task_id ON completedTasks(task_id);`
	if _, err := db.Exec(createIndexSQL); err != nil {
		return nil, fmt.Errorf("failed to create index on task_id: %w", err)
	}
	return db, nil
}

func printDB() {
	rows, err := server.DB.Query("SELECT id, task_id FROM completedTasks")
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

func addTaskToDB(taskID int) error {
	fmt.Printf("Adding to db task ID: %d to completedTasks\n", taskID)
	_, err := server.DB.Exec("INSERT INTO completedTasks (task_id) VALUES (?)", taskID)
	return err
}

func ifTaskCompleted(taskID int) bool {
	fmt.Printf("Checking if task ID: %d is completed in db\n", taskID)
	var count int
	err := server.DB.QueryRow("SELECT COUNT(*) FROM completedTasks WHERE task_id = ?", taskID).Scan(&count)
	if err != nil {
		log.Printf("Failed to check if task is completed: %v", err)
		return false
	}
	return count > 0
}

func readCompletedTasksFromDB(completedTasks *map[int]struct{}) {
	rows, err := server.DB.Query("SELECT task_id FROM completedTasks")
	if err != nil {
		log.Fatalf("Failed to query completedTasks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var taskID int
		if err := rows.Scan(&taskID); err != nil {
			log.Printf("Failed to scan row: %v", err)
		}
		(*completedTasks)[taskID] = struct{}{}
	}
}
