package task

var Statuses struct {
	Initialized string
	Ready       string
	Cancelled   string
	Received    string
	Finished    string
}

var SubStatuses struct {
	OK     string
	Failed string
}

func init() {
	Statuses.Initialized = "initialized"
	Statuses.Ready = "ready"
	Statuses.Cancelled = "cancelled"
	Statuses.Received = "received"
	Statuses.Finished = "finished"

	SubStatuses.OK = "ok"
	SubStatuses.Failed = "failed"
}
