package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func startHttpServer() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./program <async/sync> <failure threshold>")
		return
	}

	mode := os.Args[1]
	number := os.Args[2]

	serverConfig := initConfig(mode, number)
	var executor TaskExecutor
	if serverConfig.Mode == "async" {
		executor = &AsyncTaskExecutor{
			TaskChan:       make(chan Task),
			CompletedTasks: make(map[int]struct{}),
			WaitGroup:      &sync.WaitGroup{},
			Mu:             &sync.RWMutex{},
		}
	} else if serverConfig.Mode == "sync" {
		executor = &SyncTaskExecutor{
			TaskQueue:      make(map[int]struct{}),
			CompletedTasks: make(map[int]struct{}),
		}
	}

	server := Server{
		Config:       serverConfig,
		TaskExecutor: executor,
	}

	server.TaskExecutor.Start() // Start the executor
	processHttpRequests(server)
}

func initConfig(mode string, number string) ServerConfig {
	threshold, err := strconv.Atoi(number)
	if err != nil {
		fmt.Printf("Invalid failure threshold '%s': %v\n", number, err)
		os.Exit(1)
	}
	return ServerConfig{
		Mode:             mode,
		FailureThreshold: threshold,
	}
}

func processHttpRequests(server Server) {
	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		taskHandler(w, r, server)
	})

	port := "8080"
	log.Printf("Starting server in %s mode on port %s\n", server.Config.Mode, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request, server Server) {
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
		TaskId:      taskID,
		RetryCount:  0,
		Completed:   false,
		FailureRate: 0,
	}

	if _, err := server.TaskExecutor.SubmitTask(task); err != nil {
		http.Error(w, "Failed to add task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Task %d accepted for processing", taskID)
}
