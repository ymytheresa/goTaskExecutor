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
	success, err := executor.scheduleTask(task)

	if !success || err != nil {
		t.Errorf("Expected success to be true and err to be nil, got success: %v, err: %v", success, err)
	}

	if len(executor.taskQueue) != 1 {
		t.Errorf("Expected taskQueue length to be 1, got %d", len(executor.taskQueue))
	}

	if executor.taskQueue[0].TaskId != task.TaskId {
		t.Errorf("Expected task ID to be %d, got %d", task.TaskId, executor.taskQueue[0].TaskId)
	}
}

func TestPopTask(t *testing.T) {
	executor := &SyncTaskExecutor{
		taskQueue: []Task{
			{TaskId: 1},
			{TaskId: 2},
		},
	}

	task, err := executor.popTask()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if task.TaskId != 1 {
		t.Fatalf("Expected TaskId 1, got %d", task.TaskId)
	}

	task, err = executor.popTask()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if task.TaskId != 2 {
		t.Fatalf("Expected TaskId 2, got %d", task.TaskId)
	}

	_, err = executor.popTask()
	if err == nil {
		t.Fatal("Expected error, got none")
	}
	expectedError := "no tasks in queue"
	if err.Error() != expectedError {
		t.Fatalf("Expected '%s' error, got %v", expectedError, err)
	}
}

func TestCompleteTask(t *testing.T) {
	executor := &SyncTaskExecutor{
		completedTasks: make(map[int]struct{}),
	}

	resultChan := make(chan bool)
	task := Task{
		TaskId:     1,
		ResultChan: resultChan,
	}

	go func() {
		result := <-resultChan
		if !result {
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

	resultChan := make(chan bool)
	task := Task{
		TaskId:     1,
		ResultChan: resultChan,
	}

	go func() {
		result := <-resultChan
		if result {
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
