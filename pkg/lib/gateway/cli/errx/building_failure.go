package errx

import "fmt"

func FailedToBuild(entity string, err error) error {
	return fmt.Errorf("failed to build (entity='%s'): %w", entity, err)
}
