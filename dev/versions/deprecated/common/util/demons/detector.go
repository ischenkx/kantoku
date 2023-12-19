package demons

import (
	"context"
	infra2 "kantoku/framework/infra"
	"reflect"
	"unsafe"
)

type DemonDetector struct {
	source any
}

func NewDetector(source any) DemonDetector {
	return DemonDetector{source}
}

func (detector DemonDetector) Demons(ctx context.Context) []infra2.Demon {
	providers := collectProviders(detector.source, &visitor{})
	return Multi(providers).Demons()
}

func collectProviders(from any, visitor *visitor) (result []infra2.Provider) {
	if from == nil {
		return
	}
	if visitor.has(from) {
		return
	}
	visitor.add(from)

	if _, ok := from.(infra2.Provider); ok {
		result = append(result, from.(infra2.Provider))
	}

	value := reflect.ValueOf(from)

	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			result = append(result, collectProviders(value.Index(i).Interface(), visitor)...)
		}
	case reflect.Map:
		keys := value.MapKeys()
		for _, key := range keys {
			result = append(result, collectProviders(key.Interface(), visitor)...)
			result = append(result, collectProviders(value.MapIndex(key).Interface(), visitor)...)
		}
	case reflect.Struct:
		val := value
		if !val.CanAddr() {
			addressableValue := reflect.New(val.Type()).Elem()
			addressableValue.Set(val)
			val = addressableValue
		}
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
			result = append(result, collectProviders(field.Interface(), visitor)...)
		}
	}
	return
}

type visitor struct {
	data []any
}

func (v *visitor) add(x any) {
	v.data = append(v.data, x)
}

func (v *visitor) has(x any) bool {
	for _, y := range v.data {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}
