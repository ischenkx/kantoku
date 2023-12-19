package record

import "fmt"

type Entry struct {
	Name  string
	Value any
}

func (e Entry) String() string {
	return fmt.Sprintf("'%s' => %s", e.Name, e.Value)
}
