package main

import (
	"testing"
)

func TestSubmitTask(t *testing.T) {
	executor := &SyncTaskExecutor{
		completedTasks: make(map[int]struct{}),
	}

	task1 := Task{TaskId: 1}
	task2 := Task{TaskId: 1}

	success, err := executor.SubmitTask(task1)
	if !success || err != nil {
		t.Errorf("Expected success to be true, got %v, error: %v", success, err)
	}

	executor.completedTasks[task1.TaskId] = struct{}{}
	success, err = executor.SubmitTask(task2)
	if success || err == nil || err.Error() != "task already completed" {
		t.Errorf("Expected success to be false and error to indicate task already completed, got %v, error: %v", success, err)
	}
}

func TestScheduleTask(t *testing.T) {
	executor := &SyncTaskExecutor{
		taskQueue: []Task{},
	}

	task := Task{TaskId: 1}
	executor.scheduleTask(task)

	if len(executor.taskQueue) != 1 {
		t.Errorf("Expected taskQueue length to be 1, got %d", len(executor.taskQueue))
	}

	if executor.taskQueue[0].TaskId != task.TaskId {
		t.Errorf("Expected task ID to be %d, got %d", task.TaskId, executor.taskQueue[0].TaskId)
	}
}

func TestCompleteTask(t *testing.T) {
	executor := &SyncTaskExecutor{
		completedTasks: make(map[int]struct{}),
	}

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

	if _, exists := executor.completedTasks[task.TaskId]; !exists {
		t.Fatalf("Expected task ID %d to be in completed tasks", task.TaskId)
	}

}

func TestFailTask(t *testing.T) {
	executor := &SyncTaskExecutor{
		completedTasks: make(map[int]struct{}),
	}

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

	if _, exists := executor.completedTasks[task.TaskId]; exists {
		t.Fatalf("Expected task ID %d to not be in completed tasks", task.TaskId)
	}
}

func TestRetryTask(t *testing.T) {
	executor := &SyncTaskExecutor{
		taskQueue:        []Task{},
		completedTasks:   make(map[int]struct{}),
		failureThreshold: 0, //must success
		retryCount:       3,
	}

	task := Task{TaskId: 1, RetryCount: 2} // Example task with RetryCount less than failureThreshold

	success, err := executor.retryTask(task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !success {
		t.Fatalf("Expected success to be true, got false")
	}
}

func TestExecuteTask_RetryScenario(t *testing.T) {
	executor := &SyncTaskExecutor{
		failureThreshold: 101, //must failed
		taskQueue:        []Task{},
		completedTasks:   make(map[int]struct{}),
		retryCount:       3,
	}

	resultChan := make(chan int)
	task := Task{
		TaskId:     1,
		RetryCount: 2,
		ResultChan: resultChan,
	}

	// Start a goroutine to receive the completion signal
	go func() {
		result := <-resultChan
		if result != 1 {
			t.Log("Task completed successfully.")
		} else {
			t.Log("Task failed.")
		}
	}()

	success, err := executor.executeTask(task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if success {
		t.Fatalf("Expected task to fail")
	}
}
