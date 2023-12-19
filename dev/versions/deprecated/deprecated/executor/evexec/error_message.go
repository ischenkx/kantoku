package evexec

import "encoding/json"

type ErrorMessage struct {
	TaskID  string `json:"task_id,omitempty"`
	Message string `json:"message,omitempty"`
}

func (message ErrorMessage) Encode() ([]byte, error) {
	return json.Marshal(message)
}
