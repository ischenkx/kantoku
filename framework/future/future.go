package future

type ID = string

type Future struct {
	ID    ID
	Type  string
	Param []byte
}
