package main

import (
	"log"
)

// SyncTaskExecutor struct definition
type SyncTaskExecutor struct {
	TaskQueue      map[int]struct{}
	CompletedTasks map[int]struct{}
}

// Start begins executing tasks synchronously.
func (executor *SyncTaskExecutor) Start() (bool, error) {
	log.Println("Starting task execution...")
	log.Println("Task execution started.")
	return true, nil
}

// SubmitTask implements TaskExecutor.
func (executor *SyncTaskExecutor) SubmitTask(task Task) (bool, error) {
	log.Printf("Task ID: %d submitted.\n", task.TaskId)
	return true, nil
}

func (executor *SyncTaskExecutor) ScheduleTask(task Task) (bool, error) {
	log.Printf("Task ID: %d scheduled.\n", task.TaskId)
	return true, nil
}

func (executor *SyncTaskExecutor) ExecuteTask(task Task) (bool, error) {
	executor.CompletedTasks[task.TaskId] = struct{}{}
	log.Printf("Task ID: %d completed.\n", task.TaskId)
	return true, nil
}

func (executor *SyncTaskExecutor) CompleteTask(task Task) (bool, error) {
	executor.CompletedTasks[task.TaskId] = struct{}{}
	log.Printf("Task ID: %d completed.\n", task.TaskId)
	return true, nil
}

func (executor *SyncTaskExecutor) RetryTask(task Task) (bool, error) {
	log.Printf("Task ID: %d retried.\n", task.TaskId)
	return true, nil
}

func (executor *SyncTaskExecutor) FailTask(task Task) (bool, error) {
	log.Printf("Task ID: %d failed.\n", task.TaskId)
	return true, nil
}
