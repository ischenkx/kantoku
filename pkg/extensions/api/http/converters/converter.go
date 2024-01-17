package converters

import (
	"github.com/ischenkx/kantoku/pkg/extensions/api/http/oas"
	"github.com/ischenkx/kantoku/pkg/system/kernel/task"
)

func TaskToDto(t task.Task) oas.Task {
	return oas.Task{
		Id:         t.ID,
		Inputs:     t.Inputs,
		Outputs:    t.Outputs,
		Properties: PropertiesToDto(t.Properties),
	}
}

func PropertiesToDto(props task.Properties) oas.TaskProperties {
	dtoProperties := oas.TaskProperties{
		Data: map[string]any{},
	}

	for key, value := range props.Data {
		dtoProperties.Data[key] = value
	}

	for key, sub := range props.Sub {
		dtoProperties.Sub = append(dtoProperties.Sub, oas.PropertySubTree{
			Key:   key,
			Value: PropertiesToDto(sub),
		})
	}

	return dtoProperties
}

func DtoToProperties(dto oas.TaskProperties) task.Properties {
	props := task.Properties{
		Data: map[string]any{},
		Sub:  map[string]task.Properties{},
	}

	for key, value := range dto.Data {
		props.Data[key] = value
	}

	for _, sub := range dto.Sub {
		props.Sub[sub.Key] = DtoToProperties(sub.Value)
	}

	return props
}
