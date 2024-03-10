package ops

import (
	"github.com/samber/lo"
)

const (
	InOp        = "in"
	EqOp        = "eq"
	LessOp      = "lt"
	LessOrEqOp  = "le"
	GreaterOp   = "gt"
	GreaterOrEq = "ge"
	LikeOp      = "like"
	ContainsOp  = "contains"
	NotOp       = "not"
	AndOp       = "and"
	OrOp        = "or"
)

func Eq[T any](item T) Operation {
	return Operation{
		Type: EqOp,
		Data: item,
	}
}

func In[T any](items ...T) Operation {
	return Operation{
		Type: InOp,
		Data: lo.ToAnySlice(items),
	}
}

func LessThan(value any) Operation {
	return Operation{
		Type: LessOp,
		Data: value,
	}
}

func LessOrEqualThan(value any) Operation {
	return Operation{
		Type: LessOrEqOp,
		Data: value,
	}
}

func GreaterThan(value any) Operation {
	return Operation{
		Type: GreaterOp,
		Data: value,
	}
}

func GreaterOrEqualThan(value any) Operation {
	return Operation{
		Type: GreaterOrEq,
		Data: value,
	}
}

func Like(value any) Operation {
	return Operation{
		Type: LikeOp,
		Data: value,
	}
}

func Contains[T any](items ...T) Operation {
	return Operation{
		Type: ContainsOp,
		Data: lo.ToAnySlice(items),
	}
}
func Not(expr any) Operation {
	return Operation{
		Type: NotOp,
		Data: expr,
	}
}
func And(exprs ...any) Operation {
	return Operation{
		Type: AndOp,
		Data: exprs,
	}
}
func Or(exprs ...any) Operation {
	return Operation{
		Type: OrOp,
		Data: exprs,
	}
}
