package argument

import "context"

type SelfEncoder interface {
	Encode(ctx context.Context) (Argument, error)
}

type Codec interface {
	Encode(ctx context.Context, arg any) (Argument, error)
	Decode(ctx context.Context, data Argument) (any, error)
}
