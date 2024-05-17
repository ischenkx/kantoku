package typing

type Type struct {
	Name     string
	SubTypes map[string]Type
}

type Field_ struct {
	Name  string
	Value Type
}

func Field(name string, value Type) Field_ {
	return Field_{Name: name, Value: value}
}

func Number() Type {
	return Type{
		Name: "number",
	}
}

func String() Type {
	return Type{
		Name: "string",
	}
}

func Boolean() Type {
	return Type{
		Name: "boolean",
	}
}

func Struct(fields ...Field_) Type {
	fieldsMap := make(map[string]Type)
	for _, field := range fields {
		fieldsMap[field.Name] = field.Value
	}
	return Type{
		Name:     "struct",
		SubTypes: fieldsMap,
	}
}

func Map() Type {
	return Type{
		Name: "map",
	}
}

func Array(item Type) Type {
	return Type{Name: "array", SubTypes: map[string]Type{"item": item}}
}

func Ref(name string) Type {
	return Type{
		Name: "ref",
		SubTypes: map[string]Type{
			"name": {Name: name},
		},
	}
}

func Unknown() Type {
	return Type{Name: "unknown", SubTypes: make(map[string]Type)}
}
