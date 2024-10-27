package test1

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	serverConfig, err := service.initConfig(mode, number)
	if err != nil {
		log.Fatal(err)
	}
	executor := &SyncTaskExecutor{
		taskQueue:        make([]Task, 0),
		completedTasks:   make(map[int]struct{}),
		failureThreshold: serverConfig.FailureThreshold,
		stopChan:         make(chan struct{}),
		retryCount:       3,
	}
	db, err := startDB()
	server = Server{
		Config:       serverConfig,
		TaskExecutor: executor,
		DB:           db,
	}

	server.TaskExecutor.Start()
}

func TestFunctionA(t *testing.T) {
	// Test using test1 environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "postgres://user:pass@localhost/test1db" {
		t.Fatalf("unexpected DATABASE_URL: %s", dbURL)
	}
	// Rest of the test
}
