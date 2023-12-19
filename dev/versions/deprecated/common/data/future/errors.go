package future

import "errors"

var (
	ErrAlreadyResolved = errors.New("already resolved")
	ErrNotResolved     = errors.New("not resolved")
)
