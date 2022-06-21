package core

const (
	StateUnknown       string = "UNKNOWN"
	StateQueued        string = "QUEUED"
	StateInitalizing   string = "INITIALIZING"
	StateRunning       string = "RUNNING"
	StatePaused        string = "PAUSED"
	StateComplete      string = "COMPLETE"
	StateExecutorError string = "EXECUTOR_ERROR"
	StateSystemError   string = "SYSTEM_ERROR"
	StateCanceled      string = "CANCELED"
	StateCanceling     string = "CANCELING"
)
