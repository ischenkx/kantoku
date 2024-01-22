package events

func init() {
	OnTask.Created = "task:created"
	OnTask.Ready = "task:ready"
	OnTask.Received = "task:received"
	OnTask.Finished = "task:finished"
	OnTask.Cancelled = "task:cancelled"
}
