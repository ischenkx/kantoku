package record

type Mask struct {
	Operation       string
	PropertyPattern string
}

const (
	IncludeMask = "include"
	ExcludeMask = "exclude"
)

func Include(prop string) Mask {
	return Mask{
		PropertyPattern: prop,
		Operation:       IncludeMask,
	}
}

func Exclude(prop string) Mask {
	return Mask{
		PropertyPattern: prop,
		Operation:       ExcludeMask,
	}
}
