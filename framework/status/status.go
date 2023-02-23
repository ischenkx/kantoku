package status

type Status string

const (
	Blocked   Status = "blocked"
	Pending          = "pending"
	Executing        = "executing"
	Executed         = "executed"
	Complete         = "Complete"
)
