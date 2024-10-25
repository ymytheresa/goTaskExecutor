package main

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

type SyncTaskExecutor struct {
	taskQueue        []Task
	completedTasks   map[int]struct{}
	failureThreshold int
	stopChan         chan struct{}
	retryCount       int
}

func (executor *SyncTaskExecutor) Start(server Server) (bool, error) {
	executor.failureThreshold = server.Config.FailureThreshold
	executor.completedTasks = make(map[int]struct{})
	// TODO: Read completed tasks from DB

	executor.stopChan = make(chan struct{})
	go executor.processTasks()

	return true, nil
}

func (executor *SyncTaskExecutor) SubmitTask(task Task) (bool, error) {
	log.Printf("SubmitTask triggered for Task ID: %d at %s\n", task.TaskId, time.Now().Format(time.RFC3339))
	if _, ok := executor.completedTasks[task.TaskId]; ok {
		log.Printf("Task ID: %d already completed.\n", task.TaskId)
		return false, errors.New("task already completed")
	}
	executor.scheduleTask(task)
	return true, nil
}

func (executor *SyncTaskExecutor) scheduleTask(task Task) (bool, error) {
	log.Printf("ScheduleTask triggered for Task ID: %d at %s\n", task.TaskId, time.Now().Format(time.RFC3339))
	executor.taskQueue = append(executor.taskQueue, task)
	return true, nil
}

func (executor *SyncTaskExecutor) processTasks() {
	log.Println("ProcessTasks triggered at", time.Now().Format(time.RFC3339))
	for {
		select {
		case <-executor.stopChan:
			log.Println("Task processing stopped.")
			return
		default:
			if len(executor.taskQueue) > 0 {
				task, err := executor.popTask()
				if err != nil {
					log.Println("Error popping task:", err)
					continue
				}
				executor.executeTask(task)
			} else {
				// No tasks available, sleep briefly before checking again
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (executor *SyncTaskExecutor) popTask() (Task, error) {
	log.Println("PopTask triggered at", time.Now().Format(time.RFC3339))
	if len(executor.taskQueue) == 0 {
		return Task{}, errors.New("no tasks in queue")
	}
	task := executor.taskQueue[0]
	executor.taskQueue = executor.taskQueue[1:]
	return task, nil
}

func (executor *SyncTaskExecutor) executeTask(task Task) (bool, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomValue := r.Float64() * 100

	// Determine if the task execution is successful based on the threshold
	if randomValue > float64(executor.failureThreshold) {
		executor.completedTasks[task.TaskId] = struct{}{}
		executor.completeTask(task)
		return true, nil
	} else {
		task.RetryCount++
		if task.RetryCount <= executor.failureThreshold {
			executor.retryTask(task)
		} else {
			executor.failTask(task)
		}
		return false, nil
	}
}

func (executor *SyncTaskExecutor) completeTask(task Task) (bool, error) {
	log.Printf("CompleteTask triggered for Task ID: %d at %s\n", task.TaskId, time.Now().Format(time.RFC3339))
	executor.completedTasks[task.TaskId] = struct{}{}
	log.Printf("Task ID: %d completed.\n", task.TaskId)

	if task.ResultChan != nil {
		task.Completed = true
		task.ResultChan <- true
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
		task.Completed = false
		task.ResultChan <- false
	}
	return true, nil
}
