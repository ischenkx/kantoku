package plugins

import "runtime"

func Id() string {
	_, path, _, _ := runtime.Caller(2)
	return path
}
