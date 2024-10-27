package main

type ServerConfig struct {
	Mode             string
	FailureThreshold int
}

type Server struct {
	Config       ServerConfig
	TaskExecutor TaskExecutor
}

type Task struct {
	TaskId     int
	RetryCount int
	ResultChan chan int // 1: success, 2: fail, 3: repeated
}

type TaskExecutor interface {
	Start(server Server) (bool, error)
	SubmitTask(task Task) (bool, error)
	scheduleTask(task Task)
	processTasks()
	executeTask(task Task) (bool, error)
	completeTask(task Task) (bool, error)
	retryTask(task Task) (bool, error)
	failTask(task Task) (bool, error)
}
