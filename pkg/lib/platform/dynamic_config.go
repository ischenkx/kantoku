package platform

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"os"
	"reflect"
	"strings"
	"time"
)

type DynamicConfig map[string]any

func (config DynamicConfig) Kind() string {
	rawKind, ok := config["kind"]
	if !ok {
		return ""
	}

	kind, ok := rawKind.(string)
	if !ok {
		return ""
	}

	return kind
}

func (config DynamicConfig) Bind(to any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			ParseEnvVariableHookFunc(),
			ToDurationHookFunc(),
		),
		Result: to,
	})
	if err != nil {
		return fmt.Errorf("failed to create a decoder: %w", err)
	}
	return decoder.Decode(config)
}

func ToDurationHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Duration(0)) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.ParseDuration(data.(string))
		case reflect.Float64:
			return time.Duration(int64(data.(float64)) * int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Duration(data.(int64) * int64(time.Millisecond)), nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}

func ParseEnvVariableHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		//if t != reflect.TypeOf("") {
		//	return data, nil
		//}

		switch f.Kind() {
		case reflect.String:
			val := data.(string)
			if strings.HasPrefix(val, "$") {
				name := val[1:]
				if _var, ok := os.LookupEnv(name); ok {
					return _var, nil
				} else {
					return nil, fmt.Errorf("environment variable lookup failed (name='%s')", name)
				}
			}
		}
		return data, nil
	}
}
