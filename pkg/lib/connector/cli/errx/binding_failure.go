package errx

import "fmt"

func FailedToBind(err error) error {
	return fmt.Errorf("failed to bind: %w", err)
}
