package status

type Status string

const (
	Unknown   = "unknown"
	Scheduled = "scheduled"
	Received  = "received"
	Running   = "running"
	Finished  = "finished"
)
