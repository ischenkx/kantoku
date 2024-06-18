package functional

import (
	"github.com/ischenkx/kantoku/pkg/lib/tasks/specification"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/specification/typing"
	"reflect"
	"strings"
	"time"
)

func ToSpecification[Input, Output any](task Task[Input, Output]) specification.Specification {
	var (
		zeroInput  Input
		zeroOutput Output
	)

	return specification.Specification{
		ID: task.ID,
		IO: specification.IO{
			Inputs:  structToResourceSet(any(zeroInput)),
			Outputs: structToResourceSet(zeroOutput),
		},
		Executable: specification.Executable{
			Type: "functional-go",
		},
	}
}

func structToResourceSet[T any](s T) specification.ResourceSet {
	resourceSet := specification.ResourceSet{
		Naming: make(map[int]string),
		Types:  make(map[int]typing.Type),
	}

	val := reflect.ValueOf(s)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		resourceSet.Naming[i] = field.Name

		fieldType := field.Type
		if fieldType.Kind() == reflect.Struct && strings.HasPrefix(fieldType.Name(), "Future[") && fieldType.PkgPath() == "github.com/ischenkx/kantoku/pkg/lib/tasks/functional/future" {
			// If the field type is future.Future, use its inner type
			typ, ok := fieldType.FieldByName("value")
			if !ok {
				panic("no value for future")
			}
			fieldType = typ.Type.Elem()
		}
		resourceSet.Types[i] = getType(fieldType)
	}

	return resourceSet
}

func getType(t reflect.Type) typing.Type {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return typing.Number()
	case reflect.String:
		return typing.String()
	case reflect.Bool:
		return typing.Boolean()
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			return typing.TimeStamp()
		} else if t == reflect.TypeOf(time.Duration(0)) {
			return typing.Duration()
		} else {
			fields := make([]typing.Field_, 0, t.NumField())
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				fields = append(fields, typing.Field(field.Name, getType(field.Type)))
			}
			return typing.Struct(fields...)
		}
	case reflect.Map:
		return typing.Map()
	case reflect.Slice, reflect.Array:
		return typing.Array(getType(t.Elem()))
	case reflect.Ptr:
		return getType(t.Elem())
	default:
		return typing.Unknown()
	}
}
