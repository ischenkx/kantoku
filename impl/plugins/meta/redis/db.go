package redismeta

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"kantoku/framework/plugins/meta"
	"kantoku/impl/common/codec/bincodec"
	redikv "kantoku/impl/common/data/kv/redis"
)

type DB struct {
	prefix string
	client redis.UniversalClient
}

func NewDB(prefix string, client redis.UniversalClient) DB {
	return DB{
		prefix: prefix,
		client: client,
	}
}

func (db DB) Get(ctx context.Context, id string) (meta.RawMeta, error) {
	setName := fmt.Sprintf("%s_%s", db.prefix, id)
	return redikv.New[[]byte](db.client, bincodec.Codec{}, setName), nil
}
