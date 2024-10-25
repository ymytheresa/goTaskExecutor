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
	TaskId      int
	RetryCount  int
	Completed   bool
	FailureRate int
}

type TaskExecutor interface {
	Start() (bool, error)
	SubmitTask(task Task) (bool, error)
	ScheduleTask(task Task) (bool, error)
	ExecuteTask(task Task) (bool, error)
	CompleteTask(task Task) (bool, error)
	RetryTask(task Task) (bool, error)
	FailTask(task Task) (bool, error)
}

type todoTask interface {
	Do()
}
