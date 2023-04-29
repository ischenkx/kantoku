package rebimap

import (
	"context"
	"github.com/redis/go-redis/v9"
	"kantoku/common/codec"
	"kantoku/common/util"
)

type Bimap[K, V any] struct {
	keyLabel, valueLabel string
	keyCodec             codec.Codec[K, []byte]
	valueCodec           codec.Codec[V, []byte]
	client               redis.UniversalClient
}

func NewBimap[K, V any](keyLabel, valueLabel string, keyCodec codec.Codec[K, []byte], valueCodec codec.Codec[V, []byte], client redis.UniversalClient) *Bimap[K, V] {
	return &Bimap[K, V]{
		keyLabel:   keyLabel,
		valueLabel: valueLabel,
		keyCodec:   keyCodec,
		valueCodec: valueCodec,
		client:     client,
	}
}

func (bimap *Bimap[K, V]) Save(ctx context.Context, key K, value V) error {
	encodedKey, err := bimap.keyCodec.Encode(key)
	if err != nil {
		return err
	}

	encodedValue, err := bimap.valueCodec.Encode(value)
	if err != nil {
		return err
	}

	return bimap.client.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.HSet(ctx, bimap.keyLabel, string(encodedKey), encodedValue).Result()
		if err != nil {
			return err
		}

		_, err = tx.HSet(ctx, bimap.valueLabel, string(encodedValue), encodedKey).Result()
		return err
	})
}

func (bimap *Bimap[K, V]) DeleteByKey(ctx context.Context, key K) error {
	encodedKey, err := bimap.keyCodec.Encode(key)
	if err != nil {
		return err
	}

	return bimap.client.Watch(ctx, func(tx *redis.Tx) error {
		encodedValue, err := tx.HGet(ctx, bimap.keyLabel, string(encodedKey)).Result()
		if err != nil {
			return err
		}

		_, err = tx.HDel(ctx, bimap.keyLabel, string(encodedKey)).Result()
		if err != nil {
			return err
		}

		_, err = tx.HDel(ctx, bimap.valueLabel, encodedValue).Result()
		return err
	})
}

func (bimap *Bimap[K, V]) DeleteByValue(ctx context.Context, value V) error {
	encodedValue, err := bimap.valueCodec.Encode(value)
	if err != nil {
		return err
	}

	return bimap.client.Watch(ctx, func(tx *redis.Tx) error {
		encodedKey, err := tx.HGet(ctx, bimap.valueLabel, string(encodedValue)).Result()
		if err != nil {
			return err
		}

		_, err = tx.HDel(ctx, bimap.keyLabel, encodedKey).Result()
		if err != nil {
			return err
		}

		_, err = tx.HDel(ctx, bimap.valueLabel, string(encodedValue)).Result()
		return err
	})
}

func (bimap *Bimap[K, V]) ByValue(ctx context.Context, value V) (K, error) {
	encodedValue, err := bimap.valueCodec.Encode(value)
	if err != nil {
		return util.Default[K](), err
	}
	encodedKey, err := bimap.client.HGet(ctx, bimap.valueLabel, string(encodedValue)).Result()
	return bimap.keyCodec.Decode([]byte(encodedKey))
}

func (bimap *Bimap[K, V]) ByKey(ctx context.Context, key K) (V, error) {
	encodedKey, err := bimap.keyCodec.Encode(key)
	if err != nil {
		return util.Default[V](), err
	}
	encodedValue, err := bimap.client.HGet(ctx, bimap.keyLabel, string(encodedKey)).Result()
	return bimap.valueCodec.Decode([]byte(encodedValue))
}
