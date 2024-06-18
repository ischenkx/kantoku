package errx

import "fmt"

func FailedToBuild(entity string, err error) error {
	return fmt.Errorf("failed to build '%s': %w", entity, err)
}
