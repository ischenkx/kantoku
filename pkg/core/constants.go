package core

var TaskStatuses struct {
	Initialized string
	Ready       string
	Cancelled   string
	Received    string
	Finished    string
}

var TaskSubStatuses struct {
	OK     string
	Failed string
}

var OnTask struct {
	Created   string
	Ready     string
	Received  string
	Finished  string
	Cancelled string
}

var ResourceStatuses struct {
	DoesNotExist string
	Allocated    string
	Ready        string
}

func init() {
	TaskStatuses.Initialized = "initialized"
	TaskStatuses.Ready = "ready"
	TaskStatuses.Cancelled = "cancelled"
	TaskStatuses.Received = "received"
	TaskStatuses.Finished = "finished"

	TaskSubStatuses.OK = "ok"
	TaskSubStatuses.Failed = "failed"

	OnTask.Created = "task.created"
	OnTask.Ready = "task.ready"
	OnTask.Received = "task.received"
	OnTask.Finished = "task.finished"
	OnTask.Cancelled = "task.cancelled"

	ResourceStatuses.DoesNotExist = "does_not_exist"
	ResourceStatuses.Allocated = "allocated"
	ResourceStatuses.Ready = "ready"
}
