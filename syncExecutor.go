package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

/*
Sync executor will process tasks sequentially in a single thread.
Slice is used to queue tasks.
SQLite is used to persist completed tasks.
stopChan is used to signal the 'processTasks' loop to allocate new task to be executed.
*/

type SyncTaskExecutor struct {
	taskQueue        []Task
	completedTasks   map[int]struct{} //replaced by db. leave here for showing original implementation
	failureThreshold int
	stopChan         chan struct{}
	retryCount       int
}

func (executor *SyncTaskExecutor) Start() (bool, error) {
	executor.failureThreshold = server.Config.FailureThreshold

	executor.stopChan = make(chan struct{})
	go executor.processTasks() //start the task processing loop

	return true, nil
}

func (executor *SyncTaskExecutor) SubmitTask(task Task) (bool, error) {
	//check if the task is already completed
	log.Printf("SubmitTask triggered for Task ID: %d at %s\n", task.TaskId, time.Now().Format(time.RFC3339))
	if ifTaskCompleted(task.TaskId) {
		log.Printf("Task ID: %d already completed.\n", task.TaskId)
		return false, fmt.Errorf("Task ID: %d already completed.", task.TaskId)
	}
	executor.scheduleTask(task)
	return true, nil
}

func (executor *SyncTaskExecutor) scheduleTask(task Task) {
	//add task to the "todo queue"
	executor.taskQueue = append(executor.taskQueue, task)
	return
}

func (executor *SyncTaskExecutor) processTasks() {
	for {
		select {
		case <-executor.stopChan:
			task := executor.taskQueue[0]
			executor.taskQueue = executor.taskQueue[1:]
			executor.executeTask(task)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (executor *SyncTaskExecutor) executeTask(task Task) (bool, error) {
	gid := getGID()
	log.Printf("ExecuteTask triggered for Task ID: %d at %s by Goroutine %d\n", task.TaskId, time.Now().Format(time.RFC3339), gid)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomValue := r.Float64() * 100

	//determine if the task execution is successful based on the threshold
	if randomValue > float64(executor.failureThreshold) {
		time.Sleep(1 * time.Second)
		executor.completeTask(task)
		return true, nil
	} else {
		task.RetryCount++
		if task.RetryCount <= executor.retryCount {
			executor.retryTask(task)
		} else {
			executor.failTask(task)
		}
		return false, nil
	}
}

func (executor *SyncTaskExecutor) completeTask(task Task) (bool, error) {
	log.Printf("CompleteTask triggered for Task ID: %d at %s\n", task.TaskId, time.Now().Format(time.RFC3339))
	addTaskToDB(task.TaskId)

	if task.ResultChan != nil {
		task.ResultChan <- 1 //1 indicates success
	}
	return true, nil
}

func (executor *SyncTaskExecutor) retryTask(task Task) (bool, error) {
	log.Printf("Task ID: %d retried.\n", task.TaskId)
	return executor.executeTask(task)
}

func (executor *SyncTaskExecutor) failTask(task Task) (bool, error) {
	log.Printf("Task ID: %d failed.\n", task.TaskId)

	if task.ResultChan != nil {
		task.ResultChan <- 2 //2 indicates failure
	}
	return true, nil
}
