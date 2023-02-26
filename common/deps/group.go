package deps

type Group struct {
	ID           string
	Dependencies map[string]bool
}
