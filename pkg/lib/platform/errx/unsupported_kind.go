package errx

import "fmt"

func UnsupportedKind(kind string) error {
	return fmt.Errorf("unsupported kind: %s", kind)
}
