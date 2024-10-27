package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var server Server

func StartHttpServer() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./program <async/sync> <failure threshold>")
		return
	}

	mode := os.Args[1]
	number := os.Args[2]

	serverConfig, err := initConfig(mode, number)
	if err != nil {
		log.Fatal(err)
	}
	var executor TaskExecutor
	if serverConfig.Mode == "async" {
		executor = &AsyncTaskExecutor{
			taskQueue:        make(chan Task, 100),
			completedTasks:   make(map[int]struct{}),
			failureThreshold: serverConfig.FailureThreshold,
			retryCount:       3,
			wg:               sync.WaitGroup{},
			mu:               sync.RWMutex{},
		}
	} else if serverConfig.Mode == "sync" {
		executor = &SyncTaskExecutor{
			taskQueue:        make([]Task, 0),
			completedTasks:   make(map[int]struct{}),
			failureThreshold: serverConfig.FailureThreshold,
			stopChan:         make(chan struct{}),
			retryCount:       3,
		}
	}

	db, err := startDB()
	if err != nil {
		log.Fatal(err)
	}

	server = Server{
		Config:       serverConfig,
		TaskExecutor: executor,
		DB:           db,
	}

	server.TaskExecutor.Start()
	processHttpRequests()
}

func initConfig(mode string, number string) (ServerConfig, error) {
	threshold, err := strconv.Atoi(number)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid failure threshold '%s': %v", number, err)
	}
	return ServerConfig{
		Mode:             mode,
		FailureThreshold: threshold,
	}, nil
}

func processHttpRequests() {
	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		taskHandler(w, r)
	})

	port := "8080"
	log.Printf("Starting server in %s mode on port %s\n", server.Config.Mode, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var taskReq struct {
		RequestID string `json:"request_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&taskReq)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(taskReq.RequestID)
	if err != nil {
		http.Error(w, "Invalid TaskID", http.StatusBadRequest)
		return
	}

	task := Task{
		TaskId:     taskID,
		RetryCount: 0,
		ResultChan: make(chan int),
	}

	if _, err := server.TaskExecutor.SubmitTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Wait for task completion or failure
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		result := <-task.ResultChan
		defer close(task.ResultChan)

		switch result {
		case 1: // 1 is success
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Task %d completed successfully\n", taskID)
		case 2: // 2 is failure
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Task %d failed\n", taskID)
		case 3: // 3 is queue full in async mode
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Task %d is not accepted due to queue full\n", taskID)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Task %d returned an unknown result\n", taskID)
		}
		wg.Done()
		//TODO: close the result channel; might add the thread id but dont think it is necessary
	}()
	wg.Wait()
}