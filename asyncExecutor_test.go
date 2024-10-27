package main

import (
	"sync"
	"testing"
)

func setAsyncExecutor() {
	serverConfig, _ := initConfig("sync", "50")
	executor := &AsyncTaskExecutor{
		taskQueue:        make(chan Task, 100),
		completedTasks:   make(map[int]struct{}),
		failureThreshold: serverConfig.FailureThreshold,
		retryCount:       3,
		wg:               sync.WaitGroup{},
		mu:               sync.RWMutex{},
	}
	db, _ := startDB()
	server = Server{
		Config:       serverConfig,
		TaskExecutor: executor,
		DB:           db,
	}
	server.TaskExecutor.Start()
}

func TestAsyncSubmitTask(t *testing.T) {
	setAsyncExecutor()
	defer clearDB()
	executor := server.TaskExecutor

	task1 := Task{TaskId: 1}

	success, err := executor.SubmitTask(task1)
	if !success || err != nil {
		t.Errorf("Expected success to be true, got %v, error: %v", success, err)
	}
}

func TestAsyncSubmitTask_AlreadyCompleted(t *testing.T) {
	setAsyncExecutor()
	defer clearDB()
	executor := server.TaskExecutor

	task1 := Task{TaskId: 1}
	//add duplicate task to db
	addTaskToDB(task1.TaskId)

	success, err := executor.SubmitTask(task1)
	if success || err == nil {
		t.Errorf("Expected success to be true, got %v, error: %v", success, err)
	}
}

func TestAsyncCompleteTask(t *testing.T) {
	setAsyncExecutor()
	defer clearDB()
	executor := server.TaskExecutor

	resultChan := make(chan int)
	task := Task{
		TaskId:     1,
		ResultChan: resultChan,
	}

	go func() {
		result := <-resultChan
		if result != 1 {
			t.Fatal("Expected result to be true")
		}
	}()

	success, err := executor.completeTask(task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !success {
		t.Fatal("Expected task to be completed successfully")
	}

	if !ifTaskCompleted(task.TaskId) {
		t.Fatalf("Expected task ID %d to be in completed tasks", task.TaskId)
	}
}

func TestAsyncFailTask(t *testing.T) {
	setAsyncExecutor()
	defer clearDB()
	executor := server.TaskExecutor

	resultChan := make(chan int)
	task := Task{
		TaskId:     1,
		ResultChan: resultChan,
	}

	go func() {
		result := <-resultChan
		if result != 2 {
			t.Fatal("Expected result to be false")
		}
	}()

	success, err := executor.failTask(task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !success {
		t.Fatal("Expected task to be marked as failed successfully")
	}
}
