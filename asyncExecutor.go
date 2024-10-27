package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
Async executor will process tasks concurrently in multiple goroutines.
Buffered channel is used to queue tasks.
Map is used to track completed tasks.
RWMutex is used to synchronize access to the completed tasks map.
*/

type AsyncTaskExecutor struct {
	taskQueue        chan Task
	completedTasks   map[int]struct{} //replace by db. leave here for showing original implementation : will be accessed by multiple goroutines, RWMutex is used to synchronize access
	failureThreshold int
	retryCount       int
	wg               sync.WaitGroup
	mu               sync.RWMutex
}

func (executor *AsyncTaskExecutor) Start() (bool, error) {
	executor.failureThreshold = server.Config.FailureThreshold

	go executor.processTasks() //start the task processing loop

	return true, nil
}

func (executor *AsyncTaskExecutor) SubmitTask(task Task) (bool, error) {
	log.Printf("SubmitTask triggered for Task ID: %d at %s\n", task.TaskId, time.Now().Format(time.RFC3339))
	// executor.mu.RLock()
	// defer executor.mu.RUnlock()

	// if _, ok := executor.completedTasks[task.TaskId]; ok {
	if ifTaskCompleted(task.TaskId) {
		log.Printf("Task ID: %d already completed.\n", task.TaskId)
		return false, fmt.Errorf("Task ID: %d already completed.", task.TaskId)
	}
	go executor.scheduleTask(task) //concurrently schedule the task
	return true, nil
}

func (executor *AsyncTaskExecutor) scheduleTask(task Task) {
	log.Printf("ScheduleTask triggered for Task ID: %d at %s\n", task.TaskId, time.Now().Format(time.RFC3339))
	executor.taskQueue <- task //add task to the buffered channel
	return
}

func (executor *AsyncTaskExecutor) processTasks() {
	for {
		var task Task
		select {
		case task = <-executor.taskQueue:
			go executor.executeTask(task)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (executor *AsyncTaskExecutor) executeTask(task Task) (bool, error) {
	gid := getGID()
	log.Printf("ExecuteTask triggered for Task ID: %d at %s by Goroutine %d\n", task.TaskId, time.Now().Format(time.RFC3339), gid)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomValue := r.Float64() * 100

	// Determine if the task execution is successful based on the threshold
	if randomValue > float64(executor.failureThreshold) {
		go func() {
			time.Sleep(1 * time.Second)
			executor.completeTask(task)
		}()
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

func (executor *AsyncTaskExecutor) completeTask(task Task) (bool, error) {
	// executor.mu.Lock()         // Acquire write lock
	// defer executor.mu.Unlock() // Ensure the lock is released
	// executor.completedTasks[task.TaskId] = struct{}{}
	addTaskToDB(task.TaskId)
	log.Printf("CompleteTask triggered for Task ID: %d at %s\n", task.TaskId, time.Now().Format(time.RFC3339))

	task.ResultChan <- 1
	return true, nil
}

func (executor *AsyncTaskExecutor) retryTask(task Task) (bool, error) {
	log.Printf("Task ID: %d retried.\n", task.TaskId)
	return executor.executeTask(task)
}

func (executor *AsyncTaskExecutor) failTask(task Task) (bool, error) {
	log.Printf("Task ID: %d failed.\n", task.TaskId)

	task.ResultChan <- 2
	return true, nil
}

func getGID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	stack := strings.TrimPrefix(string(buf[:n]), "goroutine ")
	i := strings.Index(stack, " ")
	if i < 0 {
		return 0
	}
	id, err := strconv.ParseUint(stack[:i], 10, 64)
	if err != nil {
		return 0
	}
	return id
}
