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
	Completed  bool
	ResultChan chan bool
}

type TaskExecutor interface {
	Start(server Server) (bool, error)
	SubmitTask(task Task) (bool, error)
	scheduleTask(task Task) (bool, error)
	processTasks()
	popTask() (Task, error)
	executeTask(task Task) (bool, error)
	completeTask(task Task) (bool, error)
	retryTask(task Task) (bool, error)
	failTask(task Task) (bool, error)
}

type todoTask interface {
	Do()
}
