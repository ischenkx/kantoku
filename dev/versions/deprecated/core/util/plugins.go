package util

func FilterPlugins[T any](plugins []any) (result []T) {
	for _, plugin := range plugins {
		if x, ok := plugin.(T); ok {
			result = append(result, x)
		}
	}
	return
}
